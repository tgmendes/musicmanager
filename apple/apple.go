package apple

type Artwork struct {
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Url    string `json:"url"`
}

type ResponseMeta struct {
	Total   int     `json:"total"`
	Filters Filters `json:"filters"`
}

type Filters struct {
	ISRC map[string][]SongMeta `json:"isrc"`
}
