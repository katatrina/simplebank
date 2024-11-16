package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Store provides all functions to execute db queries and transactions.
type Store interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
	CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error)
	VerifyUserEmailTx(ctx context.Context, arg VerifyUserEmailTxParams) (VerifyUserEmailTxResult, error)
}

type SQLStore struct {
	*Queries
	connPool *sql.DB
}

// NewStore creates a new Store.
func NewStore(db *sql.DB) Store {
	return &SQLStore{
		Queries:  New(db),
		connPool: db,
	}
}

// ExecTx executes a function within a database transaction.
func (store *SQLStore) ExecTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.connPool.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	
	qTx := New(tx)
	
	err = fn(qTx)
	if err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			return fmt.Errorf("transaction error: %v, rollback error: %v", err, rbErr)
		}
		
		return err
	}
	
	return tx.Commit()
}
