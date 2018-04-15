/*
 * Created on Sun Apr 15 2018
 *
 * Copyright (c) 2018 Ozzadar.com
 * Licensed under the GNU General Public License v3.0
 */

package middlewares

import (
	"errors"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/ozzadar/monSTARS/config"
	"github.com/ozzadar/monSTARS/db"
)

//SetJWTMiddlewares add all JWT authentication middleware
func SetJWTMiddlewares(g *echo.Group) {
	if signingKey, err := config.Config.GetString("default", "jwt_secret_key"); err == nil {
		//Checks for valid JWT
		g.Use(middleware.JWTWithConfig(middleware.JWTConfig{
			SigningKey:    []byte(signingKey),
			SigningMethod: "HS512",
		}))
		//Ensures JWT is active
		g.Use(CheckJWTActive)

		//Extracts user and adds it to context for access up the stack
		g.Use(AddUserToContext)
	} else {
		panic(errors.New("jwt_secret_key defined in config"))
	}
}

//AddUserToContext :Extract user and add to context
func AddUserToContext(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if user, ok := c.Get("user").(*jwt.Token); ok {
			if claim, ok2 := user.Claims.(jwt.MapClaims); ok2 {

				if claim["username"] != "" {
					user := db.GetUserPreAuth(claim["username"].(string))

					c.Set("token_user", user)
					next(c)
					return nil
				}
			}
		}
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"message": "",
		})
	}
}

//CheckJWTActive :Extract user and add to context
func CheckJWTActive(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if user, ok := c.Get("user").(*jwt.Token); ok {
			rawToken := user.Raw

			if db.JWTExists(rawToken) {
				next(c)
				return nil
			}

			return c.JSON(http.StatusUnauthorized, map[string]string{
				"message": "JWT expired or invalid; please login again",
			})
		}
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"message": "",
		})
	}
}
