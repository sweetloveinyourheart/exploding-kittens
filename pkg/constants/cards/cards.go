package cards

const (
	ExplodingKitten    = "exploding_kitten"
	Defuse             = "defuse"
	Nope               = "nope"
	Attack             = "attack"
	Skip               = "skip"
	Favor              = "favor"
	Shuffle            = "shuffle"
	SeeTheFuture       = "see_the_future"
	TacoCat            = "taco_cat"
	Catermelon         = "catermelon"
	HairyPotatoCat     = "hairy_potato_cat"
	RainbowRalphingCat = "rainbow_ralphing_cat"
	BeardCat           = "beard_cat"
)

// Cards that must be played alone
var MustPlayAlone = map[string]bool{
	Attack:       true,
	Skip:         true,
	Favor:        true,
	Shuffle:      true,
	SeeTheFuture: true,
	Nope:         true, // special case: reaction, not turn-based play
}

// Combo cards
var ComboCards = map[string]bool{
	TacoCat:            true,
	Catermelon:         true,
	HairyPotatoCat:     true,
	RainbowRalphingCat: true,
	BeardCat:           true,
}
