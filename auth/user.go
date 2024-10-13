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

type ActivateUserParams struct {
	Username       string `json:"username"`
	ActivationCode int    `json:"activationCode"`
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

	if err := as.registerUser(params.Username, params.Password, params.Email); err != nil {
		as.logf("error in registering user: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (as *authServer) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var params ActivateUserParams

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		as.logf("error in decoding activate user params: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	valid, err := as.checkActivationCode(params.ActivationCode, params.Username)
	if err != nil {
		as.logf("error in activate user: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !valid {
		as.logf("invalid activation code: %s", params.ActivationCode)

		if err := as.increaseActivationTries(params.Username); err != nil {
			as.logf("error in increasing activation tries: %v", err)
		}

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := as.activateUser(params.Username); err != nil {
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
