package domain

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"log/slog"
)

type UserService struct {
	userStorage IUser
}

func NewUserService(storage IUser) *UserService {
	return &UserService{userStorage: storage}
}

func (u *UserService) RegistUser(l *slog.Logger, ctx context.Context, user SignUpRequestDomain) error {
	// 1. Пытаемся найти пользователя
	_, err := u.userStorage.CheckUserInAuthStorage(ctx, user.Login)

	// 2. Если ошибки НЕТ (err == nil), значит пользователь НАЙДЕН.
	// Регистрация невозможна.
	if err == nil {
		l.Debug("user already exists", "login", user.Login)
		return fmt.Errorf("пользователь уже существует")
	}

	// 3. Если ошибка ЕСТЬ, проверяем: это "не найдено" или реальный сбой?
	// Для этого нам понадобится импорт "errors" и "github.com/jackc/pgx/v5"
	if !errors.Is(err, pgx.ErrNoRows) {
		// Это НЕ "не найдено". Это сбой базы данных (нет связи, битый SQL).
		l.Error("failed to check user", "error", err)
		return err // Возвращаем наверх, пусть хендлер отдаст 500 ошибку
	}

	// 4. Если мы здесь, значит ошибка была именно pgx.ErrNoRows.
	// Пользователя НЕТ. Можно сохранять!
	err = u.userStorage.SaveUserAuthStorage(l, user)
	if err != nil {
		l.Error("cant save user to DB", "error", err)
		return fmt.Errorf("ошибка сохранения: %w", err)
	}

	l.Info("user registered successfully", "login", user.Login)

	return nil
}
