package youtubesearchapi

import (
	"fmt"
)

func GetChannelById(channelId string) ([]map[string]interface{}, error) {
	endpoint := fmt.Sprintf("%s/channel/%s", YoutubeEndpoint, channelId)

	initData, err := GetYoutubeInitData(endpoint)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении данных канала: %v", err)
	}

	var tabs []map[string]interface{}

	if contents, ok := initData.Initdata["contents"].(map[string]interface{}); ok {
		if twoColumnBrowseResultsRenderer, ok := contents["twoColumnBrowseResultsRenderer"].(map[string]interface{}); ok {
			if tabsData, ok := twoColumnBrowseResultsRenderer["tabs"].([]interface{}); ok {
				for _, tab := range tabsData {
					if tabMap, ok := tab.(map[string]interface{}); ok {
						if tabRenderer, ok := tabMap["tabRenderer"].(map[string]interface{}); ok {
							tabInfo := make(map[string]interface{})
							if title, ok := tabRenderer["title"].(string); ok {
								tabInfo["title"] = title
							}
							if content, ok := tabRenderer["content"].(map[string]interface{}); ok {
								tabInfo["content"] = content
							}
							tabs = append(tabs, tabInfo)
						}
					}
				}
			}
		}
	}

	return tabs, nil
}

func ExtractChannelData(channelRenderer map[string]interface{}) map[string]interface{} {
	channel := make(map[string]interface{})

	if id, ok := channelRenderer["channelId"].(string); ok {
		channel["id"] = id
	}
	channel["type"] = "channel"

	if thumbnail, ok := channelRenderer["thumbnail"].(map[string]interface{}); ok {
		channel["thumbnail"] = thumbnail
	}

	if title, ok := channelRenderer["title"].(map[string]interface{}); ok {
		if simpleText, ok := title["simpleText"].(string); ok {
			channel["title"] = simpleText
		}
	}

	return channel
}
