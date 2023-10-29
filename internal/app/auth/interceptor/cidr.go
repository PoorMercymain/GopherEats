package interceptor

import (
	"context"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

func CheckCIDR(CIDR string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		const registerMethodName = "/api.v1.AuthV1/RegisterV1"

		var ctxWithIsAdmin context.Context

		md, _ := metadata.FromIncomingContext(ctx) // ignoring ok because metadata can be empty

		newMD := make(metadata.MD)
		for key, values := range md {
			if key != "is_admin" { // copying all metadata except for is_admin, because user could manually set it
				newMD[key] = values
			}
		}

		if info.FullMethod == registerMethodName {
			if CIDR == "" {
				newMD["is_admin"] = []string{"false"}
				ctxWithIsAdmin = metadata.NewIncomingContext(ctx, newMD)
				return handler(ctxWithIsAdmin, req)
			}

			pr, ok := peer.FromContext(ctx)
			if !ok {
				return nil, status.Error(codes.Internal, "Failed to get peer from context")
			}

			_, subnet, err := net.ParseCIDR(CIDR)
			if err != nil {
				return nil, status.Error(codes.Internal, "Failed to parse CIDR")
			}

			host, _, err := net.SplitHostPort(pr.Addr.String())
			if err != nil { // that may happen if there are no port
				host = pr.Addr.String()
			}

			parsedIP := net.ParseIP(host)
			if parsedIP == nil {
				return nil, status.Error(codes.Internal, "Failed to parse IP")
			}

			if !subnet.Contains(parsedIP) {
				newMD["is_admin"] = []string{"false"}
				ctxWithIsAdmin = metadata.NewIncomingContext(ctx, newMD)
				return handler(ctxWithIsAdmin, req)
			}

			newMD["is_admin"] = []string{"true"}
			ctxWithIsAdmin = metadata.NewIncomingContext(ctx, newMD)
			return handler(ctxWithIsAdmin, req)
		}

		ctx = metadata.NewIncomingContext(ctx, newMD)
		return handler(ctx, req)
	}
}
