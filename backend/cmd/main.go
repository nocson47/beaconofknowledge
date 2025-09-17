package main

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	postgressql "github.com/nocson47/invoker_board/adapters/postgreSQL"
	"github.com/nocson47/invoker_board/config"
)

func main() {

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	log.Printf("Connecting to DB at %s:%d/%s", cfg.DBHost, cfg.DBPort, cfg.DBName)

	postgresConn, err := postgressql.ConnectPGX(&cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("Connected to PostgreSQL database successfully")
	defer postgresConn.Close(context.Background())

	// Initialize Fiber app
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	log.Fatal(app.Listen(":3000"))
}
