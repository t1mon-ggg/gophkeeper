package web

import (
	"context"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	mw "github.com/labstack/echo/v4/middleware"
	"github.com/mgutz/ansi"
	echoSwagger "github.com/swaggo/echo-swagger"

	"github.com/t1mon-ggg/gophkeeper/pkg/logging"
	"github.com/t1mon-ggg/gophkeeper/pkg/models"
	"github.com/t1mon-ggg/gophkeeper/pkg/server/storage"
	"github.com/t1mon-ggg/gophkeeper/pkg/server/web/auth"
	mymws "github.com/t1mon-ggg/gophkeeper/pkg/server/web/middlewares"
	"github.com/t1mon-ggg/gophkeeper/pkg/server/web/static"
)

var (
	errRegFail   = echo.NewHTTPError(http.StatusBadRequest, "registration failed")
	errBadReq    = echo.NewHTTPError(http.StatusBadRequest, "Invalid request")
	errUnAuth    = echo.NewHTTPError(http.StatusForbidden, "Unauthorized")
	errIntSrvErr = echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	errStorage   = echo.NewHTTPError(http.StatusInternalServerError, "Storage unavailable")
)

// Server - web api part of application
type Server struct {
	bind    string          // bind network address
	echo    *echo.Echo      // pointer of echo framework
	logger  logging.Logger  // instance logger
	wg      *sync.WaitGroup // pointer of waitgroup
	storage storage.Storage
}

// New - web api intialization
//		bind - string with network binding address
//		logger - application logger
func New(bind string, storage storage.Storage, logger logging.Logger) *Server {
	s := new(Server)

	s.bind = bind
	s.echo = echo.New()
	s.storage = storage
	s.logger = logger.WithPrefix("web")

	s.applyMiddlewares()

	static.ApplyStatic(s.echo)

	s.createRouter()

	return s
}

// Start - func to start web-api
//  wg - application WaitGroup
func (s *Server) Start(wg *sync.WaitGroup) {
	s.wg = wg
	err := s.echo.StartTLS(s.bind, "./ssl/server.crt", "./ssl/server.pem")
	if err != nil && !strings.Contains(err.Error(), "Server closed") {
		s.log().Fatal(err, "start web failed with error")
	}
}

// Stop - gracefull shutdown of web api with shutdown timeout in 10 seconds
func (s *Server) Stop() {
	s.log().Info(nil, "Graceful shutdown in progress...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := s.echo.Shutdown(ctx)
	if err != nil {
		s.log().Fatal(err, "shutting down ends with error")
	}
	s.log().Info(nil, "Web server stopped")
	s.wg.Done()
}

// applyMiddlewares - apply middleware set to echo framework instance
func (s *Server) applyMiddlewares() *Server {
	s.echo.Use(mw.Gzip())
	s.echo.Use(mw.Decompress())
	s.echo.Use(mw.Recover())
	s.echo.Use(mw.RequestLoggerWithConfig(mw.RequestLoggerConfig{
		LogURI:      true,
		LogStatus:   true,
		LogRemoteIP: true,
		LogHost:     true,
		LogValuesFunc: func(c echo.Context, v mw.RequestLoggerValues) error {
			s.logger.
				WithFields(logging.Fields{
					"URI":      v.URI,
					"Status":   v.Status,
					"RemoteIP": v.RemoteIP,
					"Host":     v.Host,
				}).
				Info(nil, "request")
			return nil
		},
	}))
	s.echo.Use(mymws.JSONHeader)
	return s
}

// createRouter - define http route logic
func (s *Server) createRouter() *Server {

	s.echo.GET("/api/v1/ping", s.ping)                    // db liveness probe
	s.echo.POST("/api/v1/signup", s.registration)         // signup web handler
	s.echo.POST("/api/v1/signin", s.signin)               // login web handler
	s.echo.GET("/api/swagger/*", echoSwagger.WrapHandler) // swagger web page

	restricted := s.echo.Group("/api/v1/keeper") // router group with restrected access
	jwtconfig := mw.JWTConfig{
		Claims:                  &auth.JWTClaims{},
		TokenLookup:             "cookie:token",
		SigningKey:              auth.Key(),
		ErrorHandler:            auth.JWTErrorHandler,
		ErrorHandlerWithContext: auth.JWTErrorHandlerWithContext,
	}
	restricted.Use(mw.JWTWithConfig(jwtconfig))

	return s
}

// log - return web api logger
func (s *Server) log() logging.Logger {
	return s.logger
}

// registration - handler for /signup
func (s *Server) registration(c echo.Context) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	defer c.Request().Body.Close()
	if c.Request().Body == nil {
		return errBadReq
	}
	u := new(models.User)
	if err := c.Bind(u); err != nil {
		return errBadReq
	}
	s.log().Debugf("recieved registration request for user %v", nil, ansi.Color(u.Username, "green+b"))
	err := s.storage.SignUp(u.Username, u.Password)
	if err != nil {
		s.log().Warnf("registration failed with error %v", err)
		return errRegFail
	}
	err = auth.Token(u.Username, auth.Key(), c)
	if err != nil {
		s.log().Error(err, "authorization token creation error")
		return errIntSrvErr
	}
	return c.String(http.StatusOK, "{\"status\": \"OK\"}")
}

// signin - handler for /login
func (s *Server) signin(c echo.Context) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	defer c.Request().Body.Close()
	if c.Request().Body == nil {
		return errBadReq
	}
	user := new(models.User)
	if err := c.Bind(user); err != nil {
		return errBadReq
	}
	s.log().Debugf("recieved login request for user %v", nil, ansi.Color(user.Username, "green+b"))
	err := s.storage.SignIn(*user)
	if err != nil {
		s.log().Warn(err)
		return errUnAuth
	}
	err = auth.Token(user.Username, auth.Key(), c)
	if err != nil {
		s.log().Error(err, "authorization token creation error")
		return errIntSrvErr
	}
	return c.String(http.StatusOK, "{\"status\": \"OK\"}")
}

// ping - handle for /ping
func (s *Server) ping(c echo.Context) error {
	err := s.storage.Ping()
	if err != nil {
		return errStorage
	}
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	return c.String(http.StatusOK, "{\"status\": \"OK\"}")
}
