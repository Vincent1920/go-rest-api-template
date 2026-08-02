package main

import (
	"log"

	"github.com/gin-gonic/gin"

	_ "todo-api/docs"

	"todo-api/config"
	"todo-api/internal/database"
	"todo-api/internal/routes"
	"todo-api/internal/utils"
)

// @title Todo API
// @version 1.0
// @description REST API Todo List menggunakan Gin, PostgreSQL, dan JWT.
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	gin.SetMode(cfg.AppMode)
	utils.ConfigureJWT(cfg.JWTSecret, cfg.JWTExpire)

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatal(err)
	}

	if err := database.Migrate(db); err != nil {
		log.Fatalf("database migration: %v", err)
	}

	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())
	if err := router.SetTrustedProxies(nil); err != nil {
		log.Fatalf("configure trusted proxies: %v", err)
	}

	routes.RegisterRoutes(router, db)

	log.Printf("%s listening on http://localhost:%s", cfg.AppName, cfg.AppPort)
	if err := router.Run(":" + cfg.AppPort); err != nil {
		log.Fatalf("run server: %v", err)
	}
}
