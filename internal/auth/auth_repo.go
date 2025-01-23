package auth

import (
	"context"

	"github.com/phamduytien1805/internal/platform/db"
)

type authRepo interface {
	verifyUserEmail(ctx context.Context, userEmail string) error
}

type authRepoImpl struct {
	store db.Store
}

func newAuthGatewayImpl(store db.Store) authRepo {
	return &authRepoImpl{
		store: store,
	}
}

func (gw *authRepoImpl) verifyUserEmail(ctx context.Context, userEmail string) (err error) {
	arg := db.UpdateUserByEmailParams{
		Email:         userEmail,
		EmailVerified: true,
	}
	_, err = gw.store.UpdateUserByEmail(ctx, arg)
	return err
}
