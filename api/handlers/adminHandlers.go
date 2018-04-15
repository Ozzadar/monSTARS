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

//Admin serves the admin page
func Admin(c echo.Context) error {
	return c.String(http.StatusOK, " ADMIN PAGE")
}
