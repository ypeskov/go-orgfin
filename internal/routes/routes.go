package routes

import (
	"github.com/a-h/templ"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"ypeskov/go-password-manager/cmd/web"
	"ypeskov/go-password-manager/cmd/web/components"
	"ypeskov/go-password-manager/internal/config"
	"ypeskov/go-password-manager/internal/logger"
	customMiddleware "ypeskov/go-password-manager/internal/middleware"
	"ypeskov/go-password-manager/services"
)

var log *logger.Logger
var sManager *services.ServiceManager
var cfg *config.Config

func RegisterRoutes(logger *logger.Logger, servicesManager *services.ServiceManager, configInstance *config.Config) *echo.Echo {
	cfg = configInstance
	log = logger
	log.Info("Registering routes")
	e := echo.New()

	sManager = servicesManager

	//e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(customMiddleware.CookieToHeaderMiddleware)

	fileServer := http.FileServer(http.FS(web.Files))
	e.GET("/assets/*", echo.WrapHandler(fileServer))

	e.GET("/", HomeWebHandler)

	RegisterAuthRoutes(e.Group("/auth"), cfg)

	jwtConfig := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(jwtCustomClaims)
		},
		SigningKey:   []byte(cfg.SecretKey),
		ErrorHandler: customJWTErrorHandler,
	}

	userRoutesGroup := e.Group("/user")
	userRoutesGroup.Use(echojwt.WithConfig(jwtConfig))
	RegisterUserRoutes(userRoutesGroup, cfg)

	passwordsRoutesGroup := e.Group("/passwords")
	passwordsRoutesGroup.Use(echojwt.WithConfig(jwtConfig))
	RegisterPasswordsRoutes(passwordsRoutesGroup)

	return e
}

func HomeWebHandler(c echo.Context) error {
	claims, err := getUserFromToken(c, cfg)
	if err == nil && claims != nil {
		log.Infof("Home page requested by user: %s\n", claims.Username)
		return c.Redirect(http.StatusFound, "/passwords")
	}

	component := components.HomePage()

	return Render(c, http.StatusOK, component)
}

func Render(ctx echo.Context, statusCode int, t templ.Component) error {
	buf := templ.GetBuffer()
	defer templ.ReleaseBuffer(buf)

	if err := t.Render(ctx.Request().Context(), buf); err != nil {
		log.Errorf("Error rendering component: %e\n", err)
		return err
	}

	return ctx.HTML(statusCode, buf.String())
}
