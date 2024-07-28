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
	"ypeskov/go-password-manager/internal/logger"
	"ypeskov/go-password-manager/services"
)

type Routes struct {
	logger          *logger.Logger
	Echo            *echo.Echo
	ServicesManager *services.ServiceManager
}

var log *logger.Logger
var sManager *services.ServiceManager

func RegisterRoutes(logger *logger.Logger, servicesManager *services.ServiceManager) *echo.Echo {
	log = logger
	log.Info("Registering routes")
	e := echo.New()

	sManager = servicesManager

	//e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	fileServer := http.FileServer(http.FS(web.Files))
	e.GET("/assets/*", echo.WrapHandler(fileServer))

	e.GET("/", HomeWebHandler)

	RegisterAuthRoutes(e.Group("/auth"))

	jwtConfig := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(jwtCustomClaims)
		},
		SigningKey:   []byte("secret"),
		ErrorHandler: customJWTErrorHandler,
	}

	passwordsRoutesGroup := e.Group("/passwords")
	passwordsRoutesGroup.Use(echojwt.WithConfig(jwtConfig))
	RegisterPasswordsRoutes(passwordsRoutesGroup)

	return e
}

func HomeWebHandler(c echo.Context) error {
	passwords, err := sManager.PasswordService.GetAllPasswords()
	if err != nil {
		log.Errorf("Error getting all passwords: %e\n", err)
		return err
	}

	component := components.ListOfPasswords(passwords)

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
