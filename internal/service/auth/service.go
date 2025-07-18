package auth

import (
	"context"

	"github.com/rinnothing/simple-jwt/internal/api/schema"
)

type AuthService interface {
	IssueTokens(ctx context.Context, uuid string) (schema.TokenPair, error)
	HasAccess(ctx context.Context, token schema.AccessToken) bool
	RefreshTokens(ctx context.Context, pair schema.TokenPair) (schema.TokenPair, error)
	GetUUID(ctx context.Context, token schema.AccessToken) (schema.AccessToken, error)
	Unauthorize(ctx context.Context, token schema.AccessToken) error
}
