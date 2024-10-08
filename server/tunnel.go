package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/UpTo-Space/tunnler/common"
	"github.com/coder/websocket"
	"golang.org/x/time/rate"
)

type tunnlerServer struct {
	messageBuffer  int
	messageLimiter *rate.Limiter
	serveMux       http.ServeMux
	clientMu       sync.Mutex
	client         *client

	logf func(f string, v ...interface{})
}

type client struct {
	msgs      chan []byte
	rsps      chan []byte
	closeSlow func()
}

func newtunnlerServer() *tunnlerServer {
	ts := &tunnlerServer{
		messageBuffer:  16,
		logf:           log.Printf,
		client:         nil,
		messageLimiter: rate.NewLimiter(rate.Every(time.Millisecond*100), 8),
	}

	ts.serveMux.HandleFunc("/", ts.httpHandler)
	ts.serveMux.HandleFunc("/tunnler/connection/initialize", ts.initializeHandler)

	return ts
}

func (ts *tunnlerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ts.serveMux.ServeHTTP(w, r)
}

func (ts *tunnlerServer) httpHandler(w http.ResponseWriter, r *http.Request) {
	reqData, err := common.SerializeRequest(r)
	if err != nil {
		ts.logf("error serializing request: %v", err)
	}

	ts.forwardRequest(reqData)
}

func (ts *tunnlerServer) initializeHandler(w http.ResponseWriter, r *http.Request) {
	err := ts.initialize(w, r)
	if errors.Is(err, context.Canceled) {
		return
	}
	if websocket.CloseStatus(err) == websocket.StatusNormalClosure ||
		websocket.CloseStatus(err) == websocket.StatusGoingAway {
		return
	}
	if err != nil {
		ts.logf("%v", err)
		return
	}
}

func (ts *tunnlerServer) initialize(w http.ResponseWriter, r *http.Request) error {
	var mu sync.Mutex
	var c *websocket.Conn
	var closed bool
	s := &client{
		msgs: make(chan []byte, ts.messageBuffer),
		rsps: make(chan []byte, ts.messageBuffer),
		closeSlow: func() {
			mu.Lock()
			defer mu.Unlock()
			closed = true
			if c != nil {
				c.Close(websocket.StatusPolicyViolation, "connection too slow to keep up")
			}
		},
	}

	ts.registerClient(s)

	c2, err := websocket.Accept(w, r, nil)
	if err != nil {
		return err
	}

	mu.Lock()
	if closed {
		mu.Unlock()
		return net.ErrClosed
	}

	c = c2
	mu.Unlock()
	defer c.CloseNow()

	ctx := c.CloseRead(context.Background())

	for {
		select {
		case msg := <-s.msgs:
			err := c.Write(ctx, websocket.MessageBinary, msg)
			if err != nil {
				return err
			}
		case rsps := <-s.rsps:
			ts.logf(string(rsps))
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (ts *tunnlerServer) registerClient(s *client) {
	ts.clientMu.Lock()
	ts.client = s
	ts.clientMu.Unlock()
}

func (ts *tunnlerServer) forwardRequest(msg []byte) {
	if ts.client == nil {
		return
	}

	ts.clientMu.Lock()
	defer ts.clientMu.Unlock()

	ts.messageLimiter.Wait(context.Background())

	select {
	case ts.client.msgs <- msg:
	default:
		go ts.client.closeSlow()
	}
}

func writeTimeout(ctx context.Context, timeout time.Duration, c *websocket.Conn, msg []byte) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return c.Write(ctx, websocket.MessageBinary, msg)
}
