// Package main initializes dishes service and starts it.
package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/bufbuild/protovalidate-go"
	"google.golang.org/grpc"

	"github.com/PoorMercymain/GopherEats/internal/app/dishes/config"
	"github.com/PoorMercymain/GopherEats/internal/app/dishes/handler"
	"github.com/PoorMercymain/GopherEats/internal/app/dishes/interceptor"
	"github.com/PoorMercymain/GopherEats/internal/app/dishes/repository"
	"github.com/PoorMercymain/GopherEats/internal/pkg/logger"
	pb "github.com/PoorMercymain/GopherEats/pkg/api/dishes"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {

	fmt.Println("Build version: ", buildVersion)
	fmt.Println("Build date: ", buildDate)
	fmt.Println("Build commit: ", buildCommit)

	serverConfig := config.GetServerConfig()

	if serverConfig.DatabaseDSN == "" {
		logger.Logger().Fatalln("Database DSN is missing")
	}

	if serverConfig.HostGrpc == "" {
		logger.Logger().Fatalln("GRPC Host address is missing")
	}

	repo, err := repository.NewDBStorage(serverConfig.DatabaseDSN)

	ctx := context.Background()

	listen, err := net.Listen("tcp", serverConfig.HostGrpc)

	if err != nil {
		logger.Logger().Fatalln("Failed to announce:", err)
	}

	validator, err := protovalidate.New()
	if err != nil {
		logger.Logger().Fatalln("Failed to initialize validator:", err)
	}

	s := grpc.NewServer(grpc.ChainUnaryInterceptor(interceptor.ValidateRequestUnaryServerInterceptor(validator)))

	pb.RegisterDishesServiceV1Server(s, handler.NewDishesServerV1WithStorage(repo))

	go listenAndServeGrpc(ctx, s, listen)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	select {
	case <-ctx.Done():
		err := context.Cause(ctx)
		if err != nil {
			logger.Logger().Fatalln("Context cancelled: ", err)
		}

	case <-sigChan:
		logger.Logger().Infoln("Shutdown signal")

	}

	s.GracefulStop()
}

func listenAndServeGrpc(ctx context.Context, s *grpc.Server, listen net.Listener) {
	_, cancelCtx := context.WithCancelCause(ctx)
	err := s.Serve(listen)
	if err != nil {
		logger.Logger().Infoln("listenAndServe gRPC err: ", err)
		cancelCtx(err)
	}
}
