package storage

import (
	"context"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
	"time"
)

// const urldb = "postgres://Mihail:Chupa-lupa@127.0.0.1:5432/marketplace"

type MarketStorage struct {
	User     string `env:"POSTGRES_USER" env-default:"Mihail"`
	Password string `env:"POSTGRES_PASSWORD" env-default:"Chupa-lupa"`
	Database string `env:"POSTGRES_DB" env-default:"marketplace"`
	Port     string `env-default:"5432"`
	pool     *pgxpool.Pool
}

func CreateURL(log *slog.Logger) (string, error) {
	connectStruct := MarketStorage{}

	if err := cleanenv.ReadEnv(&connectStruct); err != nil {
		log.Debug("cant make url to connect to bd", "error", err)
		return "", err
	}

	url := fmt.Sprintf(
		"postgres://%s:%s@127.0.0.1:%s/%s",
		connectStruct.User,
		connectStruct.Password,
		connectStruct.Port,
		connectStruct.Database,
	)
	fmt.Printf("---- %v---- ", connectStruct)
	return url, nil
}

// FIX:сделать подключение к базе данных
// FIX: впихнуть логгер в поключение
// HACK: разделить подключение к бд и вынести создание пула подключений в мэйн
func NewMarkerStorage(log *slog.Logger) (*MarketStorage, error) {
	// var urlStruct MarketStorage = MarketStorage{}
	urldb, err := CreateURL(log)
	if err != nil {
		return nil, err
	}
	fmt.Printf("-------- URL %v ------\n", urldb)
	cfg, err := pgxpool.ParseConfig(urldb)
	if err != nil {
		//WARNING: а должен ли я менять fmt на logger тут
		return nil, fmt.Errorf("parse pool config: %w", err)
	}

	cfg.MaxConns = 20
	cfg.MaxConnIdleTime = 5 * time.Minute
	cfg.HealthCheckPeriod = time.Minute

	conn, err := pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		// fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		return nil, fmt.Errorf("error creating database: %w", err)
	}
	if err := (conn.Ping(context.Background())); err != nil {
		conn.Close()
		return nil, fmt.Errorf("database ping failed: %w", err)
	}
	//WARNING: какие дучшие практики??
	log.Info("✅ Successfully connected to PostgreSQL")
	gs := &MarketStorage{pool: conn}
	return gs, nil
}
