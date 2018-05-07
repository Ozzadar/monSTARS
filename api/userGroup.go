/*
 * Created on Sun Apr 15 2018
 *
 * Copyright (c) 2018 Ozzadar.com
 * Licensed under the GNU General Public License v3.0
 */

package api

import (
	"github.com/labstack/echo"
	"github.com/ozzadar/monSTARS/api/handlers"
)

//UserGroup sets routes for admin group
func UserGroup(g *echo.Group) {
	g.GET("/", handlers.AppMain)
	g.GET("", handlers.AppMain)
	g.POST("/donate", handlers.Donate)
}
