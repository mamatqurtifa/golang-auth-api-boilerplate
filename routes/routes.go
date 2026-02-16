package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mamatqurtifa/golang-auth-api-boilerplate/controllers"
	"github.com/mamatqurtifa/golang-auth-api-boilerplate/middleware"
)

// SetupRoutes configures all application routes
func SetupRoutes(router *gin.Engine) {
	authController := controllers.NewAuthController()

	// API v1 group
	v1 := router.Group("/api/v1")
	{
		// Health check
		v1.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status":  "ok",
				"message": "Server is running",
			})
		})

		// Auth routes (public)
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authController.Register)
			auth.POST("/login", authController.Login)
			auth.POST("/forgot-password", authController.ForgotPassword)
			auth.POST("/reset-password", authController.ResetPassword)
			auth.GET("/verify-email", authController.VerifyEmail)
		}

		// Protected routes (require authentication)
		protected := v1.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			// User routes
			user := protected.Group("/user")
			{
				user.GET("/profile", authController.GetProfile)
				user.PUT("/profile", authController.UpdateProfile)
				user.POST("/change-password", authController.ChangePassword)
				user.POST("/logout", authController.Logout)
			}
		}
	}
}
