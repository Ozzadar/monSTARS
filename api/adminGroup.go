package api

import (
	"github.com/labstack/echo"
	"github.com/ozzadar/monSTARS/api/handlers"
)

//AdminGroup sets routes for admin group
func AdminGroup(g *echo.Group) {
	g.GET("/", handlers.Admin)
	g.GET("", handlers.Admin)
}
