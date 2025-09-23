package http

import (
	"context"
	"fmt"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/nocson47/beaconofknowledge/internal/entities"
	"github.com/nocson47/beaconofknowledge/internal/usecases"
)

type ReplyHandler struct {
	svc     usecases.ReplyService
	userSvc usecases.UserService
	cache   *redis.Client
}

func NewReplyHandler(svc usecases.ReplyService, userSvc usecases.UserService, cache *redis.Client) *ReplyHandler {
	return &ReplyHandler{svc: svc, userSvc: userSvc, cache: cache}
}

type createReplyReq struct {
	ThreadID int    `json:"thread_id"`
	ParentID *int   `json:"parent_id,omitempty"`
	Body     string `json:"body"`
}

func (h *ReplyHandler) CreateReply(c *fiber.Ctx) error {
	var req createReplyReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid payload"})
	}
	// require auth middleware to have set user id in locals
	uidVal := c.Locals("user_id")
	if uidVal == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing authorization"})
	}
	uid, ok := uidVal.(int)
	if !ok || uid == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid user id"})
	}

	rep := &entities.Reply{
		ThreadID: req.ThreadID,
		UserID:   uid,
		ParentID: req.ParentID,
		Body:     req.Body,
	}
	id, err := h.svc.CreateReply(c.UserContext(), rep)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	// invalidate thread cache
	if h.cache != nil {
		ctx := context.Background()
		key := fmt.Sprintf("thread:%d", req.ThreadID)
		h.cache.Del(ctx, key)
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"id": id})
}

func (h *ReplyHandler) GetRepliesByThread(c *fiber.Ctx) error {
	idStr := c.Params("thread_id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid thread id"})
	}
	reps, err := h.svc.GetRepliesByThread(c.UserContext(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(reps)
}

func (h *ReplyHandler) DeleteReply(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}
	// require user id from locals
	uidVal := c.Locals("user_id")
	if uidVal == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing authorization"})
	}
	uid, ok := uidVal.(int)
	if !ok || uid == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid user id"})
	}
	// determine if user is admin
	usr, uerr := h.userSvc.GetUserByID(c.UserContext(), uid)
	isAdmin := false
	if uerr == nil && usr != nil && usr.Role == "admin" {
		isAdmin = true
	}
	if err := h.svc.DeleteReply(c.UserContext(), id, uid, isAdmin); err != nil {
		// Forbidden vs other errors
		if err.Error() == "forbidden: cannot delete others' replies" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	// Note: we don't know the thread id here easily; a simple approach is to evict a generic key
	if h.cache != nil {
		// If you want precise invalidation, ReplyService.DeleteReply could return threadID.
		// For now, evict a small pattern could be used by admins, but redis doesn't support DEL by pattern efficiently without scan.
	}
	return c.SendStatus(fiber.StatusNoContent)
}
