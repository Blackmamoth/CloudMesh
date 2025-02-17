package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/blackmamoth/cloudmesh/pkg/config"
	"github.com/blackmamoth/cloudmesh/pkg/middlewares"
	"github.com/blackmamoth/cloudmesh/pkg/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type APIServer struct {
	host       string
	addr       string
	poolConfig *pgxpool.Config
}

func NewAPIServer(host, addr string, poolConfig *pgxpool.Config) *APIServer {
	return &APIServer{
		host:       host,
		addr:       addr,
		poolConfig: poolConfig,
	}
}

func (s *APIServer) Run() error {
	defer config.LOGGER.Sync()
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middlewares.HttpRequestLogger)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(middleware.Heartbeat("/ping"))
	r.Use(middleware.Compress(5, "gzip"))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		utils.SendAPIResponse(
			w,
			http.StatusOK,
			map[string]interface{}{
				"message": "Welcome to CloudMesh! Your gateway to effortlessly connecting and managing your cloud storage. Let's get started!",
			},
		)
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		utils.SendAPIErrorResponse(
			w,
			http.StatusNotFound,
			fmt.Sprintf("route not found for [%s] %s", r.Method, r.URL.Path),
		)
	})

	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		utils.SendAPIErrorResponse(
			w,
			http.StatusMethodNotAllowed,
			fmt.Sprintf("method [%s] not allowed for route %s", r.Method, r.URL.Path),
		)
	})

	config.LOGGER.Info(
		"Application started",
		zap.String("host", s.host),
		zap.String("port", s.addr),
	)

	return http.ListenAndServe(fmt.Sprintf("%s:%s", s.host, s.addr), r)
}
