CREATE TABLE IF NOT EXISTS playlists(
    playlist_id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name VARCHAR(2000) NOT NULL,
    spotify_id VARCHAR(255),
    spotify_url VARCHAR(3000),
    apple_id    VARCHAR(255),
    apple_url VARCHAR(3000),
    user_id INTEGER NOT NULL REFERENCES users(user_id),
    UNIQUE (user_id, name),
    UNIQUE (user_id, spotify_id),
    UNIQUE (user_id, spotify_url),
    UNIQUE (user_id, apple_id),
    UNIQUE (user_id, apple_url)
);

CREATE TABLE IF NOT EXISTS playlist_tracks(
    playlist_id integer NOT NULL REFERENCES playlists(playlist_id),
    track_id integer NOT NULL REFERENCES tracks(track_id),
    in_spotify bool NOT NULL DEFAULT false,
    in_apple bool NOT NULL DEFAULT false,
    UNIQUE (playlist_id, track_id)
);