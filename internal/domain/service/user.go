package service

import (
	"context"
	"fmt"
	"time"

	"git.appkode.ru/pub/go/failure"

	"github.com/Magic-Kot/store/internal/domain/entity"
	"github.com/Magic-Kot/store/internal/domain/value"
	"github.com/Magic-Kot/store/pkg/errcodes"
)

type DBUser interface {
	UserInfo(context.Context, value.PersonID) (entity.UserData, error)
	UpdateUser(context.Context, value.PersonID, int, string, string, string) error
	DeletedAtByPersonID(context.Context, value.PersonID) (*time.Time, error)
	Remove(context.Context, value.PersonID) error
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

func (u *User) RemoveUser(ctx context.Context, personID value.PersonID) error {
	deletedAt, err := u.user.DeletedAtByPersonID(ctx, personID)
	if err != nil {
		if failure.IsNotFoundError(err) {
			return failure.NewNotFoundError(
				fmt.Sprintf("user with personID %d not found", personID),
				failure.WithCode(errcodes.UserNotFound),
				failure.WithDescription(errcodes.UserNotFoundMessage.String()))
		}

		return fmt.Errorf("user.DeletedAtByPersonID: %w", err)
	}

	if deletedAt != nil {
		return failure.NewInvalidArgumentError(
			fmt.Sprintf("push token with personID: %s already deleted", personID),
			failure.WithDescription(errcodes.PushTokenAlreadyDeletedMessage))
	}

	if err = u.user.Remove(ctx, personID); err != nil {
		return fmt.Errorf("user.Remove: %w", err)
	}

	return nil
}
