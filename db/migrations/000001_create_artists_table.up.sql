CREATE TABLE IF NOT EXISTS artists(
    artist_id serial PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    popularity smallint,
    image_url VARCHAR(3000),
    spotify_id VARCHAR(62) UNIQUE NOT NULL,
    spotify_url VARCHAR(3000)
);
