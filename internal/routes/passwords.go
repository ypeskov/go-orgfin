package routes

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"ypeskov/go-password-manager/cmd/web/components"
	"ypeskov/go-password-manager/models"
)

type PasswordsRoutes struct{}

func RegisterPasswordsRoutes(g *echo.Group) {
	log.Info("Registering passwords routes")

	pr := PasswordsRoutes{}

	g.GET("/new", pr.NewPasswordWebHandler)
	g.GET("/:id", pr.PasswordDetailsWebHandler)
	g.GET("", pr.PasswordsListWeb)
	g.POST("", pr.AddPassword)
	g.POST("/:id", pr.UpdatePassword)
	g.GET("/:id/edit", pr.EditPasswordWebHandler)
	g.DELETE("/:id/delete", pr.DeletePassword)
}

func (pr *PasswordsRoutes) PasswordsListWeb(c echo.Context) error {
	passwords, err := sManager.PasswordService.GetAllPasswords()
	if err != nil {
		log.Errorf("Error getting all passwords: %e\n", err)
		return err
	}

	component := components.ListOfPasswords(passwords)

	return Render(c, http.StatusOK, component)
}

func (pr *PasswordsRoutes) PasswordDetailsWebHandler(c echo.Context) error {
	log.Infof("Password details page requested, id: [%s]\n", c.Param("id"))
	passwordId := c.Param("id")
	id, err := strconv.Atoi(passwordId)
	if err != nil {
		log.Errorf("Error converting password id to int: %e\n", err)
		return err
	}

	password, err := sManager.PasswordService.GetPasswordById(id)
	if err != nil {
		log.Errorf("Error getting password by id: %e\n", err)
		return err
	}

	component := components.PasswordDetails(*password)

	return Render(c, http.StatusOK, component)
}

func (pr *PasswordsRoutes) NewPasswordWebHandler(c echo.Context) error {
	log.Infof("New password page requested\n")
	newPassword := models.Password{}

	component := components.PasswordForm(newPassword)

	return Render(c, http.StatusOK, component)
}

func (pr *PasswordsRoutes) AddPassword(c echo.Context) error {
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

func (pr *PasswordsRoutes) EditPasswordWebHandler(c echo.Context) error {
	passwordId := c.Param("id")
	id, err := strconv.Atoi(passwordId)
	if err != nil {
		log.Errorf("Error converting password id to int: %e\n", err)
		return err
	}

	password, err := sManager.PasswordService.GetPasswordById(id)
	if err != nil {
		log.Errorf("Error getting password by id: %e\n", err)
		return err
	}

	component := components.EditPassword(*password)

	return Render(c, http.StatusOK, component)
}

func (pr *PasswordsRoutes) UpdatePassword(c echo.Context) error {
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

func (pr *PasswordsRoutes) DeletePassword(c echo.Context) error {
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
