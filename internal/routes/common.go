package routes

import (
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"net/http"
	"ypeskov/go-orgfin/cmd/web"
)

func (r *Routes) RegisterCommonRoutes(g *echo.Group) {
	g.GET("/", HelloWorldHandler)
	g.GET("/web", echo.WrapHandler(templ.Handler(web.HelloForm())))
	g.POST("/hello/", echo.WrapHandler(http.HandlerFunc(HelloWebHandler)))
}

func HelloWorldHandler(c echo.Context) error {
	resp := map[string]string{
		"message": "Hello World",
	}

	return c.JSON(http.StatusOK, resp)
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
	component := web.HelloPost(name)
	err = component.Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		routesInstance.logger.Fatalf("Error rendering in HelloWebHandler: %e", err)
	}
}
