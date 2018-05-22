package gameservice

import (
	"errors"
	"fmt"
	"log"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/websocket"
	"github.com/ozzadar/monSTARS/config"
	"github.com/ozzadar/monSTARS/db"
	"github.com/ozzadar/monSTARS/models"
)

//JWTClaims type
type JWTClaims struct {
	models.User
	jwt.StandardClaims
}

type Client struct {
	Send         chan Message
	socket       *websocket.Conn
	stopChannels map[int]chan bool
	id           string
	userName     string
}

func (c *Client) NewStopChannel(stopKey int) chan bool {
	c.StopForKey(stopKey)
	stop := make(chan bool)
	c.stopChannels[stopKey] = stop
	return stop
}

func (c *Client) StopForKey(key int) {
	if ch, found := c.stopChannels[key]; found {
		ch <- true
		delete(c.stopChannels, key)
	}
}

func (client *Client) Read() {
	var message Message
	for {
		if err := client.socket.ReadJSON(&message); err != nil {
			break
		}

		//Handle read
		HandleMessage(client, message)
	}
	client.socket.Close()
}

func (client *Client) Write() {
	for msg := range client.Send {
		if err := client.socket.WriteJSON(msg); err != nil {
			break
		}
	}
	client.socket.Close()
}

func (c *Client) Close() {
	for _, ch := range c.stopChannels {
		ch <- true
	}
	fmt.Println("Client Disconnected")
	close(c.Send)
}

func NewClient(socket *websocket.Conn) *Client {
	return &Client{
		Send:         make(chan Message),
		socket:       socket,
		stopChannels: make(map[int]chan bool),
	}
}

func (c *Client) Run(closeFunc func(*Client)) {
	defer c.Close()
	defer closeFunc(c)

	go c.Write()
	c.Read()
}

func LoginWithUserPass(username string, password string) string {
	//Check if valid login
	user := db.GetUser(username, password)

	if user != nil {
		//Create JWT Token
		token, err := CreateJwtToken(user)

		if err != nil {
			log.Printf("Failed to create token: %#v", err)
			return ""
		}

		user.LoggedIn = true
		go func() {
			db.UpdateUserLoginState(user)
			db.AddJWT(&models.JwtToken{
				Owner: user.Username,
				Token: token,
			})
		}()
		return token
	}
	return ""
}

//CreateJwtToken create a token
func CreateJwtToken(user *models.User) (string, error) {
	claims := JWTClaims{
		User: *user,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			Id:        "main_user_id",
			ExpiresAt: time.Now().Add(time.Hour * 24 * 7).Unix(),
		},
	}

	rawToken := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	if signingKey, err := config.Config.GetString("default", "jwt_secret_key"); err == nil {
		token, err := rawToken.SignedString([]byte(signingKey))

		if err != nil {
			log.Printf("Failed to create token")
			return "", err
		}

		return token, nil
	}

	panic(errors.New("jwt_secret_key defined in config"))
}
