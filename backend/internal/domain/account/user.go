package account

import (
	"strings"

	"errors"
)

type User struct {
	id        string
	firstName string
	lastName  string
	email     string
	role      Role
	accountID string
}

func (u *User) ID() string {
	return u.id
}

func (u *User) FirstName() string {
	return u.firstName
}

func (u *User) LastName() string {
	return u.lastName
}

func (u *User) Email() string {
	return u.email
}

func (u *User) Role() Role {
	return u.role
}

func (u *User) AccountID() string {
	return u.accountID
}

func NewUser(
	id string,
	firstName string,
	lastName string,
	email string,
	role string,
	accountID string,
) (*User, error) {
	firstName = strings.TrimSpace(firstName)
	lastName = strings.TrimSpace(lastName)
	email = strings.TrimSpace(email)

	userRole, err := NewRole(role)
	if err != nil {
		return nil, err
	}

	if email == "" {
		return nil, errors.New("email cannot be invalid")
	}

	return &User{
		id:        id,
		firstName: firstName,
		lastName:  lastName,
		email:     email,
		role:      userRole,
		accountID: accountID,
	}, nil
}
