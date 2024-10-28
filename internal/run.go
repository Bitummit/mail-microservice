package run

import (
	"context"
	"fmt"
	"net"
	"os/signal"
	"sync"
	"syscall"

	my_grpc "github.com/Bitummit/mail-microservice/internal/api/grpc"
	my_kafka "github.com/Bitummit/mail-microservice/internal/api/kafka"
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

	kafkaService, err := startKafka(ctx, cfg)
	if err != nil {
		log.Error("starting kafka service", err)
		return
	}
	log.Info("Kafka started")
	
	<-ctx.Done()
	kafkaService.ConsumerGroup.Close()
	kafkaService.Conn.Close()
	log.Info("kafka stopped")
	wg.Wait()
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
	server.Log.Info("Server started")
	<-ctx.Done()
	defer wg.Done()
	grpcServer.GracefulStop()
	server.Log.Info("Server stopped")
}

func startKafka(ctx context.Context, cfg *config.Config) (*my_kafka.Kafka, error){
	kafkaService, err := my_kafka.New(ctx, cfg.KafkaLeader, "emails", "user-auth-email", 0, []string{cfg.KafkaAddress})
	if err != nil {
		return nil , fmt.Errorf("starting kafka: %w", err)
	}
	kafkaService.RunConsumerWithGroup(ctx, "user-auth-email")
	return kafkaService, nil
}