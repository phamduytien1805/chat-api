package usecase

import "github.com/phamduytien1805/auth/domain"

type VerifyAccessTokenUsecase struct {
	tokenSvc domain.TokenService
}

func NewVerifyAccessTokenUsecase(tokenSvc domain.TokenService) *VerifyAccessTokenUsecase {
	return &VerifyAccessTokenUsecase{
		tokenSvc: tokenSvc,
	}
}

func (uc *VerifyAccessTokenUsecase) Exec(token string) (*domain.TokenPayload, error) {
	user, err := uc.tokenSvc.VerifyAccessToken(token)
	return user, err
}
