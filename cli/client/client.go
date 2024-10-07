package client

import (
	"context"
	"fmt"
	"log"
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

	for {
		msgType, msgByte, err := c.Read(ctx)
		if err != nil {
			log.Fatalf("error in receiving message: %v", err)
			break
		}

		fmt.Printf("Received %s // %s\n", msgType, msgByte)
	}

	defer c.CloseNow()
}
