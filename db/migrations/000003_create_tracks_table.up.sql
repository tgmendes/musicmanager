CREATE TABLE IF NOT EXISTS tracks(
    track_id serial PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    popularity smallint,
    release_date DATE,
    duration_ms integer,
    play_count integer DEFAULT 1,
    last_played TIMESTAMP,
    isrc VARCHAR(255),
    spotify_id VARCHAR(62) UNIQUE NOT NULL,
    spotify_url VARCHAR(3000),
    apple_url VARCHAR(3000) UNIQUE NOT NULL,
    album_id integer REFERENCES albums (album_id),
    artist_id integer REFERENCES artists (artist_id)
);