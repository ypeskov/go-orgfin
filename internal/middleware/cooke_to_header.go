package middleware

import "github.com/labstack/echo/v4"

func CookieToHeaderMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("auth_token")
		if err == nil {
			c.Request().Header.Set("Authorization", "Bearer "+cookie.Value)
		}
		return next(c)
	}
}
