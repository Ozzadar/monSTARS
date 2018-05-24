package gameservice

import (
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/ozzadar/monSTARS/models"
)

type Client struct {
	Send         chan Message
	socket       *websocket.Conn
	stopChannels map[int]chan bool
	authToken    *models.JwtToken
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

func (c *Client) Read() {
	var message Message
	for {
		if err := c.socket.ReadJSON(&message); err != nil {
			break
		}

		//Handle read
		if c.authToken == nil {
			HandleMessage(c, message)
		} else {
			HandleUserMessage(c, message)
		}
	}
	c.socket.Close()
}

func (c *Client) Write() {
	for msg := range c.Send {
		if err := c.socket.WriteJSON(msg); err != nil {
			break
		}
	}
	c.socket.Close()
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
		authToken:    nil,
	}
}

func (c *Client) Run(closeFunc func(*Client)) {
	defer c.Close()
	defer closeFunc(c)

	go c.Write()
	c.Read()
}
