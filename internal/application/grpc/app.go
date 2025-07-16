package grpcapp

import (
	"fmt"
	"net"

	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

type App struct {
	gRPCServer *grpc.Server
	port       string
	logger     *zerolog.Logger
}

func NewApp(port string, logger *zerolog.Logger) *App {
	gRPCServer := grpc.NewServer()

	Register(gRPCServer)

	return &App{
		gRPCServer: gRPCServer,
		port:       port,
		logger:     logger,
	}
}

func (a *App) Run(g *errgroup.Group) {
	g.Go(func() error {
		l, err := net.Listen("tcp", a.port)
		if err != nil {
			return fmt.Errorf("net.Listen: %w", err)
		}

		a.logger.Info().Str("port", a.port).Msg("grpc server started")

		if err = a.gRPCServer.Serve(l); err != nil {
			return fmt.Errorf("gRPCServer.Serve: %w", err)
		}

		return nil
	})

}

func (a *App) Stop() {
	a.logger.Info().Msg("stopping grpc server")

	a.gRPCServer.GracefulStop()
}
