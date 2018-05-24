package gameservice

import (
	"github.com/ozzadar/monSTARS/db"
)

func HandleUserMessage(client *Client, m Message) {
	if client.authToken == nil {
		SendUnauthorized(client)
	}
	switch m.Type {
	case "RequestCharacters":
		{
			HandleCharactersRequest(client)
		}
	}
}

func SendUnauthorized(client *Client) {
	client.Send <- Message{
		Type: "Unauthorized",
		Payload: map[string]interface{}{
			"message": "No Auth Token, please login again",
		},
	}
}

func HandleCharactersRequest(client *Client) {
	characters := db.GetAllCharactersForUser(client.authToken.Owner)

	client.Send <- Message{
		Type: "RequestCharacterResponse",
		Payload: map[string]interface{}{
			"message":    "success",
			"characters": characters,
		},
	}
}
