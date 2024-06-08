package account

import (
	"strings"
	"time"

	"errors"

	"github.com/dobermanndotdev/dobermann/internal/domain"
)

type User struct {
	id              domain.ID
	loginProviderID string
	email           Email
	role            Role
	accountID       domain.ID
	createdAt       time.Time
}

func NewUser(
	role Role,
	email Email,
	loginProviderID string,
	accountID domain.ID,
) (*User, error) {
	if role.IsEmpty() {
		return nil, errors.New("role cannot be invalid")
	}

	if strings.TrimSpace(loginProviderID) == "" {
		return nil, errors.New("loginProviderID cannot be empty")
	}

	return &User{
		id:              domain.NewID(),
		role:            role,
		email:           email,
		accountID:       accountID,
		createdAt:       time.Now().UTC(),
		loginProviderID: loginProviderID,
	}, nil
}

func (u *User) CreatedAt() time.Time {
	return u.createdAt
}

func (u *User) ID() domain.ID {
	return u.id
}

func (u *User) Role() Role {
	return u.role
}

func (u *User) AccountID() domain.ID {
	return u.accountID
}

func (u *User) Email() Email {
	return u.email
}

func (u *User) LoginProviderID() string {
	return u.loginProviderID
}
