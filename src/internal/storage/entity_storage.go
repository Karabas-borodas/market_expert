package storage

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MarketStorage struct {
	Conn *pgxpool.Pool
}

func NewMarketStorage(pool *pgxpool.Pool) *MarketStorage {
	return &MarketStorage{Conn: pool}
}

type SignUpRequestStorage struct {
	Uuid  uuid.UUID
	Login string
	Pass  string
}

type UserStorage struct {
	Id      uuid.UUID
	Name    string
	Sername string
	Age     uint
}
