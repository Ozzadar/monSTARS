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

//AppMain serves main logged-in user
func AppMain(c echo.Context) error {
	return c.String(http.StatusInternalServerError, "Main App")
}
