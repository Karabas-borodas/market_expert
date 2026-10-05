package storage

import (
	"github.com/google/uuid"
)

type UserRepositiry struct {
	Id      uuid.UUID
	Name    string
	Sername string
	Age     uint
}

func CreateUser()
