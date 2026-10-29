package storage

import (
	"Karabas-borodas/market_expert.git/internal/domain"
)

func UserDomainToStorage(u domain.UserDomain) *UserStorage {
	return &UserStorage{
		Phone:   u.Phone,
		Email:   u.Email,
		Id:      u.Id,
		Name:    u.Name,
		Sername: u.Sername,
		Age:     u.Age,
	}
}

func UserStorageToDomain(u UserStorage) *domain.UserDomain {
	return &domain.UserDomain{
		Phone:   u.Phone,
		Email:   u.Email,
		Id:      u.Id,
		Name:    u.Name,
		Sername: u.Sername,
		Age:     u.Age,
	}
}

func RequestDomainToStorage(u domain.SignUpRequestDomain) *SignUpRequestStorage {
	return &SignUpRequestStorage{
		Uuid:  u.Uuid,
		Login: u.Login,
		Pass:  u.Pass,
	}
}

func RquestStoregeToDomain(u SignUpRequestStorage) *domain.SignUpRequestDomain {

	return &domain.SignUpRequestDomain{
		Uuid:  u.Uuid,
		Login: u.Login,
		Pass:  u.Pass,
	}
}
