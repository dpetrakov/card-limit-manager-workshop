package main

import (
	"log"
	// "net/http" // Removed as it's not directly used

	"services/clm/internal/handlers"

	"github.com/gin-gonic/gin"
	// "services/clm/internal/storage" // Not needed if store is created in handler
)

func main() {
	router := gin.Default()

	// In a real application, you would initialize and inject the store.
	// store := storage.NewDBStore()
	// For now, the handler creates its own store instance.

	// Setup API v1 routes group
	v1 := router.Group("/api/v1")
	{
		v1.POST("/requests", handlers.CreateLimitRequestHandler)
	}

	// Start server
	port := "8080"
	log.Printf("Starting CLM service on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
