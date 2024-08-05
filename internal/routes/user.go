package routes

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"ypeskov/go-password-manager/cmd/web/components/user"
	"ypeskov/go-password-manager/internal/config"
)

type UserRoutes struct {
	cfg *config.Config
}

func RegisterUserRoutes(g *echo.Group, cfg *config.Config) {
	log.Info("Registering auth routes")

	ur := &UserRoutes{
		cfg: cfg,
	}

	g.GET("/settings", ur.Settings)
}

func (ur *UserRoutes) Settings(c echo.Context) error {
	component := user.UserSettingsPassword()

	return Render(c, http.StatusOK, component)
}
