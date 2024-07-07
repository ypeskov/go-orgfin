package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"ypeskov/go-orgfin/cmd/web"
	"ypeskov/go-orgfin/internal/logger"
)

type Routes struct {
	logger *logger.Logger
	Echo   *echo.Echo
}

var routesInstance *Routes

func RegisterRoutes(logger *logger.Logger) *Routes {
	logger.Info("Registering routes")
	e := echo.New()

	routesInstance = &Routes{
		logger: logger,
		Echo:   e,
	}

	//e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	fileServer := http.FileServer(http.FS(web.Files))
	e.GET("/assets/*", echo.WrapHandler(fileServer))

	commonGroup := e.Group("/common")
	routesInstance.RegisterCommonRoutes(commonGroup)

	return routesInstance
}
