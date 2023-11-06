package main

import (
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/PoorMercymain/GopherEats/internal/app/subscription/handler"
	"github.com/PoorMercymain/GopherEats/internal/app/subscription/interceptor"
	"github.com/PoorMercymain/GopherEats/internal/app/subscription/repository"
	"github.com/PoorMercymain/GopherEats/internal/app/subscription/service"
	"github.com/PoorMercymain/GopherEats/internal/pkg/logger"
	authApi "github.com/PoorMercymain/GopherEats/pkg/api/auth"
	subscriptionApi "github.com/PoorMercymain/GopherEats/pkg/api/subscription"
)

func main() {
	const (
		certPath    = "cert/localhost.crt"
		keyPath     = "cert/localhost.key"
		postgresDSN = "host=localhost dbname=gophereats user=gophereats password=gophereats port=5432 sslmode=disable"
	)

	creds, err := credentials.NewClientTLSFromFile(certPath, "localhost")
	if err != nil {
		logger.Logger().Errorln("failed to load credentials:", err)
		return
	}

	conn, err := grpc.Dial("localhost:8800", grpc.WithTransportCredentials(creds))
	if err != nil {
		logger.Logger().Errorln("Failed to connect:", err)
		return
	}
	defer conn.Close()

	client := authApi.NewAuthV1Client(conn)

	creds, err = credentials.NewServerTLSFromFile(certPath, keyPath)
	if err != nil {
		logger.Logger().Fatalln("Failed to setup tls:", err)
		return
	}
	grpcServer := grpc.NewServer(grpc.Creds(creds), grpc.ChainUnaryInterceptor(interceptor.ValidateRequestEmail(client)))

	pgPool, err := repository.DB(postgresDSN)
	if err != nil {
		logger.Logger().Fatalln("Failed to connect to postgres:", err)
		return
	}

	subRep := repository.New(pgPool)
	subSrv := service.New(subRep)

	subscriptionServer := handler.New(subSrv, client)

	subscriptionApi.RegisterSubscriptionV1Server(grpcServer, subscriptionServer)

	c := make(chan os.Signal, 1)
	ret := make(chan struct{}, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGQUIT, syscall.SIGTERM)
	go func() {
		<-c
		ret <- struct{}{}
	}()

	listenerGRPC, err := net.Listen("tcp", "localhost:8801")
	if err != nil {
		logger.Logger().Infoln("failed to listen:", err)
		return
	}

	go func() {
		err = grpcServer.Serve(listenerGRPC)
		if err != nil {
			logger.Logger().Errorln(err)
			ret <- struct{}{}
		}
	}()

	// waiting for a signal for shutting down or an error to occur
	<-ret

	grpcServer.GracefulStop()
}
