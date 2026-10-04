package service

import (
	"github.com/google/uuid"
)

type UserDomain struct {
	Id      uuid.UUID
	Name    string
	Sername string
	Age     uint
}

func GenerateUser() *UserDomain {
	var user UserDomain
	user.Id = uuid.New()
	return &user
}
