package main

import (
	"Karabas-borodas/market_expert.git/internal/config"
	"Karabas-borodas/market_expert.git/internal/logger"
	// "Karabas-borodas/market_expert.git/internal/storage"
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"
	// "fmt"
)

type UserService struct {
	log  *slog.Logger
	pool *pgxpool.Pool
}

func main() {

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	var userService UserService
	//NOTE: подключаем конфиг приложения(config.yaml)
	cfg := config.MustLoad()
	userService.log = logger.SetupLogger(cfg.Env)
	//NOTE: debug commands
	userService.log.Info("config file donload")
	userService.log.Info("start")
	userService.log.Debug("debug")
	storage, err := postgres.NewMarkerStorage(userService.log, ctx, cfg.Postgres)
	if err != nil {
		userService.log.Error("cant connect to DB", "error", err)
	}
	fmt.Println(storage)
}
