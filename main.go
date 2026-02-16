package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/golang-auth-api-boilerplate/config"
	"github.com/yourusername/golang-auth-api-boilerplate/database"
	"github.com/yourusername/golang-auth-api-boilerplate/middleware"
	"github.com/yourusername/golang-auth-api-boilerplate/routes"
)

func main() {
	// Load configuration
	config.LoadConfig()

	// Set Gin mode
	gin.SetMode(config.AppConfig.GinMode)

	// Connect to database
	database.ConnectDatabase()

	// Initialize Gin router
	router := gin.Default()

	// Apply global middleware
	router.Use(middleware.Logger())
	router.Use(gin.Recovery())

	// CORS middleware
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Setup routes
	routes.SetupRoutes(router)

	// Start server
	port := ":" + config.AppConfig.Port
	log.Printf("Server is running on port %s", port)
	if err := router.Run(port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
