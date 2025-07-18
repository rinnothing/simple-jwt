package authservice

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func InternalError(e echo.Context) error {
	return e.NoContent(http.StatusInternalServerError)
}

func Unauthorized(e echo.Context) error {
	return e.NoContent(http.StatusUnauthorized)
}
