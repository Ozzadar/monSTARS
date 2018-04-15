/*
 * Created on Sun Apr 15 2018
 *
 * Copyright (c) 2018 Ozzadar.com
 * Licensed under the GNU General Public License v3.0
 */

package models

//JwtToken holds jwt tokens
type JwtToken struct {
	Owner string `json:"owner" gorethink:"owner"`
	Token string `json:"token" gorethink:"token"`
}
