package routes

import (
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"ypeskov/go-orgfin/cmd/web"
	"ypeskov/go-orgfin/cmd/web/components"
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

	e.GET("/", echo.WrapHandler(templ.Handler(components.HelloForm())))

	commonGroup := e.Group("/common")
	routesInstance.RegisterCommonRoutes(commonGroup)

	passwordsGroup := e.Group("/passwords")
	routesInstance.RegisterPasswordsRoutes(passwordsGroup)

	return routesInstance
}

func HelloWebHandler(w http.ResponseWriter, r *http.Request) {
	log := routesInstance.logger
	log.Info("HelloWebHandler")
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}

	name := r.FormValue("name")
	log.Infof("Name: %s", name)
	component := components.HelloPost(name)
	err = component.Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		routesInstance.logger.Fatalf("Error rendering in HelloWebHandler: %e", err)
	}
}
