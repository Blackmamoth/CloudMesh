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
	"time"

	"github.com/blackmamoth/cloudmesh/pkg/config"
	"github.com/blackmamoth/cloudmesh/pkg/middlewares"
	"github.com/blackmamoth/cloudmesh/pkg/utils"
	"github.com/blackmamoth/cloudmesh/repository"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

const (
	DROPBOX_LIST_FOLDER_API  = "https://api.dropboxapi.com/2/files/list_folder"
	DROPBOX_AUTH_REFRESH_API = "https://api.dropbox.com/oauth2/token"
	GOOGLE_OAUTH_TOKEN_API   = "https://oauth2.googleapis.com/token"
)

type CloudStoreHandler struct {
	authMiddleware *middlewares.AuthMiddleware
	poolConfig     *pgxpool.Config
}

type CloudStoreFile struct {
	Id            string `json:"id"`
	ProviderId    string `json:"provider_id"`
	MimeType      string `json:"mime_type"`
	Name          string `json:"name"`
	Size          int64  `json:"size"`
	CreatedTime   string `json:"created_time"`
	ModifiedTime  string `json:"modified_time"`
	ThumbnailLink string `json:"thumbnail_link"`
	Extension     string `json:"extension"`
}

type DropboxListFileResponse struct {
	Entries []DropboxListFileEntries `json:"entries"`
	Cursor  string                   `json:"cursor"`
	HasMore bool                     `json:"has_more"`
}

type DropboxErrrorResponse struct {
	Error struct {
		Tag string `json:".tag"`
	}
	ErrorSummary string `json:"error_summary"`
}

type GetSynchedFilesValidation struct {
	Provider string `validate:"omitempty,oneof=google dropbox" alias:"provider" json:"provider"`
	Search   string `validate:"omitempty"                      alias:"search"   json:"search"`
}

type DropboxListFileEntries struct {
	Tag            string    `json:".tag"`
	Name           string    `json:"name"`
	PathLower      string    `json:"path_lower"`
	PathDisplay    string    `json:"path_display"`
	Id             string    `json:"id"`
	ClientModified time.Time `json:"client_modified,omitempty"`
	ServerModified time.Time `json:"server_modified,omitempty"`
	Rev            string    `json:"rev,omitempty"`
	Size           int64     `json:"size,omitempty"`
	IsDownloadable bool      `json:"is_downloadable,omitempty"`
}

type GoogleOAuthRefreshResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	IDToken     string `json:"id_token"`
}

type DropboxAuthRefreshResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

func NewCloudStoreHandler(
	authMiddleware *middlewares.AuthMiddleware,
	poolConfig *pgxpool.Config,
) *CloudStoreHandler {
	return &CloudStoreHandler{
		authMiddleware: authMiddleware,
		poolConfig:     poolConfig,
	}
}

func (h *CloudStoreHandler) RegisterRoutes() *chi.Mux {
	r := chi.NewRouter()

	r.Use(h.authMiddleware.VerifyAccessToken)

	r.Get("/sync/{provider}", h.syncFiles)
	r.Post("/get-files", h.getSynchedFiles)

	return r
}

func (h *CloudStoreHandler) syncFiles(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")

	if !slices.Contains(config.OAuthConfig.SUPPORTED_PROVIDERS, provider) {
		utils.SendAPIErrorResponse(
			w,
			http.StatusUnprocessableEntity,
			fmt.Errorf("unsupported provider"),
		)
		return
	}

	conn, err := utils.GetNewConnectionFromPool(context.Background(), h.poolConfig)
	if err != nil {
		config.LOGGER.Error(
			"an error occured while getting new connection from pool",
			zap.Error(err),
		)
		utils.SendAPIErrorResponse(w, http.StatusInternalServerError, err)
		return
	}
	defer conn.Release()

	userId := r.Context().Value(middlewares.UserKey).(pgtype.UUID)
	accessToken, refreshToken, err := getCloudAuthTokens(conn, userId, provider)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			config.LOGGER.Error("attempt to sync unlinked account", zap.Error(err))
			utils.SendAPIErrorResponse(
				w,
				http.StatusUnprocessableEntity,
				fmt.Errorf(
					"could not fetch tokens for %s, please make sure you've linked a %s account",
					provider, provider,
				),
			)
			return
		}
		utils.SendAPIErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	var count int

	switch provider {
	case "google":
		count, err = syncGoogleDriveFiles(conn, accessToken, refreshToken, userId)
		if err != nil {
			utils.SendAPIErrorResponse(w, http.StatusInternalServerError, err)
			return
		}
	case "dropbox":
		count, err = syncDropboxFiles(conn, accessToken, refreshToken, userId)
		if err != nil {
			utils.SendAPIErrorResponse(w, http.StatusInternalServerError, err)
			return
		}
	}

	utils.SendAPIResponse(w, http.StatusOK, map[string]interface{}{
		"message": fmt.Sprintf("Synched %d files", count),
	})
}

func (h *CloudStoreHandler) getSynchedFiles(w http.ResponseWriter, r *http.Request) {
	var payload GetSynchedFilesValidation

	defer r.Body.Close()

	if err := utils.ParseJSON(r, &payload); err != nil && !errors.Is(err, io.EOF) {
		config.LOGGER.Error("could not parse json payload", zap.Error(err))
		utils.SendAPIErrorResponse(
			w,
			http.StatusUnprocessableEntity,
			"please provide all the required fields",
		)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		errs := utils.GenerateValidationErrorObject(err.(validator.ValidationErrors), payload)
		utils.SendAPIErrorResponse(w, http.StatusUnprocessableEntity, errs)
		return
	}

	conn, err := utils.GetNewConnectionFromPool(context.Background(), h.poolConfig)
	if err != nil {
		config.LOGGER.Error("could not get new connection from pool", zap.Error(err))
		utils.SendAPIErrorResponse(
			w,
			http.StatusInternalServerError,
			fmt.Errorf("could not process your request, please try again later"),
		)
		return
	}

	userId := r.Context().Value(middlewares.UserKey).(pgtype.UUID)

	files, err := getSynchedFiles(conn, userId, payload.Provider, payload.Search)
	if err != nil {
		config.LOGGER.Error("could query files from the db", zap.Error(err))
		utils.SendAPIErrorResponse(
			w,
			http.StatusInternalServerError,
			fmt.Errorf("an error occured while fetching your files"),
		)
		return
	}

	data := map[string]interface{}{
		"files":   files,
		"message": "Successfully fetched files.",
	}

	utils.SendAPIResponse(w, http.StatusOK, data)
}

func getCloudAuthTokens(
	conn *pgxpool.Conn,
	userId pgtype.UUID,
	provider string,
) (string, string, error) {
	queries := repository.New(conn)

	authTokens, err := queries.GetCloudAuthTokens(
		context.Background(),
		repository.GetCloudAuthTokensParams{
			UserID:     userId,
			ProviderID: provider,
		},
	)
	if err != nil {
		config.LOGGER.Error("an error occured while fetching auth tokens from DB", zap.Error(err))
		return "", "", err
	}
	return authTokens.AccessToken.String, authTokens.RefreshToken.String, nil
}

func syncGoogleDriveFiles(conn *pgxpool.Conn,
	accessToken, refreshToken string, userId pgtype.UUID,
) (int, error) {
	client := utils.GetGoogleHttpClient(accessToken, refreshToken)

	count := 0

	driveService, err := drive.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		config.LOGGER.Error("an error occured while initializing drive service")
		return count, err
	}

	pageToken := ""

	accountId, err := getAccountId(conn, userId, "google")
	if err != nil {
		return count, err
	}

	lastSyncTime, err := getLatestSynchedFile(conn, accountId, "google")
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			config.LOGGER.Error(
				"could not fetch timestamp of the latest sync",
				zap.String("provider", "google"),
				zap.Error(err),
			)
			return count, err
		}
	}

	query := ""

	if lastSyncTime.Valid {
		query = fmt.Sprintf("modifiedTime > '%s'", lastSyncTime.Time.Format(time.RFC3339))
	}

	for {
		fileList, err := driveService.Files.List().
			Fields("files(id, name, size, mimeType, createdTime, modifiedTime, thumbnailLink, fullFileExtension)").
			PageToken(pageToken).
			PageSize(1000).
			Q(query).
			Do()
		if err != nil {
			if gErr, ok := err.(*googleapi.Error); ok {
				if gErr.Code == http.StatusUnauthorized {
					accessToken, refreshToken, err = getNewOAuthTokensGoogle(
						conn,
						userId,
						accountId,
						refreshToken,
					)
					if err != nil {
						config.LOGGER.Error(
							"an error occured while fetching google drive filed",
							zap.Error(err),
						)
						return count, err
					}

					client = utils.GetGoogleHttpClient(accessToken, refreshToken)
					driveService, err = drive.NewService(
						context.Background(),
						option.WithHTTPClient(client),
					)
					if err != nil {
						config.LOGGER.Error(
							"an error occured while initializing drive service with renewed tokens",
							zap.Error(err),
						)
						return count, err
					}

					continue

				}
				return count, err
			} else {
				return count, err
			}
		}
		for _, file := range fileList.Files {
			fileStruct := &CloudStoreFile{
				ProviderId:    file.Id,
				MimeType:      file.MimeType,
				Name:          file.Name,
				Size:          file.Size,
				CreatedTime:   file.CreatedTime,
				ModifiedTime:  file.ModifiedTime,
				ThumbnailLink: file.ThumbnailLink,
				Extension:     file.FullFileExtension,
			}
			_, err := addFilesToCloudStore(
				conn,
				fileStruct,
				accountId,
				"google",
			)
			if err != nil {
				return count, err
			}
			count += 1
		}

		if fileList.NextPageToken == "" {
			break
		}

		pageToken = fileList.NextPageToken
	}
	return count, nil
}

func getDropFolderList(
	conn *pgxpool.Conn,
	accessToken, refreshToken, cursor string,
	userId, accountId pgtype.UUID,
) (*DropboxListFileResponse, error) {
	dropboxApiUrl := DROPBOX_LIST_FOLDER_API
	reqBody := []byte(`{"path": "", "recursive": true}`)

	if cursor != "" {
		dropboxApiUrl = fmt.Sprintf("%s/continue", DROPBOX_LIST_FOLDER_API)
		reqBody = []byte(fmt.Sprintf("{\"cursor\": \"%s\"}", cursor))
	}

	httpClient := http.Client{}

	for retry := 0; retry < 2; retry++ {

		req, err := http.NewRequest(http.MethodPost, dropboxApiUrl, bytes.NewReader(reqBody))
		if err != nil {
			config.LOGGER.Error("an error occured while generating http request", zap.Error(err))
			return nil, err
		}

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
		req.Header.Set("Content-Type", "application/json")

		res, err := httpClient.Do(req)
		if err != nil {
			config.LOGGER.Error("http request to dropbox failed", zap.Error(err))
			return nil, err
		}

		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			config.LOGGER.Error("could not read response body", zap.Error(err))
			return nil, err
		}

		if res.StatusCode == http.StatusUnauthorized {
			newAccessToken, err := getNewOAuthTokensDropbox(conn, userId, accountId, refreshToken)
			if err != nil {
				config.LOGGER.Error("failed to renew token", zap.Error(err))
				return nil, err
			}
			accessToken = newAccessToken
			continue
		}

		if res.StatusCode != http.StatusOK {
			var dropboxErrorResponse DropboxErrrorResponse
			err = json.Unmarshal(body, &dropboxErrorResponse)
			if err != nil {
				config.LOGGER.Error("could not unmarshal dropbox error response", zap.Error(err))
				return nil, err
			}
			return nil, fmt.Errorf("%s", dropboxErrorResponse.ErrorSummary)
		}

		var dropboxResponse DropboxListFileResponse

		err = json.Unmarshal(body, &dropboxResponse)
		return &dropboxResponse, err
	}
	return nil, fmt.Errorf("failed after retrying request")
}

func syncDropboxFiles(
	conn *pgxpool.Conn,
	accessToken, refreshToken string,
	userId pgtype.UUID,
) (int, error) {
	accountId, err := getAccountId(conn, userId, "dropbox")
	if err != nil {
		return 0, err
	}

	cursor := ""

	count := 0

	for {
		dropboxResponse, err := getDropFolderList(
			conn,
			accessToken,
			refreshToken,
			cursor,
			userId,
			accountId,
		)
		if err != nil {
			return count, err
		}

		if !dropboxResponse.HasMore {
			break
		}

		dropboxResponse.Entries = func(entries []DropboxListFileEntries) []DropboxListFileEntries {
			var filteredEntries []DropboxListFileEntries
			for _, entry := range entries {
				if entry.Tag == "file" {
					filteredEntries = append(filteredEntries, entry)
				}
			}
			return filteredEntries
		}(dropboxResponse.Entries)

		for _, entry := range dropboxResponse.Entries {
			_, err := addFilesToCloudStore(conn, &CloudStoreFile{
				ProviderId:    "dropbox",
				MimeType:      "",
				Name:          entry.Name,
				Size:          entry.Size,
				CreatedTime:   "",
				ModifiedTime:  entry.ClientModified.String(),
				ThumbnailLink: "",
				Extension:     "",
			}, accountId, "dropbox")
			if err != nil {
				return count, err
			}
			count += 1
		}

		cursor = dropboxResponse.Cursor
	}

	return count, nil
}

func getAccountId(
	conn *pgxpool.Conn,
	userId pgtype.UUID,
	provider string,
) (pgtype.UUID, error) {
	queries := repository.New(conn)
	account, err := queries.GetAccountByUserId(
		context.Background(),
		repository.GetAccountByUserIdParams{
			UserID:     userId,
			ProviderID: provider,
		},
	)
	if err != nil {
		config.LOGGER.Error("could not fetch account id", zap.Error(err))
		return pgtype.UUID{Valid: false}, err
	}
	return account.ID, nil
}

func getLatestSynchedFile(
	conn *pgxpool.Conn,
	accountId pgtype.UUID,
	provider string,
) (pgtype.Timestamp, error) {
	queries := repository.New(conn)

	lastSyncTime, err := queries.GetLatestSynchedFile(
		context.Background(),
		repository.GetLatestSynchedFileParams{
			ProviderID: provider,
			AccountID:  accountId,
		},
	)
	return lastSyncTime, err
}

func addFilesToCloudStore(
	conn *pgxpool.Conn,
	file *CloudStoreFile,
	accountId pgtype.UUID,
	provider string,
) (pgtype.UUID, error) {
	var id pgtype.UUID

	queries := repository.New(conn)

	err := utils.WithTransaction(context.Background(), conn, func(tx pgx.Tx) error {
		qtx := queries.WithTx(tx)

		fileId, err := qtx.AddCloudStoreFile(
			context.Background(),
			repository.AddCloudStoreFileParams{
				AccountID:      accountId,
				ProviderID:     provider,
				ProviderFileID: file.ProviderId,
				FileName:       file.Name,
				FileMimeType:   file.MimeType,
				FileSize:       int32(file.Size),
				FileCreatedTime: pgtype.Text{
					String: file.CreatedTime,
					Valid:  file.CreatedTime != "",
				},
				FileModifiedTime: pgtype.Text{
					String: file.ModifiedTime,
					Valid:  file.ModifiedTime != "",
				},
				FileThumbnailLink: pgtype.Text{
					String: file.ThumbnailLink,
					Valid:  file.ThumbnailLink != "",
				},
				FileExtension: pgtype.Text{
					String: file.Extension,
					Valid:  file.Extension != "",
				},
			},
		)
		if err != nil {
			config.LOGGER.Error(
				"sql transaction failed, could not add file to cloud_store table",
				zap.Error(err),
			)
			return err
		}

		id = fileId

		return nil
	})
	if err != nil {
		return pgtype.UUID{Valid: false}, err
	}

	return id, nil
}

func getNewOAuthTokensGoogle(
	conn *pgxpool.Conn,
	userId, accountId pgtype.UUID,
	refreshToken string,
) (string, string, error) {
	reqUrl := fmt.Sprintf(
		"%s?grant_type=refresh_token&client_id=%s&client_secret=%s&refresh_token=%s",
		GOOGLE_OAUTH_TOKEN_API,
		config.OAuthConfig.GOOGLE.CLIENT_ID,
		config.OAuthConfig.GOOGLE.CLIENT_SECRET,
		refreshToken,
	)

	res, err := http.Post(reqUrl, "application/json", nil)
	if err != nil {
		config.LOGGER.Error(
			"http request for renewing auth tokens failed",
			zap.String("provider", "google"),
		)
		return "", "", err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		config.LOGGER.Error("could not read response body", zap.Error(err))
		return "", "", err
	}

	defer res.Body.Close()

	var authResponse GoogleOAuthRefreshResponse

	if err := json.Unmarshal(body, &authResponse); err != nil {
		config.LOGGER.Error("could not unmarshal response into struct", zap.Error(err))
		return "", "", err
	}

	newAccessToken, refreshToken, err := updateAuthTokens(
		authResponse.AccessToken,
		authResponse.IDToken,
		conn,
		userId,
		accountId,
	)
	if err != nil {
		return "", "", err
	}

	return newAccessToken, refreshToken, nil
}

func getNewOAuthTokensDropbox(
	conn *pgxpool.Conn,
	userId, accountId pgtype.UUID,
	refreshToken string,
) (string, error) {
	data := url.Values{}
	data.Add("grant_type", "refresh_token")
	data.Add("refresh_token", refreshToken)
	data.Add("client_id", config.OAuthConfig.DROPBOX.CLIENT_ID)
	data.Add("client_secret", config.OAuthConfig.DROPBOX.CLIENT_SECRET)

	res, err := http.Post(
		DROPBOX_AUTH_REFRESH_API,
		"application/x-www-form-urlencoded",
		bytes.NewBufferString(data.Encode()),
	)
	if err != nil {
		config.LOGGER.Error("dropbox refresh request failed", zap.Error(err))
		return "", err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		config.LOGGER.Error("could not read response body", zap.Error(err))
		return "", err
	}

	var response DropboxAuhtTokenResponse

	err = json.Unmarshal(body, &response)
	if err != nil {
		config.LOGGER.Error("could not unmarshal json response", zap.Error(err))
		return "", err
	}

	accessToken, _, err := updateAuthTokens(response.AccessToken, "", conn, userId, accountId)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func updateAuthTokens(
	newAccessToken, newIdToken string,
	conn *pgxpool.Conn,
	userId, accountId pgtype.UUID,
) (string, string, error) {
	updateOAuthTokensParams := repository.UpdateOAuthTokensParams{
		AccessToken: pgtype.Text{
			String: newAccessToken,
			Valid:  newAccessToken != "",
		},
		IDToken: pgtype.Text{
			String: newIdToken,
			Valid:  newIdToken != "",
		},
		UserID: userId,
		ID:     accountId,
	}

	var accessToken, refreshToken string

	err := utils.WithTransaction(context.Background(), conn, func(tx pgx.Tx) error {
		queries := repository.New(conn)
		qtx := queries.WithTx(tx)

		tokens, err := qtx.UpdateOAuthTokens(context.Background(), updateOAuthTokensParams)
		if err != nil {
			config.LOGGER.Error(
				"sql transaction failed, could not update oauth tokens",
				zap.Error(err),
			)
			return err
		}

		accessToken = tokens.AccessToken.String
		refreshToken = tokens.RefreshToken.String

		return nil
	})
	if err != nil {
		return "", "", nil
	}

	return accessToken, refreshToken, nil
}

func getSynchedFiles(
	conn *pgxpool.Conn,
	userId pgtype.UUID,
	provider, search string,
) ([]repository.GetSynchedFilesRow, error) {
	queries := repository.New(conn)

	return queries.GetSynchedFiles(context.Background(), repository.GetSynchedFilesParams{
		UserID:   userId,
		Provider: provider,
		Search:   search,
	})
}
