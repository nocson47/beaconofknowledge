package http

import (
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/nocson47/beaconofknowledge/adapters/jwt"
	"github.com/nocson47/beaconofknowledge/internal/usecases"
)

// RequireAuth checks Authorization Bearer token and sets user id in locals
func RequireAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		auth := c.Get("Authorization")
		if auth == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing authorization"})
		}
		// Support headers like: "Bearer <token>" (case-insensitive) and trim spaces
		token := strings.TrimSpace(strings.TrimPrefix(auth, "Bearer "))
		if token == "" {
			// try lowercase bearer
			token = strings.TrimSpace(strings.TrimPrefix(auth, "bearer "))
		}
		uid, err := jwt.ParseToken(token)
		if err != nil || uid == 0 {
			// Log parse error to help debugging token issues (don't log the token value)
			if err != nil {
				log.Printf("RequireAuth: token parse error: %v", err)
			} else {
				log.Printf("RequireAuth: token parsed but returned uid=0")
			}
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token"})
		}
		c.Locals("user_id", uid)
		return c.Next()
	}
}

// RateLimiter returns a rate limiter middleware for JWT-authenticated routes
func RateLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		// Increase Max for local development and skip OPTIONS preflight requests
		Max:        300,
		Expiration: 1 * time.Minute,
		// Don't apply rate limiting to OPTIONS (CORS preflight) or health checks
		Next: func(c *fiber.Ctx) bool {
			if c.Method() == "OPTIONS" {
				return true
			}
			// allow healthcheck root path
			if c.Path() == "/" {
				return true
			}
			return false
		},
		KeyGenerator: func(c *fiber.Ctx) string {
			// prefer user id if set, else IP
			if uid := c.Locals("user_id"); uid != nil {
				if id, ok := uid.(int); ok {
					return "user:" + strconv.Itoa(id)
				}
			}
			return c.IP()
		},
	})
}

// RateLimiterAuth is intended to be used after RequireAuth() so c.Locals("user_id") is available.
// It uses per-user keys (user:<id>) and enforces a stricter per-user limit.
func RateLimiterAuth() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        120, // 120 requests per minute per user
		Expiration: 1 * time.Minute,
		Next: func(c *fiber.Ctx) bool {
			if c.Method() == "OPTIONS" {
				return true
			}
			if c.Path() == "/" {
				return true
			}
			return false
		},
		KeyGenerator: func(c *fiber.Ctx) string {
			if uid := c.Locals("user_id"); uid != nil {
				if id, ok := uid.(int); ok {
					return "user:" + strconv.Itoa(id)
				}
			}
			// fallback to IP if not available
			return c.IP()
		},
	})
}

// RateLimiterStrict is for sensitive unauthenticated endpoints like login/register.
// It uses IP-based keys but with a low limit to prevent brute-force/spam.
func RateLimiterStrict() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        10, // e.g. 10 requests per minute per IP
		Expiration: 1 * time.Minute,
		Next: func(c *fiber.Ctx) bool {
			return c.Method() == "OPTIONS"
		},
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
	})
}

// AdminOnly ensures the authenticated user has role "admin".
// It requires a UserService to lookup the user and check role.
func AdminOnly(userSvc usecases.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		uidVal := c.Locals("user_id")
		if uidVal == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing user id"})
		}
		uid, ok := uidVal.(int)
		if !ok || uid == 0 {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid user id"})
		}
		user, err := userSvc.GetUserByID(c.UserContext(), uid)
		if err != nil || user == nil {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "forbidden"})
		}
		if user.Role != "admin" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "admin required"})
		}
		return c.Next()
	}
}

// OwnerOrAdmin allows the request to proceed only if the authenticated user
// is the owner of the resource (thread) or has an admin role.
// It requires both a UserService (to lookup roles) and a ThreadService
// (to lookup the thread owner).
func OwnerOrAdmin(userSvc usecases.UserService, threadSvc usecases.ThreadService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		uidVal := c.Locals("user_id")
		if uidVal == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing user id"})
		}
		uid, ok := uidVal.(int)
		if !ok || uid == 0 {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid user id"})
		}

		// parse thread id from params
		idStr := c.Params("id")
		if idStr == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "missing id"})
		}
		// convert to int
		tid, err := strconv.Atoi(idStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
		}

		// lookup thread to find owner
		thread, err := threadSvc.GetThreadByID(c.UserContext(), tid)
		if err != nil || thread == nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "thread not found"})
		}

		// owner allowed
		if thread.UserID == uid {
			return c.Next()
		}

		// otherwise check admin role
		user, err := userSvc.GetUserByID(c.UserContext(), uid)
		if err != nil || user == nil {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "forbidden"})
		}
		if user.Role == "admin" {
			return c.Next()
		}
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "owner or admin required"})
	}
}

// OwnerOrAdminReply checks that the authenticated user is the owner of the reply or an admin.
func OwnerOrAdminReply(userSvc usecases.UserService, replySvc usecases.ReplyService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		uidVal := c.Locals("user_id")
		if uidVal == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing user id"})
		}
		uid, ok := uidVal.(int)
		if !ok || uid == 0 {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid user id"})
		}
		idStr := c.Params("id")
		if idStr == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "missing id"})
		}
		rid, err := strconv.Atoi(idStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
		}
		rep, err := replySvc.GetReplyByID(c.UserContext(), rid)
		if err != nil || rep == nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "reply not found"})
		}
		if rep.UserID == uid {
			return c.Next()
		}
		user, err := userSvc.GetUserByID(c.UserContext(), uid)
		if err != nil || user == nil {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "forbidden"})
		}
		if user.Role == "admin" {
			return c.Next()
		}
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "owner or admin required"})
	}
}

// RequestLogger logs start and end of each request with duration and status.
func RequestLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		reqID := strconv.FormatInt(start.UnixNano(), 10)
		log.Printf("→ START %s %s id=%s", c.Method(), c.OriginalURL(), reqID)
		// let the handler run
		err := c.Next()
		duration := time.Since(start)
		status := c.Response().StatusCode()
		if err != nil {
			log.Printf("← END %s %s id=%s status=%d duration=%s err=%v", c.Method(), c.OriginalURL(), reqID, status, duration, err)
			return err
		}
		log.Printf("← END %s %s id=%s status=%d duration=%s", c.Method(), c.OriginalURL(), reqID, status, duration)
		return nil
	}
}

var requestLock sync.Mutex

// RequestLock serializes requests (global lock) so you can observe completion deterministically.
// While locked, it sets header X-Request-Locked: true. Use only for debugging/development.
func RequestLock() fiber.Handler {
	return func(c *fiber.Ctx) error {
		requestLock.Lock()
		defer requestLock.Unlock()
		// indicate lock held
		c.Set("X-Request-Locked", "true")
		// run handler while lock is held
		return c.Next()
	}
}

// Cors returns a permissive CORS middleware suitable for development.
// Adjust the config for production to restrict allowed origins, methods and headers.
func Cors() fiber.Handler {
	// Use a specific origin for development to allow credentials safely.
	// If you need multiple origins, list them comma-separated or make this value configurable.
	return cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	})
}
