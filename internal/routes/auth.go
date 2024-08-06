package routes

import (
	"errors"
	"github.com/a-h/templ"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
	"time"
	"ypeskov/go-password-manager/cmd/web/components/auth"
	"ypeskov/go-password-manager/internal/config"
	routeErrors "ypeskov/go-password-manager/internal/routes/errors"
	"ypeskov/go-password-manager/models"
)

type jwtCustomClaims struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
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
	email := c.FormValue("email")
	password := c.FormValue("password")
	log.Infof("Login attempt: %s\n", email)

	user, err := sManager.UsersService.GetUserByEmail(email)
	if err != nil {
		log.Errorf("Error getting user by email: %s\n", err)
		return echo.ErrInternalServerError
	}

	if user == nil {
		//TODO: Add component when user not found
		return echo.ErrUnauthorized
	}

	if !comparePassword(user.HashPassword, password) {
		log.Errorf("Invalid password for user: %s\n", email)
		return echo.ErrUnauthorized
	}

	log.Infof("User logged in: %+v\n", user)
	claims := &jwtCustomClaims{
		Id: user.Id,
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
		MaxAge:   -1, // Delete cookie
	})

	return c.Redirect(http.StatusFound, "/")
}

func (ar *AuthRoutes) RegisterForm(c echo.Context) error {
	formData, _ := c.Get("formData").(map[string]string)
	if formData == nil {
		formData = map[string]string{
			"email":            "",
			"password":         "",
			"confirm_password": "",
		}
	}

	errorResponse, ok := c.Get("error").(*routeErrors.UserError)
	var component templ.Component
	if ok && errorResponse != nil {
		component = auth.RegisterForm(errorResponse, formData)
	} else {
		component = auth.RegisterForm(nil, formData)
	}

	return Render(c, http.StatusOK, component)
}

func (ar *AuthRoutes) Register(c echo.Context) error {
	email := c.FormValue("email")
	password := c.FormValue("password")
	confirmPassword := c.FormValue("confirm_password")
	log.Infof("Registration attempt, username: [%s]\n", email)

	formData := map[string]string{
		"email":            email,
		"password":         password,
		"confirm_password": confirmPassword,
	}
	c.Set("formData", formData)

	if email == "" {
		log.Printf("Invalid registration attempt: missing email")
		c.Set("error", routeErrors.MissingEmail)
		return ar.RegisterForm(c)
	}
	if password == "" {
		log.Printf("Invalid registration attempt: missing password")
		c.Set("error", routeErrors.MissingPassword)
		return ar.RegisterForm(c)
	}
	if confirmPassword == "" {
		log.Printf("Invalid registration attempt: missing confirm password")
		c.Set("error", routeErrors.MissingConfirmPassword)
		return ar.RegisterForm(c)
	}

	if password != confirmPassword {
		log.Printf("Invalid registration attempt: passwords do not match")
		c.Set("error", routeErrors.UserPasswordsDoNotMatch)
		return ar.RegisterForm(c)
	}

	userExists, err := sManager.UsersService.GetUserByEmail(email)
	if err != nil {
		log.Errorf("Error getting user by email: %s\n", err)
		return echo.ErrInternalServerError
	}
	if userExists != nil {
		log.Errorf("User already exists: %s\n", email)
		c.Set("error", routeErrors.UserExists)
		return ar.RegisterForm(c)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Errorf("Error hashing password: %s\n", err)
		return echo.ErrInternalServerError
	}
	user := &models.User{
		Email:        email,
		HashPassword: string(hashedPassword),
	}

	err = user.Validate()
	if err != nil {
		log.Errorf("Error validating user: %s\n", err)
		return echo.ErrBadRequest
	}

	err = sManager.UsersService.CreateUser(user)
	if err != nil {
		log.Errorf("Error saving user to the database: %s\n", err)
		return echo.ErrInternalServerError
	}

	return c.Redirect(http.StatusFound, "/auth/login")
}

func comparePassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
