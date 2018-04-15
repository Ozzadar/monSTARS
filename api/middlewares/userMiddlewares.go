package middlewares

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/ozzadar/monSTARS/models"
)

//SetUserMiddlewares will add middleware to the main echo instance
func SetUserMiddlewares(g *echo.Group) {
	SetJWTMiddlewares(g)

	g.Use(CheckUserPrivileges)
}

//CheckUserPrivileges will make sure that the user is at least an admin
func CheckUserPrivileges(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if user, ok := c.Get("token_user").(*models.User); ok {

			if user.IsAuthenticatedForRole(models.ROLE_USER) {
				next(c)
				return nil
			}

			return c.JSON(http.StatusUnauthorized, "User does not have sufficient privileges.")
		}
		return c.JSON(http.StatusInternalServerError, "User not properly set to context")

	}
}
