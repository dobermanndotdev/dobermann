package monitor

import (
	"errors"

	"github.com/flowck/dobermann/backend/internal/domain"
)

type Subscriber struct {
	userID domain.ID
}

func NewSubscriber(userID domain.ID) (*Subscriber, error) {
	if userID.IsEmpty() {
		return nil, errors.New("userID cannot be invalid")
	}

	return &Subscriber{
		userID: userID,
	}, nil
}

func (s *Subscriber) UserID() domain.ID {
	return s.userID
}
