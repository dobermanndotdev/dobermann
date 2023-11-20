package account

import (
	"fmt"

	"github.com/nyaruka/phonenumbers"
)

type PhoneNumber struct {
	value *phonenumbers.PhoneNumber
}

func NewPhoneNumber(phoneNumber string, countryCode string) (PhoneNumber, error) {
	pn, err := phonenumbers.Parse(phoneNumber, countryCode)
	if err != nil {
		return PhoneNumber{}, fmt.Errorf("unable to parse %s with country code %s: %v", phoneNumber, countryCode, err)
	}

	return PhoneNumber{value: pn}, nil
}
