package api

import (
	"github.com/labstack/echo"
	"github.com/ozzadar/monSTARS/api/handlers"
)

//MainGroup creates routes for main group
func MainGroup(e *echo.Echo) {

	e.POST("/login", handlers.Login)
	e.POST("/register", handlers.Register)
	e.GET("/helloworld", handlers.HelloWorld)
}
