package poker

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/QQsha/cards/entity"
	"github.com/QQsha/cards/repository"
	"github.com/google/uuid"
)

func NewDeckService(rep repository.StoreRepository) *DeckService {
	return &DeckService{repository: rep}
}

type DeckService struct {
	repository repository.StoreRepository
}

var (
	errEmptyDeck       = errors.New("not enough cards in deck")
	errInvalidCardCode = errors.New("invalid card code")
	errInvalidVersion  = errors.New("invalid deck version")
)

const fullDeckSize = 52

var (
	values = []string{"ACE", "2", "3", "4", "5", "6", "7", "8", "9", "10", "JACK", "QUEEN", "KING"}
	suits  = []string{"SPADES", "DIAMONDS", "CLUBS", "HEARTS"}
	codes  = map[string]string{
		"ACE": "A", "2": "2", "3": "3", "4": "4", "5": "5", "6": "6", "7": "7",
		"8": "8", "9": "9", "10": "T", "JACK": "J", "QUEEN": "Q", "KING": "K",
		"SPADES": "S", "DIAMONDS": "D", "CLUBS": "C", "HEARTS": "H",
	}
	revCodes = map[string]string{
		"A": "ACE", "2": "2", "3": "3", "4": "4", "5": "5", "6": "6", "7": "7",
		"8": "8", "9": "9", "T": "10", "J": "JACK", "Q": "QUEEN", "K": "KING",
		"S": "SPADES", "D": "DIAMONDS", "C": "CLUBS", "H": "HEARTS",
	}
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func (d *DeckService) CreateDeck(cards []string, shuffle bool) (entity.Deck, error) {
	deck := entity.Deck{
		ID:       uuid.New().String(),
		Shuffled: shuffle,
	}

	if len(cards) > 0 {
		deck.Cards = make([]entity.Card, len(cards))

		for i, c := range cards {
			if len(c) != 2 {
				return entity.Deck{}, errInvalidCardCode
			}
			cardValue, ok := revCodes[string(c[0])]
			if !ok {
				return entity.Deck{}, errInvalidCardCode
			}
			cardSuit, ok := revCodes[string(c[1])]
			if !ok {
				return entity.Deck{}, errInvalidCardCode
			}

			deck.Cards[i] = entity.Card{Value: cardValue, Suit: cardSuit, Code: c}
		}

	} else {
		deck.Cards = make([]entity.Card, fullDeckSize)

		for i, s := range suits {
			for j, v := range values {
				deck.Cards[(i*fullDeckSize/4)+j] = entity.Card{
					Value: v,
					Suit:  s,
					Code:  fmt.Sprintf("%s%s", codes[v], codes[s]),
				}
			}
		}
	}

	if shuffle {
		rand.Shuffle(len(deck.Cards), func(i, j int) {
			deck.Cards[i], deck.Cards[j] = deck.Cards[j], deck.Cards[i]
		})
	}

	deck.Size = len(deck.Cards)

	d.repository.SaveDeck(deck)
	return deck, nil
}

func (d *DeckService) OpenDeck(id string) (entity.Deck, error) {
	deck, err := d.repository.GetDeck(id)
	if err != nil {
		return entity.Deck{}, err
	}

	if deck.Size == 0 {
		return entity.Deck{}, errEmptyDeck
	}

	return deck, nil
}

func (d *DeckService) DrawCards(deckID string, count, version int) ([]entity.Card, error) {
	deck, err := d.repository.GetDeck(deckID)
	if err != nil {
		return nil, err
	}

	if deck.Version != version {
		return nil, errInvalidVersion
	}

	if deck.Size < count {
		return nil, errEmptyDeck
	}

	cards := deck.Cards[:count]
	deck.Cards = deck.Cards[count:]
	deck.Size -= count

	deck.Version++

	d.repository.SaveDeck(deck)
	return cards, nil
}
