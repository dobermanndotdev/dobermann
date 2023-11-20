package ddd

import (
	"errors"
	"strings"

	"github.com/oklog/ulid/v2"
)

type ID struct {
	value ulid.ULID
}

func NewID() ID {
	return ID{value: ulid.Make()}
}

func NewIdFromString(id string) (ID, error) {
	id = strings.TrimSpace(id)

	if id == "" {
		return ID{}, errors.New("id cannot be empty")
	}

	value, err := ulid.Parse(id)
	if err != nil {
		return ID{}, errors.New("id cannot be invalid")
	}

	return ID{value: value}, nil
}

func (i ID) IsEmpty() bool {
	return i.value == ulid.ULID{}
}

func (i ID) String() string {
	return i.value.String()
}
