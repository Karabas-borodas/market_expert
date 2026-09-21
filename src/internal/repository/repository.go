package repository

import (
	// "context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
)

type UserRepository struct {
	log  *slog.Logger
	pool *pgxpool.Pool
}

func NewUserRepositury(log *slog.Logger, pool *pgxpool.Pool) *UserRepository {

	return &UserRepository{
		log:  log,
		pool: pool,
	}
}
