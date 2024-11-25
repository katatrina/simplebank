package db

import (
	"errors"
	
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	ForeignKeyViolationCode = "23503"
	UniqueViolationCode     = "23505"
)

const (
	UniqueUsernameConstraint = "users_pkey"
	UniqueEmailConstraint    = "users_email_key"
)

var ErrRecordNotFound = pgx.ErrNoRows

func ErrorDescription(err error) (errCode string, constraintName string) {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code, pgErr.ConstraintName
	}
	
	return
}
