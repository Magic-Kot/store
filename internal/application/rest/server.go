package rest

import (
	"context"

	"github.com/Magic-Kot/store/internal/domain/value"
)

type authService interface {
	Registration(context.Context, value.Login, value.Password) error
	Authenticate(context.Context, value.Login, value.Password) (value.TokenPair, error)
	Refresh(context.Context, value.RefreshToken) (value.TokenPair, error)
	Logout(context.Context, value.RefreshToken) error
}

type Server struct {
	auth authService
}

func NewServer(
	auth authService,
) Server {
	return Server{
		auth: auth,
	}
}
