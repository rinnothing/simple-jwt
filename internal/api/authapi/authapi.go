package authservice

import (
	"errors"
	"net/http"
	"slices"

	"github.com/labstack/echo/v4"
	"github.com/rinnothing/simple-jwt/internal/api/schema"
	"github.com/rinnothing/simple-jwt/internal/repository/postgres"
	"github.com/rinnothing/simple-jwt/internal/service/auth"
	storage "github.com/rinnothing/simple-jwt/internal/service/secure_storage"

	"go.uber.org/zap"
)

type AuthAPI interface {
	AuthorizeGUID(ctx echo.Context, guid string) error
	GetGUID(ctx echo.Context, params schema.GetGUIDParams) error
	RefreshTokens(ctx echo.Context) error
	Unauthorize(ctx echo.Context, params schema.UnauthorizeParams) error
}

type APIImpl struct {
	logger *zap.Logger

	auth    auth.AuthService
	storage storage.StorageService
}

func NewAPI(auth auth.AuthService, storage storage.StorageService, logger *zap.Logger) AuthAPI {
	return &APIImpl{
		logger:  logger,
		auth:    auth,
		storage: storage,
	}
}

func (a *APIImpl) AuthorizeGUID(e echo.Context, guid string) error {
	ctx := e.Request().Context()
	a.logRequest(e, "authorize", zap.String("guid", guid))

	uuid, err := a.storage.PutGUID(ctx, guid)
	if err != nil {
		a.logger.Error("can't put guid in storage", zap.Error(err))
		return InternalError(e)
	}

	// TODO: don't forget to set e.IPExtractor to echo.ExtractIPDirect()
	pair, err := a.auth.IssueTokens(ctx, string(uuid), e.Request().UserAgent(), e.RealIP())
	if err != nil {
		a.logger.Error("can't issue tokens", zap.Error(err))
		return InternalError(e)
	}

	return e.JSON(http.StatusCreated, pair)
}

func (a *APIImpl) GetGUID(e echo.Context, params schema.GetGUIDParams) error {
	ctx := e.Request().Context()
	a.logRequest(e, "get_guid", zap.String("access_token", params.AccessToken))

	err := a.tryAuthorize(e, params.AccessToken)
	if err != nil {
		return err
	}

	uuid, err := a.auth.GetUUID(ctx, params.AccessToken)
	if err != nil {
		a.logger.Error("can't get uuid from access token", zap.Error(err))
		return InternalError(e)
	}

	guid, err := a.storage.GetGUID(ctx, uuid)
	if err != nil {
		a.logger.Error("can't get guid from storage", zap.Error(err))
		return InternalError(e)
	}

	return e.JSON(http.StatusOK, guid)
}

func (a *APIImpl) RefreshTokens(e echo.Context) error {
	ctx := e.Request().Context()

	var pair schema.TokenPair
	err := e.Bind(&pair)
	if err != nil {
		a.logger.Error("can't unmarshal request", zap.Error(err))
		return BadRequest(e, err.Error())
	}

	a.logRequest(e, "refresh", zap.String("access_token", *pair.AccessToken), zap.String("refresh_token", *pair.RefreshToken))

	err = a.tryAuthorize(e, *pair.AccessToken)
	if err != nil {
		return err
	}

	// TODO: don't forget to set e.IPExtractor to echo.ExtractIPDirect()
	newPair, err := a.auth.RefreshTokens(ctx, pair, e.Request().UserAgent(), e.RealIP())
	if errors.Is(err, postgres.ErrWrongUserAgent) {
		a.logger.Info("refresh token denied", zap.String("access_token", string(*pair.AccessToken)),
			zap.String("refresh_token", string(*pair.RefreshToken)))
		return Unauthorized(e)
	}
	if err != nil {
		a.logger.Error("can't refresh tokens", zap.Error(err))
		return InternalError(e)
	}

	return e.JSON(http.StatusOK, newPair)
}

func (a *APIImpl) Unauthorize(e echo.Context, params schema.UnauthorizeParams) error {
	ctx := e.Request().Context()

	a.logRequest(e, "unauthorize", zap.String("access_token", params.AccessToken))

	err := a.tryAuthorize(e, params.AccessToken)
	if err != nil {
		return err
	}

	err = a.auth.Unauthorize(ctx, params.AccessToken)
	if err != nil {
		a.logger.Error("can't unauthorize user", zap.Error(err))
		return InternalError(e)
	}

	return e.NoContent(http.StatusOK)
}

func (a *APIImpl) tryAuthorize(e echo.Context, token schema.AccessToken) error {
	allow, err := a.auth.HasAccess(e.Request().Context(), token)
	if err != nil {
		a.logger.Error("can't check access", zap.Error(err))
		return InternalError(e)
	}
	if !allow {
		a.logger.Info("access denied", zap.String("access_token", string(token)))
		return e.NoContent(http.StatusUnauthorized)
	}
	return nil
}

func (a *APIImpl) logRequest(e echo.Context, name string, fields ...zap.Field) {
	a.logger.Info("got request",
		slices.Concat(
			[]zap.Field{zap.String("name", name), zap.String("user_agent", e.Request().UserAgent()), zap.String("ip", e.RealIP())},
			fields,
		)...,
	)
}
