// Package main initializes dishes service and starts it.
package main

import (
	"context"
	"database/sql"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	pb "github.com/PoorMercymain/GopherEats/api/gen/v1"
	"github.com/PoorMercymain/GopherEats/internal/pkg/logger"
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

	var err error
	var db *sql.DB

	gRPCHost, dbDSN := config.GetDishesServerConfig()

	if dbDSN != "" {
		db, err = sql.Open("pgx", dbDSN)

		if err != nil {
			logger.Logger().Fatalln("Failed to open DB:", err)
		}

		defer func() {
			err = db.Close()
			if err != nil {
				logger.Logger().Infoln("Failed to close db: ", err)
			}
		}()

		//ping
		err = db.Ping()
		if err != nil {
			logger.Logger().Fatalln("Failed to ping DB:", err)
		}
	}

	appStorage := storage.NewStorage(db, fileStoragePath)

	storageHandler := handlers.NewStorageHandler(appStorage)

	ctx := context.Background()

	listen, err := net.Listen("tcp", gRPCHost)

	if err != nil {
		log.Fatalln(err)
	}

	s := grpc.NewServer(grpc.ChainUnaryInterceptor(interceptors.ValidateRequest))
	pb.RegisterMetricsServiceV1Server(s, router.NewMetricsServerWithStorage(appStorage))

	go listenAndServeGrpc(ctx, s, listen)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	select {
	case <-ctx.Done():
		err := context.Cause(ctx)
		if err != nil {
			log.Fatalln(err)
		}

	case <-sigChan:
		fmt.Println("shutdown signal")

	}
}

func listenAndServeGrpc(ctx context.Context, s *grpc.Server, listen net.Listener) {
	_, cancelCtx := context.WithCancelCause(ctx)
	err := s.Serve(listen)
	if err != nil {
		fmt.Println("listenAndServe gRPC err", err)
		cancelCtx(err)
	}
}
