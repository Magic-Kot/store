package main

import (
	"context"
	"os"
	"time"

	"github.com/labstack/gommon/log"

	"online-store/internal/controllers"
	"online-store/internal/delivery/httpecho"
	"online-store/internal/repository/postgres"
	"online-store/internal/services/user"
	"online-store/pkg/client/postg"
	"online-store/pkg/httpserver"
)

func main() {
	serverDeps := &httpserver.ServerDeps{
		Host:    "localhost",
		Port:    ":8080",
		Timeout: 5 * time.Second,
	}

	server := httpserver.NewServer(serverDeps)

	// create client
	// postgres://postgres:12345@localhost:5438
	pool, err := postg.NewClient(context.TODO(), 5, "postgres", "12345", "localhost", "5438", "postgres")
	if err != nil {
		log.Error("failed to init storage:", err)
		os.Exit(1)
	}

	// create repository
	db := postgres.NewUserRepository(pool)

	// create service(repository, ссылка на репозиторий)
	service := user.NewUserService(db)

	contr := controllers.NewApiController(service)

	httpecho.SetUserRoutes(server.Server(), contr)

	if err := server.Start(); err != nil {
		panic(err)
	}
}
