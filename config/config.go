package config

import (
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

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")

	// Environment Variable
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	cfg := &Config{
		AppName: viper.GetString("app.name"),
		AppPort: viper.GetString("app.port"),
		AppMode: viper.GetString("app.mode"),

		DBHost:     viper.GetString("database.host"),
		DBPort:     viper.GetString("database.port"),
		DBUser:     viper.GetString("database.user"),
		DBPassword: viper.GetString("database.password"),
		DBName:     viper.GetString("database.name"),
		DBSSLMode:  viper.GetString("database.sslmode"),

		JWTSecret: viper.GetString("jwt.secret"),
		JWTExpire: viper.GetInt("jwt.expire"),
	}

	return cfg, nil
}
