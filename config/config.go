package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	AppName string
	AppPort string
	AppMode string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	JWTSecret string
	JWTExpire int
}

func LoadConfig() (*Config, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("./config")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	v.SetDefault("app.name", "Todo API")
	v.SetDefault("app.port", "8080")
	v.SetDefault("app.mode", "debug")
	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", "5433")
	v.SetDefault("database.user", "postgres")
	v.SetDefault("database.password", "postgres")
	v.SetDefault("database.name", "todo_db")
	v.SetDefault("database.sslmode", "disable")
	v.SetDefault("jwt.expire", 24)

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	cfg := &Config{
		AppName: v.GetString("app.name"),
		AppPort: v.GetString("app.port"),
		AppMode: v.GetString("app.mode"),

		DBHost:     v.GetString("database.host"),
		DBPort:     v.GetString("database.port"),
		DBUser:     v.GetString("database.user"),
		DBPassword: v.GetString("database.password"),
		DBName:     v.GetString("database.name"),
		DBSSLMode:  v.GetString("database.sslmode"),

		JWTSecret: v.GetString("jwt.secret"),
		JWTExpire: v.GetInt("jwt.expire"),
	}

	if strings.TrimSpace(cfg.JWTSecret) == "" {
		return nil, fmt.Errorf("jwt.secret must not be empty")
	}

	return cfg, nil
}
