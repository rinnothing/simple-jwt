package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/rinnothing/simple-jwt/internal/api/schema"
	"github.com/rinnothing/simple-jwt/internal/config"
	"github.com/rinnothing/simple-jwt/internal/repository/postgres"
	webhook "github.com/rinnothing/simple-jwt/internal/service/webhook_caller"
	"github.com/rinnothing/simple-jwt/utils/jwt"

	"go.uber.org/zap"
)

type AuthService interface {
	IssueTokens(ctx context.Context, uuid string, userAgent string, ip string) (schema.TokenPair, error)
	HasAccess(ctx context.Context, token schema.AccessToken) (bool, error)
	RefreshTokens(ctx context.Context, pair schema.TokenPair, userAgent, ip string) (schema.TokenPair, error)
	GetUUID(ctx context.Context, token schema.AccessToken) (schema.AccessToken, error)
	Unauthorize(ctx context.Context, token schema.AccessToken) error
}

type AuthRepo interface {
	ReviveKeys(ctx context.Context) ([3]string, error)
	StoreKeys(ctx context.Context, keys [3]string) error

	PutRefresh(ctx context.Context, uuid string, oldRefresh, newRefresh schema.RefreshToken, userAgent, IP string) (bool, error)
	FindRefresh(ctx context.Context, uuid string, refresh schema.RefreshToken) (bool, error)
	Remove(ctx context.Context, uuid string) error
}

type ServiceImpl struct {
	l *zap.Logger

	repo     AuthRepo
	authTool *jwt.Tool
	webhook  webhook.WebhookService
}

func NewService(cfg config.AuthConfig, repo AuthRepo, webhook webhook.WebhookService, l *zap.Logger) (AuthService, error) {
	keys, err := repo.ReviveKeys(context.Background())
	if err == nil {
		if cfg.AccessKey == "" {
			cfg.AccessKey = keys[0]
		}
		if cfg.RefreshKey == "" {
			cfg.RefreshKey = keys[1]
		}
		if cfg.RefreshHashKey == "" {
			cfg.RefreshHashKey = keys[2]
		}
	}
	err = repo.StoreKeys(context.Background(), [3]string{cfg.AccessKey, cfg.RefreshKey, cfg.RefreshHashKey})
	if err != nil {
		return nil, fmt.Errorf("can't store keys in database: %w", err)
	}

	return &ServiceImpl{
		l:        l,
		repo:     repo,
		webhook:  webhook,
		authTool: jwt.NewJWTTool(cfg.AccessKey, cfg.RefreshKey, cfg.RefreshHashKey),
	}, nil
}

// returns UUID which token pretends to be, first you should check it with HasAccess
func (s *ServiceImpl) GetUUID(ctx context.Context, token schema.AccessToken) (schema.AccessToken, error) {
	payload, err := jwt.AccessToken(token).GetPayload()
	if err != nil {
		return "", err
	}

	return payload.UUID, nil
}

func (s *ServiceImpl) HasAccess(ctx context.Context, token schema.AccessToken) (bool, error) {
	if !s.authTool.CheckAccess(jwt.AccessToken(token)) {
		return false, nil
	}

	payload, err := jwt.AccessToken(token).GetPayload()
	if err != nil {
		return false, fmt.Errorf("can't get uuid from access token: %w", err)
	}

	refresh := s.authTool.AccessToRefresh(jwt.AccessToken(token))
	found, err := s.repo.FindRefresh(ctx, payload.UUID, schema.RefreshToken(refresh))
	if err != nil {
		return false, fmt.Errorf("can't check if access token has expired", err)
	}

	return found, nil
}

func (s *ServiceImpl) IssueTokens(ctx context.Context, uuid string, userAgent, ip string) (schema.TokenPair, error) {
	access, refresh := s.authTool.IssueTokens(uuid)

	_, err := s.repo.PutRefresh(ctx, uuid, "", schema.RefreshToken(refresh), userAgent, ip)
	if err != nil {
		return schema.TokenPair{}, fmt.Errorf("can't update refresh token in database: %w", err)
	}

	accessToken := schema.AccessToken(access)
	refreshToken := schema.RefreshToken(refresh)
	return schema.TokenPair{
		AccessToken:  &accessToken,
		RefreshToken: &refreshToken,
	}, nil
}

func (s *ServiceImpl) RefreshTokens(ctx context.Context, pair schema.TokenPair, userAgent string, ip string) (schema.TokenPair, error) {
	payload, err := jwt.AccessToken(*pair.AccessToken).GetPayload()
	if err != nil {
		return schema.TokenPair{}, fmt.Errorf("can't get uuid from access token: %w", err)
	}

	access, refresh := s.authTool.IssueTokens(payload.UUID)
	updated, err := s.repo.PutRefresh(ctx, payload.UUID, *pair.RefreshToken, schema.RefreshToken(refresh), userAgent, ip)
	if errors.Is(err, postgres.ErrWrongUserAgent) {
		err = s.Unauthorize(ctx, *pair.AccessToken)

		return schema.TokenPair{}, err
	} else if err != nil {
		return schema.TokenPair{}, fmt.Errorf("can't update refresh token in database: %w", err)
	}

	if updated {
		err = s.webhook.CallWebhook(ip)

		// TODO: think do you really need to tell user about this or just proceed if doesn't work
		if err != nil {
			return schema.TokenPair{}, fmt.Errorf("can't call webhook on IP update: %w", err)
		}
	}

	accessToken := schema.AccessToken(access)
	refreshToken := schema.RefreshToken(refresh)
	return schema.TokenPair{
		AccessToken:  &accessToken,
		RefreshToken: &refreshToken,
	}, nil
}

func (s *ServiceImpl) Unauthorize(ctx context.Context, token schema.AccessToken) error {
	payload, err := jwt.AccessToken(token).GetPayload()
	if err != nil {
		return fmt.Errorf("can't get uuid from access token: %w", err)
	}

	err = s.repo.Remove(ctx, payload.UUID)
	if err != nil {
		return fmt.Errorf("failed to remove refresh token from database: %w", err)
	}

	return nil
}
