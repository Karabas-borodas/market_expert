package web

import (
	"Karabas-borodas/market_expert.git/internal/config"
	"Karabas-borodas/market_expert.git/internal/domain"
)

func HTTPserverDomainToWeb(serv config.HTTPserver) *HTTPserverWEB {
	return &HTTPserverWEB{
		Address:       serv.Address,
		Timeout:       serv.Timeout,
		Iddle_timeout: serv.Iddle_timeout,
	}
}

func UserWebToDomain(u UserWeb) *domain.UserDomain {
	return &domain.UserDomain{
		Id:      u.Id,
		Name:    u.Name,
		Sername: u.Sername,
		Age:     u.Age,
	}
}
func UserDomainToWeb(u domain.UserDomain) *UserWeb {
	return &UserWeb{
		Id:      u.Id,
		Name:    u.Name,
		Sername: u.Sername,
		Age:     u.Age,
	}
}

func SignUpWebToDomain(s *SignUpRequestWeb) *domain.SignUpRequestDomain {
	return &domain.SignUpRequestDomain{
		Uuid:  s.Uuid,
		Login: s.Login,
		Pass:  s.Pass,
	}
}
func SignUpDomainToWeb(s *domain.SignUpRequestDomain) *SignUpRequestWeb {
	return &SignUpRequestWeb{
		Uuid:  s.Uuid,
		Login: s.Login,
		Pass:  s.Pass,
	}
}
