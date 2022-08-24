package dto

type CardDTO struct {
	Id      int64    `json:"id"`
	Number  string   `json:"number"`
	Issuer  string   `json:"issuer"`
	Type    CardType `json:"type"`
	Balance int64    `json:"balance"`
}

type CardType string

const (
	Virtual    CardType = "VIRTUAL"
	Additional CardType = "ADDITIONAL"
)

type NewCardDTO struct {
	Id     int64    `json:"id"`
	Type   CardType `json:"type"`
	Issuer string   `json:"issuer"`
}
