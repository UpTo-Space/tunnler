package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	log.SetFlags(0)
	ctx := context.Background()

	err := run(ctx)
	if err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	handler, err := newServer()
	if err != nil {
		return err
	}

	handler.logf("listening on http://%v:%v", listenHostName, hostPort)

	s := &http.Server{
		Handler:      handler,
		Addr:         fmt.Sprintf("%s:%s", listenHostName, hostPort),
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}

	errorChannel := make(chan error, 1)
	go func() {
		errorChannel <- s.ListenAndServe()
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	select {
	case err := <-errorChannel:
		log.Printf("failed to serve: %v", err)
	case sig := <-sigs:
		log.Printf("terminating: %v", sig)
	}

	return s.Shutdown(ctx)
}
