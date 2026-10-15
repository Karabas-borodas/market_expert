package storage

import (
	"Karabas-borodas/market_expert.git/internal/domain"
	"context"
	"fmt"
	"github.com/google/uuid"
	"time"
)

func (u *UserStorage) CheckUserStorage(user domain.UserDomain) (string, uuid.UUID, error) {
	// dbUser := UserDomainToStorage(user)
	const sql = `
	SELECT password, uuid FROM users WHERE login=$1;`
	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancel()
	var password string
	var id uuid.UUID
	// err := u.Conn.QueryRow(ctx, sql, dbUser.Login).Scan(&password, &id)
	// if err != nil {
	// 	return "", uuid.Nil, l.Errorf("error querying user: %w", err)
	// }
	return password, id, nil
}

func (g *MarketStorage) CreateDB() error {
	const sql = `
	CREATE TABLE IF NOT EXISTS users(
		id SERIAL PRIMARY KEY,
		uuid UUID UNIQUE NOT NULL,
		login TEXT UNIQUE NOT NULL,
		pass TEXT NOT NULL,
		sername TEXT NOT NULL,
		age INT NULL DEFAULT 0
	);
		`

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := g.Conn.Exec(ctx, sql)
	if err != nil {
		return fmt.Errorf("create table/migration execution: %w", err)
	}
	return nil
}
