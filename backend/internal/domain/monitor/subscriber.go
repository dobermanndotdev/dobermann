package monitor

import "github.com/flowck/dobermann/backend/internal/domain"

type Subscriber struct {
	id     domain.ID
	userID domain.ID
}
