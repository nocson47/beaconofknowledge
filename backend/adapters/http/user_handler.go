package http

import (
	"log"
	"strconv"

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
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}
	return c.JSON(user)
}

func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	var user entities.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
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
	if curUser.Role != "admin" && curUser.ID != user.ID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "cannot update other users"})
	}

	if err := h.usecase.UpdateUser(c.UserContext(), &user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update user",
		})
	}
	return c.JSON(fiber.Map{
		"message": "User updated successfully",
	})
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
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}
	return c.JSON(user)
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
	return c.JSON(fiber.Map{"token": token, "expires_in": 900})
}
