package utils

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
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

func GetNewConnectionFromPool(poolConfig *pgxpool.Config) (*pgxpool.Conn, error) {
	connPool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, err
	}

	conn, err := connPool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	return conn, nil
}
