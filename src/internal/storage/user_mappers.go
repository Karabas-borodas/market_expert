package storage

import (
	"Karabas-borodas/market_expert.git/internal/domain"
)

func UserDomainToStorage(u domain.UserDomain) *UserStorage {
	return &UserStorage{
		Id:      u.Id,
		Name:    u.Name,
		Sername: u.Sername,
		Age:     u.Age,
	}
}

func UserStorageToDomain(u UserStorage) *domain.UserDomain {
	return &domain.UserDomain{
		Id:      u.Id,
		Name:    u.Name,
		Sername: u.Sername,
		Age:     u.Age,
	}
}
