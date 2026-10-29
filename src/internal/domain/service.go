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
	_, err := u.userStorage.CheckUserInAuthStorage(ctx, user.Login)

	if err == nil {
		l.Debug("user already exists", "login", user.Login)
		return fmt.Errorf("пользователь уже существует")
	}

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

func (u *UserService) SaveUserInfo(l *slog.Logger, ctx context.Context, user UserDomain) error {
	//FIX:написать функцию проверки пользователя в памяти
	userExist, err := u.userStorage.CheckUserInAuthStorageUUID(ctx, user.Id)

	if err == nil {
		l.Debug("user already exists")
		return fmt.Errorf("пользователь уже существует")
	}

	if !errors.Is(err, pgx.ErrNoRows) {
		// Это НЕ "не найдено". Это сбой базы данных (нет связи, битый SQL).
		l.Error("failed to check user", "error", err)
		return err // Возвращаем наверх, пусть хендлер отдаст 500 ошибку
	}

	// 4. Если мы здесь, значит ошибка была именно pgx.ErrNoRows.
	// Пользователя НЕТ. Можно сохранять!
	err = u.userStorage.SaveUserStorage(l, user)
	if err != nil {
		l.Error("cant save user to DB", "error", err)
		return fmt.Errorf("ошибка сохранения: %w", err)
	}

	l.Info("user registered successfully", "login", userExist.Login)

	return nil
}
