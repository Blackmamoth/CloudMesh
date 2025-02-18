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
UPDATE account SET access_token = $1, refresh_token = $2, access_token_expires_at = $3, id_token = $4, updated_at = NOW() WHERE user_id = $5; 
