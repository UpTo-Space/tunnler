package main

import (
	"log"
	"net/http"

	"github.com/UpTo-Space/tunnler/database"
)

type authServer struct {
	serveMux http.ServeMux
	logf     func(f string, v ...interface{})
	db       *database.DatabaseConnection
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
	server.registerRoutes()

	return server, nil
}

func (as *authServer) registerRoutes() {
	as.serveMux.HandleFunc("/auth/register", as.registerUserHandler)
	as.serveMux.HandleFunc("/auth/activate", as.activateUserHandler)
}

func (as *authServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	as.serveMux.ServeHTTP(w, r)
}
