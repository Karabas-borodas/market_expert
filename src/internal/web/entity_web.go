package web

import (
	"github.com/google/uuid"
)

type UserWeb struct {
	Id      uuid.UUID
	Name    string
	Sername string
	Age     uint
}

type WalletWeb struct {
	Id    uuid.UUID
	money uint
}
