package account

import (
	"errors"
	"regexp"
	"strings"
)

// loose validation
var emailRegexp = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

type Email struct {
	address string
}

func NewEmail(address string) (Email, error) {
	address = strings.TrimSpace(address)
	if address == "" {
		return Email{}, errors.New("address cannot be empty")
	}

	if !emailRegexp.MatchString(address) {
		return Email{}, errors.New("the address provided is invalid")
	}

	return Email{
		address: address,
	}, nil
}

func (e Email) Address() string {
	return e.address
}

func (e Email) IsEmpty() bool {
	return e.address == ""
}
