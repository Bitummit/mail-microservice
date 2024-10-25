package run

import (
	"context"
	"net"
	"os/signal"
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

	cfg  := config.InitConfig()
	log := logger.NewLogger()
	log.Info("logger and config init")

	server := my_grpc.New(cfg, log)
	startServer(ctx, server)
	<-ctx.Done()
	log.Info("Service stopped")
}

func startServer(ctx context.Context, server *my_grpc.Server) {
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

	<-ctx.Done()
	grpcServer.GracefulStop()
	server.Log.Info("Server stopped")
}