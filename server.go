/*
 * Created on Sun Apr 15 2018
 *
 * Copyright (c) 2018 Ozzadar.com
 * Licensed under the GNU General Public License v3.0
 */

package main

import (
	"fmt"
	"log"

	"github.com/ozzadar/monSTARS/config"
	"github.com/ozzadar/monSTARS/db"
	"github.com/ozzadar/monSTARS/router"
	"github.com/ozzadar/monSTARS/services/jwtservice"
	"github.com/ozzadar/monSTARS/services/paypalservice"
)

func main() {
	//Load config
	err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	//Config used here
	//Initialize the database
	db.Init()
	//Initialize paypal
	paypalservice.Init()

	ip, err := config.Config.GetString("default", "bind_ip")
	if err != nil {
		fmt.Println("bind_ip not defined in config; exiting.")
		return
	}
	port, err := config.Config.GetString("default", "bind_port")
	if err != nil {
		fmt.Println("bind_port not defined in config; exiting.")
		return
	}

	e := router.New()

	go jwtservice.JWTExpiryService()
	e.Start(ip + ":" + port)
}
