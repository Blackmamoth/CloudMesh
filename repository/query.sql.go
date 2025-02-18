// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: query.sql

package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createAccount = `-- name: CreateAccount :one
INSERT INTO account (
  user_id, account_id, provider_id, access_token, refresh_token, access_token_expires_at, id_token    
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
) RETURNING id, user_id, account_id, provider_id, access_token, refresh_token, access_token_expires_at, id_token, created_at, updated_at
`

type CreateAccountParams struct {
	UserID               pgtype.UUID
	AccountID            string
	ProviderID           string
	AccessToken          pgtype.Text
	RefreshToken         pgtype.Text
	AccessTokenExpiresAt pgtype.Timestamp
	IDToken              pgtype.Text
}

func (q *Queries) CreateAccount(ctx context.Context, arg CreateAccountParams) (Account, error) {
	row := q.db.QueryRow(ctx, createAccount,
		arg.UserID,
		arg.AccountID,
		arg.ProviderID,
		arg.AccessToken,
		arg.RefreshToken,
		arg.AccessTokenExpiresAt,
		arg.IDToken,
	)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.AccountID,
		&i.ProviderID,
		&i.AccessToken,
		&i.RefreshToken,
		&i.AccessTokenExpiresAt,
		&i.IDToken,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const createUser = `-- name: CreateUser :one
INSERT INTO "user" (
    name, email, image
) VALUES (
    $1, $2, $3
) RETURNING id, name, email, image, created_at, updated_at
`

type CreateUserParams struct {
	Name  string
	Email string
	Image pgtype.Text
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRow(ctx, createUser, arg.Name, arg.Email, arg.Image)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.Image,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getAccountByUserId = `-- name: GetAccountByUserId :one
SELECT id, user_id, account_id, provider_id, access_token, refresh_token, access_token_expires_at, id_token, created_at, updated_at FROM account
WHERE user_id = $1 AND provider_id = $2
`

type GetAccountByUserIdParams struct {
	UserID     pgtype.UUID
	ProviderID string
}

func (q *Queries) GetAccountByUserId(ctx context.Context, arg GetAccountByUserIdParams) (Account, error) {
	row := q.db.QueryRow(ctx, getAccountByUserId, arg.UserID, arg.ProviderID)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.AccountID,
		&i.ProviderID,
		&i.AccessToken,
		&i.RefreshToken,
		&i.AccessTokenExpiresAt,
		&i.IDToken,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, name, email, image, created_at, updated_at FROM "user"
WHERE email = $1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRow(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.Image,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getUserById = `-- name: GetUserById :one
SELECT id, name, email, image, created_at, updated_at FROM "user"
WHERE id = $1
`

func (q *Queries) GetUserById(ctx context.Context, id pgtype.UUID) (User, error) {
	row := q.db.QueryRow(ctx, getUserById, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.Image,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateAccountDetails = `-- name: UpdateAccountDetails :exec
UPDATE account SET access_token = $1, refresh_token = $2, access_token_expires_at = $3, id_token = $4, updated_at = NOW() WHERE user_id = $5
`

type UpdateAccountDetailsParams struct {
	AccessToken          pgtype.Text
	RefreshToken         pgtype.Text
	AccessTokenExpiresAt pgtype.Timestamp
	IDToken              pgtype.Text
	UserID               pgtype.UUID
}

func (q *Queries) UpdateAccountDetails(ctx context.Context, arg UpdateAccountDetailsParams) error {
	_, err := q.db.Exec(ctx, updateAccountDetails,
		arg.AccessToken,
		arg.RefreshToken,
		arg.AccessTokenExpiresAt,
		arg.IDToken,
		arg.UserID,
	)
	return err
}
