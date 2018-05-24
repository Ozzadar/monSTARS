/*
 * Created on Sun Apr 15 2018
 *
 * Copyright (c) 2018 Ozzadar.com
 * Licensed under the GNU General Public License v3.0
 */

package authservice

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/klouds/kDaemon/logging"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/ozzadar/monSTARS/config"
	"github.com/ozzadar/monSTARS/db"
	"github.com/ozzadar/monSTARS/models"
)

//JWTClaims type
type JWTClaims struct {
	models.User
	jwt.StandardClaims
}

var (
	jobRunning = false
)

/*JWTExpiryService will trigger CheckTokensForExpiry every minute
 */
func JWTExpiryService() {
	CheckTokensForExpiry()
	nextTime := time.Now().Truncate(time.Minute)
	nextTime = nextTime.Add(time.Minute)
	time.Sleep(time.Until(nextTime))
	go JWTExpiryService()
}

/*CheckTokensForExpiry checks active tokens for expiry; invalidates any found
 */
func CheckTokensForExpiry() {
	if jobRunning {
		return
	}

	currentTime := time.Now()

	logging.Log("----------")
	logging.Log("Checking for expired tokens")
	logging.Log("Time : " + currentTime.Format("Jan 2, 2006 at 15:04:05 (MST)"))

	tokens := db.GetAllActiveJWTs()

	for _, token := range tokens {
		if parsed, _ := jwt.Parse(token.Token, nil); parsed != nil {
			if claims, ok := parsed.Claims.(jwt.MapClaims); ok {
				expiresAt := int64(claims["exp"].(float64))

				if time.Now().Unix() > expiresAt {
					owner := token.Owner

					if db.DeleteJWT(&token) {
						logging.Log(fmt.Sprintf("Deleting user \"%s\"'s token. User will need to login again.", owner))
					}

				}
			}
		}
	}
	logging.Log("----------")
}

func LoginWithUserPass(username string, password string) string {
	//Check if valid login
	user := db.GetUser(username, password)

	if user != nil {
		//Create JWT Token
		token, err := CreateJwtToken(user)

		if err != nil {
			log.Printf("Failed to create token: %#v", err)
			return ""
		}

		user.LoggedIn = true

		db.UpdateUserLoginState(user)
		db.AddJWT(&models.JwtToken{
			Owner: user.Username,
			Token: token,
		})

		return token
	}
	return ""
}

//CreateJwtToken create a token
func CreateJwtToken(user *models.User) (string, error) {
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

func VerifyJWT(token string) (bool, *models.JwtToken) {
	theToken := db.GetJWT(token)
	return theToken != nil, theToken
}
