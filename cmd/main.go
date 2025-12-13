package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/student/my-kpfu-db-app/internal/config"
	"github.com/student/my-kpfu-db-app/internal/database"
	"github.com/student/my-kpfu-db-app/internal/handler"
	"github.com/student/my-kpfu-db-app/internal/repository"
)

func main() {
	// Load configuration
	cfg := config.Load()
	fmt.Printf("Connecting to database: %s\n", cfg.DBURL)

	// Connect to the database (pgx)
	dbpool, err := database.NewConnection(cfg.DBURL)
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}
	defer dbpool.Close()
	
	fmt.Println("Database connection (pgx) established successfully")

	// Connect to the database (GORM) for ORM operations
	gormDB, err := database.NewGormConnection(cfg.DBURL)
	if err != nil {
		log.Fatalf("Could not connect with GORM: %v", err)
	}
	
	fmt.Println("Database connection (GORM) established successfully")

	// Create repository and handler
	repo := repository.New(dbpool, gormDB)
	h := handler.New(repo)

	// Set up router
	r := gin.Default()
	
	// Load HTML templates
	r.LoadHTMLGlob("web/templates/*.html")
	
	// Register routes
	h.RegisterRoutes(r)

	// Start server
	fmt.Println("Starting server on :8080")
	fmt.Println("Open http://localhost:8080 in your browser")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Could not run server: %v", err)
	}
}

