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
	// 1. Проверяем, существует ли пользователь в auth
	userAuth, err := u.userStorage.CheckUserInAuthStorageUUID(ctx, user.Id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			l.Debug("user not found in auth", "uuid", user.Id)
			return fmt.Errorf("пользователь не найден, невозможно сохранить информацию")
		}
		l.Error("failed to check user auth", "error", err)
		return err
	}

	// 2. Проверяем, есть ли уже профиль в users
	isInfoExists, err := u.userStorage.CheckUserInUserStorage(ctx, userAuth.Uuid)
	if err != nil {
		// Если функция сделана через EXISTS, эта проверка не нужна,
		// но если ты делаешь SELECT - оставь её.
		l.Error("failed to check user info", "error", err)
		return err
	}

	// 3. Если профиль уже есть - выходим
	if isInfoExists {
		l.Debug("user info already exists", "uuid", user.Id)
		return fmt.Errorf("информация о пользователе уже существует")
	}

	// 4. Сохраняем, если профиля нет
	err = u.userStorage.SaveUserStorage(l, user)
	if err != nil {
		l.Error("cant save user info to DB", "error", err)
		return fmt.Errorf("ошибка сохранения: %w", err)
	}

	l.Info("user info saved successfully", "login", userAuth.Login)
	return nil
}
