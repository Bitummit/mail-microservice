package grpc

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Bitummit/mail-microservice/internal/model"
	"github.com/Bitummit/mail-microservice/internal/service"
	"github.com/Bitummit/mail-microservice/pkg/config"
	"github.com/Bitummit/mail-microservice/pkg/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)


type (
	Server struct {
		Cfg *config.Config
		Log *slog.Logger
		service Service
		proto.UnsafeMailServer
	}

	Service interface {
		SendMessage(email model.Email) error
	}
)


func New(cfg *config.Config, log *slog.Logger) *Server{
	service := service.New()
	return &Server{
		Cfg: cfg,
		Log: log,
		service: service,
	}
}

func (s *Server) Send(ctx context.Context, req *proto.EmailRequest) (*proto.EmailResponse, error) {
	email := model.Email{
		To: req.GetTo(),
		Subject: req.GetSubject(),
		Body: req.GetBody(),
	}

	err := s.service.SendMessage(email)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprint("can not send message %w", err))
	}
	
	res := proto.EmailResponse{}
	return &res, nil
}