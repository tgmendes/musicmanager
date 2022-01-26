package spotify

type TopTimeRange string

const (
	ShortTerm  TopTimeRange = "short_term"
	MediumTerm TopTimeRange = "medium_term"
	LongTerm   TopTimeRange = "long_term"
)

type Image struct {
	Url    string `json:"url"`
	Height int    `json:"height"`
	Width  int    `json:"width"`
}
