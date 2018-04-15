/*
 * Created on Sun Apr 15 2018
 *
 * Copyright (c) 2018 Ozzadar.com
 * Licensed under the GNU General Public License v3.0
 */

package db

import (
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

//HashAndSalt a given plain-text password
func HashAndSalt(pwd string) string {

	fmt.Println("Starting Hash")
	bytePwd := []byte(pwd)
	// Use GenerateFromPassword to hash & salt pwd.
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	hash, err := bcrypt.GenerateFromPassword(bytePwd, bcrypt.DefaultCost+5)
	if err != nil {
		log.Println(err)
	}
	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	fmt.Println("Ending Hash")
	return string(hash)
}

//ComparePasswords compares hashed string from db with plain text password
func ComparePasswords(hashedPwd string, plainPwd string) bool {
	// Since we'll be getting the hashed password from the DB it
	// will be a string so we'll need to convert it to a byte slice
	byteHash := []byte(hashedPwd)
	bytePwd := []byte(plainPwd)

	err := bcrypt.CompareHashAndPassword(byteHash, bytePwd)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}
