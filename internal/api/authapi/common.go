package authapi

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func InternalError(e echo.Context) error {
	return e.String(http.StatusInternalServerError, "internal error\n")
}

func Unauthorized(e echo.Context) error {
	return e.String(http.StatusUnauthorized, "unauthorized\n")
}

func BadRequest(e echo.Context, reason string) error {
	return e.String(http.StatusBadRequest, fmt.Sprintf("bad request, reason: %s\n", reason))
}
