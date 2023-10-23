package main

import (
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/PoorMercymain/GopherEats/internal/app/auth/handler"
	"github.com/PoorMercymain/GopherEats/internal/app/auth/interceptor"
	"github.com/PoorMercymain/GopherEats/internal/pkg/logger"
	"github.com/PoorMercymain/GopherEats/pkg/api"
)

func main() {
	const (
		certPath = "cert/localhost.crt"
		keyPath  = "cert/localhost.key"
	)
	creds, err := credentials.NewServerTLSFromFile(certPath, keyPath)
	if err != nil {
		logger.Logger().Fatalln("Failed to setup tls:", err)
	}
	grpcServer := grpc.NewServer(grpc.Creds(creds), grpc.ChainUnaryInterceptor(interceptor.ValidateRequest))

	authServer := handler.Server{}

	api.RegisterAuthV1Server(grpcServer, authServer)

	c := make(chan os.Signal, 1)
	ret := make(chan struct{}, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGQUIT, syscall.SIGTERM)
	go func() {
		<-c
		ret <- struct{}{}
	}()

	listenerGRPC, err := net.Listen("tcp", "localhost:8800")
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
