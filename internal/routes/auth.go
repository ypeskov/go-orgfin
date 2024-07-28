package routes

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
	"ypeskov/go-password-manager/cmd/web/components/auth"
)

type jwtCustomClaims struct {
	Name  string `json:"name"`
	Admin bool   `json:"admin"`
	jwt.RegisteredClaims
}

type AuthRoutes struct{}

// Custom error handler for JWT
func customJWTErrorHandler(c echo.Context, err error) error {
	var message string

	if err.Error() == "missing or malformed jwt" {
		message = "Missing or malformed JWT"
	} else if errors.Is(err, jwt.ErrTokenExpired) {
		message = "Expired JWT"
	} else if errors.Is(err, jwt.ErrTokenNotValidYet) {
		message = "Token not valid yet"
	} else if errors.Is(err, jwt.ErrTokenMalformed) {
		message = "Malformed JWT"
	} else if errors.Is(err, jwt.ErrSignatureInvalid) {
		message = "Invalid JWT signature"
	} else {
		message = "Invalid or missing JWT"
	}

	// Log the error
	log.Errorf("JWT error: %s\n", message)

	// Return an appropriate response
	return c.JSON(http.StatusUnauthorized, echo.Map{
		"error": message,
	})
}

func RegisterAuthRoutes(g *echo.Group) {
	log.Info("Registering auth routes")

	ar := AuthRoutes{}
	g.GET("/login", ar.LoginForm)
	g.POST("/login", ar.Login)
}

func (ar *AuthRoutes) LoginForm(c echo.Context) error {
	component := auth.LoginForm()

	return Render(c, http.StatusOK, component)
}

func (ar *AuthRoutes) Login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")
	log.Infof("Login attempt: %s\n", username)

	if username != "jon" || password != "qqq" {
		log.Warnf("Unauthorized login attempt: %s\n", username)
		return echo.ErrUnauthorized
	}

	claims := &jwtCustomClaims{
		"Jon Snow",
		true,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": t,
	})
}
