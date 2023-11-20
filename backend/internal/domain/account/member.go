package account

import (
	"net/url"
	"time"

	"github.com/flowck/doberman/internal/domain"
)

type Member struct {
	id             domain.ID
	firstName      string
	lastName       string
	email          string
	avatar         *url.URL
	timezone       time.Location
	phoneNumbers   []PhoneNumber
	holidayDetails HolidayDetails
}

type HolidayDetails struct {
	onHolidays bool
	until      time.Time
}

func (h HolidayDetails) OnHolidays() bool {
	return h.onHolidays
}

func (h HolidayDetails) Until() time.Time {
	return h.until
}

func NewMember(
	id domain.ID,
	firstName string,
	lastName string,
	email string,
	avatar *url.URL,
	timezone time.Location,
	phoneNumbers []PhoneNumber,
	holidayDetails HolidayDetails,
) (*Member, error) {
	return &Member{
		id:             id,
		firstName:      firstName,
		lastName:       lastName,
		email:          email,
		avatar:         avatar,
		timezone:       timezone,
		phoneNumbers:   phoneNumbers,
		holidayDetails: holidayDetails,
	}, nil
}
