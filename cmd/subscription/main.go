// Package main initializes subscription service and starts it.
package main

import (
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/IBM/sarama"
	"github.com/bufbuild/protovalidate-go"

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
		certPath       = "cert/localhost.crt"
		keyPath        = "cert/localhost.key"
		postgresDSN    = "host=postgres dbname=gophereats user=gophereats password=gophereats port=5432 sslmode=disable"
		baseDateString = "2023-11-02"
		smtpServer     = ""
		smtpPort       = ""
		smtpUsername   = ""
		smtpPassword   = ""
	)

	// Создаем конфигурацию для producer
	producerConfig := sarama.NewConfig()
	producerConfig.Producer.RequiredAcks = sarama.WaitForAll
	producerConfig.Producer.Retry.Max = 5
	producerConfig.Producer.Return.Successes = true

	// Создаем producer
	producer, err := sarama.NewSyncProducer([]string{"kafka:9092"}, producerConfig)
	if err != nil {
		logger.Logger().Errorln("Failed to create producer:", err)
		return
	}
	defer producer.Close()

	baseDate, err := time.Parse("2006-01-02", baseDateString)
	if err != nil {
		logger.Logger().Errorln("invalid base date format")
		return
	}

	if time.Now().Before(baseDate) {
		logger.Logger().Errorln("base date should be before current date")
		return
	}

	if baseDate.Weekday() != time.Thursday {
		logger.Logger().Errorln("base date should be thursday")
		return
	}

	weekNumber := (int(time.Since(baseDate).Hours()/24) / 7) + 1

	creds, err := credentials.NewClientTLSFromFile(certPath, "localhost")
	if err != nil {
		logger.Logger().Errorln("failed to load credentials:", err)
		return
	}

	conn, err := grpc.Dial("auth:8800", grpc.WithTransportCredentials(creds))
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

	validator, err := protovalidate.New()
	if err != nil {
		logger.Logger().Fatalln("Failed to create validator:", err)
		return
	}

	grpcServer := grpc.NewServer(grpc.Creds(creds), grpc.ChainUnaryInterceptor(interceptor.ValidateRequestEmail(client), interceptor.ValidateRequest(validator)))

	pgPool, err := repository.DB(postgresDSN)
	if err != nil {
		logger.Logger().Fatalln("Failed to connect to postgres:", err)
		return
	}

	subRep := repository.New(pgPool)
	subSrv := service.New(subRep)

	subscriptionServer := handler.New(subSrv, client, producer, &weekNumber, smtpUsername, smtpPassword, smtpServer, smtpPort)

	subscriptionApi.RegisterSubscriptionV1Server(grpcServer, subscriptionServer)

	c := make(chan os.Signal, 1)
	ret := make(chan struct{}, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGQUIT, syscall.SIGTERM)
	go func() {
		<-c
		ret <- struct{}{}
	}()

	listenerGRPC, err := net.Listen("tcp", "0.0.0.0:8801")
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

	go subscriptionServer.CountWeekAndCharge()

	// waiting for a signal for shutting down or an error to occur
	<-ret

	grpcServer.GracefulStop()
}
