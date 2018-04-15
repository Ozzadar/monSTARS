package jwtservice

import (
	"fmt"
	"time"

	"github.com/klouds/kDaemon/logging"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/ozzadar/monSTARS/db"
)

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
