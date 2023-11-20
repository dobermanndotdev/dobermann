package account

import (
	"errors"
	"strings"

	"github.com/flowck/doberman/internal/common/ddd"
)

type Account struct {
	id    ddd.ID
	name  string
	users []*User
}

func NewAccountWithOwner(id ddd.ID, name string, owner *User) (*Account, error) {
	return NewAccount(id, name, []*User{owner})
}

func NewAccount(id ddd.ID, name string, users []*User) (*Account, error) {
	name = strings.TrimSpace(name)

	if id.IsEmpty() {
		return nil, errors.New("id cannot be invalid")
	}

	if name == "" {
		return nil, errors.New("name cannot be empty")
	}

	if len(users) == 0 {
		return nil, errors.New("users cannot be empty")
	}

	hasOwner := false
	for _, user := range users {
		if user.role.value == RoleOwner.value {
			hasOwner = true
			break
		}
	}

	if !hasOwner {
		return nil, errors.New("an account must have at least one user with the Owner role")
	}

	return &Account{
		id:    id,
		name:  name,
		users: users,
	}, nil
}

func (t *Account) ID() ddd.ID {
	return t.id
}

func (t *Account) Name() string {
	return t.name
}

func (t *Account) Users() []*User {
	return t.users
}
