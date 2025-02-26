package main

import (
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/talgat-ruby/exercises-go/exercise7/blogging-platform/database"
	"github.com/talgat-ruby/exercises-go/exercise7/blogging-platform/handlers"
)

func main() {
	database.DatabaseInit()
	r := mux.NewRouter()

	r.HandleFunc("/posts", handlers.HandlePosts).Methods("GET", "POST")
	r.HandleFunc("/posts/{id}", handlers.HandlePost).Methods("GET", "PUT", "DELETE")

	port := "8080"
	slog.Info("Server is starting", "port", port)

	if err := http.ListenAndServe(":"+port, r); err != nil {
		slog.Error("Server failed to start", "error", err)
	}

	slog.Info("Server has started", "port", port)
}
