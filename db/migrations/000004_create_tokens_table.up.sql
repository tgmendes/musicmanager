CREATE TABLE IF NOT EXISTS tokens(
    token_id serial PRIMARY KEY,
    spotify_access_token VARCHAR(500) NOT NULL UNIQUE,
    spotify_refresh_token VARCHAR(500) NOT NULL UNIQUE,
    spotify_user_id VARCHAR(255) NOT NULL UNIQUE,
    apple_access_token VARCHAR(500) NOT NULL UNIQUE,
    apple_refresh_token VARCHAR(500) NOT NULL UNIQUE,
    apple_user_id VARCHAR(255) NOT NULL UNIQUE
);