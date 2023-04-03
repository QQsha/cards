package repository

import "github.com/QQsha/cards/entity"

type StoreRepository interface {
	SaveDeck(entity.Deck)
	GetDeck(string) (entity.Deck, error)
}
