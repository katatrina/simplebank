package db

import (
	"context"
	
	"github.com/jackc/pgx/v5/pgxpool"
)

// Store provides all functions to execute db queries and transactions.
type Store interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
	CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error)
	VerifyUserEmailTx(ctx context.Context, arg VerifyUserEmailTxParams) (VerifyUserEmailTxResult, error)
	Ping(ctx context.Context) error
}

type SQLStore struct {
	*Queries
	connPool *pgxpool.Pool
}

// NewStore creates a new Store.
func NewStore(db *pgxpool.Pool) Store {
	return &SQLStore{
		Queries:  New(db),
		connPool: db,
	}
}

// Ping checks if the database connection is alive.
func (store *SQLStore) Ping(ctx context.Context) error {
	return store.connPool.Ping(ctx)
}
