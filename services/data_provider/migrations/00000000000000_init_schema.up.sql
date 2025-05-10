-- CARD DEFINITION
CREATE TABLE cards (
    card_id         UUID                        DEFAULT gen_random_uuid(),
    name            VARCHAR(50)     NOT NULL,
    code            VARCHAR(50)     NOT NULL,
    description     TEXT,
    quantity        INT             NOT NULL,
    created_at      TIMESTAMP WITH  TIME ZONE   DEFAULT now(),   
    updated_at      TIMESTAMP WITH  TIME ZONE   DEFAULT now(),
    
    PRIMARY KEY (card_id)
);

CREATE TABLE card_effects (
    effect_id       UUID                        DEFAULT gen_random_uuid(),
    card_id         UUID    NOT NULL,
    effect          JSONB   NOT NULL,
    created_at      TIMESTAMP WITH  TIME ZONE   DEFAULT now(),   
    updated_at      TIMESTAMP WITH  TIME ZONE   DEFAULT now(),

    PRIMARY KEY (effect_id),
    FOREIGN KEY (card_id) REFERENCES cards(card_id) ON DELETE CASCADE
);

CREATE TABLE combo_effects (
    combo_id        UUID                            DEFAULT gen_random_uuid(),
    required_cards  INT             NOT NULL,
    effect          JSONB           NOT NULL,
    created_at      TIMESTAMP WITH  TIME ZONE       DEFAULT now(),   
    updated_at      TIMESTAMP WITH  TIME ZONE       DEFAULT now(),

    PRIMARY KEY (combo_id)
);

CREATE TABLE card_combo (
    card_id     UUID    NOT NULL,
    combo_id    UUID    NOT NULL,

    PRIMARY KEY (card_id, combo_id),
    FOREIGN KEY (card_id) REFERENCES cards(card_id) ON DELETE CASCADE,
    FOREIGN KEY (combo_id) REFERENCES combo_effects(combo_id) ON DELETE CASCADE
);

INSERT INTO cards (card_id, name, code, description, quantity) VALUES
    (gen_random_uuid(), 'Exploding Kitten', 'exploding_kitten', 'If you draw one, you explode and are out of the game unless you have a Defuse card.', 4),
    (gen_random_uuid(), 'Defuse', 'defuse', 'Allows you to prevent an explosion and secretly place the Exploding Kitten back into the deck.', 6),
    (gen_random_uuid(), 'Nope', 'nope', 'Cancels any action (except Exploding Kitten or Defuse).', 5),
    (gen_random_uuid(), 'Attack', 'attack', 'Ends your turn without drawing and forces the next player to take two turns in a row.', 4),
    (gen_random_uuid(), 'Skip', 'skip', 'Ends your turn immediately without drawing a card.', 4),
    (gen_random_uuid(), 'Favor', 'favor', 'Forces another player to give you a card of their choice.', 4),
    (gen_random_uuid(), 'Shuffle', 'shuffle', 'Shuffles the deck.', 4),
    (gen_random_uuid(), 'See the Future', 'see_the_future', 'Lets you peek at the top three cards of the deck.', 5),
    (gen_random_uuid(), 'TacoCat', 'taco_cat', 'Part of the Cat Cards.', 4),
    (gen_random_uuid(), 'Catermelon', 'catermelon', 'Part of the Cat Cards.', 4),
    (gen_random_uuid(), 'Hairy Potato Cat', 'hairy_potato_cat', 'Part of the Cat Cards.', 4),
    (gen_random_uuid(), 'Rainbow Ralphing Cat', 'rainbow_ralphing_cat', 'Part of the Cat Cards.', 4),
    (gen_random_uuid(), 'Beard Cat', 'beard_cat', 'Part of the Cat Cards.', 4);

INSERT INTO card_effects (effect_id, card_id, effect) VALUES
    (gen_random_uuid(), (SELECT card_id FROM cards WHERE code = 'exploding_kitten'), '{"type": "explode"}'),
    (gen_random_uuid(), (SELECT card_id FROM cards WHERE code = 'defuse'), '{"type": "prevent_explode"}'),
    (gen_random_uuid(), (SELECT card_id FROM cards WHERE code = 'nope'), '{"type": "cancel_action"}'),
    (gen_random_uuid(), (SELECT card_id FROM cards WHERE code = 'attack'), '{"type": "skip_turn_and_attack"}'),
    (gen_random_uuid(), (SELECT card_id FROM cards WHERE code = 'skip'), '{"type": "skip_turn"}'),
    (gen_random_uuid(), (SELECT card_id FROM cards WHERE code = 'favor'), '{"type": "steal_card"}'),
    (gen_random_uuid(), (SELECT card_id FROM cards WHERE code = 'shuffle'), '{"type": "shuffle_deck"}'),
    (gen_random_uuid(), (SELECT card_id FROM cards WHERE code = 'see_the_future'), '{"type": "peek_cards"}');

INSERT INTO combo_effects (combo_id, required_cards, effect) VALUES
    (gen_random_uuid(), 2, '{"type": "steal_random_card"}'), -- Two of a Kind
    (gen_random_uuid(), 3, '{"type": "steal_named_card"}');  -- Three of a Kind

-- Two of a Kind (All Cat Cards)
INSERT INTO card_combo (card_id, combo_id)
SELECT card_id, (SELECT combo_id FROM combo_effects WHERE required_cards = 2)
FROM cards
WHERE code IN ('taco_cat', 'catermelon', 'hairy_potato_cat', 'rainbow_ralphing_cat', 'beard_cat');

-- Three of a Kind (All Cat Cards)
INSERT INTO card_combo (card_id, combo_id)
SELECT card_id, (SELECT combo_id FROM combo_effects WHERE required_cards = 3)
FROM cards
WHERE code IN ('taco_cat', 'catermelon', 'hairy_potato_cat', 'rainbow_ralphing_cat', 'beard_cat');
