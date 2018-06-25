package gameservice

import (
	"fmt"

	"github.com/ozzadar/monSTARS/services/authservice"
)

type Message struct {
	Type    string                 `json:"type"`
	Payload map[string]interface{} `json:"payload"`
}

func HandleMessage(client *Client, m Message) {
	switch m.Type {
	case "CONNECTED":
		{
			client.Send <- Message{
				Type: "LoginRequest",
				Payload: map[string]interface{}{
					"message": "Connected to server. Please log in",
				},
			}
			break
		}
	case "LoginRequest":
		{
			username := m.Payload["username"].(string)
			password := m.Payload["password"].(string)

			fmt.Println("logging in")
			ProcessLoginRequest(username, password, client)
			break
		}
	case "LoginTokenRequest":
		{
			token := m.Payload["token"].(string)
			ProcessLoginTokenRequest(token, client)
			break
		}
	}
}

func SendLoginFailed(client *Client) {
	client.Send <- Message{
		Type: "LoginFailed",
		Payload: map[string]interface{}{
			"message": "Invalid token.",
		},
	}
}

func ProcessLoginRequest(username string, password string, client *Client) {
	token := authservice.LoginWithUserPass(username, password)

	if token != "" {

		_, tokenCheck := authservice.VerifyJWT(token)
		client.authToken = tokenCheck
		client.Send <- Message{
			Type: "LoginSuccessful",
			Payload: map[string]interface{}{
				"message":  "LOGIN_SUCCESSFUL",
				"token":    token,
				"username": username,
			},
		}
		return
	}
	SendLoginFailed(client)
}

func ProcessLoginTokenRequest(token string, client *Client) {
	isAuth, jwtToken := authservice.VerifyJWT(token)

	if isAuth {
		client.authToken = jwtToken
		client.Send <- Message{
			Type: "LoginSuccessful",
			Payload: map[string]interface{}{
				"message":  "LOGIN_SUCCESSFUL",
				"token":    token,
				"username": jwtToken.Owner,
			},
		}
		return
	}
	SendLoginFailed(client)
}
