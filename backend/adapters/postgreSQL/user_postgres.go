// ...existing code...
package postgressql

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/nocson47/beaconofknowledge/internal/entities"
	"github.com/nocson47/beaconofknowledge/internal/repositories"
)

type UserPostgres struct {
	db *pgxpool.Pool // Connection pool
}

// compile-time check: ensure UserPostgres implements repositories.UserRepository
var _ repositories.UserRepository = (*UserPostgres)(nil)

// Constructor function for UserPostgres
func NewUserPostgres(db *pgxpool.Pool) repositories.UserRepository {
	return &UserPostgres{db: db}
}

// GetAllUsers retrieves all users from the database
func (u *UserPostgres) GetAllUsers(ctx context.Context) ([]entities.User, error) {
	query := `SELECT id, username, email FROM users`
	rows, err := u.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve users: %w", err)
	}
	defer rows.Close()

	var users []entities.User
	for rows.Next() {
		var user entities.User
		if err := rows.Scan(&user.ID, &user.Username, &user.Email); err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during row iteration: %w", err)
	}

	return users, nil
}

// CreateUser inserts a new user into the database
func (u *UserPostgres) CreateUser(ctx context.Context, user *entities.User) (int, error) {
	// DB column is `pass_hash` in your schema
	query := `INSERT INTO users (username, email, pass_hash, role, created_at) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	var id int
	err := u.db.QueryRow(ctx, query, user.Username, user.Email, user.Password, user.Role, user.CreatedAt).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to create user: %w", err)
	}
	return id, nil
}

// GetUserByID retrieves a user by their ID
func (u *UserPostgres) GetUserByID(ctx context.Context, id int) (*entities.User, error) {
	// select columns including profile fields
	query := `SELECT id, username, email, pass_hash, role, created_at, updated_at, COALESCE(bio, ''), COALESCE(social, ''), COALESCE(avatar_url, '') FROM users WHERE id = $1`
	row := u.db.QueryRow(ctx, query, id)

	var user entities.User
	if err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Role, &user.CreatedAt, &user.UpdatedAt, &user.Bio, &user.Social, &user.AvatarURL); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to retrieve user by ID: %w", err)
	}
	return &user, nil
}

// GetUserByUsername retrieves a user by their username
func (u *UserPostgres) GetUserByUsername(ctx context.Context, username string) (*entities.User, error) {
	query := `SELECT id, username, email, pass_hash, role, created_at, updated_at, COALESCE(bio, ''), COALESCE(social, ''), COALESCE(avatar_url, '') FROM users WHERE username = $1`
	row := u.db.QueryRow(ctx, query, username)

	var user entities.User
	if err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Role, &user.CreatedAt, &user.UpdatedAt, &user.Bio, &user.Social, &user.AvatarURL); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to retrieve user by username: %w", err)
	}
	return &user, nil
}

// UpdateUser updates an existing user's information
func (u *UserPostgres) UpdateUser(ctx context.Context, user *entities.User) error {
	// Update profile columns including bio, social, and avatar_url
	// also update updated_at so changes are recorded even without a DB trigger
	query := `UPDATE users SET username = $1, email = $2, pass_hash = $3, role = $4, bio = $5, social = $6, avatar_url = $7, updated_at = NOW() WHERE id = $8`
	_, err := u.db.Exec(ctx, query, user.Username, user.Email, user.Password, user.Role, user.Bio, user.Social, user.AvatarURL, user.ID)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

// DeleteUser deletes a user by their ID
func (u *UserPostgres) DeleteUser(ctx context.Context, id int) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := u.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

// GetUserByEmail retrieves a user by their email
func (u *UserPostgres) GetUserByEmail(ctx context.Context, email string) (*entities.User, error) {
	query := `SELECT id, username, email, pass_hash, role, created_at, updated_at, COALESCE(bio, ''), COALESCE(social, ''), COALESCE(avatar_url, '') FROM users WHERE email = $1`
	row := u.db.QueryRow(ctx, query, email)

	var user entities.User
	if err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Role, &user.CreatedAt, &user.UpdatedAt, &user.Bio, &user.Social, &user.AvatarURL); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to retrieve user by email: %w", err)
	}
	return &user, nil
}

// ```// filepath: /Users/forson47/golang_board/backend/adapters/postgreSQL/user_postgres.go
// // ...existing code...
// package postgressql

// import (
//     "context"
//     "fmt"

//     "github.com/jackc/pgx/v4/pgxpool"
//     "github.com/nocson47/beaconofknowledge/internal/entities"
//     "github.com/nocson47/beaconofknowledge/internal/repositories"
// )

// type UserPostgres struct {
//     db *pgxpool.Pool // Connection pool
// }

// // Constructor function for UserPostgres
// func NewUserPostgres(db *pgxpool.Pool) repositories.UserRepository {
//     return &UserPostgres{db: db}
// }

// // GetAllUsers retrieves all users from the database
// func (u *UserPostgres) GetAllUsers(ctx context.Context) ([]entities.User, error) {
//     query := `SELECT id, username, email FROM users`
//     rows, err := u.db.Query(ctx, query)
//     if err != nil {
//         return nil, fmt.Errorf("failed to retrieve users: %w", err)
//     }
//     defer rows.Close()

//     var users []entities.User
//     for rows.Next() {
//         var user entities.User
//         if err := rows.Scan(&user.ID, &user.Username, &user.Email); err != nil {
//             return nil, fmt.Errorf("failed to scan user: %w", err)
//         }
//         users = append(users, user)
//     }

//     if err := rows.Err(); err != nil {
//         return nil, fmt.Errorf("error during row iteration: %w", err)
//     }

//     return users, nil
// }

// // CreateUser inserts a new user into the database
// func (u *UserPostgres) CreateUser(ctx context.Context, user entities.User) error {
//     query := `INSERT INTO users (username, email, password, role, created_at) VALUES ($1, $2, $3, $4, $5)`
//     _, err := u.db.Exec(ctx, query, user.Username, user.Email, user.Password, user.Role, user.CreateAt)
//     if err != nil {
//         return fmt.Errorf("failed to create user: %w", err)
//     }
//     return nil
// }

// // GetUserByID retrieves a user by their ID
// func (u *UserPostgres) GetUserByID(ctx context.Context, id int) (*entities.User, error) {
//     query := `SELECT id, username, email, role, created_at, updated_at FROM users WHERE id = $1`
//     row := u.db.QueryRow(ctx, query, id)

//     var user entities.User
//     if err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Role, &user.CreateAt, &user.UpdateAt); err != nil {
//         return nil, fmt.Errorf("failed to retrieve user by ID: %w", err)
//     }
//     return &user, nil
// }

// // GetUserByUsername retrieves a user by their username
// func (u *UserPostgres) GetUserByUsername(ctx context.Context, username string) (*entities.User, error) {
//     query := `SELECT id, username, email, role, created_at, updated_at FROM users WHERE username = $1`
//     row := u.db.QueryRow(ctx, query, username)

//     var user entities.User
//     if err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Role, &user.CreateAt, &user.UpdateAt); err != nil {
//         return nil, fmt.Errorf("failed to retrieve user by username: %w", err)
//     }
//     return &user, nil
// }

// // UpdateUser updates an existing user's information
// func (u *UserPostgres) UpdateUser(ctx context.Context, user entities.User) error {
//     query := `UPDATE users SET username = $1, email = $2, password = $3, role = $4, updated_at = $5 WHERE id = $6`
//     _, err := u.db.Exec(ctx, query, user.Username, user.Email, user.Password, user.Role, user.UpdateAt, user.ID)
//     if err != nil {
//         return fmt.Errorf("failed to update user: %w", err)
//     }
//     return nil
// }

// // DeleteUser deletes a user by their ID
// func (u *UserPostgres) DeleteUser(ctx context.Context, id int) error {
//     query := `DELETE FROM users WHERE id = $1`
//     _, err := u.db.Exec(ctx, query, id)
//     if err != nil {
//         return fmt.Errorf("failed to delete user: %w", err)
//     }
//     return nil
// }
