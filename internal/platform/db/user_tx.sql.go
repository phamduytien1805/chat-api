package db

import "context"

type CreateUserWithCredentialTxParams struct {
	CreateUserParams
	HashedCredential string
	Salt             string
	AfterCreate      func(User) error
}

type CreateUserWithCredentialTxResult struct {
	User User
}

func (store *SQLStore) CreateUserWithCredentialTx(ctx context.Context, arg CreateUserWithCredentialTxParams) (CreateUserWithCredentialTxResult, error) {
	var result CreateUserWithCredentialTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		result.User, err = q.CreateUser(ctx, arg.CreateUserParams)
		if err != nil {
			return err
		}

		_, err = q.CreateUserCredential(ctx, CreateUserCredentialParams{
			UserID:         result.User.ID,
			HashedPassword: arg.HashedCredential,
			Salt:           arg.Salt,
		})
		if err != nil {
			return err
		}

		return arg.AfterCreate(result.User)
	})

	return result, err
}
