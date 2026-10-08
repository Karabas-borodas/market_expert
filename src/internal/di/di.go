package di

import (
	"Karabas-borodas/market_expert.git/internal/config"
	// "Karabas-borodas/market_expert.git/internal/domain"
	"Karabas-borodas/market_expert.git/internal/logger"
	"Karabas-borodas/market_expert.git/internal/storage/postgres"
	"Karabas-borodas/market_expert.git/internal/web"
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

func StartProgramm() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	var userService UserService
	//NOTE: подключаем конфиг приложения(config.yaml)
	cfg := config.MustLoad()
	//NOTE: подкдючаю логгер в зависимости от рещима закгрузки приложения
	userService.log = logger.SetupLogger(cfg.Env)
	//NOTE: debug commands
	userService.log.Info("config file donload")
	userService.log.Info("start")
	poolStorage, err := postgres.NewPoolPostgress(userService.log, ctx, cfg.Posgress)
	if err != nil {
		//WARNING: не закрывается при остановке приложения
		userService.log.Error("cant create pool connect", "error", err)
		panic(fmt.Sprintf("failed to create pool connect: %v", err))
	}
	userService.pool = poolStorage
	// user := domain.GenerateUser()
	// fmt.Println(userService.pool)
	fmt.Println("USer CREATED:")
	fmt.Println(cfg)
	web.StartMarkerWeb(userService.log, cfg.HTTPserver)
}
