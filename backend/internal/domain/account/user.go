package account

import (
	"errors"
	"strings"

	"github.com/flowck/doberman/internal/common/ddd"
	"github.com/flowck/doberman/internal/domain"
)

type User struct {
	id               ddd.ID
	firstName        string
	lastName         string
	email            ddd.Email
	password         Password
	role             Role
	confirmationCode domain.ID
}

func NewUser(
	id ddd.ID,
	firstName string,
	lastName string,
	email ddd.Email,
	password Password,
	role Role,
	confirmationCode domain.ID,
) (*User, error) {
	firstName = strings.TrimSpace(firstName)
	lastName = strings.TrimSpace(lastName)

	if id.IsEmpty() {
		return nil, errors.New("id cannot be invalid")
	}

	if firstName == "" {
		return nil, errors.New("firstName cannot be empty")
	}

	if lastName == "" {
		return nil, errors.New("lastName cannot be empty")
	}

	if role.IsEmpty() {
		return nil, errors.New("role cannot be invalid")
	}

	if password.IsEmpty() {
		return nil, errors.New("password cannot be invalid")
	}

	return &User{
		id:               id,
		email:            email,
		role:             role,
		password:         password,
		lastName:         lastName,
		firstName:        firstName,
		confirmationCode: confirmationCode,
	}, nil
}

func (u *User) ID() ddd.ID {
	return u.id
}

func (u *User) FirstName() string {
	return u.firstName
}

func (u *User) LastName() string {
	return u.lastName
}

func (u *User) Email() ddd.Email {
	return u.email
}

func (u *User) Password() Password {
	return u.password
}

func (u *User) Role() Role {
	return u.role
}

func (u *User) ConfirmationCode() domain.ID {
	return u.confirmationCode
}
