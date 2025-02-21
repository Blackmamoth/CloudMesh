package handlers

import (
	"context"
	"fmt"
	"net/http"
	"slices"
	"time"

	"github.com/blackmamoth/cloudmesh/pkg/config"
	"github.com/blackmamoth/cloudmesh/pkg/middlewares"
	"github.com/blackmamoth/cloudmesh/pkg/utils"
	"github.com/blackmamoth/cloudmesh/repository"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

const DROPBOX_LIST_FOLDER_API = "https://api.dropboxapi.com/2/files/list_folder"

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

	userId := r.Context().Value(middlewares.UserKey).(pgtype.UUID)
	accessToken, refreshToken, err := h.getCloudAuthTokens(userId, provider)
	if err != nil {
		utils.SendAPIErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	var insertedIdList []*pgtype.UUID

	switch provider {
	case "google":
		insertedIdList, err = h.syncGoogleDriveFiles(accessToken, refreshToken, userId)
		if err != nil {
			utils.SendAPIErrorResponse(w, http.StatusInternalServerError, err)
			return
		}
		// case "dropbox":
		// 	insertedIdList, err = h.syncDropboxFiles(accessToken, userId)
		// 	if err != nil {
		// 		utils.SendAPIErrorResponse(w, http.StatusInternalServerError, err)
		// 		return
		// 	}
	}

	utils.SendAPIResponse(w, http.StatusOK, map[string]interface{}{
		"list": insertedIdList,
	})
}

func (h *CloudStoreHandler) getCloudAuthTokens(
	userId pgtype.UUID,
	provider string,
) (string, string, error) {
	conn, err := utils.GetNewConnectionFromPool(context.Background(), h.poolConfig)
	if err != nil {
		return "", "", err
	}
	defer conn.Release()

	queries := repository.New(conn)

	authTokens, err := queries.GetCloudAuthTokens(
		context.Background(),
		repository.GetCloudAuthTokensParams{
			UserID:     userId,
			ProviderID: provider,
		},
	)
	if err != nil {
		return "", "", err
	}
	return authTokens.AccessToken.String, authTokens.RefreshToken.String, nil
}

func (h *CloudStoreHandler) syncGoogleDriveFiles(
	accessToken, refreshToken string, userId pgtype.UUID,
) ([]*pgtype.UUID, error) {
	client := utils.GetGoogleHttpClient(accessToken, refreshToken)

	driveService, err := drive.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		return nil, err
	}

	list := []*pgtype.UUID{}

	pageToken := ""

	conn, err := utils.GetNewConnectionFromPool(context.Background(), h.poolConfig)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	accountId, err := h.getAccountId(conn, userId, "google")
	if err != nil {
		return nil, err
	}

	for {
		fileList, err := driveService.Files.List().
			Fields("files(id, name, size, mimeType, createdTime, modifiedTime, thumbnailLink, fullFileExtension)").
			PageToken(pageToken).
			PageSize(1000).
			Do()
		if err != nil {
			return nil, err
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
			insertedId, err := h.addFilesToCloudStore(
				conn,
				fileStruct,
				userId,
				accountId,
				"google",
			)
			if err != nil {
				return nil, err
			}
			list = append(list, &insertedId)
		}

		if fileList.NextPageToken == "" {
			break
		}

		pageToken = fileList.NextPageToken
	}
	return list, nil
}

// func (h *CloudStoreHandler) syncDropboxFiles(
// 	accessToken string,
// 	userId pgtype.UUID,
// ) ([]*pgtype.UUID, error) {
// 	dropboxApiUrl := DROPBOX_LIST_FOLDER_API
// 	requestBody := []byte(`{"path": "", "recursive": true}`)
//
// 	httpClient := http.Client{}
//
// 	for {
// 		req, err := http.NewRequest(http.MethodPost, dropboxApiUrl, bytes.NewReader(requestBody))
// 		if err != nil {
// 			return nil, err
// 		}
//
// 		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
// 		req.Header.Set("Content-Type", "application/json")
//
// 		resp, err := httpClient.Do(req)
// 		if err != nil {
// 			return nil, err
// 		}
//
// 		defer resp.Body.Close()
//
// body, err := io.ReadAll(resp.Body)
// if err != nil {
// 	return nil, err
// }
//
// if resp.StatusCode != http.StatusOK {
// 	return nil, fmt.Errorf("status: %d", resp.StatusCode)
// }
//
// var response DropboxListFileResponse
//
// err = json.Unmarshal(body, &response)
// if err != nil {
// 	return nil, err
// }
//
// response.Entries = func(entries []DropboxListFileEntries) []DropboxListFileEntries {
// 	var filteredEntries []DropboxListFileEntries
// 	for _, entry := range entries {
// 		if entry.Tag == "file" {
// 			filteredEntries = append(filteredEntries, entry)
// 		}
// 			}
// 			return filteredEntries
// 		}(response.Entries)
//
// 		if response.Cursor == "" {
// 			break
// 		}
//
// 		if !strings.Contains(dropboxApiUrl, "continue") {
// 			dropboxApiUrl = fmt.Sprintf("%s/continue", dropboxApiUrl)
// 		}
//
// 		requestBody = []byte(fmt.Sprintf("\"cursor\": \"%s\"", response.Cursor))
//
// 	}
//
// 	return nil, nil
// }

func (h *CloudStoreHandler) getAccountId(
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
		return pgtype.UUID{Valid: false}, err
	}
	return account.ID, nil
}

func (h *CloudStoreHandler) addFilesToCloudStore(
	conn *pgxpool.Conn,
	file *CloudStoreFile,
	userId, accountId pgtype.UUID,
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
