package routes

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func (r *Routes) RegisterCommonRoutes(g *echo.Group) {
	g.POST("/hello/", echo.WrapHandler(http.HandlerFunc(HelloWebHandler)))
}
