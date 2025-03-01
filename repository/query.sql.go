// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: query.sql

package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const addCloudStoreFile = `-- name: AddCloudStoreFile :one
INSERT INTO cloud_store (
  account_id, provider_id, provider_file_id, file_name, file_mime_type,
  file_size, file_created_time, file_modified_time, file_thumbnail_link, file_extension) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) ON CONFLICT (provider_file_id)
DO UPDATE SET
  file_name = $4,
  file_mime_type = $5,
  file_size = $6,
  file_created_time = $7,
  file_modified_time = $8,
  file_thumbnail_link = $9,
  file_extension = $10,
  updated_at = NOW()
RETURNING id
`

type AddCloudStoreFileParams struct {
	AccountID         pgtype.UUID
	ProviderID        string
	ProviderFileID    string
	FileName          string
	FileMimeType      string
	FileSize          int32
	FileCreatedTime   pgtype.Text
	FileModifiedTime  pgtype.Text
	FileThumbnailLink pgtype.Text
	FileExtension     pgtype.Text
}

func (q *Queries) AddCloudStoreFile(ctx context.Context, arg AddCloudStoreFileParams) (pgtype.UUID, error) {
	row := q.db.QueryRow(ctx, addCloudStoreFile,
		arg.AccountID,
		arg.ProviderID,
		arg.ProviderFileID,
		arg.FileName,
		arg.FileMimeType,
		arg.FileSize,
		arg.FileCreatedTime,
		arg.FileModifiedTime,
		arg.FileThumbnailLink,
		arg.FileExtension,
	)
	var id pgtype.UUID
	err := row.Scan(&id)
	return id, err
}

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

const getCloudAuthTokens = `-- name: GetCloudAuthTokens :one
SELECT access_token, refresh_token FROM account WHERE user_id = $1 AND provider_id = $2
`

type GetCloudAuthTokensParams struct {
	UserID     pgtype.UUID
	ProviderID string
}

type GetCloudAuthTokensRow struct {
	AccessToken  pgtype.Text
	RefreshToken pgtype.Text
}

func (q *Queries) GetCloudAuthTokens(ctx context.Context, arg GetCloudAuthTokensParams) (GetCloudAuthTokensRow, error) {
	row := q.db.QueryRow(ctx, getCloudAuthTokens, arg.UserID, arg.ProviderID)
	var i GetCloudAuthTokensRow
	err := row.Scan(&i.AccessToken, &i.RefreshToken)
	return i, err
}

const getLatestSynchedFile = `-- name: GetLatestSynchedFile :one
SELECT updated_at FROM cloud_store WHERE provider_id = $1 AND account_id = $2 ORDER BY updated_at DESC LIMIT 1
`

type GetLatestSynchedFileParams struct {
	ProviderID string
	AccountID  pgtype.UUID
}

func (q *Queries) GetLatestSynchedFile(ctx context.Context, arg GetLatestSynchedFileParams) (pgtype.Timestamp, error) {
	row := q.db.QueryRow(ctx, getLatestSynchedFile, arg.ProviderID, arg.AccountID)
	var updated_at pgtype.Timestamp
	err := row.Scan(&updated_at)
	return updated_at, err
}

const getSynchedFiles = `-- name: GetSynchedFiles :many
SELECT 
    store.file_name AS file_name, 
    store.file_size AS file_size, 
    store.provider_id AS provider,
    store.file_created_time as file_created_time, 
    store.file_modified_time as file_modified_time, 
    store.file_thumbnail_link as file_thumbnail_link 
FROM 
    cloud_store store
JOIN 
    account ON store.account_id = account.id
WHERE 
    account.user_id = $1
AND (
    $2::TEXT IS NULL 
    OR $2::TEXT = '' 
    OR store.provider_id ILIKE '%' || $2::TEXT || '%'
)
AND (
    $3::TEXT IS NULL 
    OR $3::TEXT = '' 
    OR store.file_name ILIKE '%' || $3::TEXT || '%'
    OR store.file_extension ILIKE '%s' || $3::TEXT || '%'
) ORDER BY file_name ASC LIMIT $5 OFFSET $4
`

type GetSynchedFilesParams struct {
	UserID   pgtype.UUID
	Provider string
	Search   string
	OffsetOf int32
	LimitBy  int32
}

type GetSynchedFilesRow struct {
	FileName          string
	FileSize          int32
	Provider          string
	FileCreatedTime   pgtype.Text
	FileModifiedTime  pgtype.Text
	FileThumbnailLink pgtype.Text
}

func (q *Queries) GetSynchedFiles(ctx context.Context, arg GetSynchedFilesParams) ([]GetSynchedFilesRow, error) {
	rows, err := q.db.Query(ctx, getSynchedFiles,
		arg.UserID,
		arg.Provider,
		arg.Search,
		arg.OffsetOf,
		arg.LimitBy,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetSynchedFilesRow
	for rows.Next() {
		var i GetSynchedFilesRow
		if err := rows.Scan(
			&i.FileName,
			&i.FileSize,
			&i.Provider,
			&i.FileCreatedTime,
			&i.FileModifiedTime,
			&i.FileThumbnailLink,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
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
UPDATE account SET access_token = $1, refresh_token = $2, access_token_expires_at = $3, id_token = $4, updated_at = NOW() WHERE user_id = $5 AND id = $6
`

type UpdateAccountDetailsParams struct {
	AccessToken          pgtype.Text
	RefreshToken         pgtype.Text
	AccessTokenExpiresAt pgtype.Timestamp
	IDToken              pgtype.Text
	UserID               pgtype.UUID
	ID                   pgtype.UUID
}

func (q *Queries) UpdateAccountDetails(ctx context.Context, arg UpdateAccountDetailsParams) error {
	_, err := q.db.Exec(ctx, updateAccountDetails,
		arg.AccessToken,
		arg.RefreshToken,
		arg.AccessTokenExpiresAt,
		arg.IDToken,
		arg.UserID,
		arg.ID,
	)
	return err
}

const updateOAuthTokens = `-- name: UpdateOAuthTokens :one
UPDATE account SET access_token = $1, id_token = $2, updated_at = NOW() WHERE user_id = $3 AND id = $4 RETURNING access_token, refresh_token
`

type UpdateOAuthTokensParams struct {
	AccessToken pgtype.Text
	IDToken     pgtype.Text
	UserID      pgtype.UUID
	ID          pgtype.UUID
}

type UpdateOAuthTokensRow struct {
	AccessToken  pgtype.Text
	RefreshToken pgtype.Text
}

func (q *Queries) UpdateOAuthTokens(ctx context.Context, arg UpdateOAuthTokensParams) (UpdateOAuthTokensRow, error) {
	row := q.db.QueryRow(ctx, updateOAuthTokens,
		arg.AccessToken,
		arg.IDToken,
		arg.UserID,
		arg.ID,
	)
	var i UpdateOAuthTokensRow
	err := row.Scan(&i.AccessToken, &i.RefreshToken)
	return i, err
}
