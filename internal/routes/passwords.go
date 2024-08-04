package routes

import (
	"github.com/golang-jwt/jwt/v5"
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
	userJWT := c.Get("user").(*jwt.Token)
	claims, ok := userJWT.Claims.(*jwtCustomClaims)
	if !ok {
		log.Errorf("Invalid claim type: %+v\n", userJWT.Claims)
		return c.JSON(http.StatusUnauthorized, "invalid claim type")
	}

	passwords, err := sManager.PasswordService.GetAllPasswords(claims.Id)
	if err != nil {
		log.Errorf("Error getting all passwords: %e\n", err)
		return err
	}

	component := components.ListOfPasswords(passwords)

	return Render(c, http.StatusOK, component)
}

func (pr *PasswordsRoutes) PasswordDetailsWebHandler(c echo.Context) error {
	log.Infof("EncryptedPassword Record details page requested, id: [%s]\n", c.Param("id"))
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
	newPassword := models.EncryptedPassword{}

	component := components.PasswordForm(newPassword)

	return Render(c, http.StatusOK, component)
}

func (pr *PasswordsRoutes) AddPassword(c echo.Context) error {
	password := models.EncryptedPassword{}
	if err := c.Bind(&password); err != nil {
		log.Errorf("Error binding password: %e\n", err)
		return err
	}
	log.Infof("Adding new password: %+v\n", password)

	err := sManager.PasswordService.AddPassword(&password)
	if err != nil {
		log.Errorf("Error adding password: %e\n", err)
		return err
	}

	return c.Redirect(http.StatusFound, "/")
}

func (pr *PasswordsRoutes) EditPasswordWebHandler(c echo.Context) error {
	passwordId := c.Param("id")
	log.Infof("Edit password page requested, id: [%s]\n", passwordId)
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
	log.Infof("Updating password\n")
	password := models.EncryptedPassword{}
	if err := c.Bind(&password); err != nil {
		log.Errorf("Error binding password: %e\n", err)
		return err
	}
	log.Infof("Updating password with id: %+v\n", password)
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
	log.Infof("EncryptedPassword with id %s was deleted", passwordId)

	c.Response().Header().Set("HX-Location", "/")
	return c.NoContent(http.StatusOK)
}
