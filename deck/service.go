package deck

import "github.com/QQsha/cards/entity"

type Service interface {
	CreateDeck([]string, bool) (entity.Deck, error)
	OpenDeck(string) (entity.Deck, error)
	DrawCards(string, int, int) ([]entity.Card, error)
}
