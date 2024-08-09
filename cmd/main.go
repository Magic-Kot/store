package main

import (
	"context"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"online-store/internal/config"
	"online-store/pkg/logging"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"online-store/internal/controllers"
	"online-store/internal/delivery/httpecho"
	"online-store/internal/repository/postgres"
	"online-store/internal/services/user"
	"online-store/pkg/client/postg"
	"online-store/pkg/httpserver"
)

func main() {
	// create logger
	ctx := context.Background()

	ctx, err := logging.NewLogger(ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed to init logger")
		os.Exit(1)
	}

	logger := zerolog.Ctx(ctx)

	logger.Info().Msg("Start server")
	log.Error().Err(err).Msg("failed to init logger")

	log.Debug().Msg("log level set to Debug")
	log.Info().Msg("log level set to Info")

	// read config
	//var cfg httpserver.ServerDeps
	var cfg config.Config

	err = cleanenv.ReadConfig("../config.yml", &cfg) //задан путь до файла конфигурации
	if err != nil {
		log.Error().Err(err).Msg("error initializing config")
		os.Exit(1)
	}

	// create server
	server := httpserver.NewServer(&cfg.ServerDeps)

	// create client
	//cfgRepo := postg.Config{
	//	MaxAttempts: 3,
	//	Username:    "postgres",
	//	Password:    "12345",
	//	Host:        "localhost",
	//	Port:        "5438",
	//	Database:    "postgres",
	//}

	pool, err := postg.NewClient(context.TODO(), &cfg.RepositoryConfig)
	if err != nil {
		//log.Error("failed to init storage:", err)
		os.Exit(1)
	}

	// create repository
	db := postgres.NewUserRepository(pool)

	// create service
	service := user.NewUserService(db)

	// create controller
	contr := controllers.NewApiController(service)

	// set routes
	httpecho.SetUserRoutes(server.Server(), contr)

	// start server
	if err := server.Start(); err != nil {
		panic(err)
	}
}
