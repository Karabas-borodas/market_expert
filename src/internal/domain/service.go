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

func (u *UserService) RegistUser(user SignUpRequestDomain) error {
	// _, _, err := u.userStorage.CheckUserStorage(user)
	// if err == nil {
	// 	return fmt.Errorf("пользователь уже существует")
	// }
	//
	// if user.Login== .Login{
	// 	user.Uuid = uuid.New()
	// }
	//FIX: пробросить логгер в самц функию а не создавать локально переменную
	var l *slog.Logger
	err := u.userStorage.SaveUserStorage(l, user)
	if err != nil {
		//FIX: fix error
		return nil //l.Errorf("ошибка сохранения: %w", err)
	}
	return nil
}
