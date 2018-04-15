package api

import (
	"github.com/labstack/echo"
	"github.com/ozzadar/monSTARS/api/handlers"
)

//UserGroup sets routes for admin group
func UserGroup(g *echo.Group) {
	g.GET("/", handlers.AppMain)
	g.GET("", handlers.AppMain)

}
