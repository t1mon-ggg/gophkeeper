package middlewares

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

var (
	// errRegFail   = echo.NewHTTPError(http.StatusBadRequest, "registration failed")
	errBadReq = echo.NewHTTPError(http.StatusBadRequest, "Invalid request")
	// errUnAuth    = echo.NewHTTPError(http.StatusForbidden, "Unauthorized")
	// errIntSrvErr = echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	// errStorage   = echo.NewHTTPError(http.StatusInternalServerError, "Storage unavailable")
)

// jsonHeader middleware checks application header
func JSONHeader(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		header := c.Request().Header.Get("Content-Type")
		if header != echo.MIMEApplicationJSON &&
			c.Request().RequestURI != "/" &&
			!strings.Contains(c.Request().RequestURI, "swagger") &&
			!strings.Contains(c.Request().RequestURI, "favicon") &&
			!strings.Contains(c.Request().RequestURI, "gophkeeper.png") {
			return errBadReq
		}
		return next(c)
	}
}

// jsonHeader middleware checks application header
func Logger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		return next(c)
	}
}
