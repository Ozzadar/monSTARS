package handlers

import (
	"errors"
	"log"
	"net/http"
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
	username := c.FormValue("username")
	password := c.FormValue("password")

	if username != "" && password != "" {
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

//Register a new user
func Register(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")
	email := c.FormValue("email")

	errors := make(map[string]string)

	if len(username) == 0 {
		errors["username"] = "Must provide a username"
	}

	if len(password) == 0 {
		errors["password"] = "Must provide a password"
	}

	if len(email) == 0 {
		errors["email"] = "Must provide an email"
	} else {
		formatErr := checkmail.ValidateFormat(email)
		if formatErr != nil {
			errors["email"] = "Email not in correct format"
		} else {
			hostErr := checkmail.ValidateHost(email)

			if _, ok := hostErr.(checkmail.SmtpError); !ok && hostErr != nil {
				errors["email"] = "Email host cannot be resolved; please use a real email address"
			} else if ok && hostErr != nil {
				errors["email"] = "Email account that you provided does not exist; please use a real email address"
			}
		}

	}

	if len(errors) != 0 {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Incomplete form data",
			"errors":  errors,
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
