package main

import (
	"encoding/json"
	"net/http"
	"time"
)

type RegisterUserParams struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type LoginUserParams struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
	Username    string `json:"username"`
}

func (as *authServer) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var params RegisterUserParams

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		as.logf("error in decoding register user params: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	activationCode, err := as.registerUser(params.Username, params.Password, params.Email)
	if err != nil {
		as.logf("error in registering user: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := sendActivationMain(params.Email, params.Username, activationCode); err != nil {
		as.logf("error in sending mail for user: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (as *authServer) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	userName := r.URL.Query().Get("username")
	code := r.URL.Query().Get("code")

	valid, err := as.checkActivationCode(code, userName)
	if err != nil {
		as.logf("error in activate user: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !valid {
		as.logf("invalid activation code: %s", code)

		if err := as.increaseActivationTries(userName); err != nil {
			as.logf("error in increasing activation tries: %v", err)
		}

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := as.activateUser(userName); err != nil {
		as.logf("error in activating user: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (as *authServer) loginUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var params LoginUserParams

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		as.logf("error in decoding login user params: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	success, err := as.attemptLogin(params.Username, params.Password)
	if err != nil {
		as.logf("login failed: %v", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if !success {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	tk, err := as.tokenMaker.CreateToken(params.Username, time.Hour*24)

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	response := &LoginResponse{
		AccessToken: tk,
		Username:    params.Username,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
