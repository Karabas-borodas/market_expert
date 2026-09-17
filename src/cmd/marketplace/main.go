package main

import (
	"Karabas-borodas/market_expert.git/internal/config"
	"Karabas-borodas/market_expert.git/internal/logger"
	"Karabas-borodas/market_expert.git/internal/storage"
	"fmt"
	"log/slog"
	// "fmt"
)

type UserService struct {
	log *slog.Logger
}

func main() {
	//NOTE: подключаем конфиг приложения(config.yaml)
	cfg := config.MustLoad()
	log := logger.SetupLogger(cfg.Env)
	//NOTE: debug commands
	log.Info("config file donload")
	log.Info("start")
	log.Debug("debug")
	storage, err := storage.NewMarkerStorage(log)
	if err != nil {
		log.Error("cant connect to DB", "error", err)
	}
	fmt.Println(storage)
}
