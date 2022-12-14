package main

import (
  "github.com/labstack/echo/v4"
  "github.com/labstack/echo/v4/middleware"
  "net/http"
)

const (
	frontendContentsPath        = "../public"
)

func main() {
  // Echo instance
  e := echo.New()

  // Middleware
  e.Use(middleware.Logger())
  e.Use(middleware.Recover())

  // Routes
  e.GET("/", getIndex)
  e.GET("/home", getIndex)
  e.GET("/test", hello)
  e.Static("/assets", frontendContentsPath+"/assets")

  // Start server
  e.Logger.Fatal(e.Start(":1323"))
}

// Handler
func hello(c echo.Context) error {
  return c.String(http.StatusOK, "Hello, World!")
}

func getIndex(c echo.Context) error {
	return c.File(frontendContentsPath + "/index.html")
}
