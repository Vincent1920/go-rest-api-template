package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"todo-api/config"
	"todo-api/internal/database"
	"todo-api/internal/routes"
)

// @title Todo API
// @version 1.0
// @description REST API Todo List menggunakan Gin + PostgreSQL + JWT
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatal(err)
	}

	if err := database.Migrate(db); err != nil {
		log.Fatal(err)
	}

	database.Migrate(db)

	router := gin.Default()

	routes.RegisterRoutes(router, db)

	router.Run(":" + cfg.AppPort)
}
