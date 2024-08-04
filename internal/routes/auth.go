package routes

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
	"time"
	"ypeskov/go-password-manager/cmd/web/components/auth"
	"ypeskov/go-password-manager/internal/config"
)

type jwtCustomClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type AuthRoutes struct {
	cfg *config.Config
}

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

	return c.Redirect(http.StatusFound, "/auth/login")
}

func getUserFromToken(c echo.Context, cfg *config.Config) (*jwtCustomClaims, error) {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return nil, errors.New("missing Authorization header")
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	token, err := jwt.ParseWithClaims(tokenString, &jwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.SecretKey), nil
	})
	if err != nil {
		log.Errorf("Error parsing JWT: %s\n", err)
		return nil, err
	}

	if claims, ok := token.Claims.(*jwtCustomClaims); ok && token.Valid {
		return claims, nil
	} else {
		log.Errorf("Invalid JWT: %s\n", err)
		return nil, err
	}
}

func RegisterAuthRoutes(g *echo.Group, cfg *config.Config) {
	log.Info("Registering auth routes")

	ar := AuthRoutes{
		cfg: cfg,
	}
	g.GET("/login", ar.LoginForm)
	g.POST("/login", ar.Login)
	g.GET("/register", ar.RegisterForm)
	g.POST("/register", ar.Register)
	g.GET("/logout", ar.Logout)
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
		Username: "Jon Snow",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 1)),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(ar.cfg.SecretKey))
	if err != nil {
		return err
	}

	c.SetCookie(&http.Cookie{
		Name:     "auth_token",
		Value:    t,
		HttpOnly: true,
		Secure:   false,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
	})

	return c.Redirect(http.StatusFound, "/")
}

func (ar *AuthRoutes) Logout(c echo.Context) error {
	c.SetCookie(&http.Cookie{
		Name:     "auth_token",
		Value:    "",
		HttpOnly: true,
		Secure:   false,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1, // Устанавливаем MaxAge в -1, чтобы удалить куку
	})

	return c.Redirect(http.StatusFound, "/")
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
		log.Errorf("Invalid registration attempt: missing fields: [username] or [password] or [confirm_password]")
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
