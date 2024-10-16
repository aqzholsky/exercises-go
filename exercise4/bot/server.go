package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"sync"
	"time"
)

// TODO: logic to response
// TODO: Don't start server if judge returns bad request

type readyListener struct {
	net.Listener
	ready chan struct{}
	once  sync.Once
}

func (l *readyListener) Accept() (net.Conn, error) {
	l.once.Do(func() { close(l.ready) })
	return l.Listener.Accept()
}

type RequestMove struct {
	Board []string `json:"board"`
	Token string
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func moveHandler(w http.ResponseWriter, r *http.Request) {
	var reqMove RequestMove
	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	dec.Decode(&reqMove)

	fmt.Println(reqMove.Board)
	fmt.Println(reqMove.Token)

	board := reqMove.Board
	fmt.Println(board[0])

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	js, _ := json.Marshal(map[string]int{"index": 0})
	w.Write(js)
}

func startServer(ctx context.Context) <-chan struct{} {
	port := os.Getenv("PORT")
	mux := http.NewServeMux()
	mux.HandleFunc("GET /ping", pingHandler)
	mux.HandleFunc("POST /move", moveHandler)

	server := &http.Server{
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
		BaseContext: func(_ net.Listener) context.Context {
			return ctx
		},
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		panic(err)
	}

	ready := make(chan struct{})
	list := &readyListener{Listener: listener, ready: ready}

	slog.InfoContext(
		ctx,
		"starting service",
		"port", port,
	)

	go func() {
		err := server.Serve(list)
		if !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()

	return ready
}
