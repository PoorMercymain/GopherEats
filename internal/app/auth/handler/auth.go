package handler

import (
	"context"
	"errors"
	"strconv"

	"github.com/IBM/sarama"
	"github.com/pquerna/otp/totp"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/PoorMercymain/GopherEats/internal/app/auth/domain"
	"github.com/PoorMercymain/GopherEats/internal/app/auth/email"
	authErrors "github.com/PoorMercymain/GopherEats/internal/app/auth/errors"
	"github.com/PoorMercymain/GopherEats/internal/app/auth/token"
	"github.com/PoorMercymain/GopherEats/internal/pkg/logger"
	api "github.com/PoorMercymain/GopherEats/pkg/api/auth"
)

type auth struct {
	jwtSecretKey     string
	srv              domain.AuthService
	kafkaConsumer    sarama.Consumer
	msgFromKafkaChan chan string
	api.UnimplementedAuthV1Server
}

func New(srv domain.AuthService, kafkaConsumer sarama.Consumer, jwtSecretKey string) *auth {
	return &auth{srv: srv, kafkaConsumer: kafkaConsumer, jwtSecretKey: jwtSecretKey, msgFromKafkaChan: make(chan string, 1)}
}

func (h *auth) RegisterV1(ctx context.Context, r *api.RegisterRequestV1) (*api.RegisterResponseV1, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Internal, "metadata not found in request context")
	}

	isAdminStringSlice := md.Get("is_admin")
	if len(isAdminStringSlice) == 0 {
		return nil, status.Errorf(codes.Internal, "incorrect server configuration, interceptor to check trusted subnet was not used")
	}

	isAdmin, err := strconv.ParseBool(isAdminStringSlice[0])
	if err != nil {
		return nil, status.Errorf(codes.Internal, "not a bool value provided for is_admin")
	}

	if !email.ValidateEmail(r.Email) {
		return nil, status.Errorf(codes.InvalidArgument, "wrong format of email used")
	}

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "GopherEats",
		AccountName: r.Email,
	})

	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not generate secret key for OTP")
	}

	secretKey := key.Secret()

	err = h.srv.RegisterUser(ctx, r.Email, r.Password, r.Address, secretKey, isAdmin)

	if errors.Is(err, authErrors.ErrorUserAlreadyExists) {
		return nil, status.Errorf(codes.AlreadyExists, r.Email)
	}

	if err != nil {
		return nil, status.Errorf(codes.Internal, "something went wrong on server side in auth service: %v", err)
	}

	return &api.RegisterResponseV1{OtpSecretKey: secretKey}, nil
}

func (h *auth) LoginV1(ctx context.Context, r *api.LoginRequestV1) (*api.LoginResponseV1, error) {
	if !email.ValidateEmail(r.Email) {
		return nil, status.Errorf(codes.InvalidArgument, "wrong format of email used")
	}

	err := h.srv.CheckAuthData(ctx, r.Email, r.Password)

	if errors.Is(err, authErrors.ErrorNoSuchUser) {
		return nil, status.Errorf(codes.NotFound, "no user with email %v found", r.Email)
	}

	if errors.Is(err, authErrors.ErrorWrongPassword) {
		return nil, status.Errorf(codes.Unauthenticated, "provided password is wrong")
	}

	if err != nil {
		return nil, status.Errorf(codes.Internal, "something went wrong on server side in auth service: %v", err)
	}

	hash, err := h.srv.GetPasswordHash(ctx, r.Email)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "something went wrong on server side in auth service: %v", err)
	}

	jwt, err := token.JWT(r.Email, hash, h.jwtSecretKey)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "something went wrong on server side in auth service: %v", err)
	}

	return &api.LoginResponseV1{Token: jwt}, nil
}

func (h *auth) ChangePasswordV1(ctx context.Context, r *api.ChangePasswordRequestV1) (*emptypb.Empty, error) {
	if !email.ValidateEmail(r.Email) {
		return &emptypb.Empty{}, status.Errorf(codes.InvalidArgument, "wrong format of email used")
	}

	err := h.srv.UpdatePassword(ctx, r.Email, r.OldPassword, r.NewPassword)
	if errors.Is(err, authErrors.ErrorNoSuchUser) {
		return &emptypb.Empty{}, status.Errorf(codes.NotFound, "no user with email %v found", r.Email)
	}

	if errors.Is(err, authErrors.ErrorWrongPassword) {
		return &emptypb.Empty{}, status.Errorf(codes.InvalidArgument, "the old password is wrong")
	}

	if errors.Is(err, authErrors.ErrorUserWasNotUpdated) {
		return &emptypb.Empty{}, status.Errorf(codes.Internal, "something went wrong on server side in auth service: %v", err)
	}

	if err != nil {
		return &emptypb.Empty{}, status.Errorf(codes.Internal, "something went wrong on server side in auth service: %v", err)
	}

	return &emptypb.Empty{}, nil
}

func (h *auth) LoginWithOTPV1(ctx context.Context, r *api.LoginWithOTPRequestV1) (*api.LoginResponseV1, error) {
	if !email.ValidateEmail(r.Email) {
		return nil, status.Errorf(codes.InvalidArgument, "wrong format of email used")
	}

	err := h.srv.LoginWithOTP(ctx, r.Email, r.OtpCode)

	if errors.Is(err, authErrors.ErrorNoSuchUser) {
		return nil, status.Errorf(codes.NotFound, "no user with email %v found", r.Email)
	}

	if errors.Is(err, authErrors.ErrorWrongOTP) {
		return nil, status.Errorf(codes.Unauthenticated, "provided one-time password is wrong")
	}

	if err != nil {
		return nil, status.Errorf(codes.Internal, "something went wrong on server side in auth service: %v", err)
	}

	hash, err := h.srv.GetPasswordHash(ctx, r.Email)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "something went wrong on server side in auth service: %v", err)
	}

	jwt, err := token.JWT(r.Email, hash, h.jwtSecretKey)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "something went wrong on server side in auth service: %v", err)
	}

	return &api.LoginResponseV1{Token: jwt}, nil
}

func (h *auth) ChangeAddressV1(ctx context.Context, r *api.ChangeAddressRequestV1) (*emptypb.Empty, error) {
	if !email.ValidateEmail(r.Email) {
		return &emptypb.Empty{}, status.Errorf(codes.InvalidArgument, "wrong format of email used")
	}

	err := h.srv.UpdateAddress(ctx, r.Email, r.Password, r.NewAddress)
	if errors.Is(err, authErrors.ErrorNoSuchUser) {
		return &emptypb.Empty{}, status.Errorf(codes.NotFound, "no user with email %v found", r.Email)
	}

	if errors.Is(err, authErrors.ErrorWrongPassword) {
		return &emptypb.Empty{}, status.Errorf(codes.InvalidArgument, "provided password is wrong")
	}

	if errors.Is(err, authErrors.ErrorUserWasNotUpdated) {
		return &emptypb.Empty{}, status.Errorf(codes.Internal, "something went wrong on server side in auth service: %v", err)
	}

	if err != nil {
		return &emptypb.Empty{}, status.Errorf(codes.Internal, "something went wrong on server side in auth service: %v", err)
	}

	return &emptypb.Empty{}, nil
}

func (h *auth) CheckTokenInMetadataV1(ctx context.Context, r *emptypb.Empty) (*emptypb.Empty, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return &emptypb.Empty{}, status.Errorf(codes.NotFound, "metadata not found")
	}

	tokenSlice := md.Get("authorization") // token should be under "authorization" in metadata
	if len(tokenSlice) != 1 {             // only one token allowed
		return &emptypb.Empty{}, status.Errorf(codes.InvalidArgument, "metadata should contain one token under \"authorization\" key")
	}

	providedToken := tokenSlice[0]

	emailAddress, passwordHash, err := token.GetAuthDataFromJWT(providedToken, h.jwtSecretKey)
	if err != nil {
		return &emptypb.Empty{}, status.Errorf(codes.InvalidArgument, "invalid token provided")
	}

	if !email.ValidateEmail(emailAddress) {
		return nil, status.Errorf(codes.InvalidArgument, "wrong format of email found in token %v", emailAddress)
	}

	hash, err := h.srv.GetPasswordHash(ctx, emailAddress)
	if err != nil {
		return &emptypb.Empty{}, status.Errorf(codes.Internal, "something went wrong on server side in auth service: %v", err)
	}

	if hash != passwordHash {
		return &emptypb.Empty{}, status.Errorf(codes.Internal, "something went wrong on server side in auth service: %v", authErrors.ErrorWrongToken)
	}

	return &emptypb.Empty{}, nil
}

func (h *auth) CheckIfUserIsAdminV1(ctx context.Context, r *api.CheckIfUserIsAdminRequestV1) (*api.CheckIfUserIsAdminResponseV1, error) {
	if !email.ValidateEmail(r.Email) {
		return nil, status.Errorf(codes.InvalidArgument, "wrong format of email used")
	}

	isAdmin, err := h.srv.CheckIfUserIsAdmin(ctx, r.Email)

	if errors.Is(err, authErrors.ErrorNoSuchUser) {
		return nil, status.Errorf(codes.NotFound, "no user with email %v found", r.Email)
	}

	if err != nil {
		return nil, status.Errorf(codes.Internal, "something went wrong on server side in auth service: %v", err)
	}

	return &api.CheckIfUserIsAdminResponseV1{IsAdmin: isAdmin}, nil
}

func (h *auth) GetEmailFromTokenInMetadataV1(ctx context.Context, r *emptypb.Empty) (*api.GetEmailFromTokenInMetadataResponseV1, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.NotFound, "metadata not found")
	}

	tokenSlice := md.Get("authorization") // token should be under "authorization" in metadata
	if len(tokenSlice) != 1 {             // only one token allowed
		return nil, status.Errorf(codes.InvalidArgument, "metadata should contain one token under \"authorization\" key")
	}

	providedToken := tokenSlice[0]

	emailAddress, passwordHash, err := token.GetAuthDataFromJWT(providedToken, h.jwtSecretKey)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid token provided")
	}

	if !email.ValidateEmail(emailAddress) {
		return nil, status.Errorf(codes.InvalidArgument, "wrong format of email found in token %v", emailAddress)
	}

	hash, err := h.srv.GetPasswordHash(ctx, emailAddress)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "something went wrong on server side in auth service: %v", err)
	}

	if hash != passwordHash {
		return nil, status.Errorf(codes.Internal, "something went wrong on server side in auth service: %v", authErrors.ErrorWrongToken)
	}

	return &api.GetEmailFromTokenInMetadataResponseV1{Email: emailAddress}, nil
}

func (h *auth) GetAddressV1(ctx context.Context, r *api.GetAddressRequestV1) (*api.GetAddressResponseV1, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "metadata was not provided")
	}

	c := metadata.NewOutgoingContext(ctx, md)

	emailResp, err := h.GetEmailFromTokenInMetadataV1(c, &emptypb.Empty{})
	if err != nil {
		return nil, err
	}

	c = metadata.NewOutgoingContext(ctx, md)

	isAdminResp, err := h.CheckIfUserIsAdminV1(c, &api.CheckIfUserIsAdminRequestV1{Email: emailResp.Email})
	if err != nil {
		return nil, err
	}

	if (emailResp.Email != r.Email) && !isAdminResp.IsAdmin {
		return nil, status.Errorf(codes.InvalidArgument, "email in token and request does not match")
	}

	address, err := h.srv.GetAddress(ctx, r.Email)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "something went wrong on server side in auth service: %v", err)
	}

	return &api.GetAddressResponseV1{Address: address}, nil
}

func (h *auth) ReadFromKafka(ctx context.Context, topics []string) error {
	// Подписываемся на топики
	partitionConsumer, err := h.kafkaConsumer.ConsumePartition(topics[0], 0, sarama.OffsetNewest)
	if err != nil {
		logger.Logger().Errorln("Error consuming partition:", err)
		return err
	}

	// Ждем сообщения от partitionConsumer
	go func() {
		defer partitionConsumer.Close()
		logger.Logger().Infoln("waiting for messages from kafka")
		for msg := range partitionConsumer.Messages() {
			logger.Logger().Infoln("Received message:", string(msg.Value))

			h.msgFromKafkaChan <- string(msg.Value)
		}
	}()

	return nil
}

func (h *auth) SetWarnings(ctx context.Context) error {
	err := h.ReadFromKafka(ctx, []string{"cancel-subscription"})
	if err != nil {
		return err
	}

	go func() {
		for msg := range h.msgFromKafkaChan {
			logger.Logger().Infoln("got", msg)

			err := h.srv.UpdateWarning(ctx, msg, "WARNING: not enough balance to pay for subscription!")
			if err != nil {
				logger.Logger().Errorln("Error updating warning:", err)
			}
		}
	}()

	return nil
}
