package routes

import (
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"ypeskov/go-orgfin/cmd/web"
	"ypeskov/go-orgfin/cmd/web/components"
	"ypeskov/go-orgfin/internal/logger"
	"ypeskov/go-orgfin/services"
)

type Routes struct {
	logger          *logger.Logger
	Echo            *echo.Echo
	ServicesManager *services.ServiceManager
}

var routesInstance *Routes
var log *logger.Logger

func RegisterRoutes(logger *logger.Logger, servicesManager *services.ServiceManager) *Routes {
	log = logger
	log.Info("Registering routes")
	e := echo.New()

	routesInstance = &Routes{
		logger:          logger,
		Echo:            e,
		ServicesManager: servicesManager,
	}

	//e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	fileServer := http.FileServer(http.FS(web.Files))
	e.GET("/assets/*", echo.WrapHandler(fileServer))

	e.GET("/", HomeWebHandler)

	commonGroup := e.Group("/common")
	routesInstance.RegisterCommonRoutes(commonGroup)

	passwordsGroup := e.Group("/passwords")
	routesInstance.RegisterPasswordsRoutes(passwordsGroup)

	return routesInstance
}

func HomeWebHandler(c echo.Context) error {
	log.Info("HomeWebHandler")

	passwords, err := routesInstance.ServicesManager.PasswordService.GetAllPasswords()
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
		return err
	}

	return ctx.HTML(statusCode, buf.String())
}
