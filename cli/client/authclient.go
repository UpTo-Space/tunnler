package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/UpTo-Space/tunnler/common"
)

type TunnlerAuthConnectionInfo struct {
	// Adress of the auth server
	HostAdress string
	// Port of the auth server
	HostPort string
	// Scheme of the auth server
	HostScheme string
}

type authClient struct {
	connectionInfo TunnlerAuthConnectionInfo
	logf           func(f string, v ...interface{})
}

func NewAuthClient(ci TunnlerAuthConnectionInfo) *authClient {
	ac := &authClient{
		connectionInfo: ci,
		logf:           log.Printf,
	}

	return ac
}

func (ac *authClient) Register(username, password, email string) {
	params := &common.RegisterUserParams{
		Username: username,
		Password: password,
		Email:    email,
	}

	b, err := json.Marshal(params)

	if err != nil {
		ac.logf("error in marshalling json: %v", err)
	}

	body := bytes.NewBuffer(b)

	r, err := http.Post(fmt.Sprintf("%s://%s:%s/auth/register", ac.connectionInfo.HostScheme, ac.connectionInfo.HostAdress, ac.connectionInfo.HostPort), "application/json; charset=utf-8", body)

	if err != nil {
		ac.logf("error in sending auth post request: %v", err)
	}

	defer r.Body.Close()

	if r.StatusCode != http.StatusCreated {
		ac.logf("error in registering user, received Statuscode %v", r.StatusCode)
	}

	ac.logf("check your E-Mails for activating your account")
}

func (ac *authClient) Login(username, password string) {
	params := &common.LoginUserParams{
		Username: username,
		Password: password,
	}

	b, err := json.Marshal(params)

	if err != nil {
		ac.logf("error in marshalling json: %v", err)
	}

	body := bytes.NewBuffer(b)

	r, err := http.Post(fmt.Sprintf("%s://%s:%s/auth/login", ac.connectionInfo.HostScheme, ac.connectionInfo.HostAdress, ac.connectionInfo.HostPort), "application/json; charset=utf-8", body)

	if err != nil {
		ac.logf("error in sending auth post request: %v", err)
	}

	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		ac.logf("error in logging in user, received Statuscode %v", r.StatusCode)
	}

	var loginResponse common.LoginResponse
	err = json.NewDecoder(r.Body).Decode(&loginResponse)

	if err != nil {
		ac.logf("error in decoding response: %v", err)
	}

	fmt.Println(loginResponse)
}
