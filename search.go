package youtubesearchapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strings"
)

func GetData(keyword string, withPlaylist bool, limit int, options []map[string]string) (map[string]interface{}, error) {
	endpoint := fmt.Sprintf("%s/results?search_query=%s", YoutubeEndpoint, url.QueryEscape(keyword))
	if len(options) > 0 {
		for _, opt := range options {
			if t, ok := opt["type"]; ok {
				switch strings.ToLower(t) {
				case "video":
					endpoint += "&sp=EgIQAQ%3D%3D"
				case "channel":
					endpoint += "&sp=EgIQAg%3D%3D"
				case "playlist":
					endpoint += "&sp=EgIQAw%3D%3D"
				case "movie":
					endpoint += "&sp=EgIQBA%3D%3D"
				}
				break
			}
		}
	}
	initData, err := GetYoutubeInitData(endpoint)
	if err != nil {
		return nil, err
	}
	result := make(map[string]interface{})
	items := []map[string]interface{}{}
	if contents, ok := initData.Initdata["contents"].(map[string]interface{}); ok {
		if twoColumnSearchResultsRenderer, ok := contents["twoColumnSearchResultsRenderer"].(map[string]interface{}); ok {
			if primaryContents, ok := twoColumnSearchResultsRenderer["primaryContents"].(map[string]interface{}); ok {
				if sectionListRenderer, ok := primaryContents["sectionListRenderer"].(map[string]interface{}); ok {
					if contentsList, ok := sectionListRenderer["contents"].([]interface{}); ok {
						for _, content := range contentsList {
							if contentMap, ok := content.(map[string]interface{}); ok {
								if itemSectionRenderer, ok := contentMap["itemSectionRenderer"].(map[string]interface{}); ok {
									if itemContents, ok := itemSectionRenderer["contents"].([]interface{}); ok {
										for _, item := range itemContents {
											if itemMap, ok := item.(map[string]interface{}); ok {
												if videoRenderer, ok := itemMap["videoRenderer"].(map[string]interface{}); ok {
													videoItem := ExtractVideoData(videoRenderer)
													items = append(items, videoItem)
												} else if channelRenderer, ok := itemMap["channelRenderer"].(map[string]interface{}); ok {
													channelItem := ExtractChannelData(channelRenderer)
													items = append(items, channelItem)
												} else if playlistRenderer, ok := itemMap["playlistRenderer"].(map[string]interface{}); ok && withPlaylist {
													playlistItem := ExtractPlaylistData(playlistRenderer)
													items = append(items, playlistItem)
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}

	if limit > 0 && len(items) > limit {
		items = items[:limit]
	}

	result["items"] = items
	result["nextPage"] = map[string]interface{}{
		"nextPageToken": initData.APIToken,
		"nextPageContext": map[string]interface{}{
			"context":      initData.Context,
			"continuation": ExtractContinuationToken(initData.Initdata),
		},
	}

	return result, nil
}

func NextPage(nextPage map[string]interface{}, withPlaylist bool, limit int) (map[string]interface{}, error) {
	nextPageToken, ok := nextPage["nextPageToken"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid nextPageToken format")
	}

	nextPageContext, ok := nextPage["nextPageContext"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid nextPageContext format")
	}

	endpoint := fmt.Sprintf("%s/youtubei/v1/search?key=%s", YoutubeEndpoint, nextPageToken)

	jsonData, err := json.Marshal(nextPageContext)
	if err != nil {
		return nil, fmt.Errorf("error marshaling nextPageContext: %v", err)
	}

	resp, err := HttpClient.Post(endpoint, "application/json", strings.NewReader(string(jsonData)))
	if err != nil {
		return nil, fmt.Errorf("error executing POST request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %v", err)
	}

	var responseData map[string]interface{}
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		return nil, fmt.Errorf("error parsing JSON response: %v", err)
	}

	items := []map[string]interface{}{}

	if onResponseReceivedCommands, ok := responseData["onResponseReceivedCommands"].([]interface{}); ok {
		for _, command := range onResponseReceivedCommands {
			if commandMap, ok := command.(map[string]interface{}); ok {
				if appendContinuationItemsAction, ok := commandMap["appendContinuationItemsAction"].(map[string]interface{}); ok {
					if continuationItems, ok := appendContinuationItemsAction["continuationItems"].([]interface{}); ok {
						for _, item := range continuationItems {
							if itemMap, ok := item.(map[string]interface{}); ok {
								if itemSectionRenderer, ok := itemMap["itemSectionRenderer"].(map[string]interface{}); ok {
									if contents, ok := itemSectionRenderer["contents"].([]interface{}); ok {
										for _, content := range contents {
											if contentMap, ok := content.(map[string]interface{}); ok {
												if videoRenderer, ok := contentMap["videoRenderer"].(map[string]interface{}); ok {
													videoItem := ExtractVideoData(videoRenderer)
													items = append(items, videoItem)
												} else if channelRenderer, ok := contentMap["channelRenderer"].(map[string]interface{}); ok {
													channelItem := ExtractChannelData(channelRenderer)
													items = append(items, channelItem)
												} else if playlistRenderer, ok := contentMap["playlistRenderer"].(map[string]interface{}); ok && withPlaylist {
													playlistItem := ExtractPlaylistData(playlistRenderer)
													items = append(items, playlistItem)
												}
											}
										}
									}
								} else if continuationItemRenderer, ok := itemMap["continuationItemRenderer"].(map[string]interface{}); ok {
									if continuationEndpoint, ok := continuationItemRenderer["continuationEndpoint"].(map[string]interface{}); ok {
										if continuationCommand, ok := continuationEndpoint["continuationCommand"].(map[string]interface{}); ok {
											if token, ok := continuationCommand["token"].(string); ok {
												nextPageContext["continuation"] = token
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}

	if limit > 0 && len(items) > limit {
		items = items[:limit]
	}

	result := map[string]interface{}{
		"items": items,
		"nextPage": map[string]interface{}{
			"nextPageToken":   nextPageToken,
			"nextPageContext": nextPageContext,
		},
	}

	return result, nil
}
