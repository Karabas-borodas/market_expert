package storage

import (
	// "context"
	// "fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

const urldb = "postgres://Mihail:Chupa-lupa@127.0.0.1:20005/marketplace"

type MarketStorage struct {
	User     string `yaml:"POSTGRES_USER" env-default:"Mihail"`
	Password string `yaml:"POSTGRES_PASSWORD" env-default:"Chupa-lupa"`
	Database string `yaml:"POSTGRES_DB" env-default:"marketplace"`
	Port     string `env-defaul:"1863"`
}

// FIX:сделать подключение к базе данных
// FIX: впихнуть логгер в поключение

func NewGameStorage() (*MarketStorage, error) {
	cfg, err := pgxpool.ParseConfig(urldb)
	if err != nil {
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

	fmt.Println("✅ Successfully connected to PostgreSQL")
	// gs := &MarketStorage{Conn: conn}

	// if err := gs.createDB(); err != nil {
	// 	return nil, fmt.Errorf("error creating database: %w", err)
	// }
	return gs, nil
}
