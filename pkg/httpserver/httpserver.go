package httpserver

import (
	"context"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type ConfigDeps struct {
	Host    string
	Port    string
	Timeout time.Duration
	//Logger  *zerolog.Logger
}

type Server struct {
	host    string
	port    string
	timeout time.Duration
	serv    *echo.Echo
	//logger  *zerolog.Logger
}

func NewServer(deps *ConfigDeps) *Server {
	s := echo.New()
	s.Use(middleware.Recover())

	return &Server{
		host:    deps.Host,
		port:    deps.Port,
		timeout: deps.Timeout,
		serv:    s,
	}
}

func (s *Server) Start() error {
	if err := s.serv.Start(s.host + s.port); err != nil {
		return err
	}

	return nil
}

func (s Server) Server() *echo.Echo {
	return s.serv
}

func (s Server) Shutdown(ctx context.Context) error {
	return s.serv.Shutdown(ctx)
}
