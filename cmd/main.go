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

	// Connect to the database
	dbpool, err := database.NewConnection(cfg.DBURL)
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}
	defer dbpool.Close()
	
	fmt.Println("Database connection established successfully")

	// Create repository and handler
	repo := repository.New(dbpool)
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

