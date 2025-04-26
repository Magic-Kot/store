package service

import (
	"context"
	"fmt"

	"git.appkode.ru/pub/go/failure"

	"github.com/Magic-Kot/store/internal/domain/entity"
	"github.com/Magic-Kot/store/internal/domain/value"
	"github.com/Magic-Kot/store/pkg/errcodes"
)

type DBUser interface {
	UserInfo(context.Context, value.PersonID) (entity.UserData, error)
	UpdateUser(context.Context, value.PersonID, int, string, string, string) error
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

func (u *User) UpdateUser(ctx context.Context, personID value.PersonID, user entity.UserData) error {
	if err := u.user.UpdateUser(ctx, personID, user.Age, user.Username, user.Name, user.Surname); err != nil {
		if failure.IsNotFoundError(err) {
			return failure.NewNotFoundError(
				fmt.Sprintf("user with personID %d not found", personID),
				failure.WithCode(errcodes.UserNotFound),
				failure.WithDescription(errcodes.UserNotFoundMessage.String()))
		}

		return fmt.Errorf("user.UpdateUser: %w", err)
	}

	return nil
}
