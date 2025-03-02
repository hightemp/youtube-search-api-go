package youtubesearchapi

const YoutubeEndpoint = "https://www.youtube.com"

type YoutubeInitData struct {
	Initdata map[string]interface{}
	APIToken string
	Context  map[string]interface{}
}
