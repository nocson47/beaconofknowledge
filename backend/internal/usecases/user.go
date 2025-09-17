package usecases

import (
	"github.com/nocson47/invoker_board/internal/entities"
	"github.com/nocson47/invoker_board/internal/repositories"
)

type UserUseCase struct {
	userRepo repositories.UserRepository
}

// NewUserUseCase creates a new instance of UserUseCase
func NewUserUseCase(userRepo repositories.UserRepository) UserUseCase {
	return UserUseCase{userRepo: userRepo}
}

// RegisterUser registers a new user
func (u *UserUseCase) RegisterUser(user *entities.User) error {
	// Business logic for registering a user can be added here
	return u.userRepo.CreateUser(user)
}

func (u *UserUseCase) Login(username, password string) (*entities.User, error) {
	return u.userRepo.UserLogin(username, password)
}
