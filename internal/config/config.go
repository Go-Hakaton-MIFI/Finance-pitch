package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

// DatabaseConfig содержит настройки для подключения к БД
type DatabaseConfig struct {
	Engine   string `env:"DB_ENGINE" env-default:"mysql"`
	Host     string `env:"DB_HOST" env-default:"localhost"`
	Port     int    `env:"DB_PORT" env-default:"3306"`
	User     string `env:"DB_USER" env-default:"go_clean_template_user"`
	Password string `env:"DB_PASSWORD" env-default:"password"`
	Name     string `env:"DB_NAME" env-default:"go_clean_template"`
}

func (db DatabaseConfig) DSN() string {
	if db.Engine == "mysql" {
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
			db.User, db.Password, db.Host, db.Port, db.Name)
	} else {
		return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
			db.User, db.Password, db.Host, db.Port, db.Name)

	}

}

type Server struct {
	Address string `env:"APP_ADDRESS" env-default:"0.0.0.0"`
	Port    string `env:"APP_PORT" env-default:"8089"`
}

type S3 struct {
	Url            string `env:"S3_URL" env-default:"http://minio:9000"`
	S3RootUser     string `env:"S3_ROOT_USER" env-default:"development_minio_key"`
	S3RootPassword string `env:"S3_ROOT_PASSWORD" env-default:"development_minio_secret"`
}

type Auth struct {
	PublicKey string `env:"JWT_PUBLIC"`
}

type Config struct {
	Database        DatabaseConfig
	Server          Server
	Auth            Auth
	S3              S3
	ImageBucketName string `env:"IMAGE_BUCKET_NAME" env-default:"images"`
	Debug           bool   `env:"DEBUG" env-default:"false"`
	Name            string `yaml:"name" env:"APP_NAME"`
	Version         string `yaml:"version" env:"APP_VERSION"`
}

func Load() (*Config, error) {
	// Приоритет конфига: env > yaml > defaul
	var c Config
	err := cleanenv.ReadConfig("config/config.yaml", &c)
	if err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	return &c, nil
}
