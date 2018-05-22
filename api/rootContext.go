/*
 * Created on Sun Apr 15 2018
 *
 * Copyright (c) 2018 Ozzadar.com
 * Licensed under the GNU General Public License v3.0
 */

package api

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/ozzadar/monSTARS/api/handlers"
)

//MainGroup creates routes for main group
func MainGroup(e *echo.Echo) {

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	}))
	e.GET("/ws", handlers.Hello)

	e.POST("/login", handlers.Login)
	e.POST("/register", handlers.Register)
	e.POST("/verify-jwt", handlers.VerifyJWT)
	e.GET("/helloworld", handlers.HelloWorld)
}
