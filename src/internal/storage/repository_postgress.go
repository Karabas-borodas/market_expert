package storage

import (
	"Karabas-borodas/market_expert.git/internal/domain"
	"context"
	"fmt"
	"github.com/google/uuid"
	"log/slog"
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

func (g *MarketStorage) CreateDB(l *slog.Logger) error {
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
		l.Debug("Cant connect to db", err)
		return err
	}
	return nil
}

func (m *MarketStorage) SaveUserStorage(user domain.SignUpRequestDomain) error {
	const sql = `
		INSERT INTO users (login, pass)
		VALUES ($1, $2);
		`
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := m.Conn.Exec(ctx, sql, user.Login, user.Pass)
	if err != nil {
		return fmt.Errorf("error saving user to table: %w", err)
	}
	return nil
}
