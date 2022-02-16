package spotify

const (
	// ScopePlaylistReadPrivate seeks permission to read
	// a user's private playlists.
	ScopePlaylistReadPrivate = "playlist-read-private"
	// ScopePlaylistReadCollaborative seeks permission to read
	// a user's collaborative playlists.
	ScopePlaylistReadCollaborative = "playlist-read-collaborative"
	// ScopePlaylistModifyPrivate seeks permission to modify
	// a user's private playlists.
	ScopePlaylistModifyPrivate = "playlist-modify-private"
	// ScopeUserReadPrivate seeks read access to a user's
	// subscription details (type of user account).
	ScopeUserReadPrivate = "user-read-private"
	// ScopeUserReadEmail seeks read access to a user's email address.
	ScopeUserReadEmail = "user-read-email"
	// ScopeUserReadCurrentlyPlaying seeks read access to a user's currently playing track
	ScopeUserReadCurrentlyPlaying = "user-read-currently-playing"
	// ScopeUserReadRecentlyPlayed seeks read access to a user's recently played tracks
	ScopeUserReadRecentlyPlayed = "user-read-recently-played"
	// ScopeUserReadPlaybackState seeks read access to the user's current playback state
	ScopeUserReadPlaybackState = "user-read-playback-state"
	// ScopeUserTopRead seeks read access to a user's top tracks and artists
	ScopeUserTopRead = "user-top-read"
	// ScopeStreaming seeks permission to play music and control playback on your other devices.
	ScopeStreaming = "streaming"
)

func AllScopes() []string {
	return []string{
		ScopePlaylistReadPrivate,
		ScopePlaylistReadCollaborative,
		ScopePlaylistModifyPrivate,
		ScopeUserReadPrivate,
		ScopeUserReadEmail,
		ScopeUserReadCurrentlyPlaying,
		ScopeUserReadRecentlyPlayed,
		ScopeUserReadPlaybackState,
		ScopeUserTopRead,
		ScopeStreaming,
	}
}
