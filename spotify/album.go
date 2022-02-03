package spotify

type Album struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Popularity  int         `json:"popularity"`
	TotalTracks int         `json:"total_tracks"`
	Images      []Image     `json:"images"`
	HRef        string      `json:"href"`
	Artists     []Artist    `json:"artists"`
	Tracks      []TrackItem `json:"tracks"`
}
