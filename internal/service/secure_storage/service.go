package storage

import (
	"context"

	"github.com/rinnothing/simple-jwt/internal/api/schema"

	"go.uber.org/zap"
)

type StorageService interface {
	PutGUID(ctx context.Context, guid schema.GUID) (string, error)
	GetGUID(ctx context.Context, uuid string) (schema.GUID, error)
}

type StorageRepo interface {
	PutGUID(ctx context.Context, guid schema.GUID) (string, error)
	GetGUID(ctx context.Context, uuid string) (schema.GUID, error)
}

type StorageServiceImpl struct {
	l *zap.Logger

	repo StorageRepo
}

func NewService(repo StorageRepo, l *zap.Logger) StorageService {
	return &StorageServiceImpl{
		l:    l,
		repo: repo,
	}
}

func (s *StorageServiceImpl) GetGUID(ctx context.Context, uuid string) (schema.GUID, error) {
	return s.repo.GetGUID(ctx, uuid)
}

func (s *StorageServiceImpl) PutGUID(ctx context.Context, guid schema.GUID) (string, error) {
	return s.repo.PutGUID(ctx, guid)
}
