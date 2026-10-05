package domain

import (
	"github.com/google/uuid"
)

type UserDomain struct {
	Id      uuid.UUID
	Name    string
	Sername string
	Age     uint
}

type WalletDomain struct {
	Id    uuid.UUID
	money uint
}

// func GenerateUser() *UserDomain {
// 	var user UserDomain
// 	user.Id = uuid.New()
// 	return &user
// }
// func GenerateWallet(u UserDomain) *WalletDomain {
// 	var wallet WalletDomain
// 	wallet.Id = u.Id
// 	return &wallet
// }
