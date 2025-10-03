package http

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/nocson47/beaconofknowledge/internal/usecases"
)

type AuthHandler struct {
	prUsecase *usecases.PasswordResetUsecase
}

func NewAuthHandler(pr *usecases.PasswordResetUsecase) *AuthHandler {
	return &AuthHandler{prUsecase: pr}
}

type forgotReq struct {
	Email string `json:"email"`
	Base  string `json:"base_url"`
}

func (h *AuthHandler) ForgotPassword(c *fiber.Ctx) error {
	var req forgotReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}
	// base URL: prefer header or env; for now use origin header or localhost
	base := req.Base
	if base == "" {
		base = c.Hostname()
		if base == "" {
			base = "http://localhost:3000"
		}
	}
	if err := h.prUsecase.RequestPasswordReset(c.UserContext(), req.Email, base); err != nil {
		log.Printf("ForgotPassword: %v", err)
		// do not reveal internal errors to caller
	}
	// always return 200 to avoid revealing whether email exists
	return c.JSON(fiber.Map{"message": "If that email exists, a reset link was sent"})
}

type resetReq struct {
	Token    string `json:"token"`
	Password string `json:"password"`
}

func (h *AuthHandler) ResetPassword(c *fiber.Ctx) error {
	var req resetReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}
	if req.Token == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "missing token or password"})
	}
	if err := h.prUsecase.ResetPassword(c.UserContext(), req.Token, req.Password); err != nil {
		log.Printf("ResetPassword: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid or expired token"})
	}
	return c.JSON(fiber.Map{"message": "Password reset successfully"})
}
