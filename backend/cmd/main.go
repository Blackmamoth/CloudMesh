package main

import (
	"github.com/blackmamoth/cloudmesh/cmd/api"
	"github.com/blackmamoth/cloudmesh/pkg/config"
	"github.com/blackmamoth/cloudmesh/pkg/db"
	"go.uber.org/zap"
)

func main() {
	defer config.LOGGER.Sync()

	apiServer := api.NewAPIServer(config.AppConfig.HOST, config.AppConfig.PORT, db.PoolConfig)

	if err := apiServer.Run(); err != nil {
		config.LOGGER.Fatal("Application terminated", zap.Error(err))
	}
}
