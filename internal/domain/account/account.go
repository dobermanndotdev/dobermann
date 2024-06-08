package account

import (
	"errors"
	"strings"
	"time"

	"github.com/dobermanndotdev/dobermann/internal/domain"
)

type Account struct {
	id         domain.ID
	name       string
	verifiedAt *time.Time
	users      []*User
}

func (a *Account) VerifiedAt() *time.Time {
	return a.verifiedAt
}

func (a *Account) Name() string {
	return a.name
}

func (a *Account) ID() domain.ID {
	return a.id
}

func (a *Account) Users() []*User {
	return a.users
}

func (a *Account) FirstAccountOwner() (*User, error) {
	if len(a.users) == 0 {
		return nil, errors.New("no user found")
	}

	// TODO: Implement this properly
	return a.users[0], nil
}

func NewFirstTimeAccount(name string, email Email, password Password) (*Account, error) {
	accountID := domain.NewID()

	name = strings.TrimSpace(name)
	if name == "" {
		return nil, errors.New("account name cannot be empty")
	}

	user, err := NewUser(domain.NewID(), "", "", email, RoleOwner, password, accountID, time.Now())
	if err != nil {
		return nil, err
	}

	return &Account{
		id:         accountID,
		name:       name,
		verifiedAt: nil,
		users:      []*User{user},
	}, nil
}
