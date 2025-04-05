package value

import (
	"fmt"
	"strconv"
)

type PersonID uint64

func NewPersonIDFromString(raw string) (PersonID, error) {
	personID, err := strconv.ParseUint(raw, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("strconv.ParseUint(%s): %w", raw, err)
	}

	return PersonID(personID), nil
}

func (p PersonID) String() string { return strconv.FormatUint(uint64(p), 10) }

func (p PersonID) Uint() uint {
	personID := uint(p)

	return personID
}
