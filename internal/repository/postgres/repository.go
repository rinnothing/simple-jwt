package postgres

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/rinnothing/simple-jwt/internal/api/schema"
	"github.com/rinnothing/simple-jwt/internal/config"
	"golang.org/x/crypto/bcrypt"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

var (
	ErrWrongUserAgent = errors.New("different user agent")
)

type PostgresService interface {
	ReviveKeys(ctx context.Context) ([]string, error)
	StoreKeys(ctx context.Context, keys []string) error

	PutRefresh(ctx context.Context, uuid string, oldRefresh, newRefresh schema.RefreshToken, userAgent, IP string) (bool, error)
	Remove(ctx context.Context, uuid string) error
	FindRefresh(ctx context.Context, uuid string, refresh schema.RefreshToken) (bool, error)

	PutGUID(ctx context.Context, guid schema.GUID) (string, error)
	GetGUID(ctx context.Context, uuid string) (schema.GUID, error)
}

type PostgresServiceImpl struct {
	l *zap.Logger

	pool *pgxpool.Pool

	cfg config.PostgresConfig
}

func NewRepo(cfg config.PostgresConfig, pool *pgxpool.Pool, l *zap.Logger) PostgresService {
	return &PostgresServiceImpl{
		l:    l,
		pool: pool,
		cfg:  cfg,
	}
}

func (p *PostgresServiceImpl) PutGUID(ctx context.Context, guid schema.GUID) (string, error) {
	query := `
INSERT INTO storage (guid)
VALUES ($1)
RETURNING id
`
	var uuid string
	err := p.pool.QueryRow(ctx, query, guid).Scan(&uuid)
	if err != nil {
		return "", fmt.Errorf("can't insert guid: %w", err)
	}

	return uuid, nil
}

func (p *PostgresServiceImpl) GetGUID(ctx context.Context, uuid string) (schema.GUID, error) {
	query := `
SELECT guid 
FROM storage
WHERE id = $1
`
	var guid schema.GUID
	err := p.pool.QueryRow(ctx, query, uuid).Scan(&guid)
	if err != nil {
		return "", fmt.Errorf("can't find guid for uuid %s: %w", uuid, err)
	}

	return guid, nil
}

func hashRefresh(refresh schema.RefreshToken) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(refresh), bcrypt.DefaultCost)
}

func (p *PostgresServiceImpl) PutRefresh(ctx context.Context, uuid string, oldRefresh, newRefresh schema.RefreshToken, userAgent string, IP string) (bool, error) {
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return false, fmt.Errorf("can't start transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	queryGet := `
SELECT user_agent, ip
FROM auth
WHERE uuid = $1 AND refresh_hash = $2
`
	oldRefreshHash, err := hashRefresh(oldRefresh)
	if err != nil {
		return false, fmt.Errorf("can't generate refresh hash: %w", err)
	}

	var storedUserAgent, storedIP string
	tx.QueryRow(ctx, queryGet, uuid, oldRefreshHash).Scan(&storedUserAgent, &storedIP)

	if storedUserAgent != userAgent {
		return false, ErrWrongUserAgent
	}

	querySet := `
UPDATE auth
SET refresh_hash = $1, ip = $2
`

	newRefreshHash, err := hashRefresh(newRefresh)
	if err != nil {
		return false, fmt.Errorf("can't generate refresh hash: %w", err)
	}

	_, err = tx.Exec(ctx, querySet, newRefreshHash, IP)
	if err != nil {
		return false, fmt.Errorf("can't update refresh token: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return false, fmt.Errorf("can't commit transaction: %w", err)
	}

	return storedIP != IP, nil
}

func (p *PostgresServiceImpl) FindRefresh(ctx context.Context, uuid string, refresh schema.RefreshToken) (bool, error) {
	query := `
SELECT refresh_hash
FROM auth
WHERE uuid = $1
`
	var retrievedHash []byte
	err := p.pool.QueryRow(ctx, query, uuid).Scan(&retrievedHash)
	if err != nil {
		return false, fmt.Errorf("can't find refresh_hash for uuid %s: %w", uuid, err)
	}

	refreshHash, err := hashRefresh(refresh)
	if err != nil {
		return false, fmt.Errorf("can't generate refresh hash: %w", err)
	}

	return slices.Equal(retrievedHash, refreshHash), nil
}

func (p *PostgresServiceImpl) Remove(ctx context.Context, uuid string) error {
	query := `
DELETE FROM storage
WHERE uuid = $1
`

	_, err := p.pool.Exec(ctx, query, uuid)
	if err != nil {
		fmt.Errorf("failed to remove uuid %s from authorized: %w", uuid, err)
	}

	return nil
}

func (p *PostgresServiceImpl) ReviveKeys(ctx context.Context) ([]string, error) {
	query := `
SELECT access_key, refresh_key, refresh_hash_key
FROM keys
`
	keys := make([]string, 3)
	err := p.pool.QueryRow(ctx, query).Scan(&keys[0], &keys[1], &keys[2])
	if err != nil {
		return nil, fmt.Errorf("can't revive keys: %w", err)
	}

	return keys, nil
}

func (p *PostgresServiceImpl) StoreKeys(ctx context.Context, keys []string) error {
	if len(keys) < 3 {
		return fmt.Errorf("not enough keys to store, want 3 but has %d", len(keys))
	}

	query := `
UPDATE keys
SET access_key = $1, refresh_key = $2, refresh_hash_key = $3
`
	_, err := p.pool.Exec(ctx, query, keys[0], keys[1], keys[2])
	if err != nil {
		fmt.Errorf("can't store keys: %w", err)
	}

	return nil
}
