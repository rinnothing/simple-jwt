package authservice

import (
	"github.com/labstack/echo/v4"
	"github.com/rinnothing/simple-jwt/internal/api/schema"
	"go.uber.org/zap"
)

type AuthAPI interface {
	AuthorizeGUID(ctx echo.Context, guid string) error
	GetGUID(ctx echo.Context, params schema.GetGUIDParams) error
	RefreshTokens(ctx echo.Context) error
	Unauthorize(ctx echo.Context, params schema.UnauthorizeParams) error
}

type APIImpl struct {
	logger zap.Logger
}

func NewAPI(logger zap.Logger) AuthAPI {
	return &APIImpl{
		logger: logger,
	}
}

func (a *APIImpl) AuthorizeGUID(ctx echo.Context, guid string) error {
	panic("unimplemented")
}

func (a *APIImpl) GetGUID(ctx echo.Context, params schema.GetGUIDParams) error {
	panic("unimplemented")
}

func (a *APIImpl) RefreshTokens(ctx echo.Context) error {
	panic("unimplemented")
}

func (a *APIImpl) Unauthorize(ctx echo.Context, params schema.UnauthorizeParams) error {
	panic("unimplemented")
}
