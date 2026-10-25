package domain

import (
	// "github.com/google/uuid"
	"log/slog"
)

type IUser interface {
	SaveUserStorage(l *slog.Logger, user SignUpRequestDomain) error
	// CheckUserStorage(user UserDomain) (string, uuid.UUID, error)
	// UpdateUserUUID(login string, newUuid uuid.UUID) error
	// GetUserByUUID(id uuid.UUID) (*UserDomain, error)
}
