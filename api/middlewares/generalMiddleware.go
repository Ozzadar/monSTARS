package middlewares

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

//SetMainMiddlewares will add middleware to the main echo instance
func SetMainMiddlewares(e *echo.Echo) {
	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root: "./static",
	}))
}
