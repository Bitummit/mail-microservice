package service

import "github.com/Bitummit/mail-microservice/internal/model"

type MailSender struct {

}

func New() *MailSender{
	return &MailSender{}
}

func (m *MailSender) SendMessage(model.Email) error {
	return nil
}