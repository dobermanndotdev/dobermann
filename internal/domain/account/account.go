package account

import (
	"time"

	"github.com/dobermanndotdev/dobermann/internal/domain"
)

type Account struct {
	id        domain.ID
	createdAt time.Time
}

func NewAccount() *Account {
	return &Account{
		id:        domain.NewID(),
		createdAt: time.Now(),
	}
}

func (a *Account) CreatedAt() time.Time {
	return a.createdAt
}

func (a *Account) ID() domain.ID {
	return a.id
}
