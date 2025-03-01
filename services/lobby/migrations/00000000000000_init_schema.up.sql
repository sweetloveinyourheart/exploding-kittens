
-- LOBBIES --

CREATE TABLE lobbies (
    lobby_id        UUID                        NOT NULL,
    lobby_code      VARCHAR(10)                 NOT NULL,
    lobby_name      VARCHAR(100)                NOT NULL,
    host_user_id    UUID                        NOT NULL,
    participants    JSONB                       NOT NULL    DEFAULT '[]', -- Array of user IDs
    created_at      TIMESTAMP WITH TIME ZONE    NOT NULL,
    updated_at      TIMESTAMP WITH TIME ZONE    NOT NULL,

    PRIMARY KEY (lobby_id)
);