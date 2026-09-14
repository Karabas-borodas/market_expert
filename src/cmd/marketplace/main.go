package main

import (
	"Karabas-borodas/market_expert.git/internal/config"
	"Karabas-borodas/market_expert.git/internal/logger"
	// "fmt"
)

func main() {
	//NOTE: подключаем конфиг(docker.yaml)
	cfg := config.MustLoad()
	log := logger.SetupLogger(cfg.Env)
	log.Info("start")
	log.Debug("debug")
}
