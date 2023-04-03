package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/QQsha/cards/deck"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Server struct {
	Router *chi.Mux
	Deck   deck.Service
}

func NewServer(router *chi.Mux, deckService deck.Service) *Server {
	return &Server{Router: router, Deck: deckService}
}

func (s *Server) Routes() {
	s.Router.Post("/decks/create", s.createDeckHandler())
	s.Router.Get("/decks/{id}", s.openDeckHandler())
	s.Router.Patch("/decks/{id}", s.drawCardsHandler())
}

func (s *Server) createDeckHandler() http.HandlerFunc {
	type response struct {
		ID        string `json:"deck_id"`
		Shuffled  bool   `json:"shuffled"`
		Remaining int    `json:"remaining"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			shuffled bool
			cards    []string
			err      error
		)

		if s := r.URL.Query().Get("shuffle"); s != "" {
			shuffled, err = strconv.ParseBool(s)
			if err != nil {
				http.Error(w, errors.New("invalid shuffle value").Error(), http.StatusBadRequest)
				return
			}
		}

		if c := r.URL.Query().Get("cards"); c != "" {
			cards = strings.Split(strings.ToUpper(c), ",")
		}

		deck, err := s.Deck.CreateDeck(cards, shuffled)
		if err != nil {
			http.Error(w, fmt.Errorf("error: %w", err).Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response{ID: deck.ID, Shuffled: deck.Shuffled, Remaining: deck.Size})
	}
}

func (s *Server) openDeckHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		deckID := chi.URLParam(r, "id")
		err := uuidValidator(deckID)
		if err != nil {
			http.Error(w, fmt.Errorf("error: %w", err).Error(), http.StatusNotFound)
			return
		}

		deck, err := s.Deck.OpenDeck(deckID)
		if err != nil {
			http.Error(w, fmt.Errorf("error: %w", err).Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(deck)
	}
}

func (s *Server) drawCardsHandler() http.HandlerFunc {
	type request struct {
		Draw    int `json:"draw"`
		Version int `json:"version"` // to keep this handler idempotent
	}
	return func(w http.ResponseWriter, r *http.Request) {
		deckID := chi.URLParam(r, "id")
		err := uuidValidator(deckID)
		if err != nil {
			http.Error(w, fmt.Errorf("error: %w", err).Error(), http.StatusNotFound)
			return
		}

		req := request{}
		err = json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, fmt.Errorf("error: %w", err).Error(), http.StatusBadRequest)
			return
		}

		cards, err := s.Deck.DrawCards(deckID, req.Draw, req.Version)
		if err != nil {
			http.Error(w, fmt.Errorf("error: %w", err).Error(), http.StatusBadGateway)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(cards)
	}
}

func uuidValidator(id string) error {
	_, err := uuid.Parse(id)
	return err
}
