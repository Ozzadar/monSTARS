/*
 * Created on Sun Apr 15 2018
 *
 * Copyright (c) 2018 Ozzadar.com
 * Licensed under the GNU General Public License v3.0
 */

package handlers

import (
	"net/http"

	"github.com/labstack/echo"
)

//HelloWorld says hello
func HelloWorld(c echo.Context) error {
	return c.String(http.StatusOK, "Yallo from the web side!")
}
