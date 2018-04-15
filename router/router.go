package router

import (
	"github.com/labstack/echo"
	"github.com/ozzadar/monSTARS/api"
	"github.com/ozzadar/monSTARS/api/middlewares"
)

//New will return a new instance of Echo router
func New() *echo.Echo {
	e := echo.New()

	//Create groups
	adminGroup := e.Group("/admin")
	userGroup := e.Group("/app")

	middlewares.SetMainMiddlewares(e)
	middlewares.SetAdminMiddlewares(adminGroup)
	middlewares.SetUserMiddlewares(userGroup)

	//Set routes
	api.MainGroup(e)
	api.AdminGroup(adminGroup)
	api.UserGroup(userGroup)

	return e
}
