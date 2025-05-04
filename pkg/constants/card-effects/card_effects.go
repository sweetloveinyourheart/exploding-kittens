package card_effects

const (
	// Card effects
	Explore           = "explode"
	PreventExplore    = "prevent_explode"
	CancelAction      = "cancel_action"
	SkipTurnAndAttack = "skip_turn_and_attack"
	SkipTurn          = "skip_turn"
	StealCard         = "steal_card"
	ShuffleDesk       = "shuffle_deck"
	PeekCards         = "peek_cards"

	// Combo effects
	StealRandomCard = "steal_random_card"
	StealNamedCard  = "steal_named_card"
	TradeAnyDiscard = "trade_any_discard"
)

const (
	AttackBonusCount = 1
)

var AllCardEffects = []string{
	Explore,
	PreventExplore,
	CancelAction,
	SkipTurnAndAttack,
	SkipTurn,
	StealCard,
	ShuffleDesk,
	PeekCards,
	StealRandomCard,
	StealNamedCard,
	TradeAnyDiscard,
}
