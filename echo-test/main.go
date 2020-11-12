package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Person struct {
	Name   string `json:"name"`
	Height string `json:"height"`
	Mass   string `json:"mass"`
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.GET("/", func(c echo.Context) error {
		return c.HTML(http.StatusOK, "<h1>Hello World</h1>")
	})
	e.GET("/hello", func(c echo.Context) error {
		return c.HTML(http.StatusOK, "<h1>I am not so good</h1>")
	})
	e.GET("/api", func(c echo.Context) error {
		p := Person{Height: "170", Mass: "128", Name: "Antony"}
		return c.JSON(http.StatusOK, p)
	})

	e.Logger.Fatal(e.Start(":9000"))

}
