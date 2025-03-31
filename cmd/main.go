package main

import (
	"context"
	"embed"
	"time"

	"github.com/Magic-Kot/store/internal/config"
	"github.com/Magic-Kot/store/internal/controllers"
	"github.com/Magic-Kot/store/internal/delivery/httpecho"
	"github.com/Magic-Kot/store/internal/middleware"
	"github.com/Magic-Kot/store/internal/repository/postgres"
	"github.com/Magic-Kot/store/internal/repository/redis"
	"github.com/Magic-Kot/store/internal/services/auth"
	"github.com/Magic-Kot/store/internal/services/referral"
	"github.com/Magic-Kot/store/internal/services/user"
	"github.com/Magic-Kot/store/pkg/client/postg"
	"github.com/Magic-Kot/store/pkg/client/reds"
	"github.com/Magic-Kot/store/pkg/httpserver"
	"github.com/Magic-Kot/store/pkg/logging"
	"github.com/Magic-Kot/store/pkg/ossignal"
	"github.com/Magic-Kot/store/pkg/utils/jwt_token"

	"github.com/go-playground/validator/v10"
	"github.com/ilyakaznacheev/cleanenv"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/speakeasy-api/goose/v3"
	"golang.org/x/sync/errgroup"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func main() {
	// read config
	var cfg config.Config

	err := cleanenv.ReadConfig("internal/config/config.yml", &cfg) // Local: internal/config/config.yml Docker: config.yml
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

	// migrations
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}

	if err := goose.Up(pool, "migrations"); err != nil {
		panic(err)
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

	rds := redis.NewAuthRepository(clientRedis)
	middlewareUser := middleware.NewMiddleware(logger, tokenJWT)

	// Auth
	authRepository := postgres.NewAuthPostgresRepository(pool)
	authService := auth.NewAuthService(authRepository, rds, tokenJWT)
	authController := controllers.NewApiAuthController(authService, logger, validate)
	httpecho.SetAuthRoutes(server.Server(), authController)

	// User
	userRepository := postgres.NewUserRepository(pool)
	userService := user.NewUserService(userRepository, rds)
	userController := controllers.NewApiController(userService, logger, validate)
	httpecho.SetUserRoutes(server.Server(), userController, middlewareUser)

	// Referral
	referralRepository := postgres.NewReferralRepository(pool)
	redisURL := redis.NewReferralRepository(clientRedis)
	referralService := referral.NewReferralService(referralRepository, redisURL)
	referralController := controllers.NewApiReferralController(referralService, logger, validate)
	httpecho.SetReferralRoutes(server.Server(), referralController, middlewareUser)

	runner, ctx := errgroup.WithContext(ctx)

	// start server
	logger.Info().Msg("starting server")
	runner.Go(func() error {
		if err := server.Start(); err != nil {
			logger.Fatal().Msgf("%v", err)
		}

		return nil
	})

	runner.Go(func() error {
		if err := ossignal.DefaultSignalWaiter(ctx); err != nil {
			return errors.Wrap(err, "waiting os signal")
		}

		return nil
	})

	runner.Go(func() error {
		<-ctx.Done()

		ctxSignal, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		defer cancel()

		if err := server.Shutdown(ctxSignal); err != nil {
			logger.Error().Err(err).Msg("shutdown http server")
		}

		return nil
	})

	if err := runner.Wait(); err != nil {
		switch {
		case ossignal.IsExitSignal(err):
			logger.Info().Msg("exited by exit signal")
		default:
			logger.Fatal().Msgf("exiting with error: %v", err)
		}
	}
}
