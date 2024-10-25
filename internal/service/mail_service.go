package service

import (
	"os"

	"github.com/Bitummit/mail-microservice/internal/model"
	"github.com/Bitummit/mail-microservice/pkg/config"
)

type MailSender struct {
	From string
	Password string
	Cfg *config.Config
}

func New() *MailSender{
	return &MailSender{
		From: os.Getenv("EMAIL_ACCOUNT"),
		Password: os.Getenv("EMAIL_PASSWORD"),
	}
}

func (m *MailSender) SendMessage(model.Email) error {
	
	return nil
}