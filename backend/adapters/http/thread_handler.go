// ...existing code...
package http

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/nocson47/beaconofknowledge/internal/entities"
	"github.com/nocson47/beaconofknowledge/internal/usecases"
)

// ThreadHandler depends on the ThreadService port (usecase) and optional Redis client for caching.
type ThreadHandler struct {
	svc   usecases.ThreadService
	cache *redis.Client
}

func NewThreadHandler(svc usecases.ThreadService, cache *redis.Client) *ThreadHandler {
	return &ThreadHandler{svc: svc, cache: cache}
}

type createThreadReq struct {
	UserID int      `json:"user_id"`
	Title  string   `json:"title"`
	Body   string   `json:"body"`
	Tags   []string `json:"tags,omitempty"`
}

func (h *ThreadHandler) CreateThread(c *fiber.Ctx) error {
	var req createThreadReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid payload"})
	}
	// If authentication middleware provided a user id, prefer that over the body value.
	userID := req.UserID
	if uidVal := c.Locals("user_id"); uidVal != nil {
		if uid, ok := uidVal.(int); ok && uid != 0 {
			userID = uid
		}
	}
	thread := &entities.Thread{
		UserID: userID,
		Title:  req.Title,
		Body:   req.Body,
		Tags:   req.Tags,
	}
	id, err := h.svc.CreateThread(c.UserContext(), thread)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	// invalidate cache for this thread id if present
	if h.cache != nil {
		key := fmt.Sprintf("thread:%d", id)
		ctx := context.Background()
		h.cache.Del(ctx, key)
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"id": id})
}

func (h *ThreadHandler) GetThreadByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}
	// try cache-aside: look up Redis first
	if h.cache != nil {
		key := fmt.Sprintf("thread:%d", id)
		ctx := context.Background()
		if data, err := h.cache.Get(ctx, key).Result(); err == nil {
			var t entities.Thread
			if jerr := json.Unmarshal([]byte(data), &t); jerr == nil {
				return c.JSON(t)
			}
			// if unmarshal fails, fallthrough to DB fetch
		}
	}
	thread, err := h.svc.GetThreadByID(c.UserContext(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}
	// cache the result (best-effort)
	if h.cache != nil {
		key := fmt.Sprintf("thread:%d", id)
		ctx := context.Background()
		if b, jerr := json.Marshal(thread); jerr == nil {
			// set a short TTL because votes/replies may update frequently
			h.cache.Set(ctx, key, b, 60*time.Second)
		}
	}
	return c.JSON(thread)
}

func (h *ThreadHandler) GetAllThreads(c *fiber.Ctx) error {
	threads, err := h.svc.GetAllThreads(c.UserContext())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(threads)
}

type updateThreadReq struct {
	Title     string   `json:"title"`
	Body      string   `json:"body"`
	Tags      []string `json:"tags,omitempty"`
	IsLocked  bool     `json:"is_locked"`
	IsDeleted bool     `json:"is_deleted"`
}

func (h *ThreadHandler) UpdateThread(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}
	var req updateThreadReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid payload"})
	}
	thread := &entities.Thread{
		ID:        id,
		Title:     req.Title,
		Body:      req.Body,
		Tags:      req.Tags,
		IsLocked:  req.IsLocked,
		IsDeleted: req.IsDeleted,
	}
	if err := h.svc.UpdateThread(c.UserContext(), thread); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	// invalidate cache
	if h.cache != nil {
		key := fmt.Sprintf("thread:%d", id)
		ctx := context.Background()
		h.cache.Del(ctx, key)
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *ThreadHandler) DeleteThread(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}
	if err := h.svc.DeleteThread(c.UserContext(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if h.cache != nil {
		key := fmt.Sprintf("thread:%d", id)
		ctx := context.Background()
		h.cache.Del(ctx, key)
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// ...existing code...
