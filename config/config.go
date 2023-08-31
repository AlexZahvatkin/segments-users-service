package config

import (
	"log"
	"os"
	"time"
)

type Config struct {
	Env        string `yaml:"env" env:"ENV" env-default:"local"`
	HTTPServer `yaml:"http_server"`
	Database   `yaml:"database"`
}

type HTTPServer struct {
	Port        string        `yaml:"port" env-default:"8080"`
	Host        string        `yaml:"address" env-default:"localhost"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

type Database struct {
	Host     string `yaml:"host" env-required:"true"`
	Port     string `yaml:"port" env-required:"true"`
	User     string `yaml:"user" env-required:"true"`
	Password string `yaml:"password"`
	Name     string `yaml:"name" env-required:"true"`
	SSLMode  string `yaml:"ssl_mode" env-default:"disable"`
}

func MustLoad() *Config {
	var cfg Config
	cfg.Env = os.Getenv("ENV_TYPE")
	cfg.HTTPServer.Port = os.Getenv("SERVER_PORT")
	cfg.HTTPServer.Host = os.Getenv("SERVER_HOST")
	timeout := os.Getenv("SERVER_TIMEOUT")
	timeoutDur, err := time.ParseDuration(timeout)
	if err != nil {
		log.Fatal("Wrong timeout format")
	}
	cfg.HTTPServer.Timeout = timeoutDur
	idleTimeout := os.Getenv("SERVER_IDLE_TIMEOUT")
	idleTimeoutDur, err := time.ParseDuration(idleTimeout)
	if err != nil {
		log.Fatal("Wrong idle timeout format")
	}
	cfg.HTTPServer.IdleTimeout = idleTimeoutDur

	cfg.Database.Host = os.Getenv("POSTGRES_HOST")
	cfg.Database.Port = os.Getenv("POSTGRES_PORT")
	cfg.Database.User = os.Getenv("POSTGRES_USER")
	cfg.Database.Password = os.Getenv("POSTGRES_PASSWORD")
	cfg.Database.Name = os.Getenv("POSTGRES_DB")
	cfg.Database.SSLMode = os.Getenv("POSTGRES_SSLMODE")

	return &cfg
}
