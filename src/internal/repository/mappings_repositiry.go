package repository

import (
	// "github.com/google/uuid"
	"Karabas-borodas/market_expert.git/internal/service"
)

func UserDomainToStorage(u service.UserDomain) *service.UserDomain {
	return &service.UserDomain{
		Id:      u.Id,
		Name:    u.Name,
		Sername: u.Sername,
		Age:     u.Age,
	}
}
