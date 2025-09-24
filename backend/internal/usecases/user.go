package usecases

import (
	"context"
	"errors"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/nocson47/beaconofknowledge/internal/entities"
	"github.com/nocson47/beaconofknowledge/internal/repositories"
)

// UserService is the application port for user operations.
type UserService interface {
	CreateUser(ctx context.Context, u *entities.User) (int, error)
	GetAllUsers(ctx context.Context) ([]entities.User, error)
	GetUserByID(ctx context.Context, id int) (*entities.User, error)
	GetUserByUsername(ctx context.Context, username string) (*entities.User, error)
	UpdateUser(ctx context.Context, u *entities.User) error
	DeleteUser(ctx context.Context, id int) error
}

// unexported implementation to enforce interface usage
type userService struct {
	userRepo repositories.UserRepository
}

func NewUserUseCase(userRepo repositories.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) CreateUser(ctx context.Context, u *entities.User) (int, error) {
	if u == nil {
		return 0, errors.New("user is nil")
	}
	u.Username = strings.TrimSpace(u.Username)
	u.Email = strings.TrimSpace(u.Email)
	u.Password = strings.TrimSpace(u.Password)
	if u.Username == "" || u.Email == "" || u.Password == "" {
		return 0, errors.New("username, email and password are required")
	}
	// Hash password before storing
	hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}
	u.Password = string(hashed)
	// Set defaults
	if strings.TrimSpace(u.Role) == "" {
		u.Role = "user"
	}
	if u.CreatedAt.IsZero() {
		u.CreatedAt = time.Now()
	}
	return s.userRepo.CreateUser(ctx, u)
}

func (s *userService) GetAllUsers(ctx context.Context) ([]entities.User, error) {
	return s.userRepo.GetAllUsers(ctx)
}

func (s *userService) GetUserByID(ctx context.Context, id int) (*entities.User, error) {
	return s.userRepo.GetUserByID(ctx, id)
}

func (s *userService) GetUserByUsername(ctx context.Context, username string) (*entities.User, error) {
	return s.userRepo.GetUserByUsername(ctx, username)
}

func (s *userService) UpdateUser(ctx context.Context, u *entities.User) error {
	if u == nil {
		return errors.New("user is nil")
	}
	u.Username = strings.TrimSpace(u.Username)
	u.Email = strings.TrimSpace(u.Email)
	u.Bio = strings.TrimSpace(u.Bio)
	u.Social = strings.TrimSpace(u.Social)
	u.AvatarURL = strings.TrimSpace(u.AvatarURL)
	return s.userRepo.UpdateUser(ctx, u)
}

func (s *userService) DeleteUser(ctx context.Context, id int) error {
	return s.userRepo.DeleteUser(ctx, id)
}
