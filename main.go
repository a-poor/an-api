package main

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	port      = "1323"
	version   = "0.0.0"
	jwtSecret = []byte("secret")
)

type User struct {
	Name  string `json:"name"`
	Admin bool   `json:"admin"`
	jwt.StandardClaims
}

func main() {
	// Create the app
	e := echo.New()

	// Add middleware
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Logger())
	e.HideBanner = true
	e.HidePort = true

	// Add routes
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"success": true,
			"version": version,
		})
	})
	e.GET("/greet", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"success": true,
			"version": version,
			"message": "Hello, you!",
		})
	})
	e.GET("/greet/:name", func(c echo.Context) error {
		name := c.Param("name")
		if name == "" {
			name = "World"
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"success": true,
			"version": version,
			"message": "Hello, " + name + "!",
		})
	})
	e.POST("/auth", func(c echo.Context) error {
		username := c.FormValue("username")
		password := c.FormValue("password")

		var user User
		if username == "admin" && password == "admin" {
			user = User{
				username,
				true,
				jwt.StandardClaims{
					ExpiresAt: time.Now().Add(time.Hour).Unix(),
				},
			}
		} else if username == "user" && password == "user" {
			user = User{
				username,
				false,
				jwt.StandardClaims{
					ExpiresAt: time.Now().Add(time.Hour).Unix(),
				},
			}
		} else {
			return echo.ErrUnauthorized
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, user)
		t, err := token.SignedString(jwtSecret)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"success": true,
			"token":   t,
		})
	})
	gu := e.Group("/users", middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: jwtSecret,
		Claims:     &User{},
	}))
	gu.GET("/all", func(c echo.Context) error {
		c.Logger().Warnf("User: %+v", c.Get("user"))
		return c.JSON(http.StatusOK, map[string]interface{}{"user": c.Get("user")})
	})

	// Start the server
	if err := e.Start(":" + port); err != nil {
		e.Logger.Fatal(err)
	}
}
