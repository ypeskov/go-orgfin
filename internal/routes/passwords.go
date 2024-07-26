package routes

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"ypeskov/go-orgfin/cmd/web/components"
	"ypeskov/go-orgfin/models"
)

func RegisterPasswordsRoutes(g *echo.Group) {
	log.Info("Registering passwords routes")
	g.GET("/new", NewPasswordWebHandler)
	g.GET("/:id", PasswordDetailsWebHandler)
	g.POST("", AddPassword)
	g.POST("/:id", UpdatePassword)
	g.GET("/:id/edit", EditPasswordWebHandler)
	g.DELETE("/:id/delete", DeletePassword)
}

func PasswordDetailsWebHandler(c echo.Context) error {
	passwordId := c.Param("id")
	password, err := sManager.PasswordService.GetPasswordById(passwordId)
	if err != nil {
		log.Errorf("Error getting password by id: %e\n", err)
		return err
	}

	component := components.PasswordDetails(*password)

	return Render(c, http.StatusOK, component)
}

func NewPasswordWebHandler(c echo.Context) error {
	newPassword := models.Password{}

	component := components.PasswordForm(newPassword)

	return Render(c, http.StatusOK, component)
}

func AddPassword(c echo.Context) error {
	password := models.Password{}
	if err := c.Bind(&password); err != nil {
		log.Errorf("Error binding password: %e\n", err)
		return err
	}

	err := sManager.PasswordService.AddPassword(&password)
	if err != nil {
		log.Errorf("Error adding password: %e\n", err)
		return err
	}

	return c.Redirect(http.StatusFound, "/")
}

func EditPasswordWebHandler(c echo.Context) error {
	passwordId := c.Param("id")
	password, err := sManager.PasswordService.GetPasswordById(passwordId)
	if err != nil {
		log.Errorf("Error getting password by id: %e\n", err)
		return err
	}

	component := components.EditPassword(*password)

	return Render(c, http.StatusOK, component)
}

func UpdatePassword(c echo.Context) error {
	password := models.Password{}
	if err := c.Bind(&password); err != nil {
		log.Errorf("Error binding password: %e\n", err)
		return err
	}

	err := sManager.PasswordService.UpdatePassword(&password)
	if err != nil {
		log.Errorf("Error updating password: %e\n", err)
		return err
	}

	return c.Redirect(http.StatusFound, "/")
}

func DeletePassword(c echo.Context) error {
	passwordId := c.Param("id")
	err := sManager.PasswordService.DeletePassword(passwordId)
	if err != nil {
		log.Errorf("Error deleting password: %e\n", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error deleting password"})
	}
	log.Infof("Password with id %s was deleted", passwordId)

	c.Response().Header().Set("HX-Location", "/")
	return c.NoContent(http.StatusOK)
}
