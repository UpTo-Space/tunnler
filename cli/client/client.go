package client

import (
	"context"
	"fmt"
	"time"

	"github.com/coder/websocket"
)

type tunnelClient struct {
}

func Connect(adress string, port string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	c, _, err := websocket.Dial(ctx, fmt.Sprintf("ws://%s:%s/tunnler/connection/subscribe", adress, port), nil)
	if err != nil {
		c.Close(websocket.StatusInternalError, "")
	}

	defer c.CloseNow()
}
