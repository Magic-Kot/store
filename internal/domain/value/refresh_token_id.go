package value

import (
	"fmt"

	"github.com/rs/xid"
)

type RefreshTokenID struct{ xid.ID }

func NewRefreshTokenID(id xid.ID) RefreshTokenID {
	return RefreshTokenID{id}
}

func ParseRefreshTokenID(raw string) (RefreshTokenID, error) {
	id, err := xid.FromString(raw)
	if err != nil {
		return RefreshTokenID{}, fmt.Errorf("xid.FromString(%s): %w", raw, err)
	}

	return RefreshTokenID{id}, nil
}
