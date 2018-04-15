package middlewares

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/ozzadar/monSTARS/models"
)

//SetAdminMiddlewares will add middleware to the main echo instance
func SetAdminMiddlewares(g *echo.Group) {
	SetJWTMiddlewares(g)

	g.Use(CheckAdminPrivileges)
}

//CheckAdminPrivileges will make sure that the user is at least an admin
func CheckAdminPrivileges(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if user, ok := c.Get("token_user").(*models.User); ok {

			if user.IsAuthenticatedForRole(models.ROLE_ADMIN) {
				next(c)
				return nil
			}

			return c.JSON(http.StatusUnauthorized, "User does not have sufficient privileges.")
		}
		return c.JSON(http.StatusInternalServerError, "User not properly set to context")

	}
}
