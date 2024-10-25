package run

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/Bitummit/mail-microservice/pkg/config"
	"github.com/Bitummit/mail-microservice/pkg/logger"
)

func Run() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg  := config.InitConfig()
	log := logger.NewLogger()
	log.Info("logger and config init")

	
}