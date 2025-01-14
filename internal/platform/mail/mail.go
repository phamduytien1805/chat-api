package mail

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"os"

	"html/template"

	"github.com/phamduytien1805/package/config"
	gomail "gopkg.in/mail.v2"
)

type SendMailParams struct {
	To         string
	VerifyLink string
}

type MailService interface {
	SendVerificationMail(ctx context.Context, payload SendMailParams) error
}

type MailServiceImpl struct {
	dialer       *gomail.Dialer
	origin       string
	logger       *slog.Logger
	templatePath string
}

func NewMailService(configEmail *config.MailConfig, logger *slog.Logger) MailService {
	dialer := gomail.NewDialer(configEmail.Host, configEmail.Port, configEmail.Username, configEmail.Password)
	return &MailServiceImpl{
		dialer:       dialer,
		origin:       configEmail.Origin,
		logger:       logger,
		templatePath: "internal/platform/mail",
	}
}

func (m *MailServiceImpl) SendVerificationMail(ctx context.Context, payload SendMailParams) error {
	m.logger.Info("sending email", "to", payload.To)
	message := gomail.NewMessage()
	message.SetHeader("From", m.origin)
	message.SetHeader("To", payload.To)
	message.SetHeader("Subject", "Verify Your Account")
	data := struct {
		VerifyLink string
	}{
		VerifyLink: payload.VerifyLink,
	}

	dir, err := os.Getwd()
	if err != nil {
		slog.Error("error while getting working directory", "detail", err.Error())
		return err
	}

	// Parse the email template
	tmpl, err := template.ParseFiles(fmt.Sprintf("%s/%s/verify_account.html", dir, m.templatePath))
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
	if err := m.dialer.DialAndSend(message); err != nil {
		slog.Error("error while sending email", "detail", err.Error())
		return err
	}
	return nil
}
