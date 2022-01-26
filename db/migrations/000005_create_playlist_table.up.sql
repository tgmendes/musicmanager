CREATE TYPE playlist_type AS ENUM ('short_term', 'medium_term', 'long_term', 'generic');
CREATE TABLE IF NOT EXISTS playlists(
    playlist_id serial PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    internal_identifier VARCHAR(255) NOT NULL UNIQUE,
    playlist_type playlist_type NOT NULL,
    created_date DATE NOT NULL,
    spotify_id VARCHAR(255) NOT NULL UNIQUE,
    spotify_url VARCHAR(3000),
    user_id VARCHAR(255) NOT NULL UNIQUE,
    UNIQUE (name, user_id)
);