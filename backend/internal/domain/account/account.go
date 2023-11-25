package account

import "errors"

type Account struct {
	id    string
	users []*User
}

func (a *Account) ID() string {
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
