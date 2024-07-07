package routes

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (r *Routes) RegisterPasswordsRoutes(g *echo.Group) {
	routesInstance.logger.Info("Registering passwords routes")
	g.GET("/", GetAllPasswordsHandler)
}

func GetAllPasswordsHandler(c echo.Context) error {
	passwords, _ := routesInstance.ServicesManager.PasswordService.GetAllPasswords()
	fmt.Println(passwords)

	return c.JSON(http.StatusOK, passwords)
}
