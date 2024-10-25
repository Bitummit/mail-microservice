package grpc

import (
	"context"
	"log/slog"

	"github.com/Bitummit/mail-microservice/internal/service"
	"github.com/Bitummit/mail-microservice/pkg/config"
	"github.com/Bitummit/mail-microservice/pkg/proto"
)


type (
	Server struct {
		cfg *config.Config
		log *slog.Logger
		service *service.MailSender
		proto.UnsafeMailServer
	}
)


func New(cfg *config.Config, log *slog.Logger) *Server{
	service := service.New()
	return &Server{
		cfg: cfg,
		log: log,
		service: service,
	}
}

func (s *Server) Send(ctx context.Context, req *proto.EmailRequest) *proto.EmailResponse {
	res := proto.EmailResponse{}
	return &res
}