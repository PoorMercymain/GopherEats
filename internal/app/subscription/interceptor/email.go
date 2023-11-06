package interceptor

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/PoorMercymain/GopherEats/pkg/api/auth"
)

type emailProvider interface {
	GetEmail() string
}

func ValidateRequestEmail(authGRPCClient auth.AuthV1Client) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if ep, ok := req.(emailProvider); ok {
			md, ok := metadata.FromIncomingContext(ctx)
			if !ok {
				return nil, status.Errorf(codes.InvalidArgument, "no metadata provided")
			}

			c := metadata.NewOutgoingContext(ctx, md)

			mdEmail, err := authGRPCClient.GetEmailFromTokenInMetadataV1(c, &emptypb.Empty{})
			if err != nil {
				return nil, err
			}

			c = metadata.NewOutgoingContext(ctx, md)

			isAdmin, err := authGRPCClient.CheckIfUserIsAdminV1(c, &auth.CheckIfUserIsAdminRequestV1{Email: mdEmail.Email})
			if err != nil {
				return nil, err
			}

			if (mdEmail.Email != ep.GetEmail()) && !isAdmin.IsAdmin {
				return nil, status.Errorf(codes.InvalidArgument, "email in token and request does not match")
			}

			ctx = metadata.NewOutgoingContext(ctx, md)
		}

		return handler(ctx, req)
	}
}
