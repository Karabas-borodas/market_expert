package main

import (
	"Karabas-borodas/market_expert.git/internal/config"
	"Karabas-borodas/market_expert.git/internal/logger"
	"Karabas-borodas/market_expert.git/internal/storage"
	"fmt"
	// "fmt"
)

func main() {
	//NOTE: подключаем конфиг(docker.yaml)
	cfg := config.MustLoad()
	log := logger.SetupLogger(cfg.Env)
	log.Info("start")
	log.Debug("debug")
	storage, err := storage.NewGameStorage()
	if err != nil {
		log.Error("cant connect to DB %v", err)
	}
	fmt.Println(storage)
}
