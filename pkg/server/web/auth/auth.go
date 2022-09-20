package auth

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"

	"github.com/t1mon-ggg/gophkeeper/pkg/helpers"
	"github.com/t1mon-ggg/gophkeeper/pkg/logging/zerolog"
)

var signingKey []byte

type JWTClaims struct {
	Name string `json:"name"`
	jwt.StandardClaims
}

func init() {
	signingKey, _ = helpers.GenSecretKey(256)
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

	zerolog.New().WithPrefix("auth").Error(err, "jwt error")
	return echo.NewHTTPError(http.StatusForbidden, "Unauthorized")
}

func JWTErrorHandlerWithContext(err error, c echo.Context) error {
	zerolog.New().WithPrefix("auth").Error(err, "jwt error")
	zerolog.New().WithPrefix("auth").Tracef("method=%s URI=%s RemoteAddr=%s", err, c.Request().Method, c.Request().RequestURI, c.Request().RemoteAddr)
	return echo.NewHTTPError(http.StatusForbidden, "Unauthorized")
}

func Key() []byte {
	return signingKey
}
