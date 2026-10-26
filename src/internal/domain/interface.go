package domain

import (
	"context"
	"log/slog"
)

type IUser interface {
	// SaveUserStorage(l *slog.Logger, user SignUpRequestDomain) error
	SaveUserAuthStorage(l *slog.Logger, user SignUpRequestDomain) error
	// CheckUserAuthStorage(ctx context.Context, login string) (bool, error)
	CheckUserInAuthStorage(ctx context.Context, login string) (SignUpRequestDomain, error)
	// CheckUserStorage(user UserDomain) (string, uuid.UUID, error)
	// UpdateUserUUID(login string, newUuid uuid.UUID) error
	// GetUserByUUID(id uuid.UUID) (*UserDomain, error)
}
