package postgres

import (
	"Karabas-borodas/market_expert.git/internal/config"
	"context"
	"fmt"
	// "github.com/ilyakaznacheev/cleanenv"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
	"time"
)

// создание пула подключений к постгресс
func NewPoolPostgress(u *slog.Logger, ctx context.Context, cfg config.PostgresConfig) (*pgxpool.Pool, error) {

	urldb := fmt.Sprintf("postgres://%s:%s@127.0.0.1:%s/%s",
		cfg.User,
		cfg.Password,
		cfg.Port,
		cfg.Database,
	)

	postgresPool, err := pgxpool.ParseConfig(urldb)
	if err != nil {
		//FIX: fig log fatal error
		// log.Fatalf("fale read config %s", err)
		return nil, err
	}

	postgresPool.MaxConns = 20
	postgresPool.MaxConnIdleTime = 5 * time.Minute
	postgresPool.HealthCheckPeriod = time.Minute

	conn, err := pgxpool.NewWithConfig(ctx, postgresPool)
	if err != nil {
		u.Debug("error sreating database %s", err)
		// fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		return nil, fmt.Errorf("error creating database: %w", err)
	}
	if err := (conn.Ping(ctx)); err != nil {
		conn.Close()
		u.Debug("database ping is fale ", "error", err)
		return nil, fmt.Errorf("database ping failed: %w", err)
	}
	u.Info("✅ Successfully connected to PostgreSQL")
	return conn, nil
}
