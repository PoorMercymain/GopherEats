package main

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/IBM/sarama"
	"github.com/bufbuild/protovalidate-go"

	"github.com/PoorMercymain/GopherEats/internal/app/auth/handler"
	"github.com/PoorMercymain/GopherEats/internal/app/auth/interceptor"
	"github.com/PoorMercymain/GopherEats/internal/app/auth/repository"
	"github.com/PoorMercymain/GopherEats/internal/app/auth/service"
	"github.com/PoorMercymain/GopherEats/internal/pkg/logger"
	api "github.com/PoorMercymain/GopherEats/pkg/api/auth"
)

func main() {
	const (
		certPath      = "cert/localhost.crt"
		keyPath       = "cert/localhost.key"
		mongoURI      = "mongodb://mongodb:27017"
		trustedSubnet = ""
		jwtSecretKey  = "somesecretkey"
	)

	// Создаем конфигурацию для consumer
	consumerConfig := sarama.NewConfig()
	consumerConfig.Consumer.Offsets.Initial = sarama.OffsetNewest

	// Создаем consumer
	consumer, err := sarama.NewConsumer([]string{"kafka:9092"}, consumerConfig)
	if err != nil {
		logger.Logger().Errorln("Failed to create consumer:", err)
		return
	}
	defer consumer.Close()

	creds, err := credentials.NewServerTLSFromFile(certPath, keyPath)
	if err != nil {
		logger.Logger().Fatalln("Failed to setup tls:", err)
	}

	validator, err := protovalidate.New()
	if err != nil {
		logger.Logger().Fatalln("Failed to create validator:", err)
		return
	}

	grpcServer := grpc.NewServer(grpc.Creds(creds), grpc.ChainUnaryInterceptor(interceptor.ValidateRequest(validator), interceptor.CheckCIDR(trustedSubnet)))

	clientOptions := options.Client().ApplyURI(mongoURI)
	ctx := context.Background()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		logger.Logger().Errorln(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		logger.Logger().Errorln(err)
	}

	collection := client.Database("GopherEats").Collection("users")

	indexModel := mongo.IndexModel{
		Keys:    bson.M{"email": 1},
		Options: options.Index().SetUnique(true),
	}

	_, err = collection.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		logger.Logger().Errorln(err)
	}

	ar := repository.New(collection)
	as := service.New(ar)

	authServer := handler.New(as, consumer, jwtSecretKey)

	api.RegisterAuthV1Server(grpcServer, authServer)

	c := make(chan os.Signal, 1)
	ret := make(chan struct{}, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGQUIT, syscall.SIGTERM)
	go func() {
		<-c
		ret <- struct{}{}
	}()

	listenerGRPC, err := net.Listen("tcp", "0.0.0.0:8800")
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

	go func() {
		err := authServer.SetWarnings(context.Background())
		if err != nil {
			logger.Logger().Errorln(err)
		}
	}()

	// waiting for a signal for shutting down or an error to occur
	<-ret

	grpcServer.GracefulStop()
}
