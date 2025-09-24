package http

import (
	"log"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/gofiber/fiber/v2"
	"github.com/nocson47/beaconofknowledge/adapters/jwt"
	"github.com/nocson47/beaconofknowledge/internal/entities"
	"github.com/nocson47/beaconofknowledge/internal/usecases"
)

type UserHandler struct {
	usecase usecases.UserService
}

func NewUserHandler(usecase usecases.UserService) *UserHandler {
	return &UserHandler{usecase: usecase}
}

func (h *UserHandler) GetAllUsers(c *fiber.Ctx) error {
	log.Println("Handler: Fetching all users")
	users, err := h.usecase.GetAllUsers(c.UserContext())
	if err != nil {
		log.Printf("Handler Error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve users",
		})
	}
	log.Println("Handler: Successfully fetched users")
	return c.JSON(users)
}

func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var user entities.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	id, err := h.usecase.CreateUser(c.UserContext(), &user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user",
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"id": id,
	})
}

func (h *UserHandler) GetUserByID(c *fiber.Ctx) error {
	id := c.Params("id")

	// Convert id from string to int
	intID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	user, err := h.usecase.GetUserByID(c.UserContext(), intID)
	if err != nil || user == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	// build public view - hide sensitive fields like password and email unless owner/admin
	type publicUser struct {
		ID        int       `json:"id"`
		Username  string    `json:"username"`
		AvatarURL string    `json:"avatar_url,omitempty"`
		Bio       string    `json:"bio,omitempty"`
		Social    string    `json:"social,omitempty"`
		CreatedAt time.Time `json:"created_at,omitempty"`
		UpdatedAt time.Time `json:"updated_at,omitempty"`
		Email     string    `json:"email,omitempty"`
	}

	pu := publicUser{
		ID:        user.ID,
		Username:  user.Username,
		AvatarURL: user.AvatarURL,
		Bio:       user.Bio,
		Social:    user.Social,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	// include email only if requester is owner or admin
	uidVal := c.Locals("user_id")
	if uidVal != nil {
		if uid, ok := uidVal.(int); ok {
			if uid == user.ID {
				pu.Email = user.Email
			} else {
				// if not owner, check if requester is admin
				curUser, _ := h.usecase.GetUserByID(c.UserContext(), uid)
				if curUser != nil && curUser.Role == "admin" {
					pu.Email = user.Email
				}
			}
		}
	}

	return c.JSON(pu)
}

func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	var user entities.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// get id from URL param and use that as target
	idParam := c.Params("id")
	intID, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user id in path"})
	}

	// require auth: user id in locals
	uidVal := c.Locals("user_id")
	if uidVal == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing authorization"})
	}
	uid, ok := uidVal.(int)
	if !ok || uid == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid user id"})
	}
	// only allow owner or admin
	curUser, err := h.usecase.GetUserByID(c.UserContext(), uid)
	if err != nil || curUser == nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "forbidden"})
	}
	if curUser.Role != "admin" && curUser.ID != intID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "cannot update other users"})
	}

	// fetch existing user to merge fields (avoid overwriting pass_hash/role unintentionally)
	existing, gerr := h.usecase.GetUserByID(c.UserContext(), intID)
	if gerr != nil || existing == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
	}

	// ensure ID is set to URL id
	user.ID = intID
	// preserve fields not supplied in request
	if strings.TrimSpace(user.Username) == "" {
		user.Username = existing.Username
	}
	if strings.TrimSpace(user.Email) == "" {
		user.Email = existing.Email
	}
	// preserve password hash and role
	user.Password = existing.Password
	user.Role = existing.Role
	// preserve avatar if not provided
	if strings.TrimSpace(user.AvatarURL) == "" {
		user.AvatarURL = existing.AvatarURL
	}

	if err := h.usecase.UpdateUser(c.UserContext(), &user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update user",
		})
	}

	// fetch updated user and return a sanitized view
	refreshed, rerr := h.usecase.GetUserByID(c.UserContext(), intID)
	if rerr != nil || refreshed == nil {
		log.Printf("UpdateUser: updated but failed to fetch refreshed user: %v", rerr)
		return c.JSON(fiber.Map{"message": "User updated successfully"})
	}
	log.Printf("User %d updated_at: %v", refreshed.ID, refreshed.UpdatedAt)

	type publicUser struct {
		ID        int       `json:"id"`
		Username  string    `json:"username"`
		AvatarURL string    `json:"avatar_url,omitempty"`
		Bio       string    `json:"bio,omitempty"`
		Social    string    `json:"social,omitempty"`
		CreatedAt time.Time `json:"created_at,omitempty"`
		UpdatedAt time.Time `json:"updated_at,omitempty"`
		Email     string    `json:"email,omitempty"`
	}

	pu := publicUser{
		ID:        refreshed.ID,
		Username:  refreshed.Username,
		AvatarURL: refreshed.AvatarURL,
		Bio:       refreshed.Bio,
		Social:    refreshed.Social,
		CreatedAt: refreshed.CreatedAt,
		UpdatedAt: refreshed.UpdatedAt,
		Email:     refreshed.Email, // owner or admin only - UpdateUser enforces that
	}
	return c.JSON(pu)
}

func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	intID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}
	// require auth
	uidVal := c.Locals("user_id")
	if uidVal == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing authorization"})
	}
	uid, ok := uidVal.(int)
	if !ok || uid == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid user id"})
	}
	curUser, err := h.usecase.GetUserByID(c.UserContext(), uid)
	if err != nil || curUser == nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "forbidden"})
	}
	if curUser.Role != "admin" && curUser.ID != intID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "cannot delete other users"})
	}

	if err := h.usecase.DeleteUser(c.UserContext(), intID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete user",
		})
	}
	return c.JSON(fiber.Map{
		"message": "User deleted successfully",
	})
}

func (h *UserHandler) GetUserByUsername(c *fiber.Ctx) error {
	username := c.Params("username")
	user, err := h.usecase.GetUserByUsername(c.UserContext(), username)
	if err != nil || user == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	type publicUser struct {
		ID        int       `json:"id"`
		Username  string    `json:"username"`
		AvatarURL string    `json:"avatar_url,omitempty"`
		Bio       string    `json:"bio,omitempty"`
		Social    string    `json:"social,omitempty"`
		CreatedAt time.Time `json:"created_at,omitempty"`
		UpdatedAt time.Time `json:"updated_at,omitempty"`
		Email     string    `json:"email,omitempty"`
	}

	pu := publicUser{
		ID:        user.ID,
		Username:  user.Username,
		AvatarURL: user.AvatarURL,
		Bio:       user.Bio,
		Social:    user.Social,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	uidVal := c.Locals("user_id")
	if uidVal != nil {
		if uid, ok := uidVal.(int); ok {
			if uid == user.ID {
				pu.Email = user.Email
			} else {
				curUser, _ := h.usecase.GetUserByID(c.UserContext(), uid)
				if curUser != nil && curUser.Role == "admin" {
					pu.Email = user.Email
				}
			}
		}
	}

	return c.JSON(pu)
}

type loginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Login validates credentials and returns a JWT token
func (h *UserHandler) Login(c *fiber.Ctx) error {
	var req loginReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid payload"})
	}
	user, err := h.usecase.GetUserByUsername(c.UserContext(), req.Username)
	if err != nil || user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid credentials"})
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid credentials"})
	}
	token, terr := jwt.GenerateToken(user.ID)
	if terr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to generate token"})
	}
	// expires_in is in seconds (15 minutes)
	// return a sanitized user object along with token so frontend can store role and show admin UI
	type loginUser struct {
		ID        int       `json:"id"`
		Username  string    `json:"username"`
		AvatarURL string    `json:"avatar_url,omitempty"`
		Bio       string    `json:"bio,omitempty"`
		Social    string    `json:"social,omitempty"`
		CreatedAt time.Time `json:"created_at,omitempty"`
		UpdatedAt time.Time `json:"updated_at,omitempty"`
		Email     string    `json:"email,omitempty"`
		Role      string    `json:"role,omitempty"`
	}
	lu := loginUser{
		ID:        user.ID,
		Username:  user.Username,
		AvatarURL: user.AvatarURL,
		Bio:       user.Bio,
		Social:    user.Social,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
		Role:      user.Role,
	}
	return c.JSON(fiber.Map{"token": token, "expires_in": 900, "user": lu})
}

// GetMe returns the currently authenticated user's full info (requires auth)
func (h *UserHandler) GetMe(c *fiber.Ctx) error {
	uidVal := c.Locals("user_id")
	if uidVal == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing authorization"})
	}
	uid, ok := uidVal.(int)
	if !ok || uid == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid user id"})
	}
	user, err := h.usecase.GetUserByID(c.UserContext(), uid)
	if err != nil || user == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
	}

	type meUser struct {
		ID        int       `json:"id"`
		Username  string    `json:"username"`
		AvatarURL string    `json:"avatar_url,omitempty"`
		Bio       string    `json:"bio,omitempty"`
		Social    string    `json:"social,omitempty"`
		CreatedAt time.Time `json:"created_at,omitempty"`
		UpdatedAt time.Time `json:"updated_at,omitempty"`
		Email     string    `json:"email,omitempty"`
		Role      string    `json:"role,omitempty"`
	}
	mu := meUser{
		ID:        user.ID,
		Username:  user.Username,
		AvatarURL: user.AvatarURL,
		Bio:       user.Bio,
		Social:    user.Social,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
		Role:      user.Role,
	}
	return c.JSON(mu)
}
