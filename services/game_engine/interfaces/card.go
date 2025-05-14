package interfaces

type CardEffect struct {
	Type string `json:"type"`
}

type CardComboEffect struct {
	Type          string `json:"type"`
	RequiredCards int    `json:"required_cards"`
}
