package mail

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"log/slog"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/phamduytien1805/auth/domain"
	"github.com/phamduytien1805/package/config"
	redis_engine "github.com/phamduytien1805/package/redis"
	gomail "gopkg.in/mail.v2"
)

const (
	emailVerificationKey = "verify_email"
)

type MailService struct {
	dialer *gomail.Dialer
	redis  redis_engine.RedisQuerier
	logger *slog.Logger

	emailDuration time.Duration
	verifyLink    string
	origin        string
	templatePath  string

	taskq   domain.TaskProducer
	usersvc domain.UserService
}

func NewMailService(config *config.MailConfig, logger *slog.Logger, redis redis_engine.RedisQuerier, taskq domain.TaskProducer, usersvc domain.UserService) domain.MailService {
	return &MailService{
		dialer:        gomail.NewDialer(config.Host, config.Port, config.Username, config.Password),
		redis:         redis,
		logger:        logger,
		emailDuration: config.Expired,
		verifyLink:    config.VerifyEmailUrl,
		origin:        config.Origin,
		templatePath:  "internal/platform/mail",
		taskq:         taskq,
		usersvc:       usersvc,
	}
}

func (svc *MailService) SendEmailAsync(ctx context.Context, userEmail string) error {
	token, err := uuid.NewRandom()
	if err != nil {
		svc.logger.Error("failed to generate uuid SendEmailAsync", "detail", err.Error())
		return err
	}
	verifyKey := fmt.Sprintf("%s:%s", emailVerificationKey, token)
	if err := svc.redis.SetTx(ctx, verifyKey, userEmail, svc.emailDuration); err != nil {
		svc.logger.Error("failed to set verifyKey", "detail", err.Error())
		return err
	}
	payload := domain.SendMailParams{
		To:         userEmail,
		VerifyLink: fmt.Sprintf("%s?token=%s", svc.verifyLink, token),
	}
	if err := svc.taskq.EnqueueSendMailTask(ctx, payload); err != nil {
		return err
	}
	return nil
}

func (svc *MailService) VerifyEmail(ctx context.Context, token string) (string, error) {
	verifyKey := fmt.Sprintf("%s:%s", emailVerificationKey, token)
	userEmail, err := svc.redis.GetRaw(ctx, verifyKey)
	if err != nil {
		return "", err
	}
	if err := svc.redis.Delete(ctx, verifyKey); err != nil {
		return "", err
	}
	if err := svc.usersvc.VerifyUserEmail(ctx, userEmail); err != nil {
		return "", err
	}
	return userEmail, nil
}

func (svc *MailService) SendVerificationMail(ctx context.Context, params domain.SendMailParams) error {
	svc.logger.Info("sending email", "to", params.To)
	message := gomail.NewMessage()
	message.SetHeader("From", svc.origin)
	message.SetHeader("To", params.To)
	message.SetHeader("Subject", "Verify Your Account")
	data := struct {
		VerifyLink string
	}{
		VerifyLink: params.VerifyLink,
	}

	dir, err := os.Getwd()
	if err != nil {
		slog.Error("error while getting working directory", "detail", err.Error())
		return err
	}

	// Parse the email template
	tmpl, err := template.ParseFiles(fmt.Sprintf("%s/%s/verify_account.html", dir, svc.templatePath))
	if err != nil {
		slog.Error("error while parsing email template", "detail", err.Error())
		return err
	}
	// Render the template into a string
	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		slog.Error("error while parsing email template", "detail", err.Error())
		return err
	}

	message.SetBody("text/html", body.String())
	if err := svc.dialer.DialAndSend(message); err != nil {
		slog.Error("error while sending email", "detail", err.Error())
		return err
	}
	return nil
}
