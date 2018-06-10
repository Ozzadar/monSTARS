/*
 * Created on Sun Apr 15 2018
 *
 * Copyright (c) 2018 Ozzadar.com
 * Licensed under the GNU General Public License v3.0
 */

package models

//TODO: FINISH THIS

//Character type
type Character struct {
	Name     string   `json:"name" gorethink:"name"`
	SpriteID string   `json:"spriteid" gorethink:"spriteid"`
	Gender   string   `json:"gender" gorethink:"gender"`
	MapID    string   `json:"mapid" gorethink:"mapid"`
	Location Position `json:"position" gorethink:"position"`
	OwnerID  string   `json:"ownerid" gorethink:"ownerid"`
}
