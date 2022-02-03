CREATE TABLE IF NOT EXISTS tracks(
    track_id serial PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    duration_ms integer,
    isrc VARCHAR(255) UNIQUE,
    spotify_id VARCHAR(62) UNIQUE,
    spotify_url VARCHAR(3000),
    apple_url VARCHAR(3000) UNIQUE,
    album_id integer REFERENCES albums (album_id),
    artist_id integer REFERENCES artists (artist_id)
);