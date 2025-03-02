package youtubesearchapi

const YoutubeEndpoint = "https://www.youtube.com"

type YoutubeInitData struct {
	Initdata map[string]interface{}
	APIToken string
	Context  map[string]interface{}
}

type Channel struct {
	ID        string
	Type      string
	Thumbnail map[string]interface{}
	Title     string
}

type Tab struct {
	Title   string
	Content map[string]interface{}
}

type Video struct {
	ID           string
	Type         string
	Thumbnail    map[string]interface{}
	Title        string
	ChannelTitle string
	Length       string
	IsLive       bool
}

type Playlist struct {
	ID         string
	Type       string
	Thumbnail  []interface{}
	Title      string
	Length     string
	VideoCount string
	IsLive     bool
}

type SearchResult struct {
	Items    []interface{}
	NextPage map[string]interface{}
}

type VideoDetails struct {
	ID          string
	Title       string
	Channel     string
	Description string
	Thumbnail   map[string]interface{}
	Keywords    []string
	Suggestion  []interface{}
}

type PlaylistData struct {
	Items    []interface{}
	Metadata map[string]interface{}
}

type SuggestData struct {
	Items []interface{}
}
