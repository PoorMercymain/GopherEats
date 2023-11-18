package interceptor

import (
	"context"
	"errors"

	"github.com/bufbuild/protovalidate-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

// ValidateRequestUnaryServerInterceptor validates fields in gRPC messages.
func ValidateRequestUnaryServerInterceptor(validator *protovalidate.Validator) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		switch msg := req.(type) {
		case proto.Message:
			if err := validator.Validate(msg); err != nil {
				return nil, status.Error(codes.InvalidArgument, err.Error())
			}
		default:
			return nil, errors.New("unsupported message type")
		}
		return handler(ctx, req)
	}
}
