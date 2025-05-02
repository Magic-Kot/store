package rest

import (
	"context"

	"github.com/nats-io/nats.go"

	"github.com/Magic-Kot/store/internal/domain/entity"
	"github.com/Magic-Kot/store/internal/domain/value"
)

type authService interface {
	Registration(context.Context, value.Login, value.Password) error
	Authenticate(context.Context, value.Login, value.Password) (value.TokenPair, error)
	Refresh(context.Context, value.RefreshToken) (value.TokenPair, error)
	Logout(context.Context, value.RefreshToken) error
}

type userService interface {
	UserInfo(context.Context, value.PersonID) (entity.UserData, error)
	UpdateUser(context.Context, value.PersonID, entity.UserData) error
	RemoveUser(context.Context, value.PersonID) error
}

type Server struct {
	auth       authService
	user       userService
	natsClient *nats.Conn
}

func NewServer(
	auth authService,
	user userService,
	natsClient *nats.Conn,
) Server {
	return Server{
		auth:       auth,
		user:       user,
		natsClient: natsClient,
	}
}
