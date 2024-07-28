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
	g.GET("/register", ar.RegisterForm)
	g.POST("/register", ar.Register)
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

func (ar *AuthRoutes) RegisterForm(c echo.Context) error {
	component := auth.RegisterForm()

	return Render(c, http.StatusOK, component)
}

func (ar *AuthRoutes) Register(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")
	confirmPassword := c.FormValue("confirm_password")
	log.Infof("Registration attempt, username: [%s]\n", username)

	if username == "" || password == "" || confirmPassword == "" {
		log.Warnf("Invalid registration attempt: missing fields: [username] or [password] or [confirm_password]")
		return echo.ErrBadRequest
	}

	if password != confirmPassword {
		log.Printf("Invalid registration attempt: passwords do not match")
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Passwords do not match",
		})
	}

	// Add logic to save the new user to the database
	// Example: user, err := saveUserToDatabase(username, password)

	// For demonstration purposes, we just return a success message
	return c.JSON(http.StatusOK, echo.Map{
		"message": "User registered successfully",
	})
}
