package service

import (
	// "crypto/tls"
	"fmt"
	"log/slog"
	// "net/smtp"
	"os"

	"github.com/Bitummit/mail-microservice/internal/model"
	"github.com/Bitummit/mail-microservice/pkg/config"
	mail "gopkg.in/mail.v2"
)

type MailSender struct {
	from string
	password string
	Cfg *config.Config
}

func New(cfg *config.Config) *MailSender{
	return &MailSender{
		from: os.Getenv("EMAIL_ACCOUNT"),
		password: os.Getenv("EMAIL_PASSWORD"),
		Cfg: cfg,
	}
}


func (m *MailSender) SendMessage(email model.Email) error {
	message := mail.NewMessage()
	message.SetHeader("From", m.from)
	message.SetHeader("To", email.To...)
	message.SetHeader("Subject", email.Subject)
	message.SetBody("text/plain", email.Body)
	slog.Info(m.Cfg.Email.Server, m.Cfg.Email.Port)
	dealer := mail.NewDialer(m.Cfg.Email.Server, m.Cfg.Email.Port, m.from, m.password)
	// dealer.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	if err := dealer.DialAndSend(message); err != nil {
		return fmt.Errorf("send email: %w", err)
	}
	return nil
}


// func (m *MailSender) SendMessage(message model.Email) error {
// 	auth := smtp.PlainAuth("", m.from, m.password, m.Cfg.Email.Server)
// 	err := smtp.SendMail(m.Cfg.Email.Server+string(m.Cfg.Email.Port), auth, m.from, message.To, []byte(message.Body))
// 	if err != nil {
// 		return fmt.Errorf("send email: %w", err)
// 	}
// 	return nil
// }

// {
//     "to": ["ignat.001@mail.ru"],
//     "subject": "Privet",
//     "body": "TEstovoy zapusk"
// }