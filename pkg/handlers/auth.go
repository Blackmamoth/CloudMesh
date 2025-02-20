package handlers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"slices"

	"github.com/blackmamoth/cloudmesh/pkg/config"
	"github.com/blackmamoth/cloudmesh/pkg/utils"
	"github.com/blackmamoth/cloudmesh/repository"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/dropbox"
	"github.com/markbates/goth/providers/google"
	"go.uber.org/zap"
)

type AuthHandler struct {
	poolConfig *pgxpool.Config
}

func NewAuthHandler(poolConfig *pgxpool.Config) *AuthHandler {
	return &AuthHandler{
		poolConfig: poolConfig,
	}
}

func (h *AuthHandler) RegisterRoutes() *chi.Mux {
	r := chi.NewRouter()

	goth.UseProviders(
		google.New(
			config.OAuthConfig.GOOGLE.CLIENT_ID,
			config.OAuthConfig.GOOGLE.CLIENT_SECRET,
			config.OAuthConfig.GOOGLE.CALLBACK_URL,
			config.OAuthConfig.GOOGLE.SCOPES...,
		),
		dropbox.New(
			config.OAuthConfig.DROPBOX.CLIENT_ID,
			config.OAuthConfig.DROPBOX.CLIENT_SECRET,
			config.OAuthConfig.DROPBOX.CALLBACK_URL,
			config.OAuthConfig.DROPBOX.SCOPES...,
		),
	)

	r.Get("/{provider}", h.auth)
	r.Get("/{provider}/callback", h.authCallback)

	return r
}

func (h *AuthHandler) auth(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")

	if !slices.Contains(config.OAuthConfig.SUPPORTED_PROVIDERS, provider) {
		utils.SendAPIErrorResponse(
			w,
			http.StatusInternalServerError,
			fmt.Errorf("invalid provider"),
		)
		return
	}

	r = r.WithContext(context.WithValue(r.Context(), "provider", provider))

	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		gothic.BeginAuthHandler(w, r)
		return
	}

	utils.SendAPIErrorResponse(w, http.StatusOK, user)
}

func (h *AuthHandler) authCallback(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")

	conn, err := utils.GetNewConnectionFromPool(context.Background(), h.poolConfig)
	if err != nil {
		utils.SendAPIErrorResponse(w, http.StatusInternalServerError, err)
		return
	}
	defer conn.Release()

	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		// h.errorRedirect(w, r, "Gothic Error", err)
		utils.SendAPIErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	var userId pgtype.UUID
	queries := repository.New(conn)

	existingUser, err := queries.GetUserByEmail(context.Background(), user.Email)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		newUser, err := h.createNewUser(user)
		if err != nil {
			utils.SendAPIErrorResponse(w, http.StatusInternalServerError, err)
			return
		}

		userId = newUser.ID
		_, err = h.createNewAccount(user, userId)
		if err != nil {
			utils.SendAPIErrorResponse(w, http.StatusInternalServerError, err)
			return
		}
	} else {
		userId = existingUser.ID
		_, err = queries.GetAccountByUserId(context.Background(), repository.GetAccountByUserIdParams{
			UserID:     userId,
			ProviderID: provider,
		})

		if err != nil && errors.Is(err, sql.ErrNoRows) {
			_, err = h.createNewAccount(user, userId)
			if err != nil {
				utils.SendAPIErrorResponse(w, http.StatusInternalServerError, err)
				return
			}
		} else {
			err = h.updateAccount(user, userId)
			if err != nil {
				utils.SendAPIErrorResponse(w, http.StatusInternalServerError, err)
				return
			}
		}
	}

	accessTokenCookie, refreshTokenCookie, err := h.generateAuthTokenCookies(userId)
	if err != nil {
		utils.SendAPIErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	http.SetCookie(w, accessTokenCookie)
	http.SetCookie(w, refreshTokenCookie)

	utils.SendAPIResponse(
		w,
		http.StatusOK,
		map[string]interface{}{
			"user": user,
		},
	)
}

func (h *AuthHandler) createNewUser(user goth.User) (*repository.User, error) {
	newUserParams := repository.CreateUserParams{
		Name:  user.Name,
		Email: user.Email,
		Image: pgtype.Text{
			String: user.AvatarURL,
			Valid:  user.AvatarURL != "",
		},
	}
	conn, err := utils.GetNewConnectionFromPool(context.Background(), h.poolConfig)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	var newUser repository.User

	err = utils.WithTransaction(context.Background(), conn, func(tx pgx.Tx) error {
		queries := repository.New(conn)
		qtx := queries.WithTx(tx)

		createdUser, err := qtx.CreateUser(context.Background(), newUserParams)
		if err != nil {
			return err
		}

		newUser = createdUser
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &newUser, nil
}

func (h *AuthHandler) createNewAccount(
	user goth.User,
	userId pgtype.UUID,
) (*repository.Account, error) {
	newAccountParams := repository.CreateAccountParams{
		AccountID:  user.UserID,
		ProviderID: user.Provider,
		UserID:     userId,
		AccessToken: pgtype.Text{
			String: user.AccessToken,
			Valid:  user.AccessToken != "",
		},
		RefreshToken: pgtype.Text{
			String: user.RefreshToken,
			Valid:  user.RefreshToken != "",
		},
		AccessTokenExpiresAt: pgtype.Timestamp{
			Time:  user.ExpiresAt,
			Valid: !user.ExpiresAt.IsZero(),
		},
		IDToken: pgtype.Text{
			String: user.IDToken,
			Valid:  user.IDToken != "",
		},
	}

	conn, err := utils.GetNewConnectionFromPool(context.Background(), h.poolConfig)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	var newAccount repository.Account

	err = utils.WithTransaction(context.Background(), conn, func(tx pgx.Tx) error {
		queries := repository.New(conn)
		qtx := queries.WithTx(tx)

		createdAccount, err := qtx.CreateAccount(context.Background(), newAccountParams)
		if err != nil {
			return err
		}

		newAccount = createdAccount
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &newAccount, nil
}

func (h *AuthHandler) updateAccount(user goth.User, userId pgtype.UUID) error {
	updateAccountParams := repository.UpdateAccountDetailsParams{
		AccessToken: pgtype.Text{
			String: user.AccessToken,
			Valid:  user.AccessToken != "",
		},
		RefreshToken: pgtype.Text{
			String: user.RefreshToken,
			Valid:  user.RefreshToken != "",
		},
		AccessTokenExpiresAt: pgtype.Timestamp{
			Time:  user.ExpiresAt,
			Valid: !user.ExpiresAt.IsZero(),
		},
		IDToken: pgtype.Text{
			String: user.IDToken,
			Valid:  user.IDToken != "",
		},
		UserID: userId,
	}

	conn, err := utils.GetNewConnectionFromPool(context.Background(), h.poolConfig)
	if err != nil {
		return err
	}
	defer conn.Release()

	err = utils.WithTransaction(context.Background(), conn, func(tx pgx.Tx) error {
		queries := repository.New(conn)
		qtx := queries.WithTx(tx)

		return qtx.UpdateAccountDetails(context.Background(), updateAccountParams)
	})
	return err
}

func (h *AuthHandler) generateAuthTokenCookies(
	userId pgtype.UUID,
) (*http.Cookie, *http.Cookie, error) {
	accessToken, err := utils.SignJWTToken(userId.String(), utils.ACCESS_TOKEN)
	if err != nil {
		return nil, nil, err
	}

	refreshToken, err := utils.SignJWTToken(userId.String(), utils.REFRESH_TOKEN)
	if err != nil {
		return nil, nil, err
	}

	accessTokenExpiration := utils.GetAccessTokenExpirationTime()
	accessTokenCookie := &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Path:     "/",
		MaxAge:   int(accessTokenExpiration.Unix()),
		Secure:   config.AppConfig.ENVIRONMENT != "DEVELOPMENT",
		HttpOnly: false,
		SameSite: http.SameSiteLaxMode,
		Expires:  accessTokenExpiration,
	}

	refreshTokenExpiration := utils.GetRefreshTokenExpirationTime()
	refreshTokenCookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/",
		MaxAge:   int(refreshTokenExpiration.Unix()),
		Secure:   config.AppConfig.ENVIRONMENT != "DEVELOPMENT",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Expires:  refreshTokenExpiration,
	}

	return accessTokenCookie, refreshTokenCookie, nil
}

func (h *AuthHandler) errorRedirect(w http.ResponseWriter, r *http.Request, msg string, err error) {
	defer config.LOGGER.Sync()
	config.LOGGER.Error(msg, zap.Error(err))
	http.Redirect(w, r, "http://localhost:3000/error", http.StatusSeeOther)
}
