package account

import (
	"strings"
	"time"

	"errors"

	"github.com/flowck/dobermann/backend/internal/domain"
)

type User struct {
	id                domain.ID
	firstName         string
	lastName          string
	email             string
	password          string
	role              Role
	verificationToken *domain.ID
	verifiedAt        *time.Time
	accountID         domain.ID
	createdAt         time.Time
}

func (u *User) VerificationToken() *domain.ID {
	return u.verificationToken
}

func (u *User) VerifiedAt() *time.Time {
	return u.verifiedAt
}

func (u *User) CreatedAt() time.Time {
	return u.createdAt
}

func (u *User) ID() domain.ID {
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

func (u *User) AccountID() domain.ID {
	return u.accountID
}

func NewUser(
	id domain.ID,
	firstName string,
	lastName string,
	email string,
	role Role,
	password string,
	accountID domain.ID,
) (*User, error) {
	firstName = strings.TrimSpace(firstName)
	lastName = strings.TrimSpace(lastName)
	email = strings.TrimSpace(email)
	password = strings.TrimSpace(password)

	if email == "" {
		return nil, errors.New("email cannot be invalid")
	}

	if password == "" {
		return nil, errors.New("password cannot be empty")
	}

	return &User{
		id:        id,
		firstName: firstName,
		lastName:  lastName,
		email:     email,
		password:  password,
		role:      role,
		accountID: accountID,
	}, nil
}
