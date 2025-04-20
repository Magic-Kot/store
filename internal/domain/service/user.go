package service

import (
	"context"
	"fmt"

	"github.com/Magic-Kot/store/internal/domain/entity"
	"github.com/Magic-Kot/store/internal/domain/value"
)

type DBUser interface {
	UserInfo(context.Context, value.PersonID) (entity.UserData, error)
}

type User struct {
	user DBUser
}

func NewUser(
	user DBUser,
) *User {
	return &User{
		user: user,
	}
}

func (u *User) UserInfo(ctx context.Context, personID value.PersonID) (entity.UserData, error) {
	user, err := u.user.UserInfo(ctx, personID)
	if err != nil {
		return entity.UserData{}, fmt.Errorf("user.UserInfo: %w", err)
	}

	return user, nil
}
