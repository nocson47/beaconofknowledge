package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nocson47/beaconofknowledge/internal/usecases"
)

func SetupRouter(app *fiber.App, userHandler *UserHandler, userSvc usecases.UserService, threadHandler *ThreadHandler, voteHandler *VoteHandler, replyHandler *ReplyHandler) {
	// apply CORS globally
	app.Use(Cors())
	// User routes
	users := app.Group("/users")
	users.Get("/", userHandler.GetAllUsers)
	users.Get("/:id", userHandler.GetUserByID)
	// use a separate path for username lookup to avoid conflicting with the :id route
	users.Get("/username/:username", userHandler.GetUserByUsername)
	users.Post("/", userHandler.CreateUser)
	users.Post("/login", userHandler.Login)
	users.Put(":id", RequireAuth(), userHandler.UpdateUser)
	// Delete user: admin-only (or owner; handler also enforces ownership)
	users.Delete(":id", RequireAuth(), AdminOnly(userSvc), userHandler.DeleteUser)

	// Thread routes
	threads := app.Group("/threads")
	threads.Get("/", threadHandler.GetAllThreads)                  // GET /threads
	threads.Post("/", RequireAuth(), threadHandler.CreateThread)   // POST /threads
	threads.Get("/:id", threadHandler.GetThreadByID)               // GET /threads/:id
	threads.Put("/:id", RequireAuth(), threadHandler.UpdateThread) // PUT /threads/:id
	// Thread deletion is admin-only
	threads.Delete("/:id", RequireAuth(), AdminOnly(userSvc), threadHandler.DeleteThread) // DELETE /threads/:id

	// Vote routes
	votes := app.Group("/votes")
	votes.Post("/", RequireAuth(), voteHandler.CreateVote) // POST /votes

	// Thread vote counts
	threads.Get("/:id/votes", voteHandler.GetThreadCounts) // GET /threads/:id/votes

	// Reply routes
	replies := app.Group("/replies")
	replies.Post("/", RequireAuth(), replyHandler.CreateReply)
	replies.Get("/thread/:thread_id", replyHandler.GetRepliesByThread)
	replies.Delete(":id", RequireAuth(), replyHandler.DeleteReply)
}
