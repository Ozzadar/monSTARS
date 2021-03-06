/*
 * Created on Sun Apr 15 2018
 *
 * Copyright (c) 2018 Ozzadar.com
 * Licensed under the GNU General Public License v3.0
 */

package db

import (
	r "github.com/dancannon/gorethink"
	"github.com/ozzadar/monSTARS/models"
)

//GetJWT returns a jwt
func GetJWT(token string) *models.JwtToken {
	var JWT *models.JwtToken

	rows, err := r.Table(ActiveJWTDB).GetAllByIndex("token", token).Run(Session)

	if err != nil {
		return nil
	}

	err = rows.One(&JWT)

	if err != nil {
		return nil
	}

	return JWT
}

/*AddJWT adds jwt as active in DB
 */
func AddJWT(jwt *models.JwtToken) (bool, string) {
	success := false

	res, err := r.Table(ActiveJWTDB).Insert(jwt, r.InsertOpts{
		Conflict: "update",
	}).RunWrite(Session)

	if err != nil {
		if res.Inserted == 0 {
			return false, "Only one JWT per owner allowed"
		}
		return false, "Error occurred"
	}

	success = true

	return success, "User created successfully"
}

/*JWTExists returns the requested user; this function assumes that the method has already authenticated the user*/
func JWTExists(jwt string) bool {

	return GetJWT(jwt) != nil
}

/*GetAllActiveJWTs :returns all active JWTs in system
 */
func GetAllActiveJWTs() []models.JwtToken {
	res, err := r.Table(ActiveJWTDB).Run(Session)

	if err != nil {
		return nil
	}

	var tokens []models.JwtToken

	err = res.All(&tokens)

	if err != nil {
		return nil
	}

	if len(tokens) <= 0 {
		return nil
	}

	return tokens
}

/*DeleteJWT :removes JWT from database
 */
func DeleteJWT(token *models.JwtToken) bool {
	err := r.Table(ActiveJWTDB).Get(token.Owner).Delete().Exec(Session)

	if err != nil {
		return false
	}

	return true
}
