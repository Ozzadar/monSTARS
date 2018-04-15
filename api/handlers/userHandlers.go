package handlers

import (
	"net/http"

	"github.com/labstack/echo"
)

//AppMain serves main logged-in user
func AppMain(c echo.Context) error {
	return c.String(http.StatusInternalServerError, "Main App")
}
