package youtubesearchapi

import (
	"fmt"
)

func GetPlaylistData(playlistId string, limit int) (map[string]interface{}, error) {
	endpoint := fmt.Sprintf("%s/playlist?list=%s", YoutubeEndpoint, playlistId)

	initData, err := GetYoutubeInitData(endpoint)
	if err != nil {
		return nil, fmt.Errorf("error getting playlist data: %v", err)
	}

	result := make(map[string]interface{})
	items := []map[string]interface{}{}

	if contents, ok := initData.Initdata["contents"].(map[string]interface{}); ok {
		if twoColumnBrowseResultsRenderer, ok := contents["twoColumnBrowseResultsRenderer"].(map[string]interface{}); ok {
			if tabs, ok := twoColumnBrowseResultsRenderer["tabs"].([]interface{}); ok && len(tabs) > 0 {
				if tab, ok := tabs[0].(map[string]interface{}); ok {
					if tabRenderer, ok := tab["tabRenderer"].(map[string]interface{}); ok {
						if content, ok := tabRenderer["content"].(map[string]interface{}); ok {
							if sectionListRenderer, ok := content["sectionListRenderer"].(map[string]interface{}); ok {
								if sectionContents, ok := sectionListRenderer["contents"].([]interface{}); ok && len(sectionContents) > 0 {
									if itemSection, ok := sectionContents[0].(map[string]interface{}); ok {
										if itemSectionRenderer, ok := itemSection["itemSectionRenderer"].(map[string]interface{}); ok {
											if sectionContents, ok := itemSectionRenderer["contents"].([]interface{}); ok && len(sectionContents) > 0 {
												if playlistVideoListRenderer, ok := sectionContents[0].(map[string]interface{}); ok {
													if contents, ok := playlistVideoListRenderer["playlistVideoListRenderer"].(map[string]interface{}); ok {
														if videoItems, ok := contents["contents"].([]interface{}); ok {
															for _, item := range videoItems {
																if videoRenderer, ok := item.(map[string]interface{}); ok {
																	if playlistVideoRenderer, ok := videoRenderer["playlistVideoRenderer"].(map[string]interface{}); ok {
																		videoItem := ExtractVideoData(playlistVideoRenderer)
																		items = append(items, videoItem)
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
					}
				}
			}
		}
	}

	if limit > 0 && len(items) > limit {
		items = items[:limit]
	}

	result["items"] = items
	result["metadata"] = extractPlaylistMetadata(initData.Initdata)

	return result, nil
}

func ExtractPlaylistData(playlistRenderer map[string]interface{}) map[string]interface{} {
	playlist := make(map[string]interface{})

	if id, ok := playlistRenderer["playlistId"].(string); ok {
		playlist["id"] = id
	}
	playlist["type"] = "playlist"

	if thumbnail, ok := playlistRenderer["thumbnails"].([]interface{}); ok {
		playlist["thumbnail"] = thumbnail
	}

	if title, ok := playlistRenderer["title"].(map[string]interface{}); ok {
		if simpleText, ok := title["simpleText"].(string); ok {
			playlist["title"] = simpleText
		}
	}

	if videoCount, ok := playlistRenderer["videoCount"].(string); ok {
		playlist["length"] = videoCount
		playlist["videoCount"] = videoCount
	}

	playlist["isLive"] = false

	return playlist
}

func extractPlaylistMetadata(initdata map[string]interface{}) map[string]interface{} {
	metadata := make(map[string]interface{})

	if microformat, ok := initdata["microformat"].(map[string]interface{}); ok {
		if microformatDataRenderer, ok := microformat["microformatDataRenderer"].(map[string]interface{}); ok {
			if title, ok := microformatDataRenderer["title"].(string); ok {
				metadata["title"] = title
			}
			if description, ok := microformatDataRenderer["description"].(string); ok {
				metadata["description"] = description
			}
			if thumbnail, ok := microformatDataRenderer["thumbnail"].(map[string]interface{}); ok {
				if thumbnails, ok := thumbnail["thumbnails"].([]interface{}); ok && len(thumbnails) > 0 {
					metadata["thumbnail"] = thumbnails[0]
				}
			}
		}
	}

	return metadata
}
