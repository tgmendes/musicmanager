CREATE TYPE playlist_type AS ENUM ('short_term', 'medium_term', 'long_term', 'generic');
CREATE TABLE IF NOT EXISTS playlists(
    playlist_id serial PRIMARY KEY,
    name VARCHAR(2000) NOT NULL,
    internal_identifier VARCHAR(255) NOT NULL UNIQUE,
    playlist_type playlist_type NOT NULL,
    created_date DATE NOT NULL,
    spotify_id VARCHAR(255) UNIQUE,
    spotify_url VARCHAR(3000) UNIQUE,
    apple_id VARCHAR(255) UNIQUE,
    apple_url VARCHAR(3000) UNIQUE,
    user_id VARCHAR(255) NOT NULL UNIQUE,
    UNIQUE (name, user_id)
);