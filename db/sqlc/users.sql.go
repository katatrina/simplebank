// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: users.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (username, hashed_password, full_name, email)
VALUES ($1, $2, $3, $4) RETURNING username, role, hashed_password, full_name, email, is_email_verified, password_changed_at, created_at
`

type CreateUserParams struct {
	Username       string `json:"username"`
	HashedPassword string `json:"-"`
	FullName       string `json:"full_name"`
	Email          string `json:"email"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRow(ctx, createUser,
		arg.Username,
		arg.HashedPassword,
		arg.FullName,
		arg.Email,
	)
	var i User
	err := row.Scan(
		&i.Username,
		&i.Role,
		&i.HashedPassword,
		&i.FullName,
		&i.Email,
		&i.IsEmailVerified,
		&i.PasswordChangedAt,
		&i.CreatedAt,
	)
	return i, err
}

const getUser = `-- name: GetUser :one
SELECT username, role, hashed_password, full_name, email, is_email_verified, password_changed_at, created_at
FROM users
WHERE username = $1
`

func (q *Queries) GetUser(ctx context.Context, username string) (User, error) {
	row := q.db.QueryRow(ctx, getUser, username)
	var i User
	err := row.Scan(
		&i.Username,
		&i.Role,
		&i.HashedPassword,
		&i.FullName,
		&i.Email,
		&i.IsEmailVerified,
		&i.PasswordChangedAt,
		&i.CreatedAt,
	)
	return i, err
}

const updateUser = `-- name: UpdateUser :one
UPDATE users
SET full_name           = COALESCE($1, full_name),
    email               = COALESCE($2, email),
    hashed_password     = COALESCE($3, hashed_password),
    password_changed_at = COALESCE($4, password_changed_at),
    is_email_verified   = COALESCE($5, is_email_verified)
WHERE username = $6 RETURNING username, role, hashed_password, full_name, email, is_email_verified, password_changed_at, created_at
`

type UpdateUserParams struct {
	FullName          pgtype.Text        `json:"full_name"`
	Email             pgtype.Text        `json:"email"`
	HashedPassword    pgtype.Text        `json:"-"`
	PasswordChangedAt pgtype.Timestamptz `json:"password_changed_at"`
	IsEmailVerified   pgtype.Bool        `json:"is_email_verified"`
	Username          string             `json:"username"`
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error) {
	row := q.db.QueryRow(ctx, updateUser,
		arg.FullName,
		arg.Email,
		arg.HashedPassword,
		arg.PasswordChangedAt,
		arg.IsEmailVerified,
		arg.Username,
	)
	var i User
	err := row.Scan(
		&i.Username,
		&i.Role,
		&i.HashedPassword,
		&i.FullName,
		&i.Email,
		&i.IsEmailVerified,
		&i.PasswordChangedAt,
		&i.CreatedAt,
	)
	return i, err
}
