package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

const TASK_COMPLETE = "Task Complete"

func longHandler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(3 * time.Second)
	fmt.Fprintf(w, TASK_COMPLETE)
}

func routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /long", longHandler)
	return mux
}

func NewServer() *http.Server {
	srv := http.Server{
		Addr:    ":8080",
		Handler: routes(),
	}

	return &srv
}

func RunNotGracefulServer() *http.Server {
	srv := NewServer()
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server execution error: %v", err)
		}
	}()

	return srv
}

func gracefulShutdown(server *http.Server, done chan bool, shutdownChan <-chan struct{}) {
	<-shutdownChan
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Fail to shutdown server: %v", err)
	}
	done <- true
}

func RunGracefulServer(done chan bool, shutdownChan <-chan struct{}) (*http.Server, error) {
	srv := NewServer()
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server execution error: %v", err)
		}
	}()
	go gracefulShutdown(srv, done, shutdownChan)
	return srv, nil
}
