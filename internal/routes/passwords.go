package routes

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"ypeskov/go-orgfin/cmd/web/components"
)

func (r *Routes) RegisterPasswordsRoutes(g *echo.Group) {
	routesInstance.logger.Info("Registering passwords routes")
	g.GET("/:id", PasswordDetailsWebHandler)
}

func PasswordDetailsWebHandler(c echo.Context) error {
	log := routesInstance.logger
	log.Info("PasswordDetailsWebHandler")

	passwordId := c.Param("id")
	password, err := routesInstance.ServicesManager.PasswordService.GetPasswordById(passwordId)
	if err != nil {
		log.Errorf("Error getting password by id: %e\n", err)
		return err
	}

	component := components.PasswordDetails(*password)

	return Render(c, http.StatusOK, component)
}
