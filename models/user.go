package models

type Role int

const (
	ROLE_GUEST Role = iota
	ROLE_USER
	ROLE_PREMIUM
	ROLE_ADMIN
)

//User struct for holding user data
type User struct {
	Username      string      `json:"username" gorethink:"username"`
	PasswordHash  string      `json:"-" gorethink:"password_hash"`
	Email         string      `json:"email" gorethink:"email"`
	Characters    []Character `json:"characters" gorethink:"characters"`
	LoggedIn      bool        `json:"logged_in" gorethink:"logged_in"`
	Role          Role        `json:"role" gorethink:"role"`
	CurrentToken  string      `json:"-" gorethink:"current_token"`
	CurrentClient int         `json:"-" gorethink:"-"`
}

//NewUser returns a new user with default settings
func NewUser(username string, passwordhash string, email string) User {
	newUser := User{
		Username:      username,
		PasswordHash:  passwordhash,
		Email:         email,
		Characters:    []Character{},
		LoggedIn:      false,
		Role:          ROLE_ADMIN,
		CurrentToken:  "",
		CurrentClient: 0,
	}

	return newUser
}

/*IsAuthenticatedForRole returns true if the user is AT OR ABOVE
the requested security level*/
func (u *User) IsAuthenticatedForRole(r Role) bool {
	return u.Role >= r
}
