package domain

type UserService struct {
	userStorage IUser
}

func (u *UserService) RegistUser(user SignUpRequestDomain) error {
	// _, _, err := u.userStorage.CheckUserStorage(user)
	// if err == nil {
	// 	return fmt.Errorf("пользователь уже существует")
	// }
	//
	// if user.Uuid == uuid.Nil {
	// 	user.Uuid = uuid.New()
	// }
	//
	// err = u.userStorage.SaveUserStorage(user)
	// if err != nil {
	// 	return fmt.Errorf("ошибка сохранения: %w", err)
	// }
	return nil
}
