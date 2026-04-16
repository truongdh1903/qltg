package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"time-tracker/internal/config"
	"time-tracker/internal/database"
	"time-tracker/internal/handler"
	"time-tracker/internal/middleware"
	"time-tracker/internal/scheduler"
)

func main() {
	// Load .env file (chỉ ở local dev; production dùng env vars thật)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Load config
	cfg := config.Load()

	// Connect database
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Run migrations (auto-run khi start)
	if err := database.RunMigrations(cfg); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// Setup Gin
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Load HTML templates
	r.LoadHTMLGlob("web/templates/**/*.html")

	// Serve static files
	r.Static("/static", "./web/static")

	// Register middleware
	r.Use(middleware.RequestID())

	// Register routes
	handler.RegisterRoutes(r, db, cfg)

	// Start scheduler (Trash cleanup job)
	s := scheduler.New(db)
	s.Start()
	defer s.Stop()

	// Start server
	log.Printf("Server running on :%s (env: %s)", cfg.AppPort, cfg.AppEnv)
	if err := r.Run(":" + cfg.AppPort); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
