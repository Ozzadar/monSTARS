package handlers

import (
	"net/http"

	"github.com/labstack/echo"
)

//Admin serves the admin page
func Admin(c echo.Context) error {
	return c.String(http.StatusOK, " ADMIN PAGE")
}
