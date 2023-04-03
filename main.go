package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/QQsha/cards/deck/poker"
	"github.com/QQsha/cards/handlers"
	"github.com/QQsha/cards/repository/memorystore"
	"github.com/go-chi/chi/v5"
)

func main() {
	store := memorystore.NewInMemoryStore()
	deckService := poker.NewDeckService(store)
	rt := chi.NewRouter()
	
	srv := handlers.NewServer(rt, deckService)

	srv.Routes()
	http.ListenAndServe(":8080", srv.Router)

	sig := make(chan os.Signal, 1)

	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
}
