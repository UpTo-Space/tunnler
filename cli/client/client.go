package client

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/UpTo-Space/tunnler/common"
	"github.com/coder/websocket"
)

type TunnlerConnectionInfo struct {
	// Adress of the local server to forward requests to
	HostAdress string
	// Port of the local server
	HostPort string

	// Tunnler IP adress
	TunnlerAdress string
	// Tunnler Port
	TunnlerPort string
}

type tunnlerClient struct {
	connectionInfo TunnlerConnectionInfo
	logf           func(f string, v ...interface{})
}

func NewTunnlerClient(ci TunnlerConnectionInfo) *tunnlerClient {
	tc := &tunnlerClient{
		connectionInfo: ci,
		logf:           log.Printf,
	}

	return tc
}

func (tc *tunnlerClient) Connect() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	c, _, err := websocket.Dial(
		ctx, fmt.Sprintf("ws://%s:%s/tunnler/connection/subscribe",
			tc.connectionInfo.TunnlerAdress, tc.connectionInfo.TunnlerPort), nil)

	if err != nil {
		c.Close(websocket.StatusInternalError, "")
	}

	for {
		_, msgByte, err := c.Read(ctx)
		if err != nil {
			log.Fatalf("error in receiving message: %v", err)
			break
		}

		req, err := common.DeserializeRequest(msgByte)
		if err != nil {
			log.Fatalf("error in deserializing request: %v", err)
		}

		resp, err := tc.TunnelRequest(req)
		if err != nil {
			log.Fatalf("error in tunneling request: %v", err)
		}

		fmt.Println(resp)
	}

	defer c.CloseNow()
}

func (tc *tunnlerClient) TunnelRequest(req *http.Request) (*http.Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	treq := req.Clone(ctx)
	treq.RequestURI = ""
	treq.URL = &url.URL{
		Scheme:      "http",
		Opaque:      req.URL.Opaque,
		User:        req.URL.User,
		Host:        fmt.Sprintf("%s:%s", tc.connectionInfo.HostAdress, tc.connectionInfo.HostPort),
		Path:        req.URL.Path,
		RawPath:     req.URL.RawPath,
		OmitHost:    req.URL.OmitHost,
		ForceQuery:  req.URL.ForceQuery,
		RawQuery:    req.URL.RawQuery,
		Fragment:    req.URL.Fragment,
		RawFragment: req.URL.RawFragment,
	}

	client := &http.Client{}
	resp, err := client.Do(treq)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (tc *tunnlerClient) TunnelResponse(resp *http.Response) {

}
