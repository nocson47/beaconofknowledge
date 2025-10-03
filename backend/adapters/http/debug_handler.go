package http

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/gofiber/fiber/v2"
)

// DebugHandler exposes simple debug endpoints such as avatar upload
type DebugHandler struct{}

func NewDebugHandler() *DebugHandler { return &DebugHandler{} }

// UploadAvatar accepts multipart/form-data with fields: avatar (file) and user_id
// Saves file to ./public/avatars and returns JSON { url: "/public/avatars/..." }
func (h *DebugHandler) UploadAvatar(c *fiber.Ctx) error {
	// Parse multipart form
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid multipart form"})
	}
	files := form.File["avatar"]
	if len(files) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "missing avatar file"})
	}
	userIDs := form.Value["user_id"]
	userID := "unknown"
	if len(userIDs) > 0 {
		userID = userIDs[0]
	}

	// ensure public/avatars dir exists
	outDir := filepath.Join("public", "avatars")
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to create avatars dir"})
	}

	// take the first file only
	fh := files[0]
	src, err := fh.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to open uploaded file"})
	}
	defer src.Close()

	ext := filepath.Ext(fh.Filename)
	name := fmt.Sprintf("%s-%d%s", userID, time.Now().UnixNano(), ext)
	outPath := filepath.Join(outDir, name)
	outFile, err := os.Create(outPath)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to create file"})
	}
	defer outFile.Close()

	if _, err := io.Copy(outFile, src); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to save file"})
	}

	// return URL path (the frontend can request from same host /avatars/..)
	url := "/avatars/" + name
	return c.JSON(fiber.Map{"url": url})
}
