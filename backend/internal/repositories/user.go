package repositories

import "github.com/nocson47/invoker_board/internal/entities"

// UserRepository defines the interface for user-related database operations
type UserRepository interface {
	CreateUser(user *entities.User) error
	GetUserByID(id int) (*entities.User, error)
	GetUserByUsername(username string) (*entities.User, error)
	GetAllUsers() ([]*entities.User, error) // use pointer to avoid copying and null values
	UpdateUser(user *entities.User) error
	DeleteUser(id int) error
	UserLogin(username, password string) (*entities.User, error)
}
