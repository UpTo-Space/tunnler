package main

import (
	"context"
	"errors"
	"io"
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
	subscriberMu   sync.Mutex
	subscriber     *subscriber

	logf func(f string, v ...interface{})
}

type subscriber struct {
	msgs      chan []byte
	closeSlow func()
}

func newtunnlerServer() *tunnlerServer {
	ts := &tunnlerServer{
		messageBuffer:  16,
		logf:           log.Printf,
		subscriber:     nil,
		messageLimiter: rate.NewLimiter(rate.Every(time.Millisecond*100), 8),
	}

	ts.serveMux.HandleFunc("/", ts.httpHandler)
	ts.serveMux.HandleFunc("/tunnler/connection/subscribe", ts.subscribeHandler)
	ts.serveMux.HandleFunc("/tunnler/connection/publish", ts.publishHandler)

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

	ts.publish(reqData)
}

func (ts *tunnlerServer) subscribeHandler(w http.ResponseWriter, r *http.Request) {
	err := ts.subscribe(w, r)
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

func (ts *tunnlerServer) publishHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}

	body := http.MaxBytesReader(w, r.Body, 8192)
	msg, err := io.ReadAll(body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusRequestEntityTooLarge), http.StatusRequestEntityTooLarge)
	}

	ts.publish(msg)

	w.WriteHeader(http.StatusAccepted)
}

func (ts *tunnlerServer) subscribe(w http.ResponseWriter, r *http.Request) error {
	var mu sync.Mutex
	var c *websocket.Conn
	var closed bool
	s := &subscriber{
		msgs: make(chan []byte, ts.messageBuffer),
		closeSlow: func() {
			mu.Lock()
			defer mu.Unlock()
			closed = true
			if c != nil {
				c.Close(websocket.StatusPolicyViolation, "connection too slow to keep up")
			}
		},
	}

	ts.registerSubscriber(s)

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
			err := writeTimeout(ctx, time.Second*5, c, msg)
			if err != nil {
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (ts *tunnlerServer) registerSubscriber(s *subscriber) {
	ts.subscriberMu.Lock()
	ts.subscriber = s
	ts.subscriberMu.Unlock()
}

func (ts *tunnlerServer) publish(msg []byte) {
	if ts.subscriber == nil {
		return
	}

	ts.subscriberMu.Lock()
	defer ts.subscriberMu.Unlock()

	ts.messageLimiter.Wait(context.Background())

	select {
	case ts.subscriber.msgs <- msg:
	default:
		go ts.subscriber.closeSlow()
	}
}

func writeTimeout(ctx context.Context, timeout time.Duration, c *websocket.Conn, msg []byte) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return c.Write(ctx, websocket.MessageBinary, msg)
}
