CREATE TABLE IF NOT EXISTS artists(
    artist_id serial PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    image_url VARCHAR(3000),
    spotify_id VARCHAR(62) UNIQUE,
    spotify_url VARCHAR(3000) UNIQUE,
    apple_url VARCHAR(3000) UNIQUE
);
