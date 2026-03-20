package service

import (
	"context"
	"fmt"
	"log"
	"net/smtp"

	"github.com/rado31/rabbit/notification/internal/config"
	"github.com/rado31/rabbit/notification/internal/repository"
)

type NotificationService struct {
	cfg config.Config
	mq  *repository.RabbitMQRepository
}

func New(cfg config.Config, mq *repository.RabbitMQRepository) *NotificationService {
	return &NotificationService{cfg: cfg, mq: mq}
}

func (s *NotificationService) Start(ctx context.Context) error {
	return s.mq.Consume(ctx, func(e repository.ClientCreatedEvent) {
		if err := s.sendWelcomeEmail(e); err != nil {
			log.Printf("notification: send email to %s: %v", e.Email, err)
		}
	})
}

func (s *NotificationService) sendWelcomeEmail(e repository.ClientCreatedEvent) error {
	subject := "Welcome!"

	body := fmt.Sprintf(
		"Hi %s,\n\nYou have been successfully registered.\n\nBest regards,",
		e.Name,
	)

	msg := []byte(fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s",
		s.cfg.SMTPFrom, e.Email, subject, body,
	))

	addr := fmt.Sprintf("%s:%d", s.cfg.SMTPHost, s.cfg.SMTPPort)
	auth := smtp.PlainAuth("", s.cfg.SMTPUser, s.cfg.SMTPPass, s.cfg.SMTPHost)

	return smtp.SendMail(addr, auth, s.cfg.SMTPFrom, []string{e.Email}, msg)
}
