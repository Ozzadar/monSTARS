package models

//JwtToken holds jwt tokens
type JwtToken struct {
	Owner string `json:"owner" gorethink:"owner"`
	Token string `json:"token" gorethink:"token"`
}
