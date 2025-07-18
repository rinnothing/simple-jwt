package storage

import (
	"context"

	"github.com/rinnothing/simple-jwt/internal/api/schema"
)

type StorageService interface {
	PutGUID(ctx context.Context, guid schema.GUID) (string, error)
	GetGUID(ctx context.Context, uuid string) (schema.GUID, error)
}
