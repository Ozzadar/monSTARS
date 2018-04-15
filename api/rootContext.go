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

//MainGroup creates routes for main group
func MainGroup(e *echo.Echo) {

	e.POST("/login", handlers.Login)
	e.POST("/register", handlers.Register)
	e.GET("/helloworld", handlers.HelloWorld)
}
