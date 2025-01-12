
-- USERS --

CREATE TABLE public.users
(
    user_id         UUID                        NOT NULL,
    username        VARCHAR(255)                NOT NULL,
    full_name       VARCHAR(255)                NOT NULL,
    status          INT                         NOT NULL,
    created_at      TIMESTAMP WITH TIME ZONE    NOT NULL,
    updated_at      TIMESTAMP WITH TIME ZONE    NOT NULL,
    PRIMARY KEY (user_id)
);

-- USER CREDENTIALS --

CREATE TABLE public.user_credentials
(
    user_id         UUID                        NOT NULL,
    auth_provider   VARCHAR(255)                NOT NULL,
    meta            JSONB,
    created_at      TIMESTAMP WITH TIME ZONE    NOT NULL,
    updated_at      TIMESTAMP WITH TIME ZONE    NOT NULL,
    PRIMARY KEY (user_id),
    CONSTRAINT fk_user_credentials_users FOREIGN KEY (user_id) REFERENCES users (user_id)
);

-- USER SESSIONS --

CREATE TABLE public.user_sessions
(
    session_id          BIGSERIAL                   NOT NULL,
    user_id             UUID                        NOT NULL,
    token               VARCHAR(255)                NOT NULL,
    session_start       TIMESTAMP WITH TIME ZONE    NOT NULL,
    last_updated        TIMESTAMP WITH TIME ZONE    NOT NULL,
    session_expiration  TIMESTAMP WITH TIME ZONE,
    session_end         TIMESTAMP WITH TIME ZONE,
    PRIMARY KEY (session_id),
    CONSTRAINT fk_user_sessions_users FOREIGN KEY (user_id) REFERENCES users (user_id)
);

CREATE INDEX ON user_sessions (token);