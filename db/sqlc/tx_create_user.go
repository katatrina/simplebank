package db

import (
	"context"
)

type CreateUserTxParams struct {
	CreateUserParams
	AfterCreate func(user User) error // function to be called after creating the user, inside the same transaction.
}

type CreateUserTxResult struct {
	User User
}

func (store *SQLStore) CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error) {
	var result CreateUserTxResult
	
	err := store.ExecTx(ctx, func(qTx *Queries) error {
		var err error
		
		result.User, err = qTx.CreateUser(ctx, arg.CreateUserParams)
		if err != nil {
			return err
		}
		
		return arg.AfterCreate(result.User)
	})
	
	return result, err
}
