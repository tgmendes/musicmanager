CREATE TABLE IF NOT EXISTS spotify_tokens(
    token_id serial PRIMARY KEY,
    access_token VARCHAR(500) NOT NULL UNIQUE,
    refresh_token VARCHAR(500) NOT NULL UNIQUE,
    user_id VARCHAR(255) NOT NULL UNIQUE
);