package account

import (
	"strings"
	"time"

	"errors"

	"github.com/flowck/dobermann/backend/internal/common/hashing"
	"github.com/flowck/dobermann/backend/internal/domain"
)

type User struct {
	id                domain.ID
	firstName         string
	lastName          string
	email             Email
	password          Password
	role              Role
	verificationToken *domain.ID
	verifiedAt        *time.Time
	accountID         domain.ID
	createdAt         time.Time
}

func NewUser(
	id domain.ID,
	firstName,
	lastName string,
	email Email,
	role Role,
	password Password,
	accountID domain.ID,
	createdAt time.Time,
) (*User, error) {
	firstName = strings.TrimSpace(firstName)
	lastName = strings.TrimSpace(lastName)

	if id.IsEmpty() {
		return nil, errors.New("id cannot be invalid")
	}

	if email.IsEmpty() {
		return nil, errors.New("email cannot be invalid")
	}

	if password.IsEmpty() {
		return nil, errors.New("password cannot be invalid")
	}

	if role.IsEmpty() {
		return nil, errors.New("role cannot be invalid")
	}

	if time.Now().Before(createdAt) {
		return nil, errors.New("createdAt cannot be set in the future")
	}

	return &User{
		id:        id,
		firstName: firstName,
		lastName:  lastName,
		email:     email,
		password:  password,
		role:      role,
		accountID: accountID,
		createdAt: createdAt.UTC(),
	}, nil
}

func (u *User) Password() Password {
	return u.password
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

func (u *User) Email() Email {
	return u.email
}

func (u *User) Role() Role {
	return u.role
}

func (u *User) AccountID() domain.ID {
	return u.accountID
}

func (u *User) Authenticate(plainTextPassword string) error {
	err := hashing.Compare(plainTextPassword, u.password.hash)
	if err != nil {
		return ErrAuthenticationFailed
	}

	return nil
}
