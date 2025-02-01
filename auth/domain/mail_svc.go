package domain

import "context"

type MailService interface {
	SendEmailAsync(ctx context.Context, userEmail string) error
	VerifyEmail(ctx context.Context, token string) (string, error)
	SendVerificationMail(ctx context.Context, params SendMailParams) error
}

type SendMailParams struct {
	To         string
	VerifyLink string
}
