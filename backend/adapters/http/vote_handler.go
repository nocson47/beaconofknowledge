package http

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/nocson47/beaconofknowledge/adapters/jwt"
	"github.com/nocson47/beaconofknowledge/internal/entities"
	"github.com/nocson47/beaconofknowledge/internal/usecases"
)

type VoteHandler struct {
	svc   usecases.VoteService
	cache *redis.Client
}

func NewVoteHandler(svc usecases.VoteService, cache *redis.Client) *VoteHandler {
	return &VoteHandler{svc: svc, cache: cache}
}

type createVoteReq struct {
	ThreadID *int        `json:"thread_id,omitempty"`
	ReplyID  *int        `json:"reply_id,omitempty"`
	Value    interface{} `json:"value"` // accept string "up"/"down" or numeric 1/-1
}

func (h *VoteHandler) CreateVote(c *fiber.Ctx) error {
	var req createVoteReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid payload"})
	}
	// derive user from Authorization header (JWT)
	auth := c.Get("Authorization")
	if auth == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing authorization"})
	}
	token := strings.TrimSpace(strings.TrimPrefix(auth, "Bearer"))
	token = strings.TrimSpace(token)
	userID, err := jwt.ParseToken(token)
	if err != nil || userID == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token"})
	}

	// parse value (accept string or numeric)
	var valInt int
	switch v := req.Value.(type) {
	case string:
		// may be "up" or "down" or numeric as string
		vi, perr := entities.ParseVoteValueString(v)
		if perr != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": perr.Error()})
		}
		valInt = vi
	case float64:
		// JSON numbers decode as float64
		if int(v) == 1 || int(v) == -1 {
			valInt = int(v)
		} else {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid numeric vote value"})
		}
	case json.Number:
		si := string(v)
		vi, perr := entities.ParseVoteValueString(si)
		if perr != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": perr.Error()})
		}
		valInt = vi
	default:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid vote value type"})
	}

	vote := &entities.Vote{
		UserID:   userID,
		ThreadID: req.ThreadID,
		ReplyID:  req.ReplyID,
		Value:    valInt,
	}

	id, err := h.svc.CreateVote(c.UserContext(), vote)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// invalidate thread cache if vote affected a thread
	if h.cache != nil && vote.ThreadID != nil {
		ctx := context.Background()
		key := fmt.Sprintf("thread:%d", *vote.ThreadID)
		h.cache.Del(ctx, key)
	}

	// If vote is for a thread, return updated counts
	resp := fiber.Map{"id": id, "value": vote.ValueString()}
	if vote.ThreadID != nil {
		up, down, gerr := h.svc.GetVoteCountsForThread(c.UserContext(), *vote.ThreadID)
		if gerr != nil {
			// non-fatal for client; still return created id
			resp["warning"] = fmt.Sprintf("failed to fetch counts: %v", gerr)
		} else {
			resp["upvotes"] = up
			resp["downvotes"] = down
		}
	}

	return c.Status(fiber.StatusCreated).JSON(resp)
}

func (h *VoteHandler) GetThreadCounts(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}
	up, down, err := h.svc.GetVoteCountsForThread(c.UserContext(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"upvotes": up, "downvotes": down})
}
