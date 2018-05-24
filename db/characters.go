package db

import (
	r "github.com/dancannon/gorethink"
	"github.com/ozzadar/monSTARS/models"
)

func GetAllCharactersForUser(username string) []*models.Character {
	var characters []*models.Character

	res, err := r.Table(CharactersDB).GetAllByIndex("ownerid", username).Run(Session)

	if err != nil {
		return nil
	}

	err = res.All(&characters)

	if err != nil {
		return nil
	}

	return characters
}

func CreateNewCharacter(character *models.Character) bool {

	_, err := r.Table(CharactersDB).Insert(character).RunWrite(Session)

	if err != nil {

		return false
	}

	return true
}
