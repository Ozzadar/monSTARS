/*
 * Created on Sun Apr 15 2018
 *
 * Copyright (c) 2018 Ozzadar.com
 * Licensed under the GNU General Public License v3.0
 */

package db

import (
	r "github.com/dancannon/gorethink"
	"github.com/klouds/kDaemon/logging"
	"github.com/ozzadar/monSTARS/config"
)

var (
	//Session singleton variable for accessing the db
	Session *r.Session
	//UserDB the name of the user table
	UserDB = "users"
	//ActiveJWTDB : the name of the active jwt table
	ActiveJWTDB = "active_tokens"
	//TransactionsDB : name of the transactions table
	TransactionsDB = "transactions"
	//CharactersDB : name of the characters table
	CharactersDB = "characters"
)

//Init supporting database functions
func Init() {
	InitDB()

}

//InitDB : initializes connection to the rethink database
func InitDB() {

	logging.Log("Initializing Database connection.")

	rethinkdbhost, err := config.Config.GetString("default", "rethinkdb_host")
	if err != nil {
		logging.Log("Problem with config file! (rethinkdb_host)")
	}

	rethinkdbport, err := config.Config.GetString("default", "rethinkdb_port")
	if err != nil {
		logging.Log("Problem with config file! (rethinkdb_port)")
	}

	rethinkdbname, err := config.Config.GetString("default", "rethinkdb_dbname")

	if err != nil {
		logging.Log("Problem with config file! (rethinkdb_dbname)")
	}

	session, err := r.Connect(r.ConnectOpts{
		Address: rethinkdbhost + ":" + rethinkdbport,
	})

	if err != nil {

	}

	session, err = r.Connect(r.ConnectOpts{
		Address: rethinkdbhost + ":" + rethinkdbport,
	})

	if err != nil {
		logging.Log("rethinkdb not found at given address: ", rethinkdbhost, ":", rethinkdbport)
		panic(true)
	}

	_, err = r.DBCreate(rethinkdbname).RunWrite(session)

	if err != nil {
		logging.Log("Unable to create DB, probably already exists.")

	}

	//Create user table
	_, err = r.DB(rethinkdbname).TableCreate(UserDB, r.TableCreateOpts{
		PrimaryKey: "username",
	}).RunWrite(session)

	if err != nil {
		logging.Log("Failed in creating users table")
	}

	//Create active jwt table
	_, err = r.DB(rethinkdbname).TableCreate(ActiveJWTDB, r.TableCreateOpts{
		PrimaryKey: "owner",
	}).RunWrite(session)

	if err != nil {
		logging.Log("Failed in creating Active JWT table")
	}

	//Create transactions table
	_, err = r.DB(rethinkdbname).TableCreate(TransactionsDB).RunWrite(session)

	if err != nil {
		logging.Log("Failed in creating Transactions table")
	}

	//Create active jwt table
	_, err = r.DB(rethinkdbname).TableCreate(CharactersDB, r.TableCreateOpts{
		PrimaryKey: "name",
	}).RunWrite(session)

	if err != nil {
		logging.Log("Failed in creating Characters table")
	}

	_, err = r.DB(rethinkdbname).Table(CharactersDB).IndexCreate("ownerid").Run(session)
	if err != nil {
		logging.Log(err)
	}

	session, err = r.Connect(r.ConnectOpts{
		Address:  rethinkdbhost + ":" + rethinkdbport,
		Database: rethinkdbname,
	})

	Session = session

	// character := &models.Character{
	// 	Name:    "Character",
	// 	OwnerID: "pmauviel",
	// 	MapID:   "1",
	// 	Location: models.Position{
	// 		X: 0,
	// 		Y: 0,
	// 	},
	// 	SpriteID: "male_123",
	// }

	// for i := 0; i < 10; i++ {
	// 	character.Name += strconv.Itoa(i)

	// 	CreateNewCharacter(character)
	// }
}
