package handlers

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
	"github.com/ozzadar/monSTARS/services/gameservice"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func Hello(c echo.Context) error {

	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}

	//Register online
	gameservice.RegisterClientOnline(gameservice.NewClient(ws))

	return nil
}
