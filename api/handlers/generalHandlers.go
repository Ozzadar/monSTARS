package handlers

import (
	"net/http"

	"github.com/labstack/echo"
)

//HelloWorld says hello
func HelloWorld(c echo.Context) error {
	return c.String(http.StatusOK, "Yallo from the web side!")
}
