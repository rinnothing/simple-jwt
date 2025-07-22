package authservice

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func InternalError(e echo.Context) error {
	return e.String(http.StatusInternalServerError, "internal error")
}

func Unauthorized(e echo.Context) error {
	return e.String(http.StatusUnauthorized, "unauthorized")
}

func BadRequest(e echo.Context, reason string) error {
	return e.String(http.StatusBadRequest, fmt.Sprintf("bad request, reason: %s", reason))
}
