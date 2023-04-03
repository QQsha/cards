package memorystore

import (
	"errors"
	"sync"

	"github.com/QQsha/cards/entity"
)

var errInvalidID = errors.New("invalid deck id")

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		map[string]entity.Deck{},
		sync.RWMutex{},
	}
}

type InMemoryStore struct {
	store map[string]entity.Deck
	lock  sync.RWMutex
}

func (i *InMemoryStore) SaveDeck(deck entity.Deck) {
	i.lock.Lock()
	defer i.lock.Unlock()
	i.store[deck.ID] = deck
}

func (i *InMemoryStore) GetDeck(id string) (entity.Deck, error) {
	i.lock.RLock()
	defer i.lock.RUnlock()

	deck, ok := i.store[id]
	if !ok {
		return entity.Deck{}, errInvalidID
	}

	return deck, nil
}
