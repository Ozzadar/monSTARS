/*
 * Created on Sun Apr 15 2018
 *
 * Copyright (c) 2018 Ozzadar.com
 * Licensed under the GNU General Public License v3.0
 */

package config

import (
	"github.com/dlintw/goconf"
)

var Config *goconf.ConfigFile

func LoadConfig() error {

	err := error(nil)
	Config, err = goconf.ReadConfigFile("config/app.conf")

	return err

}
