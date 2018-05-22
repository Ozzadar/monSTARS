package gameservice

import (
	"fmt"
)

type Message struct {
	Type    string      `json:"type"`
	Message interface{} `json:"message"`
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func HandleMessage(client *Client, m Message) {
	fmt.Println("Received message: " + m.Type)
	switch m.Type {
	case "CONNECTED":
		{
			client.Send <- Message{
				Type:    "LoginRequest",
				Message: "Welcome to monSTARS, please log in.",
			}
			break
		}
	case "LoginRequest":
		{

			if creds, ok := m.Message.(map[string]interface{}); ok {
				username := creds["username"].(string)
				password := creds["password"].(string)

				token := LoginWithUserPass(username, password)

				if token != "" {
					client.Send <- Message{
						Type:    "LoginSuccessful",
						Message: token,
					}
					return
				}
			}
			client.Send <- Message{
				Type:    "LoginFailed",
				Message: nil,
			}
			break
		}
	}
}
