package db

import (
	"context"
	"database/sql"
)

type VerifyUserEmailTxParams struct {
	EmailID    int64
	SecretCode string
}

type VerifyUserEmailTxResult struct {
	User        User
	VerifyEmail VerifyEmail
}

func (store *SQLStore) VerifyUserEmailTx(ctx context.Context, arg VerifyUserEmailTxParams) (VerifyUserEmailTxResult, error) {
	var result VerifyUserEmailTxResult
	
	err := store.ExecTx(ctx, func(qTx *Queries) error {
		var err error
		
		result.VerifyEmail, err = qTx.UpdateVerifyEmail(ctx, UpdateVerifyEmailParams{
			ID:         arg.EmailID,
			SecretCode: arg.SecretCode,
		})
		if err != nil {
			return err
		}
		
		result.User, err = qTx.UpdateUser(ctx, UpdateUserParams{
			Username: result.VerifyEmail.Username,
			IsEmailVerified: sql.NullBool{
				Bool:  true,
				Valid: true,
			},
		})
		if err != nil {
			return err
		}
		
		return nil
	})
	
	return result, err
}
