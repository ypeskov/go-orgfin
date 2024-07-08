package routes

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"ypeskov/go-orgfin/cmd/web/components"
	"ypeskov/go-orgfin/models"
)

func (r *Routes) RegisterPasswordsRoutes(g *echo.Group) {
	log.Info("Registering passwords routes")
	g.GET("/new", NewPasswordWebHandler)
	g.GET("/:id", PasswordDetailsWebHandler)
	g.POST("", AddPassword)
	g.POST("/:id", UpdatePassword)
	g.GET("/:id/edit", EditPasswordWebHandler)
}

func PasswordDetailsWebHandler(c echo.Context) error {
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

func NewPasswordWebHandler(c echo.Context) error {
	log.Info("NewPasswordWebHandler")

	newPassword := models.Password{}

	component := components.EditPassword(newPassword)

	return Render(c, http.StatusOK, component)
}

func AddPassword(c echo.Context) error {
	log.Info("AddPassword")

	password := models.Password{}
	if err := c.Bind(&password); err != nil {
		log.Errorf("Error binding password: %e\n", err)
		return err
	}

	err := routesInstance.ServicesManager.PasswordService.AddPassword(&password)
	if err != nil {
		log.Errorf("Error adding password: %e\n", err)
		return err
	}

	return c.Redirect(http.StatusFound, "/")
}

func EditPasswordWebHandler(c echo.Context) error {
	log.Info("EditPasswordWebHandler")

	passwordId := c.Param("id")
	password, err := routesInstance.ServicesManager.PasswordService.GetPasswordById(passwordId)
	if err != nil {
		log.Errorf("Error getting password by id: %e\n", err)
		return err
	}

	component := components.EditPassword(*password)

	return Render(c, http.StatusOK, component)
}

func UpdatePassword(c echo.Context) error {
	log.Info("UpdatePassword")

	password := models.Password{}
	if err := c.Bind(&password); err != nil {
		log.Errorf("Error binding password: %e\n", err)
		return err
	}

	err := routesInstance.ServicesManager.PasswordService.UpdatePassword(&password)
	if err != nil {
		log.Errorf("Error updating password: %e\n", err)
		return err
	}

	return c.Redirect(http.StatusFound, "/")
}
