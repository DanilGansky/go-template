CREATE TABLE IF NOT EXISTS users (
    id         SERIAL NOT NULL PRIMARY KEY,
    username   VARCHAR NOT NULL UNIQUE,
    password   VARCHAR NOT NULL,
    joined_at  TIMESTAMP NOT NULL,
    last_login TIMESTAMP NOT NULL
);