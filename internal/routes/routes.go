package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"

	"todo-api/internal/handler"
	"todo-api/internal/middleware"
	"todo-api/internal/repository"
	"todo-api/internal/service"
)

func RegisterRoutes(router *gin.Engine, db *gorm.DB) {
	// Repository
	userRepo := repository.NewUserRepository(db)
	todoRepo := repository.NewTodoRepository(db)

	// Service
	authService := service.NewAuthService(userRepo)
	todoService := service.NewTodoService(todoRepo)

	// Handler
	authHandler := handler.NewAuthHandler(authService)
	todoHandler := handler.NewTodoHandler(todoService)

	// API version group
	api := router.Group("/api/v1")

	// Authentication routes
	auth := api.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)

		auth.GET(
			"/profile",
			middleware.JWTMiddleware(),
			authHandler.GetProfile,
		)
	}

	// Protected Todo routes
	todos := api.Group("/todos")
	todos.Use(middleware.JWTMiddleware())
	{
		todos.POST("", todoHandler.Create)
		todos.GET("", todoHandler.GetAll)
		todos.GET("/:id", todoHandler.GetByID)
		todos.PUT("/:id", todoHandler.Update)
		todos.DELETE("/:id", todoHandler.Delete)
	}

	// Swagger documentation
	router.GET(
		"/swagger/*any",
		ginSwagger.WrapHandler(swaggerFiles.Handler),
	)

	// Health check
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Todo API Running",
		})
	})
}
