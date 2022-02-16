CREATE TABLE IF NOT EXISTS albums(
    album_id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name VARCHAR(255) NOT NULL,
    total_tracks smallint,
    image_url VARCHAR(3000),
    spotify_id VARCHAR(62) UNIQUE,
    spotify_url VARCHAR(3000) UNIQUE,
    apple_url VARCHAR(3000) UNIQUE,
    artist_id integer REFERENCES artists (artist_id)
);