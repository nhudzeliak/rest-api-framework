CREATE TABLE users (
    id              SERIAL          PRIMARY KEY
    , fullname      VARCHAR(255)    NOT NULL
    , updated_at    TIMESTAMP       DEFAULT CURRENT_TIMESTAMP
    , created_at    TIMESTAMP       DEFAULT CURRENT_TIMESTAMP
);
