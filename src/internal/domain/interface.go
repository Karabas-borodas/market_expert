package domain

import (
	"context"
	"github.com/google/uuid"
	"log/slog"
)

type IUser interface {
	SaveUserAuthStorage(l *slog.Logger, user SignUpRequestDomain) error
	CheckUserInAuthStorage(ctx context.Context, login string) (SignUpRequestDomain, error)
	CheckUserInAuthStorageUUID(ctx context.Context, Id uuid.UUID) (SignUpRequestDomain, error)
	SaveUserStorage(l *slog.Logger, user UserDomain) error
}
