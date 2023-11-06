package handler

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/PoorMercymain/GopherEats/internal/app/subscription/domain"
	subErrors "github.com/PoorMercymain/GopherEats/internal/app/subscription/errors"
	"github.com/PoorMercymain/GopherEats/pkg/api/auth"
	api "github.com/PoorMercymain/GopherEats/pkg/api/subscription"
)

var _ api.SubscriptionV1Server = (*subscription)(nil)

type subscription struct {
	srv    domain.SubscriptionService
	client auth.AuthV1Client
	api.UnimplementedSubscriptionV1Server
}

func New(srv domain.SubscriptionService, client auth.AuthV1Client) *subscription {
	return &subscription{srv: srv, client: client}
}

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

func (h *subscription) ReadSubscriptionV1(ctx context.Context, r *api.ReadSubscriptionRequestV1) (*api.ReadSubscriptionResponseV1, error) {
	bundleID, err := h.srv.ReadSubscription(ctx, r.Email)

	if errors.Is(err, subErrors.ErrorNoRowsWhileReading) {
		return nil, status.Errorf(codes.NotFound, subErrors.ErrorNoRowsWhileReading.Error())
	}

	if err != nil {
		return nil, status.Errorf(codes.Internal, "something went wrong in subscription service: %v", err)
	}

	return &api.ReadSubscriptionResponseV1{BundleId: bundleID}, nil
}

func (h *subscription) ChangeSubscriptionV1(ctx context.Context, r *api.ChangeSubscriptionRequestV1) (*emptypb.Empty, error) {
	err := h.srv.UpdateSubscription(ctx, r.Email, r.BundleId)

	if errors.Is(err, subErrors.ErrorNoRowsUpdated) {
		return nil, status.Errorf(codes.NotFound, subErrors.ErrorNoRowsUpdated.Error())
	}

	if err != nil {
		return nil, status.Errorf(codes.Internal, "something went wrong in subscription service: %v", err)
	}

	return &emptypb.Empty{}, nil
}

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

func (h *subscription) AddBalanceV1(ctx context.Context, r *api.AddBalanceRequestV1) (*emptypb.Empty, error) {
	err := h.srv.AddBalance(ctx, r.Email, r.Amount)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "something went wrong in subscription service: %v", err)
	}

	return &emptypb.Empty{}, nil
}

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
