package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"todo-api/config"
)

func Connect(cfg *config.Config) (*gorm.DB, error) {

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBPort,
		cfg.DBSSLMode,
	)

	// ==========================
	// DEBUG CONFIG
	// ==========================
	fmt.Println("========== DATABASE CONFIG ==========")
	fmt.Println("Host     :", cfg.DBHost)
	fmt.Println("Port     :", cfg.DBPort)
	fmt.Println("User     :", cfg.DBUser)
	fmt.Println("Password :", cfg.DBPassword)
	fmt.Println("Database :", cfg.DBName)
	fmt.Println("SSLMode  :", cfg.DBSSLMode)
	fmt.Println("=====================================")

	fmt.Println("DSN:")
	fmt.Println(dsn)

	// ==========================
	// CONNECT DATABASE
	// ==========================
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	fmt.Println("✅ Database connected successfully")

	return db, nil
}
