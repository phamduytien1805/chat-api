package auth

type EmailVerificationForm struct {
	Token string `json:"token" validate:"required"`
}
