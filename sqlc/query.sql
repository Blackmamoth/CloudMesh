-- name: CreateUser :one
INSERT INTO "user" (
    name, email, image
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM "user"
WHERE email = $1;

-- name: GetUserById :one
SELECT * FROM "user"
WHERE id = $1;

-- name: CreateAccount :one
INSERT INTO account (
  user_id, account_id, provider_id, access_token, refresh_token, access_token_expires_at, id_token    
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: GetAccountByUserId :one
SELECT * FROM account
WHERE user_id = $1 AND provider_id = $2;

-- name: UpdateAccountDetails :exec
UPDATE account SET access_token = $1, refresh_token = $2, access_token_expires_at = $3, id_token = $4, updated_at = NOW() WHERE user_id = $5 AND id = $6;

-- name: UpdateOAuthTokens :one
UPDATE account SET access_token = $1, id_token = $2, updated_at = NOW() WHERE user_id = $3 AND id = $4 RETURNING access_token, refresh_token;

-- name: GetCloudAuthTokens :one
SELECT access_token, refresh_token FROM account WHERE user_id = $1 AND provider_id = $2;

-- name: AddCloudStoreFile :one
INSERT INTO cloud_store (
  account_id, provider_id, provider_file_id, file_name, file_mime_type,
  file_size, file_created_time, file_modified_time, file_thumbnail_link, file_extension
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) ON CONFLICT (provider_file_id)
DO UPDATE SET
  file_name = $4,
  file_mime_type = $5,
  file_size = $6,
  file_created_time = $7,
  file_modified_time = $8,
  file_thumbnail_link = $9,
  file_extension = $10,
  updated_at = NOW()
RETURNING id;

-- name: GetLatestSynchedFile :one
SELECT updated_at FROM cloud_store WHERE provider_id = $1 AND account_id = $2 ORDER BY updated_at DESC LIMIT 1;
