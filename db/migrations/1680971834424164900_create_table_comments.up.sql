CREATE TABLE comments (
    id              SERIAL          PRIMARY KEY
    , post_id       INTEGER         REFERENCES posts(id)
    , content       TEXT            NOT NULL
    , updated_at    TIMESTAMP       DEFAULT CURRENT_TIMESTAMP
    , created_at    TIMESTAMP       DEFAULT CURRENT_TIMESTAMP
);