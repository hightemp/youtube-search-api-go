package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	youtubesearchapi "github.com/hightemp/youtube-search-api-go"
)

func main() {
	proxyURL := os.Getenv("HTTP_PROXY")
	if proxyURL == "" {
		fmt.Println("HTTP_PROXY is not set. Using direct connection.")
	} else {
		fmt.Printf("Using proxy server: %s\n", proxyURL)
	}

	result, err := youtubesearchapi.GetData("golang tutorial", true, 5, nil)
	if err != nil {
		log.Fatalf("Error getting data: %v", err)
	}
	printJSON("GetData Result", result)

	videoDetails, err := youtubesearchapi.GetVideoDetails("dQw4w9WgXcQ")
	if err != nil {
		log.Fatalf("Error getting video details: %v", err)
	}
	printJSON("GetVideoDetails Result", videoDetails)

	playlistData, err := youtubesearchapi.GetPlaylistData("PL0Zuz27SZ-6Oi6xNtL_fwCrwpuqylMsgT", 5)
	if err != nil {
		log.Fatalf("Error getting playlist data: %v", err)
	}
	printJSON("GetPlaylistData Result", playlistData)

	suggestData, err := youtubesearchapi.GetSuggestData(5)
	if err != nil {
		log.Fatalf("Error getting suggestions: %v", err)
	}
	printJSON("GetSuggestData Result", suggestData)

	channelData, err := youtubesearchapi.GetChannelById("UCCezIgC97PvUuR4_gbFUs5g")
	if err != nil {
		log.Fatalf("Error getting channel data: %v", err)
	}
	printJSON("GetChannelById Result", channelData)

	shortVideos, err := youtubesearchapi.GetShortVideo()
	if err != nil {
		log.Fatalf("Error getting short videos: %v", err)
	}
	printJSON("GetShortVideo Result", shortVideos)
}

func printJSON(title string, data interface{}) {
	fmt.Printf("\n%s:\n", title)
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling JSON: %v", err)
	}
	fmt.Println(string(jsonData))
}
