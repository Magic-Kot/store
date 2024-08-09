package httpserver

import (
	"net/http"
	"online-store/internal/config"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

//type ServerDeps struct {
//	Host    string        `yaml:"host" env:"HOST" env-default:"localhost"`
//	Port    string        `yaml:"port" env:"PORT" env-default:":8000"`
//	Timeout time.Duration `yaml:"timeout" env:"TIMEOUT" env-default:"5s"`
//}

type Server struct {
	host    string
	port    string
	timeout time.Duration
	serv    *echo.Echo
}

// func NewServer(deps *ServerDeps) *Server {
func NewServer(deps *config.ServerDeps) *Server {
	s := echo.New()
	s.Use(middleware.Recover())

	s.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

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

// getter
func (s *Server) Server() *echo.Echo {
	return s.serv
}

// setter
func (s *Server) SetTimeout(timeout time.Duration) {
	s.timeout = timeout
}
