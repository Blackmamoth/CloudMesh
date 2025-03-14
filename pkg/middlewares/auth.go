package middlewares

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/blackmamoth/cloudmesh/pkg/utils"
	"github.com/blackmamoth/cloudmesh/repository"
	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthMiddleware struct {
	poolConfig *pgxpool.Config
}

func NewAuthMiddleware(poolConfig *pgxpool.Config) *AuthMiddleware {
	return &AuthMiddleware{
		poolConfig: poolConfig,
	}
}

type UserKeyType string

var UserKey UserKeyType = "userId"

var (
	ErrNoToken             = errors.New("unauthorized, no token")
	ErrInvalidAccessToken  = errors.New("unauthorized, invalid access token")
	ErrNoUser              = errors.New("unauthorized, user not found")
	ErrNoCookies           = errors.New("unauthorized, no cookies sent with request")
	ErrNoRefreshToken      = errors.New("unauthorized, no refresh token cookie")
	ErrInvalidRefreshToken = errors.New("unauthorized, invalid refresh token")
)

func (m *AuthMiddleware) VerifyAccessToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bearer := r.Header.Get("Authorization")

		var accessToken string

		if len(bearer) > 7 && strings.ToUpper(bearer[0:6]) == "BEARER" {
			accessToken = bearer[7:]
		} else {
			utils.SendAPIErrorResponse(w, http.StatusUnauthorized, ErrNoToken)
			return
		}

		token, err := jwtauth.VerifyToken(utils.AccessTokenAuth, accessToken)
		if err != nil {
			utils.SendAPIErrorResponse(w, http.StatusUnauthorized, ErrInvalidAccessToken)
			return
		}

		jwtSub, ok := token.Get("sub")

		if !ok {
			utils.SendAPIErrorResponse(w, http.StatusUnauthorized, ErrInvalidAccessToken)
			return
		}

		uuid, err := uuid.Parse(jwtSub.(string))
		if err != nil {
			utils.SendAPIErrorResponse(w, http.StatusUnauthorized, ErrInvalidAccessToken)
			return
		}

		userId := pgtype.UUID{
			Bytes: uuid,
			Valid: true,
		}

		// For now assume any type of error means no user exists
		if err := m.checkUserExists(userId); err != nil {
			utils.SendAPIErrorResponse(w, http.StatusUnauthorized, ErrNoUser)
			return
		}

		ctx := r.Context()

		ctx = context.WithValue(ctx, UserKey, userId)

		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func (m *AuthMiddleware) VerifyRefreshToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookies := r.Cookies()
		if len(cookies) == 0 {
			utils.SendAPIErrorResponse(w, http.StatusUnauthorized, ErrNoCookies)
			return
		}

		refreshToken, err := r.Cookie("refresh_token")
		if err != nil || refreshToken == nil || refreshToken.Value == "" {
			utils.SendAPIErrorResponse(w, http.StatusUnauthorized, ErrNoRefreshToken)
			return
		}

		token, err := jwtauth.VerifyToken(utils.RefreshTokenAuth, refreshToken.Value)
		if err != nil {
			utils.SendAPIErrorResponse(w, http.StatusUnauthorized, ErrInvalidRefreshToken)
			return
		}

		jwtSub, ok := token.Get("sub")

		if !ok {
			utils.SendAPIErrorResponse(w, http.StatusUnauthorized, ErrInvalidRefreshToken)
			return
		}

		uuid, err := uuid.Parse(jwtSub.(string))
		if err != nil {
			utils.SendAPIErrorResponse(w, http.StatusUnauthorized, ErrInvalidAccessToken)
			return
		}

		userId := pgtype.UUID{
			Bytes: uuid,
			Valid: true,
		}

		// For now assume any type of error means no user exists
		if err := m.checkUserExists(userId); err != nil {
			utils.SendAPIErrorResponse(w, http.StatusUnauthorized, ErrNoUser)
			return
		}

		ctx := r.Context()

		ctx = context.WithValue(ctx, UserKey, userId)

		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func (m *AuthMiddleware) checkUserExists(userId pgtype.UUID) error {
	conn, err := utils.GetNewConnectionFromPool(context.Background(), m.poolConfig)
	if err != nil {
		return err
	}
	defer conn.Release()
	queries := repository.New(conn)
	_, err = queries.GetUserById(context.Background(), userId)
	return err
}
