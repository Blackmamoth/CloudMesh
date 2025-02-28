package handlers

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"slices"
	"strings"

	"github.com/blackmamoth/cloudmesh/pkg/config"
	"github.com/blackmamoth/cloudmesh/pkg/middlewares"
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

const (
	DROPBOX_AUTH_URL            = "https://www.dropbox.com/oauth2/authorize"
	DROPBOX_AUTH_TOKEN_API      = "https://api.dropboxapi.com/oauth2/token"
	DROPBOX_CURRENT_ACCOUNT_API = "https://api.dropboxapi.com/2/users/get_current_account"
)

type AuthHandler struct {
	authMiddleware *middlewares.AuthMiddleware
	poolConfig     *pgxpool.Config
}

type DropboxAuhtTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	UId          string `json:"uid"`
	AccountId    string `json:"account_id"`
}

type GetDropboxAccountResponse struct {
	AccountId string `json:"account_id"`
	Name      struct {
		DisplayName string `json:"display_name"`
	} `json:"name"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
}

func NewAuthHandler(
	authMiddleware *middlewares.AuthMiddleware,
	poolConfig *pgxpool.Config,
) *AuthHandler {
	return &AuthHandler{
		authMiddleware: authMiddleware,
		poolConfig:     poolConfig,
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

	googleProvider, err := goth.GetProvider("google")
	if err != nil {
		config.LOGGER.Fatal("could not access google provider", zap.Error(err))
	}
	googleProvider.(*google.Provider).SetAccessType("offline")
	googleProvider.(*google.Provider).SetPrompt("consent")

	r.Get("/{provider}", h.auth)
	r.Get("/{provider}/callback", h.authCallback)

	r.Group(func(r chi.Router) {
		r.Use(h.authMiddleware.VerifyRefreshToken)
		r.Get("/refresh", h.refreshTokens)
	})

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

	if strings.ToLower(provider) == "dropbox" {
		http.Redirect(w, r, getDropboxAuthUrl(), http.StatusFound)
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
		config.LOGGER.Error("could not get new connection from pool", zap.Error(err))
		utils.SendAPIErrorResponse(w, http.StatusInternalServerError, err)
		return
	}
	defer conn.Release()

	var userId pgtype.UUID

	var user goth.User

	if provider == "dropbox" {
		dropboxResponse, err := handleDropboxAuthCallbaack(r)
		if err != nil {
			utils.SendAPIErrorResponse(w, http.StatusInternalServerError, err)
			return
		}
		currentAccountResponse, err := getDropboxCurrentAccount(dropboxResponse.AccessToken)
		if err != nil {
			utils.SendAPIErrorResponse(w, http.StatusInternalServerError, err)
			return
		}

		user = goth.User{
			Name:         currentAccountResponse.Name.DisplayName,
			Provider:     "dropbox",
			Email:        currentAccountResponse.Email,
			AccessToken:  dropboxResponse.AccessToken,
			RefreshToken: dropboxResponse.RefreshToken,
		}
	} else {
		user, err = gothic.CompleteUserAuth(w, r)
		if err != nil {
			utils.SendAPIErrorResponse(w, http.StatusInternalServerError, err)
			return
		}
	}

	queries := repository.New(conn)

	existingUser, err := queries.GetUserByEmail(context.Background(), user.Email)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		newUser, err := createNewUser(conn, user)
		if err != nil {
			utils.SendAPIErrorResponse(w, http.StatusInternalServerError, err)
			return
		}

		userId = newUser.ID
		_, err = createNewAccount(conn, user, userId)
		if err != nil {
			utils.SendAPIErrorResponse(w, http.StatusInternalServerError, err)
			return
		}
	} else {
		userId = existingUser.ID
		account, err := queries.GetAccountByUserId(context.Background(), repository.GetAccountByUserIdParams{
			UserID:     userId,
			ProviderID: provider,
		})

		if err != nil && errors.Is(err, sql.ErrNoRows) {
			_, err = createNewAccount(conn, user, userId)
			if err != nil {
				utils.SendAPIErrorResponse(w, http.StatusInternalServerError, err)
				return
			}
		} else {
			err = updateAccount(conn, user, userId, account.ID)
			if err != nil {
				utils.SendAPIErrorResponse(w, http.StatusInternalServerError, err)
				return
			}
		}
	}

	if err := handleAuthTokenCookies(w, userId); err != nil {

		utils.SendAPIErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	utils.SendAPIResponse(
		w,
		http.StatusOK,
		map[string]interface{}{
			"user": user,
		},
	)
}

func getDropboxAuthUrl() string {
	return fmt.Sprintf(
		"%s?client_id=%s&token_access_type=offline&response_type=code",
		DROPBOX_AUTH_URL,
		config.OAuthConfig.DROPBOX.CLIENT_ID,
	)
}

func handleDropboxAuthCallbaack(
	r *http.Request,
) (*DropboxAuhtTokenResponse, error) {
	code := r.URL.Query().Get("code")

	if code == "" {
		return nil, fmt.Errorf("exchange code not found")
	}

	httpClient := http.Client{}

	data := url.Values{}
	data.Set("code", code)
	data.Set("grant_type", "authorization_code")

	req, err := http.NewRequest(
		http.MethodPost,
		DROPBOX_AUTH_TOKEN_API,
		bytes.NewBufferString(data.Encode()),
	)
	if err != nil {
		config.LOGGER.Error("could not generate new http request", zap.Error(err))
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(config.OAuthConfig.DROPBOX.CLIENT_ID, config.OAuthConfig.DROPBOX.CLIENT_SECRET)

	res, err := httpClient.Do(req)
	if err != nil {
		config.LOGGER.Error("dropbox oauth request failed", zap.Error(err))
		return nil, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		config.LOGGER.Error("error reading dropbox oauth response body", zap.Error(err))
		return nil, err
	}

	var dropboxResponse DropboxAuhtTokenResponse

	err = json.Unmarshal(body, &dropboxResponse)

	return &dropboxResponse, err
}

func getDropboxCurrentAccount(
	accessToken string,
) (*GetDropboxAccountResponse, error) {
	httpClient := http.Client{}

	req, err := http.NewRequest(http.MethodPost, DROPBOX_CURRENT_ACCOUNT_API, nil)
	if err != nil {
		config.LOGGER.Error(
			"could not generate new dropbox current_account request",
			zap.Error(err),
		)
		return nil, err
	}

	// req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	res, err := httpClient.Do(req)
	if err != nil {
		config.LOGGER.Error("dropbox current_account request failed", zap.Error(err))
		return nil, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		config.LOGGER.Error("could not read dropbox current_account response body", zap.Error(err))
		return nil, err
	}

	var response GetDropboxAccountResponse

	err = json.Unmarshal(body, &response)

	return &response, err
}

func (h *AuthHandler) refreshTokens(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(middlewares.UserKey).(pgtype.UUID)

	accessToken, refreshToken, err := generateAuthTokens(userId)
	if err != nil {
		utils.SendAPIErrorResponse(
			w,
			http.StatusInternalServerError,
			fmt.Errorf("could not generate tokens"),
		)
		return
	}
	utils.SetHTTPCookie(w, refreshToken, utils.REFRESH_TOKEN)
	utils.SendAPIResponse(w, http.StatusOK, map[string]string{
		"access_token": accessToken,
	})
}

func createNewUser(conn *pgxpool.Conn, user goth.User) (*repository.User, error) {
	newUserParams := repository.CreateUserParams{
		Name:  user.Name,
		Email: user.Email,
		Image: pgtype.Text{
			String: user.AvatarURL,
			Valid:  user.AvatarURL != "",
		},
	}

	var newUser repository.User

	err := utils.WithTransaction(context.Background(), conn, func(tx pgx.Tx) error {
		queries := repository.New(conn)
		qtx := queries.WithTx(tx)

		createdUser, err := qtx.CreateUser(context.Background(), newUserParams)
		if err != nil {
			config.LOGGER.Error("sql transaction failed while creating new user", zap.Error(err))
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

func createNewAccount(
	conn *pgxpool.Conn,
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

	var newAccount repository.Account

	err := utils.WithTransaction(context.Background(), conn, func(tx pgx.Tx) error {
		queries := repository.New(conn)
		qtx := queries.WithTx(tx)

		createdAccount, err := qtx.CreateAccount(context.Background(), newAccountParams)
		if err != nil {
			config.LOGGER.Error("sql transaction failed while creating new account", zap.Error(err))
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

func updateAccount(
	conn *pgxpool.Conn,
	user goth.User,
	userId, accountId pgtype.UUID,
) error {
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
		ID:     accountId,
	}

	err := utils.WithTransaction(context.Background(), conn, func(tx pgx.Tx) error {
		queries := repository.New(conn)
		qtx := queries.WithTx(tx)

		return qtx.UpdateAccountDetails(context.Background(), updateAccountParams)
	})
	if err != nil {
		config.LOGGER.Error("sql transaction failed while updating account details", zap.Error(err))
	}
	return err
}

func generateAuthTokens(userId pgtype.UUID) (string, string, error) {
	accessToken, err := utils.SignJWTToken(userId.String(), utils.ACCESS_TOKEN)
	if err != nil {
		config.LOGGER.Error("could not generate access token", zap.Error(err))
		return "", "", err
	}

	refreshToken, err := utils.SignJWTToken(userId.String(), utils.REFRESH_TOKEN)
	if err != nil {
		config.LOGGER.Error("could not generate refresh token", zap.Error(err))
		return "", "", err
	}
	return accessToken, refreshToken, nil
}

func handleAuthTokenCookies(
	w http.ResponseWriter, userId pgtype.UUID,
) error {
	accessToken, refreshToken, err := generateAuthTokens(userId)
	if err != nil {
		return err
	}

	if err := utils.SetHTTPCookie(w, accessToken, utils.ACCESS_TOKEN); err != nil {
		return err
	}

	if err := utils.SetHTTPCookie(w, refreshToken, utils.REFRESH_TOKEN); err != nil {
		return nil
	}

	return nil
}

func getUserDetails(
	conn *pgxpool.Conn,
	userId pgtype.UUID,
) (*repository.User, error) {
	queries := repository.New(conn)

	user, err := queries.GetUserById(context.Background(), userId)

	return &user, err
}

func errorRedirect(w http.ResponseWriter, r *http.Request, msg string, err error) {
	defer config.LOGGER.Sync()
	config.LOGGER.Error(msg, zap.Error(err))
	http.Redirect(w, r, "http://localhost:3000/error", http.StatusSeeOther)
}
