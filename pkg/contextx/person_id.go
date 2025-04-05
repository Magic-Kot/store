package contextx

import (
	"context"
	"fmt"

	"github.com/Magic-Kot/store/internal/domain/value"
)

type contextKeyPersonID struct{}

func WithPersonID(ctx context.Context, personID value.PersonID) context.Context {
	return context.WithValue(ctx, contextKeyPersonID{}, personID)
}

func PersonIDFromContext(ctx context.Context) (value.PersonID, error) {
	personID, ok := ctx.Value(contextKeyPersonID{}).(value.PersonID)
	if !ok {
		return personID, fmt.Errorf("personID: %w", ErrNoValue)
	}

	return personID, nil
}
