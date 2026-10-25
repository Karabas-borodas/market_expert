package domain

import (
	"log/slog"
)

type UserService struct {
	userStorage IUser
}

func NewUserService(storage IUser) *UserService {
	return &UserService{userStorage: storage}
}

func (u *UserService) RegistUser(l *slog.Logger, user SignUpRequestDomain) error {
	_, _, err := u.userStorage.CheckUserAuthStorage(l, user.Login)
	if err != nil {
		l.Debug("user is exist %s", err)
		return err
	}
	// if user.Login== .Login{
	// 	user.Uuid = uuid.New()
	// }
	err = u.userStorage.SaveUserAuthStorage(l, user)
	if err != nil {
		l.Debug("Cant save usert to DB", err)
		return nil //l.Errorf("ошибка сохранения: %w", err)
	}
	return nil
}
