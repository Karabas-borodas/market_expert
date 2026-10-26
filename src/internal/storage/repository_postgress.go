package storage

import (
	"Karabas-borodas/market_expert.git/internal/domain"
	"context"
	"fmt"
	"github.com/google/uuid"
	"log/slog"
	"time"
)

// FIX: заменить возсратны параметры на структуру
//
//	func (m *MarketStorage) CheckUserAuthStorage(l *slog.Logger, user domain.SignUpRequestDomain) (string, uuid.UUID, error) {
//		dbUser := SignUpRequestStorage(user)
//		const sql = `
//		SELECT pass, uuid FROM users_auth WHERE login=$1;`
//		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
//		defer cancel()
//		var password string
//		var id uuid.UUID
//		err := m.Conn.QueryRow(ctx, sql, dbUser.Login).Scan(&password, &id)
//		if err != nil {
//			l.Info("user not found")
//			return "", uuid.Nil, err
//		}
//		return password, id, nil
//	}
//
// WARNING: решить проблему с функйией
// func (m *MarketStorage) CheckUserAuthStorage(ctx context.Context, login string) (bool, error) {
// 	const sql = `SELECT EXISTS(SELECT 1 FROM users WHERE login = $1)`
//
// 	var exists bool
// 	err := m.Conn.QueryRow(ctx, sql, login).Scan(&exists)
// 	if err != nil {
// 		return false, fmt.Errorf("query user exists: %w", err)
// 	}
// 	return exists, nil
// }

// NOTE: таблица пользовательский аутентификации
func (g *MarketStorage) CreateUserAuth(l *slog.Logger) error {
	const sql = `
	CREATE TABLE IF NOT EXISTS users_auth(
		id SERIAL PRIMARY KEY,
		uuid UUID UNIQUE NOT NULL,
		login TEXT UNIQUE NOT NULL,
		pass TEXT NOT NULL
	);`
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := g.Conn.Exec(ctx, sql)
	if err != nil {
		l.Debug("Cant create user_auth table", err)
		return err
	}
	return nil
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

func (u *MarketStorage) CheckUserInAuthStorage(ctx context.Context, dbUser string) (domain.SignUpRequestDomain, error) {
	const sql = `SELECT pass, uuid FROM users_auth WHERE login=$1`

	var userStorage SignUpRequestStorage

	// Используем переданный ctx, не создаем новый!
	err := u.Conn.QueryRow(ctx, sql, dbUser).Scan(&userStorage.Pass, &userStorage.Uuid)
	if err != nil {
		// Если ошибка - возвращаем пустую структуру и ошибку (будь то ErrNoRows или сбой БД)
		return domain.SignUpRequestDomain{}, fmt.Errorf("error querying user: %w", err)
	}

	// Сюда мы попадаем ТОЛЬКО если err == nil (пользователь найден)
	userStorage.Login = dbUser
	userDomain := RquestStoregeToDomain(userStorage)

	return *userDomain, nil
}

// func (m *MarketStorage) SaveUserStorage(l *slog.Logger, user domain.SignUpRequestDomain) error {
// 	const sql = `
// 		INSERT INTO users (uuid, login, pass, sername)
// 		VALUES ($1, $2, $3, $4);
// 		`
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()
//
// 	newUUID := uuid.New() // генерируем uuid здесь
// 	_, err := m.Conn.Exec(ctx, sql, newUUID, user.Login, user.Pass, "")
// 	if err != nil {
// 		l.Info("error savin user to table ", err)
// 		return err
// 	}
// 	return nil
// }

func (m *MarketStorage) SaveUserAuthStorage(l *slog.Logger, user domain.SignUpRequestDomain) error {
	const sql = `
		INSERT INTO users_auth (uuid, login, pass)
		VALUES ($1, $2, $3);
		`
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	newUUID := uuid.New() // генерируем uuid здесь
	_, err := m.Conn.Exec(ctx, sql, newUUID, user.Login, user.Pass)
	if err != nil {
		l.Info("error savin user to table ", err)
		return err
	}
	return nil
}
