package main

import (
	"context"
	"errors"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
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

	return ts
}

func (ts *tunnlerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ts.serveMux.ServeHTTP(w, r)
}

func (ts *tunnlerServer) httpHandler(w http.ResponseWriter, r *http.Request) {
	// First of all, check if we have a websocket connection request
	if r.Header.Get("Upgrade") == "websocket" {
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
	} else {
		// Otherwise handle it as a forwardable http request
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		reqData, err := common.SerializeRequest(r)
		if err != nil {
			ts.logf("error serializing request: %v", err)
		}

		ts.forwardRequest(reqData)

		for {
			select {
			case rsp := <-ts.client.rsps:
				ts.handleResponse(w, rsp)
				return
			case <-ctx.Done():
				return
			}
		}
	}
}

func (ts *tunnlerServer) handleResponse(w http.ResponseWriter, b []byte) {
	msg, err := common.DeserializeResponse(b)
	if err != nil {
		ts.logf("error in deserializing response: %v", err)
	}

	body, err := io.ReadAll(msg.Body)
	if err != nil {
		ts.logf("error in reading response object: %v", err)
	}

	for k, v := range msg.Header {
		w.Header().Set(k, strings.Join(v, ","))
	}
	w.WriteHeader(msg.StatusCode)
	w.Write(body)
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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		for {
			_, msg, err := c.Read(ctx)
			if err != nil {
				ts.logf("error in receiving: %v", err)
			}
			s.rsps <- msg
		}
	}()

	for {
		select {
		case msg := <-s.msgs:
			err := c.Write(ctx, websocket.MessageBinary, msg)
			if err != nil {
				return err
			}
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
