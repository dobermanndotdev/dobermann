package monitor

import (
	"net/url"

	"github.com/friendsofgo/errors"
)

type Endpoint struct {
	value *url.URL
}

func NewEndpoint(value string) (Endpoint, error) {
	u, err := url.Parse(value)
	if err != nil {
		return Endpoint{}, errors.New("the endpoint address provided cannot be invalid")
	}

	return Endpoint{value: u}, nil
}

func (e Endpoint) IsValid() bool {
	return e.value != nil && e.value.Scheme != ""
}

func (e Endpoint) String() string {
	return e.value.String()
}
