package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	port    = "1323"
	version = "0.1.0"
)

func main() {
	// Create the app
	e := echo.New()

	// Add middleware
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Logger())

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

	// Start the server
	if err := e.Start(":" + port); err != nil {
		e.Logger.Fatal(err)
	}
}
