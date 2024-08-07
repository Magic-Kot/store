package httpserver

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type ServerDeps struct {
	Host    string
	Port    string
	Timeout time.Duration
}

type Server struct {
	host    string
	port    string
	timeout time.Duration
	serv    *echo.Echo
}

func NewServer(deps *ServerDeps) *Server {
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
