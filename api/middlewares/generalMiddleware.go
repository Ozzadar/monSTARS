/*
 * Created on Sun Apr 15 2018
 *
 * Copyright (c) 2018 Ozzadar.com
 * Licensed under the GNU General Public License v3.0
 */

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
