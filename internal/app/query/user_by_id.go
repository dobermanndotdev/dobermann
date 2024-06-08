package query

import (
	"context"

	"github.com/dobermanndotdev/dobermann/internal/domain"
	"github.com/dobermanndotdev/dobermann/internal/domain/account"
)

type UserByID struct {
	ID domain.ID
}

type UserByIdHandler struct {
	userRepository account.UserRepository
}

func NewUserByIdHandler(userRepository account.UserRepository) UserByIdHandler {
	return UserByIdHandler{
		userRepository: userRepository,
	}
}

func (h UserByIdHandler) Execute(ctx context.Context, q UserByID) (*account.User, error) {
	return h.userRepository.FindByID(ctx, q.ID)
}
