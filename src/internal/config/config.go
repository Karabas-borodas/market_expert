package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

// FIX: захардкожен путь конфига чтол бы не вызвать каждый раз в командной строке при запуске:
const configPath = "config/config.yaml"

type Config struct {
	Env        string `yaml:"env" env-default:"local"`
	HTTPserver string `yaml:"http_server"`
	// Posgress  strign`yaml:"`
}

type HTTPserver struct {
	Address       string        `yaml:"address" env-default:"localhost:20005"`
	Timeout       time.Duration `yaml:"timeout" env-default:"10s"`
	Iddle_timeout time.Duration `yaml:"iddle_timeout" env-default:"80s"`
}

type PostgresConfig struct {
	User     string `yaml:"POSTGRES_USER" env-default:"Mihail"`
	Password string `yaml:"POSTGRES_PASSWORD" env-default:"Chupa-lupa"`
	Host     string `yaml:"POSTGRESS_HOST" env-default:"127.0.0.1"`
	Database string `yaml:"POSTGRES_DB" env-default:"marketplace"`
	Port     string `yaml-default:"5432"`
}

func MustLoad() Config {
	//FIX: захардкожен путь конфига чтол бы не вызвать каждый раз:
	//FIX:CONFIG_PATH=config/docker.yaml go run cmd/marketplace/main.go
	// configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatalf("config file is not set")
	}
	//проверяю существоание файла
	if _, err := os.Stat(configPath); os.IsNotExist(err) {

		log.Fatalf("config file is not exist: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {

		log.Fatalf("fale read config %s", err)
	}
	return cfg
}
