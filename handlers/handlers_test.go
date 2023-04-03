package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/QQsha/cards/deck/poker"
	"github.com/QQsha/cards/entity"
	"github.com/QQsha/cards/repository/memorystore"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestCreateDeckHandler(t *testing.T) {
	store := memorystore.NewInMemoryStore()
	deckService := poker.NewDeckService(store)
	rt := chi.NewRouter()
	srv := NewServer(rt, deckService)
	srv.Routes()
	t.Run("create deck successful request", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/decks/create", nil)
		w := httptest.NewRecorder()
		srv.Router.ServeHTTP(w, req)

		deck := entity.Deck{}
		json.NewDecoder(w.Body).Decode(&deck)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, 52, deck.Size)
	})

	t.Run("create deck invalid shuffle", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/decks/create?shuffle=go", nil)
		w := httptest.NewRecorder()
		srv.Router.ServeHTTP(w, req)
		res := strings.TrimSpace(w.Body.String())
		assert.Equal(t, res, "invalid shuffle value")
	})
}

func TestOpenDeckHandler(t *testing.T) {
	store := memorystore.NewInMemoryStore()
	deckService := poker.NewDeckService(store)
	rt := chi.NewRouter()
	srv := NewServer(rt, deckService)
	srv.Routes()
	expectedDeck, err := deckService.CreateDeck(nil, false)
	assert.NoError(t, err)

	t.Run("open deck successful request", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/decks/%s", expectedDeck.ID), nil)
		w := httptest.NewRecorder()
		srv.Router.ServeHTTP(w, req)

		deck := entity.Deck{}
		json.NewDecoder(w.Body).Decode(&deck)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, expectedDeck, deck)
	})

	t.Run("open deck invalid id", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/decks/id123", nil)
		w := httptest.NewRecorder()
		srv.Router.ServeHTTP(w, req)

		res := strings.TrimSpace(w.Body.String())
		assert.Equal(t, res, "error: invalid UUID length: 5")
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestDrawCardsHandler(t *testing.T) {
	type request struct {
		Draw    int `json:"draw"`
		Version int `json:"version"`
	}

	store := memorystore.NewInMemoryStore()
	deckService := poker.NewDeckService(store)
	rt := chi.NewRouter()
	srv := NewServer(rt, deckService)
	srv.Routes()
	expectedDeck, err := deckService.CreateDeck(nil, false)
	assert.NoError(t, err)

	t.Run("draw card, successful request", func(t *testing.T) {
		f := request{Draw: 1, Version: 0}
		var buf bytes.Buffer
		err := json.NewEncoder(&buf).Encode(f)
		assert.NoError(t, err)
		req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/decks/%s", expectedDeck.ID), &buf)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		srv.Router.ServeHTTP(w, req)

		cards := []entity.Card{}
		json.NewDecoder(w.Body).Decode(&cards)

		expect := []entity.Card{{Value: "ACE", Suit: "SPADES", Code: "AS"}}
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, expect, cards)
	})

	t.Run("draw card, invalid id", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPatch, "/decks/id123", nil)
		w := httptest.NewRecorder()
		srv.Router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Equal(t, w.Body.String(), "error: invalid UUID length: 5\n")
	})
}
