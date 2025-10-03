package http

import (
	"context"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/nocson47/beaconofknowledge/adapters/jwt"
	"github.com/nocson47/beaconofknowledge/internal/entities"
	"github.com/nocson47/beaconofknowledge/internal/usecases"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ReportHandler struct {
	svc    usecases.ReportService
	logCol *mongo.Collection // optional mongo collection for logs
}

// NewReportHandler creates a handler. Pass nil for logCol if you don't have Mongo logging.
func NewReportHandler(svc usecases.ReportService, logCol *mongo.Collection) *ReportHandler {
	return &ReportHandler{svc: svc, logCol: logCol}
}

type createReportReq struct {
	Kind     string `json:"kind"`
	TargetID int    `json:"target_id"`
	Reason   string `json:"reason,omitempty"`
}

// CreateReport accepts POST /reports
// reporters may be anonymous; if Authorization header present, we attach reporter_id
func (h *ReportHandler) CreateReport(c *fiber.Ctx) error {
	var req createReportReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid payload"})
	}
	var reporterID *int
	auth := c.Get("Authorization")
	if auth != "" {
		token := auth
		token = strings.TrimSpace(strings.TrimPrefix(token, "Bearer"))
		token = strings.TrimSpace(token)
		if token != "" {
			if uid, err := jwt.ParseToken(token); err == nil && uid != 0 {
				reporterID = &uid
			}
		}
	}
	rep := &entities.Report{
		ReporterID: reporterID,
		Kind:       req.Kind,
		TargetID:   req.TargetID,
		Reason:     req.Reason,
		Status:     "open",
		CreatedAt:  time.Now(),
	}
	id, err := h.svc.CreateReport(context.Background(), rep)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	// write an audit/log entry to mongo if available
	if h.logCol != nil {
		_, _ = h.logCol.InsertOne(context.Background(), bson.M{
			"type":        "report_created",
			"report_id":   id,
			"kind":        rep.Kind,
			"target_id":   rep.TargetID,
			"reporter_id": rep.ReporterID,
			"reason":      rep.Reason,
			"created_at":  time.Now(),
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"id": id})
}

// GetReports returns reports (admin only)
func (h *ReportHandler) GetReports(c *fiber.Ctx) error {
	kind := c.Query("kind")
	var k *string
	if kind != "" {
		k = &kind
	}
	reps, err := h.svc.GetReports(context.Background(), k)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"reports": reps})
}

// UpdateReport allows admin to change status (resolve/dismiss)
type updateReportReq struct {
	Status string `json:"status"`
}

func (h *ReportHandler) UpdateReport(c *fiber.Ctx) error {
	id := c.Params("id")
	var req updateReportReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid payload"})
	}
	// derive admin user id
	auth := c.Get("Authorization")
	var adminID *int
	if auth != "" {
		token := strings.TrimSpace(strings.TrimPrefix(auth, "Bearer"))
		token = strings.TrimSpace(token)
		if token != "" {
			if uid, err := jwt.ParseToken(token); err == nil && uid != 0 {
				adminID = &uid
			}
		}
	}
	if adminID == nil {
		// should be guarded by AdminOnly middleware, but double-check
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}
	if err := h.svc.UpdateReportStatus(context.Background(), id, req.Status, adminID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	// log to mongo
	if h.logCol != nil {
		_, _ = h.logCol.InsertOne(context.Background(), bson.M{
			"type":        "report_update",
			"report_id":   id,
			"status":      req.Status,
			"resolved_by": adminID,
			"updated_at":  time.Now(),
		})
	}
	return c.JSON(fiber.Map{"ok": true})
}
