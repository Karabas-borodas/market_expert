package storage

import (
	"Karabas-borodas/market_expert.git/internal/domain"
)

func UserDomainToStorage(u domain.UserDomain) *UserRepositiry {
	return &UserRepositiry{
		Id:      u.Id,
		Name:    u.Name,
		Sername: u.Sername,
		Age:     u.Age,
	}
}

func UserStorageToDomain(u UserRepositiry) *domain.UserDomain {
	return &domain.UserDomain{
		Id:      u.Id,
		Name:    u.Name,
		Sername: u.Sername,
		Age:     u.Age,
	}
}
