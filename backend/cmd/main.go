package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/nocson47/beaconofknowledge/adapters/http"
	postgressql "github.com/nocson47/beaconofknowledge/adapters/postgreSQL"
	redisadapters "github.com/nocson47/beaconofknowledge/adapters/redis"
	"github.com/nocson47/beaconofknowledge/config"
	"github.com/nocson47/beaconofknowledge/internal/usecases"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	log.Printf("Connecting to DB at %s:%d/%s", cfg.DBHost, cfg.DBPort, cfg.DBName)

	// Connect to PostgreSQL
	postgresConn, err := postgressql.ConnectPGX(&cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("Connected to PostgreSQL database successfully")
	defer postgresConn.Close()

	// Initialize repository, use case, and handler
	userRepo := postgressql.NewUserPostgres(postgresConn)
	userService := usecases.NewUserUseCase(userRepo)
	userHandler := http.NewUserHandler(userService)

	threadRepo := postgressql.NewThreadPostgres(postgresConn)
	threadService := usecases.NewThreadService(threadRepo) // returns usecases.ThreadService (interface)
	// connect to redis and wire cache to handlers
	redisClient := redisadapters.ConnectRedis(&cfg)
	if err := redisadapters.PingRedis(redisClient); err != nil {
		log.Printf("Warning: failed to connect to redis: %v", err)
	} else {
		log.Println("Connected to Redis")
	}
	threadHandler := http.NewThreadHandler(threadService, redisClient)

	// Votes
	voteRepo := postgressql.NewVotePostgres(postgresConn)
	voteService := usecases.NewVoteService(voteRepo)
	voteHandler := http.NewVoteHandler(voteService, redisClient)

	// Replies
	replyRepo := postgressql.NewReplyPostgres(postgresConn)
	replyService := usecases.NewReplyService(replyRepo)
	replyHandler := http.NewReplyHandler(replyService, userService, redisClient)

	// Initialize Fiber app
	app := fiber.New()
	// Apply CORS before rate limiter so preflight (OPTIONS) get CORS headers
	app.Use(http.Cors())
	// Apply rate limiter globally (you can scope it per-route as needed)
	app.Use(http.RateLimiter())

	// Admin-protected actions: wire AdminOnly middleware where needed
	// Note: router mounts unprotected routes; we'll add admin-protected handlers below
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	// Set up routes (router config will use auth middleware where needed)
	http.SetupRouter(app, userHandler, userService, threadHandler, threadService, voteHandler, replyHandler)

	// Start the server
	log.Fatal(app.Listen(":3000"))
}
