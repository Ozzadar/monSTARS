package gameservice

import (
	"fmt"
	"strconv"
)

var (
	OnlineUsers []*Client
)

func RegisterClientOnline(client *Client) {

	OnlineUsers = append(OnlineUsers, client)
	go client.Run(clientDisconnected)

	fmt.Println("User online #" + strconv.Itoa(len(OnlineUsers)))
}

func clientDisconnected(client *Client) {
	for i, c := range OnlineUsers {
		if c == client {
			OnlineUsers = append(OnlineUsers[:i], OnlineUsers[i+1:]...)
			fmt.Println("Client #" + strconv.Itoa(i+1) + " disconnected. Online users remaining: " + strconv.Itoa(len(OnlineUsers)))
		}
	}
}
