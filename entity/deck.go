package entity

type Deck struct {
	ID       string `json:"deck_id"`
	Shuffled bool   `json:"shuffled"`
	Size     int    `json:"remaining"`
	Cards    []Card `json:"cards"`
	Version  int    `json:"version"`
}

type Card struct {
	Value string `json:"value"`
	Suit  string `json:"suit"`
	Code  string `json:"code"`
}
