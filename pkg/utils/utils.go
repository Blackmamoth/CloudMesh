package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/blackmamoth/cloudmesh/pkg/config"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var Validate = validator.New()

var AccessTokenAuth *jwtauth.JWTAuth = jwtauth.New(
	"HS256",
	[]byte(config.JWTConfig.ACCESS_TOKEN_SECRET),
	nil,
)

var RefreshTokenAuth *jwtauth.JWTAuth = jwtauth.New(
	"HS256",
	[]byte(config.JWTConfig.REFRESH_TOKEN_SECRET),
	nil,
)

type JWTTokenType string

var (
	ACCESS_TOKEN  JWTTokenType = "ACCESS_TOKEN"
	REFRESH_TOKEN JWTTokenType = "REFRESH_TOKEN"
)

func SendAPIResponse(w http.ResponseWriter, status int, data any, cookies ...*http.Cookie) error {
	if len(cookies) > 0 {
		for _, cookie := range cookies {
			http.SetCookie(w, cookie)
		}
	}

	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(generateAPIResponseBody(status, data))
}

func SendAPIRedirect(
	w http.ResponseWriter,
	r *http.Request,
	url string,
	status int,
	cookies ...*http.Cookie,
) {
	if len(cookies) > 0 {
		for _, cookie := range cookies {
			http.SetCookie(w, cookie)
		}
	}

	http.Redirect(w, r, url, status)
}

func SendAPIErrorResponse(w http.ResponseWriter, status int, err interface{}) {
	if e, ok := err.(error); ok {
		SendAPIResponse(w, status, map[string]interface{}{"message": e.Error()})
	} else {
		SendAPIResponse(w, status, map[string]interface{}{"message": err})
	}
}

func generateAPIResponseBody(status int, data any) map[string]any {
	if status >= 400 {
		return map[string]any{"status": status, "error": data}
	}
	return map[string]any{"status": status, "data": data}
}

func ParseJSON(r *http.Request, v interface{}) error {
	if r.Body == nil {
		return fmt.Errorf("request body should not be empty")
	}
	return json.NewDecoder(r.Body).Decode(v)
}

func generateMsgForField(fe validator.FieldError, v interface{}) (string, string) {
	t := reflect.TypeOf(v)

	field, _ := t.FieldByName(fe.StructField())

	aliasTag := field.Tag.Get("alias")

	switch fe.Tag() {
	case "required":
		return aliasTag, fmt.Sprintf("\"%s\" is required", aliasTag)
	case "email":
		return aliasTag, fmt.Sprintf("\"%s\" must be a valid email address", aliasTag)
	case "min":
		return aliasTag, fmt.Sprintf(
			"\"%s\" should contain at least %s characters",
			aliasTag,
			fe.Param(),
		)
	case "max":
		return aliasTag, fmt.Sprintf(
			"\"%s\" should contain at most %s characters",
			aliasTag,
			fe.Param(),
		)
	case "dive":
		return aliasTag, fmt.Sprintf("\"%s\" should be in an array", aliasTag)
	case "oneof":
		return aliasTag, fmt.Sprintf("\"%s\" should be one of [%s]", aliasTag, fe.Param())
	case "alphanum":
		return aliasTag, fmt.Sprintf("\"%s\" should be alpha numerical", aliasTag)
	case "lowercase":
		return aliasTag, fmt.Sprintf("\"%s\" should be all lower case", aliasTag)
	case "uuid":
		return aliasTag, fmt.Sprintf("\"%s\" should be a valid UUID", aliasTag)
	}

	return fe.Field(), fe.Error()
}

func GenerateValidationErrorObject(ve validator.ValidationErrors, v interface{}) map[string]string {
	errs := map[string]string{}
	for _, fe := range ve {
		key, value := generateMsgForField(fe, v)
		errs[key] = value
	}
	return errs
}

func PingPostgresConnection(poolConfig *pgxpool.Config) error {
	connPool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return err
	}
	conn, err := connPool.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()
	return conn.Ping(context.Background())
}

func GetNewConnectionFromPool(
	ctx context.Context,
	poolConfig *pgxpool.Config,
) (*pgxpool.Conn, error) {
	conn, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, err
	}

	return conn.Acquire(ctx)
}

func WithTransaction(ctx context.Context, conn *pgxpool.Conn, fn func(pgx.Tx) error) error {
	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if err := fn(tx); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func GetAccessTokenExpirationTime() time.Time {
	return time.Now().
		Add(time.Hour * time.Duration(config.JWTConfig.ACCESS_TOKEN_EXPIRATION_IN_HOURS))
}

func GetRefreshTokenExpirationTime() time.Time {
	return time.Now().
		Add(time.Duration(config.JWTConfig.REFRESH_TOKEN_EXPIRATION_IN_DAYS) * 24 * time.Hour)
}

func SignJWTToken(userId string, tokenType JWTTokenType) (string, error) {
	claims := map[string]interface{}{"sub": userId}
	jwtauth.SetIssuedNow(claims)
	switch tokenType {
	case ACCESS_TOKEN:
		jwtauth.SetExpiry(claims, GetAccessTokenExpirationTime())
		_, tokenString, err := AccessTokenAuth.Encode(claims)
		return tokenString, err
	case REFRESH_TOKEN:
		jwtauth.SetExpiry(claims, GetRefreshTokenExpirationTime())
		_, tokenString, err := RefreshTokenAuth.Encode(claims)
		return tokenString, err
	default:
		return "", fmt.Errorf("invalid token type")
	}
}

func SetHTTPCookie(w http.ResponseWriter, token string, tokenType JWTTokenType) error {
	switch tokenType {
	case ACCESS_TOKEN:
		accessTokenExpiration := GetAccessTokenExpirationTime()
		http.SetCookie(w, &http.Cookie{
			Name:     "access_token",
			Value:    token,
			Path:     "/",
			MaxAge:   int(accessTokenExpiration.Unix()),
			Secure:   config.AppConfig.ENVIRONMENT != "DEVELOPMENT",
			HttpOnly: false,
			SameSite: http.SameSiteLaxMode,
			Expires:  accessTokenExpiration,
		})
		return nil
	case REFRESH_TOKEN:
		refreshTokenExpiration := GetRefreshTokenExpirationTime()
		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    token,
			Path:     "/",
			MaxAge:   int(refreshTokenExpiration.Unix()),
			Secure:   config.AppConfig.ENVIRONMENT != "DEVELOPMENT",
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
			Expires:  refreshTokenExpiration,
		})
		return nil
	default:
		return fmt.Errorf("invalid token type")
	}
}

func GetGoogleHttpClient(accessToken, refreshToken string) *http.Client {
	token := &oauth2.Token{AccessToken: accessToken, RefreshToken: refreshToken}

	conf := &oauth2.Config{
		ClientID:     config.OAuthConfig.GOOGLE.CLIENT_ID,
		ClientSecret: config.OAuthConfig.GOOGLE.CLIENT_SECRET,
		Scopes:       config.OAuthConfig.GOOGLE.SCOPES,
		RedirectURL:  config.OAuthConfig.GOOGLE.CALLBACK_URL,
		Endpoint:     google.Endpoint,
	}

	tokenSource := conf.TokenSource(context.Background(), token)

	reusableTokenSource := oauth2.ReuseTokenSource(nil, tokenSource)

	return oauth2.NewClient(context.Background(), reusableTokenSource)
}
