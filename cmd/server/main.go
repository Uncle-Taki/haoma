package main

import (
	"log"
	"strconv"

	"haoma/internal/adapters/http"
	"haoma/internal/config"
	"haoma/internal/infrastructure/persistence"
	"haoma/internal/infrastructure/web"

	_ "haoma/api/docs" // Swagger docs registration
)

// @title Haoma - Black-Box Carnival API
// @version 1.0
// @description Persian god meets cyber trials. A 7-node security quiz for ELECOMP 1404.
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	// Initialize persistence layer
	db, err := persistence.NewDatabase()
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// Initialize HTTP server
	router := web.NewRouter()

	// Initialize handlers
	http.RegisterRoutes(router, db)

	// Start the carnival
	port := ":" + strconv.Itoa(config.DEFAULT_PORT)
	log.Println("🎪 Haoma's carnival opens at " + port)
	if err := router.Run(port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
