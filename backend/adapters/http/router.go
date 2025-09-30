package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nocson47/beaconofknowledge/internal/usecases"
)

func SetupRouter(app *fiber.App, userHandler *UserHandler, userSvc usecases.UserService, threadHandler *ThreadHandler, threadSvc usecases.ThreadService, voteHandler *VoteHandler, replyHandler *ReplyHandler, reportHandler *ReportHandler, authHandler *AuthHandler) {
	// Debug endpoints (file upload for avatar)
	dbg := NewDebugHandler()
	app.Post("/debug/avatar", dbg.UploadAvatar)

	// Serve avatars directly at /avatars from the public/avatars folder
	app.Static("/avatars", "public/avatars")
	// Also keep serving any other public assets under /public if needed
	app.Static("/public", "public")

	// User routes
	users := app.Group("/users")
	users.Get("/", userHandler.GetAllUsers)
	// current user info
	users.Get("/me", RequireAuth(), userHandler.GetMe)
	users.Get(":id", userHandler.GetUserByID)

	// use a separate path for username lookup to avoid conflicting with the :id route
	users.Get("/username/:username", userHandler.GetUserByUsername)
	users.Post("/", userHandler.CreateUser)
	// apply stricter limiter to login to reduce brute force attempts
	users.Post("/login", RateLimiterStrict(), userHandler.Login)
	// protect update with auth + per-user rate limiter
	users.Put(":id", RequireAuth(), RateLimiterAuth(), userHandler.UpdateUser)

	// Delete user: admin-only (or owner; handler also enforces ownership)
	users.Delete(":id", RequireAuth(), AdminOnly(userSvc), userHandler.DeleteUser)

	// Thread routes
	threads := app.Group("/threads")
	threads.Get("/", threadHandler.GetAllThreads)                                   // GET /threads
	threads.Post("/", RequireAuth(), RateLimiterAuth(), threadHandler.CreateThread) // POST /threads
	threads.Get("/:id", threadHandler.GetThreadByID)                                // GET /threads/:id
	// Only owner or admin may update/delete a thread
	threads.Put("/:id", RequireAuth(), RateLimiterAuth(), OwnerOrAdmin(userSvc, threadSvc), threadHandler.UpdateThread)    // PUT /threads/:id
	threads.Delete("/:id", RequireAuth(), RateLimiterAuth(), OwnerOrAdmin(userSvc, threadSvc), threadHandler.DeleteThread) // DELETE /threads/:id

	// Vote routes
	votes := app.Group("/votes")
	votes.Post("/", RequireAuth(), RateLimiterAuth(), voteHandler.CreateVote) // POST /votes

	// Thread vote counts
	threads.Get("/:id/votes", voteHandler.GetThreadCounts) // GET /threads/:id/votes

	// Reply routes
	replies := app.Group("/replies")
	replies.Post("/", RequireAuth(), RateLimiterAuth(), replyHandler.CreateReply)
	replies.Get("/thread/:thread_id", replyHandler.GetRepliesByThread)
	// Allow owner or admin to update or delete a reply
	replies.Put(":id", RequireAuth(), RateLimiterAuth(), OwnerOrAdminReply(userSvc, replyHandler.svc), replyHandler.UpdateReply)
	replies.Delete(":id", RequireAuth(), RateLimiterAuth(), OwnerOrAdminReply(userSvc, replyHandler.svc), replyHandler.DeleteReply)

	// Reports
	app.Post("/reports", reportHandler.CreateReport)
	app.Get("/reports", RequireAuth(), AdminOnly(userSvc), reportHandler.GetReports)
	app.Put("/reports/:id", RequireAuth(), AdminOnly(userSvc), reportHandler.UpdateReport)

	// Auth endpoints (password reset)
	// Note: PasswordResetUsecase and its handler must be constructed/wired in cmd/main.go and passed in when SetupRouter is called.
	// For now, register paths if handlers are present in globals (constructed elsewhere)
	if authHandler != nil {
		app.Post("/auth/forgot", authHandler.ForgotPassword)
		app.Post("/auth/reset", authHandler.ResetPassword)
	}
}
