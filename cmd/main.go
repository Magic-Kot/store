package main

import (
	"context"

	"github.com/Magic-Kot/store/internal/config"
	"github.com/Magic-Kot/store/internal/controllers"
	"github.com/Magic-Kot/store/internal/delivery/httpecho"
	"github.com/Magic-Kot/store/internal/repository/postgres"
	"github.com/Magic-Kot/store/internal/repository/redis"
	"github.com/Magic-Kot/store/internal/services/referral"
	"github.com/Magic-Kot/store/internal/services/user"
	"github.com/Magic-Kot/store/pkg/client/postg"
	"github.com/Magic-Kot/store/pkg/client/reds"
	"github.com/Magic-Kot/store/pkg/httpserver"
	"github.com/Magic-Kot/store/pkg/logging"
	"github.com/Magic-Kot/store/pkg/utils/jwt_token"

	"github.com/go-playground/validator/v10"
	"github.com/ilyakaznacheev/cleanenv"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

func main() {
	// read config
	//var cfg httpserver.ServerDeps
	var cfg config.Config

	err := cleanenv.ReadConfig("internal/config/config.yml", &cfg)
	//err := cleanenv.ReadConfig("config.yml", &cfg) // for docker
	//err := cleanenv.ReadConfig("internal/config/config.env", &cfg)
	//err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("error initializing config")
	}

	// create logger
	logCfg := logging.LoggerDeps{
		LogLevel: cfg.LoggerDeps.LogLevel,
	}

	logger, err := logging.NewLogger(&logCfg)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to init logger")
	}

	logger.Info().Msg("init logger")

	ctx := context.Background()
	ctx = logger.WithContext(ctx)

	logger.Debug().Msgf("config: %+v", cfg)

	// create server
	serv := httpserver.ConfigDeps{
		Host:    cfg.ServerDeps.Host,
		Port:    cfg.ServerDeps.Port,
		Timeout: cfg.ServerDeps.Timeout,
		//Logger:  logger,
	}

	server := httpserver.NewServer(&serv)

	// create client Postgres
	repo := postg.ConfigDeps{
		MaxAttempts: cfg.PostgresDeps.MaxAttempts,
		Delay:       cfg.PostgresDeps.Delay,
		Username:    cfg.PostgresDeps.Username,
		Password:    cfg.PostgresDeps.Password,
		Host:        cfg.PostgresDeps.Host,
		Port:        cfg.PostgresDeps.Port,
		Database:    cfg.PostgresDeps.Database,
		SSLMode:     cfg.PostgresDeps.SSLMode,
	}

	pool, err := postg.NewClient(ctx, &repo)
	if err != nil {
		logger.Fatal().Err(err).Msgf("NewClient: %s", err)
	}

	// create client Redis for refresh tokens
	redisCfg := reds.ConfigDeps{
		Username: cfg.RedisDeps.Username,
		Password: cfg.RedisDeps.Password,
		Host:     cfg.RedisDeps.Host,
		Port:     cfg.RedisDeps.Port,
		Database: cfg.RedisDeps.Database,
	}

	clientRedis, err := reds.NewClientRedis(ctx, &redisCfg)
	if err != nil {
		logger.Fatal().Err(err).Msgf("redis refresh tokens: %s", err)
	}

	// create client Redis for referral urls
	//redisUrl := reds.ConfigDeps{
	//	Username: "reds",
	//	Password: "",
	//	Host:     "127.0.0.1",
	//	Port:     "6385",
	//	Database: "0",
	//}
	//
	//clientRedisURL, err := reds.NewClientRedis(ctx, &redisUrl)
	//if err != nil {
	//	logger.Fatal().Err(err).Msgf("redis URL: %s", err)
	//}

	// create tokenJWT
	tokenCfg := jwt_token.TokenJWTDeps{
		SigningKey:      cfg.AuthDeps.SigningKey,
		AccessTokenTTL:  cfg.AuthDeps.AccessTokenTTL,
		RefreshTokenTTL: cfg.AuthDeps.RefreshTokenTTL,
	}

	tokenJWT, err := jwt_token.NewTokenJWT(&tokenCfg)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to init tokenJWT")
	}

	// create validator
	validate := validator.New()

	// User
	userRepository := postgres.NewUserRepository(pool)                                      // create user repository
	rds := redis.NewAuthRepository(clientRedis)                                             // create auth repository
	userService := user.NewUserService(userRepository, rds, tokenJWT)                       // create service
	userController := controllers.NewApiController(userService, logger, validate, tokenJWT) // create controller
	httpecho.SetUserRoutes(server.Server(), userController)                                 // set routes

	// Referral
	referralRepository := postgres.NewReferralRepository(pool)
	redisURL := redis.NewReferralRepository(clientRedis)                                          // clientRedisURL
	referralService := referral.NewReferralService(referralRepository, redisURL)                  //
	referralController := controllers.NewApiReferralController(referralService, logger, validate) //
	httpecho.SetReferralRoutes(server.Server(), userController, referralController)

	// start server
	logger.Info().Msg("starting server")

	if err := server.Start(); err != nil {
		logger.Fatal().Msgf("serverStart: %v", err)
	}
}
