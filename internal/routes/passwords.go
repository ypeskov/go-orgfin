package routes

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func (r *Routes) RegisterPasswordsRoutes(g *echo.Group) {
	routesInstance.logger.Info("Registering passwords routes")
	g.GET("/", GetAllPasswordsHandler)
}

func GetAllPasswordsHandler(c echo.Context) error {
	routesInstance.logger.Info("Getting all passwords handler")
	passwords, _ := routesInstance.ServicesManager.PasswordService.GetAllPasswords()

	return c.JSON(http.StatusOK, passwords)
}
