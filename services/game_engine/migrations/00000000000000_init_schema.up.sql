-- CARD DEFINITION
CREATE TABLE card_types (
    id              SERIAL,
    name            VARCHAR(50)     NOT NULL    UNIQUE,
    description     TEXT,
    created_at      TIMESTAMP WITH TIME ZONE    DEFAULT now(),
    updated_at      TIMESTAMP WITH TIME ZONE    DEFAULT now(),

    PRIMARY KEY (id)
);

CREATE TABLE cards (
    id              UUID                        DEFAULT gen_random_uuid(),
    type_id         INT             NOT NULL,
    name            VARCHAR(50)     NOT NULL,
    description     TEXT,
    effect          JSONB,
    quantity        INT             NOT NULL,
    created_at      TIMESTAMP WITH TIME ZONE    DEFAULT now(),   
    updated_at      TIMESTAMP WITH TIME ZONE    DEFAULT now(),
    
    PRIMARY KEY(id),
    FOREIGN KEY (type_id) REFERENCES card_types(id) ON DELETE CASCADE
);

INSERT INTO card_types (name, description) VALUES
('Action', 'Cards with special effects that change gameplay.'),
('Cat Card', 'Cards that do not have abilities unless played in pairs or combos.');

INSERT INTO cards (type_id, name, description, effect, quantity) VALUES
-- Action Cards
(1, 'Exploding Kitten', 'If you draw one, you explode and are out of the game unless you have a Defuse card.', '{"instant_death": true}', 4),
(1, 'Defuse', 'Prevents an Exploding Kitten from eliminating you.', '{"prevent_explode": true, "place_back": "secret"}', 6),
(1, 'Nope', 'Cancels any action except Exploding Kitten or Defuse.', '{"cancel_action": true}', 5),
(1, 'Attack', 'Ends your turn and forces the next player to take two turns.', '{"skip_turns": 1, "force_next_player": 2}', 4),
(1, 'Skip', 'Ends your turn immediately without drawing a card.', '{"skip_turns": 1}', 4),
(1, 'Favor', 'Forces another player to give you a card of their choice.', '{"steal_card": 1}', 4),
(1, 'Shuffle', 'Shuffles the deck.', '{"shuffle_deck": true}', 4),
(1, 'See the Future', 'Peek at the top three cards of the deck.', '{"peek_cards": 3}', 5),

-- Cat Cards (No Special Ability)
(2, 'TacoCat', 'A mysterious cat with no ability unless played in combos.', NULL, 4),
(2, 'Catermelon', 'A watermelon-cat hybrid with no standalone ability.', NULL, 4),
(2, 'Hairy Potato Cat', 'A fluffy potato disguised as a cat.', NULL, 4),
(2, 'Rainbow Ralphing Cat', 'A cat that barfs rainbows. No effect alone.', NULL, 4),
(2, 'Beard Cat', 'A cat with a majestic beard. Used in combos.', NULL, 4);
