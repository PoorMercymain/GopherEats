package handler

import (
	"context"
	"errors"
	"fmt"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/IBM/sarama"

	"github.com/PoorMercymain/GopherEats/internal/app/subscription/domain"
	subErrors "github.com/PoorMercymain/GopherEats/internal/app/subscription/errors"
	"github.com/PoorMercymain/GopherEats/internal/pkg/logger"
	"github.com/PoorMercymain/GopherEats/pkg/api/auth"
	api "github.com/PoorMercymain/GopherEats/pkg/api/subscription"
)

var _ api.SubscriptionV1Server = (*subscription)(nil)

type subscription struct {
	srv                      domain.SubscriptionService
	client                   auth.AuthV1Client
	emailSender              smtpSender
	weekNumber               *int
	kafkaProducer            sarama.SyncProducer
	notEnoughFundsEmailsChan chan string
	api.UnimplementedSubscriptionV1Server
}

// New returns pointer to new instance of subscription struct with service, auth and Kafka producer.
func New(srv domain.SubscriptionService, client auth.AuthV1Client, kafkaProducer sarama.SyncProducer, weekNumber *int, smtpUsername string, smtpPassword string, smtpServer string, smtpPort string) *subscription {
	return &subscription{srv: srv, client: client, kafkaProducer: kafkaProducer, weekNumber: weekNumber,
		emailSender: smtpSender{username: smtpUsername, password: smtpPassword, server: smtpServer, port: smtpPort}, notEnoughFundsEmailsChan: make(chan string, 1)}
}

// CreateSubscriptionV1 creates new subscription.
func (h *subscription) CreateSubscriptionV1(ctx context.Context, r *api.CreateSubscriptionRequestV1) (*emptypb.Empty, error) {
	err := h.srv.CreateSubscription(ctx, r.Email, r.BundleId)

	if errors.Is(err, subErrors.ErrorUniqueViolationWhileCreating) {
		return &emptypb.Empty{}, status.Errorf(codes.AlreadyExists, "the user already have a subscription, to change it, use another endpoint")
	}

	if err != nil {
		return &emptypb.Empty{}, status.Errorf(codes.Internal, "something went wrong in subscription service: %v", err)
	}

	return &emptypb.Empty{}, nil
}

// ReadSubscriptionV1 returns subscription info.
func (h *subscription) ReadSubscriptionV1(ctx context.Context, r *api.ReadSubscriptionRequestV1) (*api.ReadSubscriptionResponseV1, error) {
	bundleID, isDeleted, err := h.srv.ReadSubscription(ctx, r.Email)

	if errors.Is(err, subErrors.ErrorNoRowsWhileReading) {
		return nil, status.Errorf(codes.NotFound, subErrors.ErrorNoRowsWhileReading.Error())
	}

	if err != nil {
		return nil, status.Errorf(codes.Internal, "something went wrong in subscription service: %v", err)
	}

	return &api.ReadSubscriptionResponseV1{BundleId: bundleID, IsDeleted: isDeleted}, nil
}

// ChangeSubscriptionV1 updates subscription.
func (h *subscription) ChangeSubscriptionV1(ctx context.Context, r *api.ChangeSubscriptionRequestV1) (*emptypb.Empty, error) {
	err := h.srv.UpdateSubscription(ctx, r.Email, r.BundleId, r.IsDeleted)

	if errors.Is(err, subErrors.ErrorNoRowsUpdated) {
		return nil, status.Errorf(codes.NotFound, subErrors.ErrorNoRowsUpdated.Error())
	}

	if err != nil {
		return nil, status.Errorf(codes.Internal, "something went wrong in subscription service: %v", err)
	}

	return &emptypb.Empty{}, nil
}

// CancelSubscriptionV1 allows to cancel subscription.
func (h *subscription) CancelSubscriptionV1(ctx context.Context, r *api.CancelSubscriptionRequestV1) (*emptypb.Empty, error) {
	err := h.srv.CancelSubscription(ctx, r.Email)

	if errors.Is(err, subErrors.ErrorNoRowsUpdated) {
		return nil, status.Errorf(codes.NotFound, subErrors.ErrorNoRowsUpdated.Error())
	}

	if err != nil {
		return nil, status.Errorf(codes.Internal, "something went wrong in subscription service: %v", err)
	}

	return &emptypb.Empty{}, nil
}

// AddBalanceV1 adds funds to balance.
func (h *subscription) AddBalanceV1(ctx context.Context, r *api.AddBalanceRequestV1) (*emptypb.Empty, error) {
	err := h.srv.AddBalance(ctx, r.Email, r.Amount)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "something went wrong in subscription service: %v", err)
	}

	return &emptypb.Empty{}, nil
}

// ReadUserDataV1 returns user data.
func (h *subscription) ReadUserDataV1(ctx context.Context, r *api.ReadUserDataRequestV1) (*api.ReadUserDataResponseV1, error) {
	addressResp, err := h.client.GetAddressV1(ctx, &auth.GetAddressRequestV1{Email: r.Email})
	if err != nil {
		return nil, err
	}

	userData, err := h.srv.ReadUserData(ctx, r.Email)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "something went wrong in subscription service: %v", err)
	}

	return &api.ReadUserDataResponseV1{Address: addressResp.Address, BundleId: userData.BundleID, Balance: userData.Balance}, nil
}

// ReadBalanceHistoryV1 returns history of balance funding and purchases.
func (h *subscription) ReadBalanceHistoryV1(ctx context.Context, r *api.ReadBalanceHistoryRequestV1) (*api.ReadBalanceHistoryResponseV1, error) {
	history, err := h.srv.ReadBalanceHistory(ctx, r.Email, r.Page)
	if errors.Is(err, subErrors.ErrorNoRowsWhileReading) {
		return nil, status.Errorf(codes.NotFound, subErrors.ErrorNoRowsWhileReadingHistory.Error())
	}

	if err != nil {
		return nil, status.Errorf(codes.Internal, "something went wrong in subscription service: %v", err)
	}

	return &api.ReadBalanceHistoryResponseV1{History: history}, nil
}

// SendEmail sends email to user.
func (h *subscription) SendEmail(ctx context.Context, to string, subject string, message string) error {
	return h.emailSender.SendEmail(ctx, to, subject, message)
}

// CountWeekAndCharge charges user for current week subscription.
func (h *subscription) CountWeekAndCharge() {
	currentTime := time.Now()

	for currentTime.Weekday() != time.Thursday {
		<-time.After(24 * time.Hour)
		currentTime = currentTime.Add(24 * time.Hour)
		if currentTime.Weekday() == time.Thursday {
			*h.weekNumber += 1
		}
	}

	ticker := time.NewTicker(7 * 24 * time.Hour)

	var email string

	go func() {
		for email = range h.notEnoughFundsEmailsChan {
			logger.Logger().Infoln("got", email)
			err := h.SendToKafka(context.Background(), "cancel-subscription", email)
			if err != nil {
				logger.Logger().Errorln("sending data to auth service failed:", err)
			}
		}
	}()

	for range ticker.C {
		*h.weekNumber += 1
		logger.Logger().Infoln("new week:", *h.weekNumber) // TODO: use func to charge for all subscriptions, send messages to kafka and send emails if not enough funds
		err := h.srv.ChargeForSubscription(context.Background(), h.notEnoughFundsEmailsChan)

		if errors.Is(err, subErrors.ErrorNotEnoughFunds) {
			logger.Logger().Errorln("not enough funds on balance and could not delete", err)
		}

		if !errors.Is(err, subErrors.ErrorNotEnoughFunds) && err != nil {
			logger.Logger().Errorln(err)
		}
	}
}

// SendToKafka sends message to Kafka.
func (h *subscription) SendToKafka(ctx context.Context, topic string, message string) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}

	// Отправляем сообщение
	partition, offset, err := h.kafkaProducer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("delivery failed: %v", err)
	}

	logger.Logger().Infoln("Delivered message to topic", topic, "[", partition, "] at offset", offset)

	return nil
}
