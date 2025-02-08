package db

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/phamduytien1805/package/common"
	"github.com/phamduytien1805/user/domain"
)

// SQLStore provides all functions to execute SQL queries and transactions
type SQLStore struct {
	connPool *pgxpool.Pool
	q        *Queries
}

// NewStore creates a new store
func NewStore(connPool *pgxpool.Pool) domain.UserRepo {
	return &SQLStore{
		connPool: connPool,
		q:        New(connPool),
	}
}

func (store *SQLStore) CreateUserWithCredential(ctx context.Context, userParams domain.User, userCredential domain.UserCredential) (domain.User, error) {
	var result domain.User
	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		result, err := q.CreateUser(ctx, CreateUserParams(userParams))
		if err != nil {
			return err
		}

		_, err = q.CreateUserCredential(ctx, CreateUserCredentialParams{
			UserID:         result.ID,
			HashedPassword: userCredential.HashedPassword,
		})
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == common.UNIQUE_CONSTRAINT_VIOLATION {
				return domain.User{}, common.ErrorUserResourceConflict
			}
		}
	}

	return result, err
}

func (store *SQLStore) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	userEntity, err := store.q.GetUserByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return mapToUser(userEntity), nil
}

func (store *SQLStore) GetUserById(ctx context.Context, userID uuid.UUID) (domain.User, error) {
	userEntity, err := store.q.GetUserById(ctx, userID)
	if err != nil {
		return domain.User{}, err
	}
	return mapToUser(userEntity), nil
}

func (store *SQLStore) GetUserCredentialByUserId(ctx context.Context, userID uuid.UUID) (domain.UserCredential, error) {
	userCredentialEntity, err := store.q.GetUserCredentialByUserId(ctx, userID)
	if err != nil {
		return domain.UserCredential{}, err
	}
	return mapToUserCredential(userCredentialEntity), nil
}

func (store *SQLStore) GetUserByEmailOrUsername(ctx context.Context, emailOrUsername string) (domain.User, error) {
	userEntity, err := store.q.GetUserByEmailOrUsername(ctx, emailOrUsername)
	if err != nil {
		return domain.User{}, err
	}
	return mapToUser(userEntity), nil
}

func (store *SQLStore) UpdateUser(ctx context.Context, user domain.User) (domain.User, error) {
	userEntity, err := store.q.UpdateUser(ctx, UpdateUserParams(user))
	if err != nil {
		return domain.User{}, err
	}
	return mapToUser(userEntity), nil
}
