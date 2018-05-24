package models

type Position struct {
	X int `json:"x" gorethink:"x"`
	Y int `json:"y" gorethink:"y"`
}
