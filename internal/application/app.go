package application

import (
	"context"
	"fmt"
	"github.com/Magic-Kot/store/pkg/masker"
	"net"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/benbjohnson/clock"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/jackc/pgx/v5/stdlib" // postgres driver
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
	"golang.org/x/sync/errgroup"

	"github.com/Magic-Kot/store/internal/application/rest"
	"github.com/Magic-Kot/store/internal/config"
	"github.com/Magic-Kot/store/internal/domain/component"
	"github.com/Magic-Kot/store/internal/domain/service"
	"github.com/Magic-Kot/store/internal/infrastructure/persistence"
	"github.com/Magic-Kot/store/pkg/logging"
	"github.com/Magic-Kot/store/pkg/middlewarex"
)

type App struct {
	name     string
	version  string
	cfg      config.Config
	deferred []func()

	httpServer        *http.Server
	privateHTTPServer *http.Server
	authRepository    *persistence.AuthPostgresRepository
	authService       service.Auth
	postgresClient    *sqlx.DB

	redisClient   *redis.Client
	refreshTokens persistence.RefreshTokens

	clock clock.Clock

	dbUser      persistence.DBUser
	userService *service.User
}

func New(name, version string, cfg config.Config) *App {
	return &App{ //nolint:exhaustruct
		name:    name,
		version: version,
		cfg:     cfg,
	}
}

func (app *App) Run() error {
	defer app.shutdown()

	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	defer stop()

	logger, err := logging.NewLogger(&logging.LoggerDeps{Level: app.cfg.Logger.Level})
	if err != nil {
		log.Fatal().Err(err).Msg("failed to init logger")
	}

	ctx = logger.WithContext(ctx)

	g, ctx := errgroup.WithContext(ctx)

	app.clock = clock.New()

	postgresClient, err := sqlx.ConnectContext(ctx, "pgx", app.cfg.Postgres.DSN)
	if err != nil {
		return fmt.Errorf("sqlx.ConnectContext: %w", err)
	}

	postgresClient.SetMaxOpenConns(app.cfg.Postgres.MaxOpenConns)
	postgresClient.SetMaxIdleConns(app.cfg.Postgres.MaxIdleConns)
	postgresClient.SetConnMaxLifetime(app.cfg.Postgres.ConnMaxLifetime)

	logger.Info().Str("database", app.cfg.Postgres.DSN).Msg("connected to postgres")

	app.postgresClient = postgresClient

	app.redisClient = lo.Must(
		ConnectRedis(
			ctx,
			app.cfg.Redis.Address,
			app.cfg.Redis.Username,
			app.cfg.Redis.Password,
			app.cfg.Redis.DatabaseNumber,
			app.cfg.Redis.PoolSize,
			app.cfg.Redis.MinIdleConnections,
			app.cfg.Redis.MaxIdleConnections,
		),
	)

	app.refreshTokens = persistence.NewRefreshTokens(app.redisClient)

	// auth
	app.authRepository = persistence.NewAuthPostgresRepository(app.postgresClient)
	app.authService = service.NewAuth(
		component.NewRefreshTokenParser(app.cfg.JWT.PublicKey),
		component.NewAccessTokenGenerator(app.cfg.JWT.PrivateKey, app.cfg.JWT.AccessTokenTTL, app.clock),
		component.NewRefreshTokenGenerator(app.cfg.JWT.PrivateKey, app.cfg.JWT.RefreshTokenTTL, app.clock),
		app.refreshTokens,
		app.cfg.JWT.RefreshTokenTTL,
		app.authRepository,
	)

	// user
	app.dbUser = persistence.NewDBUser(app.postgresClient)
	app.userService = service.NewUser(app.dbUser)

	// Referral
	//referralRepository := postgres.NewReferralRepository(app.postgresClient)
	//redisURL := redis.NewReferralRepository(app.redisClient)
	//referralService := referral.NewReferralService(referralRepository, redisURL)
	//referralController := rest.NewApiReferralController(referralService, logger, validate)

	app.runHTTPServer(ctx, g)

	if err = g.Wait(); err != nil {
		return fmt.Errorf("g.Wait: %w", err)
	}

	return nil
}

func (app *App) shutdown() {
	for _, fn := range app.deferred {
		fn()
	}
}

func (app *App) runHTTPServer(ctx context.Context, g *errgroup.Group) {
	app.httpServer = app.newHTTPServer(ctx)

	g.Go(func() error {
		go func() {
			<-ctx.Done()

			ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), app.cfg.ServerHTTP.ShutdownTimeout) //nolint:govet
			defer cancel()

			if err := app.httpServer.Shutdown(ctx); err != nil {
				zerolog.Ctx(ctx).Error().Err(err).Msg("httpServer.Shutdown")
			}
		}()

		zerolog.Ctx(ctx).Info().Str("address", app.cfg.ServerHTTP.ListenAddress).Msg("http server started")

		if err := app.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("httpServer.ListenAndServe: %w", err)
		}

		zerolog.Ctx(ctx).Info().Msg("http server stopped")

		return nil
	})
}

func (app *App) newHTTPServer(ctx context.Context) *http.Server {
	router := chi.NewRouter()

	router.Use(
		middleware.RealIP,
		middlewarex.TraceID,
		middlewarex.Logger,
		middlewarex.RequestLogging(masker.NewNopSensitiveDataMasker(), app.cfg.Logger.FieldMaxLen),
		middlewarex.ResponseLogging(masker.NewNopSensitiveDataMasker(), app.cfg.Logger.FieldMaxLen),
		middlewarex.Recovery,
	)

	restServer := rest.NewServer(app.authService, app.userService)

	authMiddleware := rest.NewBearerAuth(middlewarex.NewHeaderAuthorizationBearerTokenFinder(), component.NewAccessTokenParser(app.cfg.JWT.PublicKey))
	rest.RegisterRoutes(router, restServer, authMiddleware)

	return &http.Server{ //nolint:exhaustruct
		Addr:              app.cfg.ServerHTTP.ListenAddress,
		WriteTimeout:      app.cfg.ServerHTTP.WriteTimeout,
		ReadTimeout:       app.cfg.ServerHTTP.ReadTimeout,
		ReadHeaderTimeout: app.cfg.ServerHTTP.ReadTimeout,
		IdleTimeout:       app.cfg.ServerHTTP.IdleTimeout,
		Handler:           router,
		BaseContext: func(net.Listener) context.Context {
			return ctx
		},
	}
}

func ConnectRedis(
	ctx context.Context,
	address,
	username,
	password string,
	databaseNumber,
	poolSize,
	minIdleConnections,
	maxIdleConnections int,
) (*redis.Client, error) {
	redisClient := redis.NewClient(&redis.Options{ //nolint:exhaustruct
		Network:      "tcp",
		Addr:         address,
		Username:     username,
		Password:     password,
		DB:           databaseNumber,
		PoolSize:     poolSize,
		MinIdleConns: minIdleConnections,
		MaxIdleConns: maxIdleConnections,
	})

	if err := redisClient.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redisClient.Ping: %w", err)
	}

	zerolog.Ctx(ctx).Info().Str("address", address).Int("database-number", databaseNumber).Msg("connected to redis")

	return redisClient, nil
}
