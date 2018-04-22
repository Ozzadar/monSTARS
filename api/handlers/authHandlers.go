/*
 * Created on Sun Apr 15 2018
 *
 * Copyright (c) 2018 Ozzadar.com
 * Licensed under the GNU General Public License v3.0
 */

package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/badoux/checkmail"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/ozzadar/monSTARS/config"
	"github.com/ozzadar/monSTARS/db"
	"github.com/ozzadar/monSTARS/models"
)

//JWTClaims type
type JWTClaims struct {
	models.User
	jwt.StandardClaims
}

//Login logs user in
func Login(c echo.Context) error {
	jsonMap := make(map[string]interface{})

	err := json.NewDecoder(c.Request().Body).Decode(&jsonMap)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "No credentials provided.",
		})
	}

	username, okUsername := jsonMap["username"].(string)
	password, okPassword := jsonMap["password"].(string)

	if okUsername && okPassword && username != "" && password != "" {
		//Check if valid login
		user := db.GetUser(username, password)

		if user != nil {
			//Create JWT Token
			token, err := CreateJwtToken(c, user)

			if err != nil {
				log.Printf("Failed to create token: %#v", err)
				return c.String(http.StatusInternalServerError, "something")
			}

			user.LoggedIn = true
			go func() {
				db.UpdateUserLoginState(user)
				db.AddJWT(&models.JwtToken{
					Owner: user.Username,
					Token: token,
				})
			}()

			return c.JSON(http.StatusOK, map[string]string{
				"message": "Logged in!",
				"token":   token,
			})
		}
	}

	return c.JSON(http.StatusBadRequest, map[string]string{
		"message": "Invalid username/password",
	})
}

//VerifyJWT checks to see if JWT is valid and active
func VerifyJWT(c echo.Context) error {
	jsonMap := make(map[string]interface{})

	err := json.NewDecoder(c.Request().Body).Decode(&jsonMap)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "No JWT provided.",
		})
	}

	token, okJWT := jsonMap["token"].(string)

	if okJWT && token != "" {
		//Check if valid login
		isValid := db.JWTExists(token)

		if isValid {
			return c.JSON(http.StatusOK, map[string]string{
				"message": "Valid and Active JWT",
				"token":   token,
			})
		}
	}

	return c.JSON(http.StatusUnauthorized, map[string]string{
		"message": "JWT expired or inactive",
	})
}

//Register a new user
func Register(c echo.Context) error {
	jsonMap := make(map[string]interface{})

	err := json.NewDecoder(c.Request().Body).Decode(&jsonMap)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "No credentials provided.",
		})
	}

	username, okUsername := jsonMap["username"].(string)
	password, okPassword := jsonMap["password"].(string)
	email, okEmail := jsonMap["email"].(string)

	errors := make(map[string]string)

	if okUsername && len(username) == 0 {
		errors["username"] = "Must provide a username"
	}

	if okPassword && len(password) == 0 {
		errors["password"] = "Must provide a password"
	}

	if okEmail && len(email) == 0 {
		errors["email"] = "Must provide an email"
	} else {
		formatErr := checkmail.ValidateFormat(email)
		if formatErr != nil {
			errors["email"] = "Email not in correct format"
		}
	}

	if len(errors) != 0 {

		messages := ""

		for _, message := range errors {
			messages += message + ", "
		}

		messages = strings.TrimSuffix(messages, ", ")

		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Incomplete/Invalid form data: " + messages,
		})
	}

	hashedPwd := db.HashAndSalt(password)
	newUser := models.NewUser(username, hashedPwd, email)

	success, message := db.RegisterUser(&newUser)

	if success {
		return c.JSON(http.StatusOK, map[string]string{
			"message": message,
		})
	}

	return c.JSON(http.StatusConflict, map[string]string{
		"message": message,
	})

}

//CreateJwtToken create a token
func CreateJwtToken(c echo.Context, user *models.User) (string, error) {
	claims := JWTClaims{
		User: *user,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			Id:        "main_user_id",
			ExpiresAt: time.Now().Add(time.Hour * 24 * 7).Unix(),
		},
	}

	rawToken := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	if signingKey, err := config.Config.GetString("default", "jwt_secret_key"); err == nil {
		token, err := rawToken.SignedString([]byte(signingKey))

		if err != nil {
			log.Printf("Failed to create token")
			return "", err
		}

		return token, nil
	}

	panic(errors.New("jwt_secret_key defined in config"))
}
