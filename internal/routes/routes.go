package routes

import (
	"todo-api/internal/handler"
	"todo-api/internal/middleware"
	"todo-api/internal/repository"
	"todo-api/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(router *gin.Engine, db *gorm.DB) {

	// Repository
	userRepo := repository.NewUserRepository(db)

	// Service
	authService := service.NewAuthService(userRepo)

	// Handler
	authHandler := handler.NewAuthHandler(authService)

	// API Group
	api := router.Group("/api/v1")

	// Authentication Routes
	auth := api.Group("/auth")
	{
		// Public Routes
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)

		// Protected Routes
		auth.GET(
			"/profile",
			middleware.JWTMiddleware(),
			authHandler.GetProfile,
		)
	}

	// Health Check
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Todo API Running",
		})
	})
}
