package main

import (
	"flag"
	"fmt"
	"log"

	youtubesearchapi "github.com/hightemp/youtube-search-api-go/pkg/youtubesearchapi"
)

func main() {
	query := flag.String("query", "", "Search query")
	limit := flag.Int("limit", 5, "Number of results to display")
	includePlaylist := flag.Bool("playlist", false, "Include playlists in search results")
	flag.Parse()

	if *query == "" {
		log.Fatal("Please specify a search query using the -query flag")
	}

	runCLI(*query, *limit, *includePlaylist)

	fmt.Println("Program completed.")
}

func runCLI(query string, limit int, includePlaylist bool) {
	result, err := youtubesearchapi.GetData(query, includePlaylist, limit, nil)
	if err != nil {
		log.Fatalf("Error searching for videos: %v", err)
	}

	fmt.Printf("Search results for query: %s\n\n", query)
	for _, item := range result.Items {
		switch v := item.(type) {
		case youtubesearchapi.Video:
			printVideo(v)
		case youtubesearchapi.Channel:
			printChannel(v)
		case youtubesearchapi.Playlist:
			printPlaylist(v)
		}
		fmt.Println("---")
	}
}

func printVideo(video youtubesearchapi.Video) {
	fmt.Printf("Type: Video\n")
	fmt.Printf("Title: %s\n", video.Title)
	fmt.Printf("Channel: %s\n", video.ChannelTitle)
	fmt.Printf("URL: https://www.youtube.com/watch?v=%s\n", video.ID)
	fmt.Printf("Length: %s\n", video.Length)
	if video.IsLive {
		fmt.Println("Live: Yes")
	}
}

func printChannel(channel youtubesearchapi.Channel) {
	fmt.Printf("Type: Channel\n")
	fmt.Printf("Name: %s\n", channel.Title)
	fmt.Printf("URL: https://www.youtube.com/channel/%s\n", channel.ID)
}

func printPlaylist(playlist youtubesearchapi.Playlist) {
	fmt.Printf("Type: Playlist\n")
	fmt.Printf("Title: %s\n", playlist.Title)
	fmt.Printf("URL: https://www.youtube.com/playlist?list=%s\n", playlist.ID)
	fmt.Printf("Video Count: %s\n", playlist.VideoCount)
	if playlist.IsLive {
		fmt.Println("Live: Yes")
	}
}
