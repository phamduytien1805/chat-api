package user

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"github.com/phamduytien1805/internal/platform/db"
	"github.com/phamduytien1805/package/common"
)

type userRepo interface {
	createUserWithCredential(ctx context.Context, userParams *User, userCredential *UserCredential, afterCreateFn func(*User) error) (*User, error)
	getUserByEmail(ctx context.Context, email string) (*User, error)
	getUserCredentialByUserId(ctx context.Context, userID uuid.UUID) (*UserCredential, error)
}

type userRepoImpl struct {
	store db.Store
}

func newUserGatewayImpl(store db.Store) userRepo {
	return &userRepoImpl{
		store: store,
	}
}

func (gw *userRepoImpl) createUserWithCredential(ctx context.Context, userParams *User, userCredential *UserCredential, afterCreateFn func(*User) error) (*User, error) {
	arg := db.CreateUserWithCredentialTxParams{
		CreateUserParams: db.CreateUserParams{
			ID:            userParams.ID,
			Username:      userParams.Username,
			Email:         userParams.Email,
			EmailVerified: userParams.EmailVerified,
		},
		HashedCredential: userCredential.HashedPassword,
		Salt:             userCredential.Salt,
		AfterCreate: func(u db.User) error {
			return afterCreateFn(mapToUser(u))
		},
	}
	txResult, err := gw.store.CreateUserWithCredentialTx(ctx, arg)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == common.UNIQUE_CONSTRAINT_VIOLATION {
				return nil, ErrorUserResourceConflict
			}

		}
		return nil, err
	}
	return mapToUser(txResult.User), nil
}

func (gw *userRepoImpl) getUserByEmail(ctx context.Context, email string) (*User, error) {
	u, err := gw.store.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return mapToUser(u), nil
}

func (gw *userRepoImpl) getUserCredentialByUserId(ctx context.Context, userID uuid.UUID) (*UserCredential, error) {
	uc, err := gw.store.GetUserCredentialByUserId(ctx, userID)
	if err != nil {
		return nil, err
	}
	return mapToUserCredential(uc), nil
}

func mapToUser(u db.User) *User {
	return &User{
		ID:            u.ID,
		Username:      u.Username,
		Email:         u.Email,
		EmailVerified: u.EmailVerified,
	}
}

func mapToUserCredential(uc db.UserCredential) *UserCredential {
	return &UserCredential{
		HashedPassword: uc.HashedPassword,
		Salt:           uc.Salt,
	}
}
