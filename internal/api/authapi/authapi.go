package authservice

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rinnothing/simple-jwt/internal/api/schema"
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

	uuid, err := a.storage.PutGUID(ctx, guid)
	if err != nil {
		return InternalError(e)
	}

	// TODO: don't forget to set e.IPExtractor to echo.ExtractIPDirect()
	pair, err := a.auth.IssueTokens(ctx, string(uuid), e.Request().UserAgent(), e.RealIP())
	if err != nil {
		return InternalError(e)
	}

	return e.JSON(http.StatusCreated, pair)
}

func (a *APIImpl) GetGUID(e echo.Context, params schema.GetGUIDParams) error {
	ctx := e.Request().Context()

	err := a.tryAuthorize(e, params.AccessToken)
	if err != nil {
		return err
	}

	uuid, err := a.auth.GetUUID(ctx, params.AccessToken)
	if err != nil {
		return InternalError(e)
	}

	guid, err := a.storage.GetGUID(ctx, uuid)
	if err != nil {
		return InternalError(e)
	}

	return e.JSON(http.StatusOK, guid)
}

func (a *APIImpl) RefreshTokens(e echo.Context) error {
	ctx := e.Request().Context()

	var pair schema.TokenPair
	e.Bind(&pair)

	err := a.tryAuthorize(e, *pair.AccessToken)
	if err != nil {
		return err
	}

	// TODO: don't forget to set e.IPExtractor to echo.ExtractIPDirect()
	newPair, err := a.auth.RefreshTokens(ctx, pair, e.Request().UserAgent(), e.RealIP())
	if err != nil {
		return Unauthorized(e)
	}

	return e.JSON(http.StatusOK, newPair)
}

func (a *APIImpl) Unauthorize(e echo.Context, params schema.UnauthorizeParams) error {
	ctx := e.Request().Context()

	err := a.tryAuthorize(e, params.AccessToken)
	if err != nil {
		return err
	}

	err = a.auth.Unauthorize(ctx, params.AccessToken)
	if err != nil {
		return InternalError(e)
	}

	return e.NoContent(http.StatusOK)
}

func (a *APIImpl) tryAuthorize(e echo.Context, token schema.AccessToken) error {
	allow, err := a.auth.HasAccess(e.Request().Context(), token)
	if err != nil {
		return InternalError(e)
	}
	if !allow {
		return e.NoContent(http.StatusUnauthorized)
	}
	return nil
}
