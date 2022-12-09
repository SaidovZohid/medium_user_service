package config

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	GrpcPort      string
	Postgres      PostgresConfig
	Authorization string
	Redis         Redis

	NotificationServiceHost     string
	NotificationServiceGrpcPort string
}

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

type Redis struct {
	Addr string
}

func Load(path string) Config {
	godotenv.Load(path + "/.env")

	conf := viper.New()
	conf.AutomaticEnv()

	cfg := Config{
		GrpcPort: conf.GetString("USER_SERVICE_GRPC_PORT"),
		Postgres: PostgresConfig{
			Host:     conf.GetString("POSTGRES_HOST"),
			Port:     conf.GetString("POSTGRES_PORT"),
			User:     conf.GetString("POSTGRES_USER"),
			Password: conf.GetString("POSTGRES_PASSWORD"),
			Database: conf.GetString("POSTGRES_DATABASE"),
		},
		Authorization: conf.GetString("SECRET_KEY"),
		Redis: Redis{
			Addr: conf.GetString("REDIS_ADDR"),
		},
		NotificationServiceHost:     conf.GetString("NOTIFICATION_SERVICE_HOST"),
		NotificationServiceGrpcPort: conf.GetString("NOTIFICATION_SERVICE_USER_SERVICE_GRPC_PORT"),
	}
	return cfg
}
