package run

import (
	"context"
	"net"
	"os/signal"
	"sync"
	"syscall"

	my_grpc "github.com/Bitummit/mail-microservice/internal/api/grpc"
	"github.com/Bitummit/mail-microservice/pkg/config"
	"github.com/Bitummit/mail-microservice/pkg/logger"
	"github.com/Bitummit/mail-microservice/pkg/proto"
	"google.golang.org/grpc"
)

func Run() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	wg := &sync.WaitGroup{}

	cfg  := config.InitConfig()
	log := logger.NewLogger()
	log.Info("logger and config init")

	server := my_grpc.New(cfg, log)
	wg.Add(1)
	go startServer(ctx, wg, server)
	<-ctx.Done()
	log.Info("Service stopped")
}

func startServer(ctx context.Context, wg *sync.WaitGroup, server *my_grpc.Server) {
	listener, err := net.Listen("tcp", server.Cfg.GrpcAddress)
	if err != nil {
		server.Log.Error("Error trying to listen: %w", logger.Err(err))
	}
	opts := []grpc.ServerOption{}
	grpcServer := grpc.NewServer(opts...)
	proto.RegisterMailServer(grpcServer, server)
	
	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			server.Log.Error("Server error: %w", logger.Err(err))
		}
	}()
	server.Log.Info("Server started5")
	<-ctx.Done()
	defer wg.Done()
	grpcServer.GracefulStop()
	server.Log.Info("Server stopped")
}