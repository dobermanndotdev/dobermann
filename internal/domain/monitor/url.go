package monitor

import (
	"fmt"
	"net/url"
	"strings"
)

type URL struct {
	value *url.URL
}

func NewURL(value string) (URL, error) {
	value = strings.TrimSpace(value)

	parsedURL, err := url.Parse(value)
	if err != nil {
		return URL{}, fmt.Errorf("%s is not a valid URL", value)
	}

	if parsedURL.Scheme == "" {
		return URL{}, fmt.Errorf("the url must start with http/https/tcp instead of '%s'", parsedURL.Scheme)
	}

	return URL{
		value: parsedURL,
	}, nil
}

func (u URL) String() string {
	return u.value.String()
}

func (u URL) IsEmpty() bool {
	return u.value == nil || u.value.Host == ""
}
