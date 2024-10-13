package main

import (
	"log"
	"net/http"

	"github.com/UpTo-Space/tunnler/common"
	"github.com/UpTo-Space/tunnler/database"
)

type authServer struct {
	serveMux   http.ServeMux
	logf       func(f string, v ...interface{})
	db         *database.DatabaseConnection
	tokenMaker *common.PasetoMaker
}

func newServer() (*authServer, error) {
	server := &authServer{
		logf: log.Printf,
	}

	db, err := database.NewDatabaseConnection()
	if err != nil {
		server.logf("error in setting up database connection: %v", err)
		return nil, err
	}
	server.db = db

	tokenMaker, err := common.NewPaseto(symmetricKey)
	if err != nil {
		server.logf("error in setting up the token maker: %v", err)
	}
	server.tokenMaker = tokenMaker

	server.registerRoutes()

	return server, nil
}

func (as *authServer) registerRoutes() {
	as.serveMux.HandleFunc("/auth/register", as.registerUserHandler)
	as.serveMux.HandleFunc("/auth/activate", as.activateUserHandler)
	as.serveMux.HandleFunc("/auth/login", as.loginUserHandler)
}

func (as *authServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	as.serveMux.ServeHTTP(w, r)
}
