package handler

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/PoorMercymain/GopherEats/internal/app/auth/domain"
	authErrors "github.com/PoorMercymain/GopherEats/internal/app/auth/errors"
	"github.com/PoorMercymain/GopherEats/pkg/api"
)

type auth struct {
	srv domain.AuthService
	api.UnimplementedAuthV1Server
}

func New(srv domain.AuthService) *auth {
	return &auth{srv: srv}
}

func (h *auth) RegisterV1(ctx context.Context, r *api.RegisterRequestV1) (*emptypb.Empty, error) {
	err := h.srv.RegisterUser(ctx, r.Email, r.Password)

	if errors.Is(err, authErrors.ErrorUserAlreadyExists) {
		return &emptypb.Empty{}, status.Errorf(codes.AlreadyExists, r.Email)
	}

	if err != nil {
		return &emptypb.Empty{}, status.Errorf(codes.Internal, "something went wrong on server side in auth service: %v", err)
	}

	return &emptypb.Empty{}, nil
}

func (h *auth) LoginV1(ctx context.Context, r *api.LoginRequestV1) (*emptypb.Empty, error) {
	err := h.srv.CheckAuthData(ctx, r.Email, r.Password)

	if errors.Is(err, authErrors.ErrorNoSuchUser) {
		return &emptypb.Empty{}, status.Errorf(codes.NotFound, "no user with email %v found", r.Email)
	}

	if errors.Is(err, authErrors.ErrorWrongPassword) {
		return &emptypb.Empty{}, status.Errorf(codes.Unauthenticated, "provided password is wrong")
	}

	if err != nil {
		return &emptypb.Empty{}, status.Errorf(codes.Internal, "something went wrong on server side in auth service: %v", err)
	}

	return &emptypb.Empty{}, nil
}
