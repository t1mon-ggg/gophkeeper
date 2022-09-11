package auth

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"

	"github.com/t1mon-ggg/gophkeeper/pkg/helpers"
	"github.com/t1mon-ggg/gophkeeper/pkg/logging"
	"github.com/t1mon-ggg/gophkeeper/pkg/logging/zerolog"
)

var signingKey []byte
var log logging.Logger

type JWTClaims struct {
	Name string `json:"name"`
	jwt.StandardClaims
}

func init() {
	signingKey, _ = helpers.GenSecretKey(256)
	log = zerolog.New().WithPrefix("web-auth")
}

func Token(name string, key []byte, c echo.Context) error {
	expiration := time.Now().Add(time.Hour * 2)
	claims := &JWTClaims{
		name,
		jwt.StandardClaims{
			ExpiresAt: expiration.Unix(),
		},
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := t.SignedString(key)
	if err != nil {
		return err
	}
	cookie := new(http.Cookie)
	cookie.Name = "token"
	cookie.Value = token
	cookie.Expires = expiration
	c.SetCookie(cookie)
	return nil
}

func JWTErrorHandler(err error) error {
	log.Error(err, "jwt error")
	return echo.NewHTTPError(http.StatusForbidden, "Unauthorized")
}

func JWTErrorHandlerWithContext(err error, c echo.Context) error {
	log.Error(err, "jwt error")
	log.Tracef("method=%s URI=%s RemoteAddr=%s", err, c.Request().Method, c.Request().RequestURI, c.Request().RemoteAddr)
	return echo.NewHTTPError(http.StatusForbidden, "Unauthorized")
}

func Key() []byte {
	return signingKey
}
