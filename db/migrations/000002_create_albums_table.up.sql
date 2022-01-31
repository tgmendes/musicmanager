CREATE TABLE IF NOT EXISTS albums(
    album_id serial PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    popularity smallint,
    release_date DATE,
    total_tracks smallint,
    image_url VARCHAR(3000),
    spotify_id VARCHAR(62) UNIQUE NOT NULL,
    spotify_url VARCHAR(3000),
    apple_url VARCHAR(3000) UNIQUE NOT NULL,
    artist_id integer REFERENCES artists (artist_id)
);