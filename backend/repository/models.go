// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package repository

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Account struct {
	ID                   pgtype.UUID
	UserID               pgtype.UUID
	AccountID            string
	ProviderID           string
	AccessToken          pgtype.Text
	RefreshToken         pgtype.Text
	AccessTokenExpiresAt pgtype.Timestamp
	IDToken              pgtype.Text
	CreatedAt            pgtype.Timestamp
	UpdatedAt            pgtype.Timestamp
}

type CloudStore struct {
	ID                pgtype.UUID
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
	CreatedAt         pgtype.Timestamp
	UpdatedAt         pgtype.Timestamp
}

type User struct {
	ID        pgtype.UUID
	Name      string
	Email     string
	Image     pgtype.Text
	CreatedAt pgtype.Timestamp
	UpdatedAt pgtype.Timestamp
}
