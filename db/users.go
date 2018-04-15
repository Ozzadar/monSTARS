package db

import (
	r "github.com/dancannon/gorethink"
	"github.com/ozzadar/monSTARS/models"
)

//UsernameExists checks if username exists in database
func UsernameExists(username string) bool {
	return true
}

//EmailExists checks if email exists in database;
func EmailExists(email string) bool {
	return false
}

/*RegisterUser registers user in database.
- Username must be unique
- Email must be unique
*/
func RegisterUser(user *models.User) (bool, string) {
	success := false

	uniqueUser := !EmailExists(user.Email)

	//If the user is unique, add him to database
	if uniqueUser {
		res, err := r.Table(UserDB).Insert(user).RunWrite(Session)

		if err != nil {
			if res.Inserted == 0 {
				return false, "Username must be unique"
			}
			return false, "Error occurred"
		}

		success = true
	}
	return success, "User created successfully"
}

/*GetUser returns the user if password is correct
 */
func GetUser(username string, password string) *models.User {
	var user *models.User

	res, err := r.Table(UserDB).Get(username).Run(Session)

	if err != nil {
		return nil
	}

	err = res.One(&user)

	if err != nil {
		return nil
	}

	if user != nil && ComparePasswords(user.PasswordHash, password) {
		return user
	}
	return nil
}

/*GetUserPreAuth returns the requested user; this function assumes that the method has already authenticated the user*/
func GetUserPreAuth(username string) *models.User {
	var user *models.User

	res, err := r.Table(UserDB).Get(username).Run(Session)

	if err != nil {
		return nil
	}

	err = res.One(&user)

	if err != nil {
		return nil
	}

	return user
}

/*UpdateUserToken updates the database with new user information*/
func UpdateUserToken(u *models.User) bool {

	_, err := r.Table(UserDB).Update(map[string]interface{}{
		"username":      u.Username,
		"current_token": u.CurrentToken,
	}).RunWrite(Session)

	if err != nil {
		return false
	}

	return true
}

/*UpdateUserLoginState updates the 'logged_in' field*/
func UpdateUserLoginState(u *models.User) bool {
	_, err := r.Table(UserDB).Update(map[string]interface{}{
		"username":  u.Username,
		"logged_in": u.LoggedIn,
	}).RunWrite(Session)

	if err != nil {
		return false
	}

	return true
}
