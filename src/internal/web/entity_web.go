package web

import (
	"github.com/google/uuid"
	"time"
)

type HTTPserverWEB struct {
	Address       string        `yaml:"address" env-default:"localhost:20005"`
	Timeout       time.Duration `yaml:"timeout" env-default:"10s"`
	Iddle_timeout time.Duration `yaml:"iddle_timeout" env-default:"80s"`
}

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

type SignUpRequestWeb struct {
	Uuid  uuid.UUID
	Login string `json:"login"`
	Pass  string `json:"password"`
}
