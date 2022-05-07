package main

import (
	"context"
	"math/rand"
	"net"
	"os"
	"time"

	"github.com/Snikimonkd/tinyUrl/internal/pkg/tinyUrl/delivery"
	server "github.com/Snikimonkd/tinyUrl/internal/pkg/tinyUrl/delivery/server"
	"github.com/Snikimonkd/tinyUrl/internal/pkg/tinyUrl/repository"
	"github.com/Snikimonkd/tinyUrl/internal/pkg/tinyUrl/usecase"
	"github.com/Snikimonkd/tinyUrl/internal/tinyUrl"
	"github.com/Snikimonkd/tinyUrl/internal/tinyUrl/utils"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
)

type ServerInterceptor struct {
	Logger *utils.Logger
}

func (s *ServerInterceptor) logger(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	start := time.Now()
	md, _ := metadata.FromIncomingContext(ctx)

	reqId := rand.Uint64()

	s.Logger.Logger = s.Logger.Logger.WithFields(logrus.Fields{
		"requestId": reqId,
		"method":    info.FullMethod,
		"context":   md,
		"request":   req,
		"response":  resp,
		"error":     err,
		"work_time": time.Since(start),
	})

	s.Logger.LogInfo("Entry Point")

	reply, err := handler(ctx, req)

	s.Logger.LogInfo("USER Interceptor")
	return reply, err
}

func main() {
	rand.Seed(time.Now().UnixNano())
	listener, err := net.Listen("tcp", ":5000")

	if err != nil {
		grpclog.Fatalf("failed to listen: %v", err)
	}

	utils.MainLogger = &utils.Logger{Logger: logrus.NewEntry(logrus.StandardLogger())}
	utils.MainLogger.Logger.Logger.SetFormatter(&logrus.TextFormatter{ForceColors: true})
	ServerInterceptor := ServerInterceptor{utils.MainLogger}
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(ServerInterceptor.logger))

	var currentDB repository.TinyUrlRepositoryInterface

	buf := os.Getenv("DB")
	utils.MainLogger.LogInfo(buf)

	if os.Getenv("DB") == "NOSQL" {
		currentDB = &repository.TinyUrlInMemoryRepository{DB: make(map[string]string)}
		utils.MainLogger.LogInfo("Current build uses NO SQL database")
	} else {
		currentDB = &repository.TinyUrlSQLRepository{DB: tinyUrl.Init()}
		utils.MainLogger.LogInfo("Current build uses SQL database")
	}

	tinyUrlServer := delivery.TinyUrlHandler{
		Usecase: &usecase.TinyUrlUseCase{Repository: currentDB},
	}

	server.RegisterTinyUrlServerServer(grpcServer, &tinyUrlServer)
	utils.MainLogger.LogInfo("Tiny URL server start at 5000")
	reflection.Register(grpcServer)
	err = grpcServer.Serve(listener)
	if err != nil {
		grpclog.Fatalf("Failed to server: %v", err)
	}
}
