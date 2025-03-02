package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	youtubesearchapi "github.com/hightemp/youtube-search-api-go"
)

func main() {
	// Пример установки прокси-сервера
	proxyURL := os.Getenv("HTTP_PROXY")
	if proxyURL == "" {
		fmt.Println("HTTP_PROXY не установлен. Используется прямое соединение.")
	} else {
		fmt.Printf("Используется прокси-сервер: %s\n", proxyURL)
	}

	// Вы также можете установить прокси программно:
	// os.Setenv("HTTP_PROXY", "http://your-proxy-server:port")
	// Пример использования GetData
	result, err := youtubesearchapi.GetData("golang tutorial", true, 5, nil)
	if err != nil {
		log.Fatalf("Ошибка при получении данных: %v", err)
	}
	printJSON("GetData Result", result)

	// Пример использования GetVideoDetails
	videoDetails, err := youtubesearchapi.GetVideoDetails("dQw4w9WgXcQ")
	if err != nil {
		log.Fatalf("Ошибка при получении деталей видео: %v", err)
	}
	printJSON("GetVideoDetails Result", videoDetails)

	// Пример использования GetPlaylistData
	playlistData, err := youtubesearchapi.GetPlaylistData("PL0Zuz27SZ-6Oi6xNtL_fwCrwpuqylMsgT", 5)
	if err != nil {
		log.Fatalf("Ошибка при получении данных плейлиста: %v", err)
	}
	printJSON("GetPlaylistData Result", playlistData)

	// Пример использования GetSuggestData
	suggestData, err := youtubesearchapi.GetSuggestData(5)
	if err != nil {
		log.Fatalf("Ошибка при получении рекомендаций: %v", err)
	}
	printJSON("GetSuggestData Result", suggestData)

	// Пример использования GetChannelById
	channelData, err := youtubesearchapi.GetChannelById("UCCezIgC97PvUuR4_gbFUs5g")
	if err != nil {
		log.Fatalf("Ошибка при получении данных канала: %v", err)
	}
	printJSON("GetChannelById Result", channelData)

	// Пример использования GetShortVideo
	shortVideos, err := youtubesearchapi.GetShortVideo()
	if err != nil {
		log.Fatalf("Ошибка при получении коротких видео: %v", err)
	}
	printJSON("GetShortVideo Result", shortVideos)
}

func printJSON(title string, data interface{}) {
	fmt.Printf("\n%s:\n", title)
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Fatalf("Ошибка при маршалинге JSON: %v", err)
	}
	fmt.Println(string(jsonData))
}
