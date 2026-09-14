package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

// FIX: захардкожен путь конфига чтол бы не вызвать каждый раз:
const configPath = "config/docker.yaml"

type Config struct {
	Env        string `yaml:"env" env-default:"local"`
	HTTPserver `yaml:"http_server"`
}

type HTTPserver struct {
	Address       string        `yaml: "address" env-default:"localhost:20005"`
	Timeout       time.Duration `yaml:"timeout" env-default:"10s"`
	Iddle_timeout time.Duration `yaml:"iddle_timeout" env-default:"80s"`
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
