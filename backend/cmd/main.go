package main

import (
	"context"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/nocson47/beaconofknowledge/adapters/email"
	"github.com/nocson47/beaconofknowledge/adapters/http"
	mongoadapters "github.com/nocson47/beaconofknowledge/adapters/mongo"
	postgressql "github.com/nocson47/beaconofknowledge/adapters/postgreSQL"
	redisadapters "github.com/nocson47/beaconofknowledge/adapters/redis"
	"github.com/nocson47/beaconofknowledge/config"
	"github.com/nocson47/beaconofknowledge/internal/repositories"
	"github.com/nocson47/beaconofknowledge/internal/usecases"
	"go.mongodb.org/mongo-driver/mongo"
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

	// Reports: prefer Mongo if available, otherwise Postgres
	var reportRepoUse repositories.ReportRepository
	var reportService usecases.ReportService
	var reportHandler *http.ReportHandler
	var logCol *mongo.Collection
	mongoClient, mErr := mongoadapters.ConnectMongo(&cfg)
	if mErr == nil {
		log.Println("Connected to MongoDB, using Mongo report repository")
		reportRepoUse = mongoadapters.NewMongoReportRepo(mongoClient, cfg.MongoDBName)
		// prepare a logs collection for optional audit writes
		logCol = mongoClient.Database(cfg.MongoDBName).Collection("logs")
		// ensure indexes and collections
		if err := mongoadapters.EnsureIndexes(context.Background(), mongoClient, cfg.MongoDBName); err != nil {
			log.Printf("Warning: failed to ensure mongo indexes: %v", err)
		}
	} else {
		log.Printf("Mongo not available: %v — using Postgres reports", mErr)
		reportRepoUse = postgressql.NewReportPostgres(postgresConn)
		logCol = nil
	}
	reportService = usecases.NewReportService(reportRepoUse)
	reportHandler = http.NewReportHandler(reportService, logCol)

	// Password reset: wire Postgres password reset repo and usecase. Use SMTP if configured, otherwise default to console sender (dev)
	prRepo := postgressql.NewPasswordResetPostgres(postgresConn)
	var emailSender usecases.EmailSender
	if cfg.SMTPHost != "" && cfg.SMTPPort != 0 {
		log.Printf("Using SMTP email sender %s:%d", cfg.SMTPHost, cfg.SMTPPort)
		smtpSender := email.NewSMTPEmailSender(cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUser, cfg.SMTPPassword, cfg.SMTPFrom)
		emailSender = smtpSender
	} else {
		log.Printf("Using console email sender (development)")
		emailSender = email.NewConsoleEmailSender()
	}
	prUsecase := usecases.NewPasswordResetUsecase(userRepo, prRepo, emailSender, time.Hour*24)
	authHandler := http.NewAuthHandler(prUsecase)

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
	http.SetupRouter(app, userHandler, userService, threadHandler, threadService, voteHandler, replyHandler, reportHandler, authHandler)

	// Start the server
	log.Fatal(app.Listen(":3000"))
}
