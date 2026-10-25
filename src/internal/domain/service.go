package domain

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

	err := u.userStorage.SaveUserStorage(user)
	if err != nil {
		//FIX: fix error
		return nil //l.Errorf("ошибка сохранения: %w", err)
	}
	return nil
}
