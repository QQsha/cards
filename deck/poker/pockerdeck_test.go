package poker_test

import (
	"testing"

	"github.com/QQsha/cards/deck/poker"
	"github.com/QQsha/cards/entity"
	"github.com/QQsha/cards/repository/memorystore"
	"github.com/stretchr/testify/assert"
)

func TestDeckCreate(t *testing.T) {
	store := memorystore.NewInMemoryStore()
	deckService := poker.NewDeckService(store)
	t.Run("full deck, no shuffle", func(t *testing.T) {
		deck, err := deckService.CreateDeck(nil, false)
		assert.NoError(t, err)
		assert.Equal(t, 52, deck.Size)
		assert.Equal(t, false, deck.Shuffled)
		assert.Equal(t, entity.Card{Value: "ACE", Suit: "SPADES", Code: "AS"}, deck.Cards[0])
		assert.Equal(t, entity.Card{Value: "KING", Suit: "HEARTS", Code: "KH"}, deck.Cards[51])
	})

	t.Run("full deck, shuffle", func(t *testing.T) {
		deck, err := deckService.CreateDeck(nil, true)
		assert.NoError(t, err)

		assert.Equal(t, 52, deck.Size)
		assert.Equal(t, true, deck.Shuffled)
		assert.NotEqual(t, entity.Card{Value: "ACE", Suit: "SPADES", Code: "AS"}, deck.Cards[0])
		assert.NotEqual(t, entity.Card{Value: "KING", Suit: "HEARTS", Code: "KH"}, deck.Cards[51])
	})

	t.Run("custom deck, no shuffle", func(t *testing.T) {
		expected := []entity.Card{
			{
				Value: "ACE",
				Suit:  "SPADES",
				Code:  "AS",
			},
			{
				Value: "KING",
				Suit:  "DIAMONDS",
				Code:  "KD",
			},
			{
				Value: "ACE",
				Suit:  "CLUBS",
				Code:  "AC",
			},
		}

		deck, err := deckService.CreateDeck([]string{"AS", "KD", "AC"}, false)
		if !assert.NoError(t, err) {
			t.Error(err)
		}

		assert.Equal(t, 3, deck.Size)
		assert.Equal(t, false, deck.Shuffled)
		assert.Equal(t, expected, deck.Cards)
	})

	t.Run("custom deck, shuffle", func(t *testing.T) {
		expected := []entity.Card{
			{
				Value: "ACE",
				Suit:  "SPADES",
				Code:  "AS",
			},
			{
				Value: "KING",
				Suit:  "DIAMONDS",
				Code:  "KD",
			},
			{
				Value: "ACE",
				Suit:  "CLUBS",
				Code:  "AC",
			},
			{
				Value: "2",
				Suit:  "CLUBS",
				Code:  "2C",
			},
			{
				Value: "KING",
				Suit:  "HURTS",
				Code:  "KH",
			},
		}

		deck, err := deckService.CreateDeck([]string{"AS", "KD", "AC", "2C", "KH"}, true)
		assert.NoError(t, err)

		assert.Equal(t, 5, deck.Size)
		assert.Equal(t, true, deck.Shuffled)
		assert.NotEqual(t, expected, deck.Cards)
	})

	t.Run("invalid card code error", func(t *testing.T) {
		_, err := deckService.CreateDeck([]string{"11", "22", "33"}, true)
		assert.EqualError(t, err, "invalid card code")

		_, err = deckService.CreateDeck([]string{"AAAAA", "22", "33"}, true)
		assert.EqualError(t, err, "invalid card code")
	})
}

func TestOpenDeck(t *testing.T) {
	store := memorystore.NewInMemoryStore()
	deckService := poker.NewDeckService(store)
	t.Run("open full deck", func(t *testing.T) {
		deck, err := deckService.CreateDeck(nil, false)
		assert.NoError(t, err)

		deck, err = deckService.OpenDeck(deck.ID)
		assert.NoError(t, err)

		assert.Equal(t, 52, deck.Size)
		assert.Equal(t, false, deck.Shuffled)
		assert.Equal(t, entity.Card{Value: "ACE", Suit: "SPADES", Code: "AS"}, deck.Cards[0])
	})

	t.Run("open empty deck", func(t *testing.T) {
		deck, err := deckService.CreateDeck(nil, false)
		assert.NoError(t, err)

		_, err = deckService.DrawCards(deck.ID, 52, deck.Version)
		assert.NoError(t, err)

		_, err = deckService.OpenDeck(deck.ID)
		assert.EqualError(t, err, "not enough cards in deck")
	})
}

func TestDeckDraw(t *testing.T) {
	store := memorystore.NewInMemoryStore()
	t.Run("full deck, no shuffle, draw a card", func(t *testing.T) {
		deckService := poker.NewDeckService(store)
		deck, err := deckService.CreateDeck(nil, false)
		assert.NoError(t, err)
		version := deck.Version

		cards, err := deckService.DrawCards(deck.ID, 1, deck.Version)
		assert.NoError(t, err)

		assert.Equal(t, entity.Card{Value: "ACE", Suit: "SPADES", Code: "AS"}, cards[0])

		deck, err = deckService.OpenDeck(deck.ID)
		assert.NoError(t, err)

		assert.Equal(t, 51, deck.Size)
		assert.Equal(t, version+1, deck.Version)
	})
	t.Run("draw 2 cards in one size deck", func(t *testing.T) {
		deckService := poker.NewDeckService(store)
		deck, err := deckService.CreateDeck([]string{"AS"}, false)
		assert.NoError(t, err)

		_, err = deckService.DrawCards(deck.ID, 2, deck.Version)

		assert.EqualError(t, err, "not enough cards in deck")
	})

	t.Run("draw cards with old deck version", func(t *testing.T) {
		deckService := poker.NewDeckService(store)
		deck, err := deckService.CreateDeck([]string{"AS"}, false)
		assert.NoError(t, err)

		_, err = deckService.DrawCards(deck.ID, 2, deck.Version-1)

		assert.EqualError(t, err, "invalid deck version")
	})
}
