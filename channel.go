package youtubesearchapi

import (
	"fmt"
)

func GetChannelById(channelId string) ([]Tab, error) {
	endpoint := fmt.Sprintf("%s/channel/%s", YoutubeEndpoint, channelId)

	initData, err := GetYoutubeInitData(endpoint)
	if err != nil {
		return nil, fmt.Errorf("error getting channel data: %v", err)
	}

	var tabs []Tab

	if contents, ok := initData.Initdata["contents"].(map[string]interface{}); ok {
		if twoColumnBrowseResultsRenderer, ok := contents["twoColumnBrowseResultsRenderer"].(map[string]interface{}); ok {
			if tabsData, ok := twoColumnBrowseResultsRenderer["tabs"].([]interface{}); ok {
				for _, tab := range tabsData {
					if tabMap, ok := tab.(map[string]interface{}); ok {
						if tabRenderer, ok := tabMap["tabRenderer"].(map[string]interface{}); ok {
							var tabInfo Tab
							if title, ok := tabRenderer["title"].(string); ok {
								tabInfo.Title = title
							}
							if content, ok := tabRenderer["content"].(map[string]interface{}); ok {
								tabInfo.Content = content
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

func ExtractChannelData(channelRenderer map[string]interface{}) Channel {
	channel := Channel{
		Type: "channel",
	}

	if id, ok := channelRenderer["channelId"].(string); ok {
		channel.ID = id
	}

	if thumbnail, ok := channelRenderer["thumbnail"].(map[string]interface{}); ok {
		channel.Thumbnail = thumbnail
	}

	if title, ok := channelRenderer["title"].(map[string]interface{}); ok {
		if simpleText, ok := title["simpleText"].(string); ok {
			channel.Title = simpleText
		}
	}

	return channel
}
