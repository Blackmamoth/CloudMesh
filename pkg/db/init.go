package db

import (
	"fmt"
	"time"

	"github.com/blackmamoth/cloudmesh/pkg/config"
	"github.com/blackmamoth/cloudmesh/pkg/utils"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

var PoolConfig *pgxpool.Config

func init() {
	config.LOGGER.Sync()
	poolConfig, err := connectPostgres()
	if err != nil {
		config.LOGGER.Fatal("Application disconnected from PostgreSQL Server", zap.Error(err))
	}

	if err := utils.PingPostgresConnection(poolConfig); err != nil {
		config.LOGGER.Fatal("Application disconnected from PostgreSQL Server", zap.Error(err))
	}

	config.LOGGER.Info("Application connected to PostgreSQL Server")

	PoolConfig = poolConfig
}

func connectPostgres() (*pgxpool.Config, error) {
	dsn := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=%s",
		config.PostgresConfig.POSTGRES_USER,
		config.PostgresConfig.POSTGRES_PASS,
		config.PostgresConfig.POSTGRES_HOST,
		config.PostgresConfig.POSTGRES_PORT,
		config.PostgresConfig.POSTGRES_DBNAME,
		config.PostgresConfig.POSTGRES_SSLMODE,
	)

	conn, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	conn.MaxConns = int32(10)
	conn.MinConns = int32(0)
	conn.MaxConnLifetime = time.Hour
	conn.MaxConnIdleTime = time.Minute * 30
	conn.HealthCheckPeriod = time.Minute
	conn.ConnConfig.ConnectTimeout = time.Second * 30

	return conn, nil
}
