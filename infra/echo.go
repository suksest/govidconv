package infra

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

//NewServer ...
func NewServer() (e *echo.Echo) {
	e = echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	return
}
