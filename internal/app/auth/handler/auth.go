package handler

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/PoorMercymain/GopherEats/pkg/api"
)

type Server struct {
	api.UnimplementedAuthV1Server
}

func (h Server) RegisterV1(ctx context.Context, r *api.RegisterRequestV1) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, status.Errorf(codes.OK, r.Email)
}
