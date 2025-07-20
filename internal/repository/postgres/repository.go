package postgres

import (
	"context"
	"errors"

	"github.com/rinnothing/simple-jwt/internal/api/schema"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

var (
	ErrWrongUserAgent = errors.New("different user agent")
)

type PostgresService interface {
	ReviveKeys(ctx context.Context) ([3]string, error)
	StoreKeys(ctx context.Context, keys [3]string) error

	PutRefresh(ctx context.Context, uuid string, refresh schema.RefreshToken, userAgent, IP string) (bool, error)
	RemoveRefresh(ctx context.Context, uuid string, refresh schema.RefreshToken) error
	FindRefresh(ctx context.Context, uuid string, refresh schema.RefreshToken) (bool, error)

	PutGUID(ctx context.Context, guid schema.GUID) (string, error)
	GetGUID(ctx context.Context, uuid string) (schema.GUID, error)
}

type PostgresServiceImpl struct {
	l *zap.Logger

	pool *pgxpool.Pool
}

func NewRepo(pool *pgxpool.Pool, l *zap.Logger) PostgresService {
	return &PostgresServiceImpl{
		l:    l,
		pool: pool,
	}
}

func (p *PostgresServiceImpl) FindRefresh(ctx context.Context, uuid string, refresh schema.RefreshToken) (bool, error) {
	panic("unimplemented")
}

func (p *PostgresServiceImpl) GetGUID(ctx context.Context, uuid string) (schema.GUID, error) {
	panic("unimplemented")
}

func (p *PostgresServiceImpl) PutGUID(ctx context.Context, guid schema.GUID) (string, error) {
	panic("unimplemented")
}

func (p *PostgresServiceImpl) PutRefresh(ctx context.Context, uuid string, refresh schema.RefreshToken, userAgent string, IP string) (bool, error) {
	panic("unimplemented")
}

func (p *PostgresServiceImpl) RemoveRefresh(ctx context.Context, uuid string, refresh schema.RefreshToken) error {
	panic("unimplemented")
}

func (p *PostgresServiceImpl) ReviveKeys(ctx context.Context) ([3]string, error) {
	panic("unimplemented")
}

func (p *PostgresServiceImpl) StoreKeys(ctx context.Context, keys [3]string) error {
	panic("unimplemented")
}
