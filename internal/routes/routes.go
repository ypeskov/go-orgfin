package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"ypeskov/go-orgfin/cmd/web"
	"ypeskov/go-orgfin/internal/logger"
	"ypeskov/go-orgfin/services"
)

type Routes struct {
	logger          *logger.Logger
	Echo            *echo.Echo
	ServicesManager *services.ServiceManager
}

var routesInstance *Routes

func RegisterRoutes(logger *logger.Logger, servicesManager *services.ServiceManager) *Routes {
	logger.Info("Registering routes")
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

	commonGroup := e.Group("/common")
	routesInstance.RegisterCommonRoutes(commonGroup)

	passwordsGroup := e.Group("/passwords")
	routesInstance.RegisterPasswordsRoutes(passwordsGroup)

	return routesInstance
}
