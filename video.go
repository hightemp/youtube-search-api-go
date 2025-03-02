package youtubesearchapi

import (
	"encoding/json"
	"fmt"
	"io"
	"regexp"
)

func GetSuggestData(limit int) (map[string]interface{}, error) {
	initData, err := GetYoutubeInitData(YoutubeEndpoint)
	if err != nil {
		return nil, fmt.Errorf("error getting recommendation data: %v", err)
	}

	result := make(map[string]interface{})
	items := []map[string]interface{}{}

	if contents, ok := initData.Initdata["contents"].(map[string]interface{}); ok {
		if twoColumnBrowseResultsRenderer, ok := contents["twoColumnBrowseResultsRenderer"].(map[string]interface{}); ok {
			if tabs, ok := twoColumnBrowseResultsRenderer["tabs"].([]interface{}); ok && len(tabs) > 0 {
				if tab, ok := tabs[0].(map[string]interface{}); ok {
					if tabRenderer, ok := tab["tabRenderer"].(map[string]interface{}); ok {
						if content, ok := tabRenderer["content"].(map[string]interface{}); ok {
							if richGridRenderer, ok := content["richGridRenderer"].(map[string]interface{}); ok {
								if contents, ok := richGridRenderer["contents"].([]interface{}); ok {
									for _, item := range contents {
										if richItemRenderer, ok := item.(map[string]interface{}); ok {
											if content, ok := richItemRenderer["content"].(map[string]interface{}); ok {
												if videoRenderer, ok := content["videoRenderer"].(map[string]interface{}); ok {
													videoItem := ExtractVideoData(videoRenderer)
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

	if limit > 0 && len(items) > limit {
		items = items[:limit]
	}

	result["items"] = items

	return result, nil
}

func GetVideoDetails(videoId string) (map[string]interface{}, error) {
	endpoint := fmt.Sprintf("%s/watch?v=%s", YoutubeEndpoint, videoId)

	initData, err := GetYoutubeInitData(endpoint)
	if err != nil {
		return nil, fmt.Errorf("error getting video data: %v", err)
	}

	playerData, err := GetYoutubePlayerDetail(endpoint)
	if err != nil {
		return nil, fmt.Errorf("error getting player data: %v", err)
	}

	result := make(map[string]interface{})

	if contents, ok := initData.Initdata["contents"].(map[string]interface{}); ok {
		if twoColumnWatchNextResults, ok := contents["twoColumnWatchNextResults"].(map[string]interface{}); ok {
			if results, ok := twoColumnWatchNextResults["results"].(map[string]interface{}); ok {
				if resultsContents, ok := results["results"].(map[string]interface{}); ok {
					if contents, ok := resultsContents["contents"].([]interface{}); ok && len(contents) > 1 {
						if videoPrimaryInfoRenderer, ok := contents[0].(map[string]interface{}); ok {
							if primaryInfo, ok := videoPrimaryInfoRenderer["videoPrimaryInfoRenderer"].(map[string]interface{}); ok {
								result["title"] = ExtractTextFromRuns(primaryInfo["title"])
								result["isLive"] = IsLiveVideo(primaryInfo["viewCount"])
							}
						}
						if videoSecondaryInfoRenderer, ok := contents[1].(map[string]interface{}); ok {
							if secondaryInfo, ok := videoSecondaryInfoRenderer["videoSecondaryInfoRenderer"].(map[string]interface{}); ok {
								if owner, ok := secondaryInfo["owner"].(map[string]interface{}); ok {
									if videoOwnerRenderer, ok := owner["videoOwnerRenderer"].(map[string]interface{}); ok {
										result["channel"] = ExtractTextFromRuns(videoOwnerRenderer["title"])
									}
								}
							}
						}
					}
				}
			}
		}
	}

	result["id"] = playerData["videoId"]
	result["thumbnail"] = playerData["thumbnail"]
	result["channelId"] = playerData["channelId"]
	result["description"] = playerData["shortDescription"]
	result["keywords"] = playerData["keywords"]
	suggestions := ExtractSuggestions(initData.Initdata)
	result["suggestion"] = suggestions
	return result, nil
}

func GetYoutubePlayerDetail(url string) (map[string]interface{}, error) {
	resp, err := HttpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error executing request: %v", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %v", err)
	}
	bodyStr := string(body)
	ytInitialPlayerResponse := regexp.MustCompile(`var ytInitialPlayerResponse = (.+?);</script>`).FindStringSubmatch(bodyStr)
	if len(ytInitialPlayerResponse) < 2 {
		return nil, fmt.Errorf("failed to find ytInitialPlayerResponse")
	}
	var playerData map[string]interface{}
	err = json.Unmarshal([]byte(ytInitialPlayerResponse[1]), &playerData)
	if err != nil {
		return nil, fmt.Errorf("error parsing ytInitialPlayerResponse: %v", err)
	}
	result := make(map[string]interface{})
	if videoDetails, ok := playerData["videoDetails"].(map[string]interface{}); ok {
		result["videoId"] = videoDetails["videoId"]
		result["title"] = videoDetails["title"]
		result["channelId"] = videoDetails["channelId"]
		result["shortDescription"] = videoDetails["shortDescription"]
		result["thumbnail"] = videoDetails["thumbnail"]
		result["keywords"] = videoDetails["keywords"]
	}
	return result, nil
}

func GetShortVideo() ([]map[string]interface{}, error) {
	initData, err := GetYoutubeInitData(YoutubeEndpoint)
	if err != nil {
		return nil, fmt.Errorf("error getting short video data: %v", err)
	}
	var shortVideos []map[string]interface{}
	if contents, ok := initData.Initdata["contents"].(map[string]interface{}); ok {
		if twoColumnBrowseResultsRenderer, ok := contents["twoColumnBrowseResultsRenderer"].(map[string]interface{}); ok {
			if tabs, ok := twoColumnBrowseResultsRenderer["tabs"].([]interface{}); ok && len(tabs) > 0 {
				if tab, ok := tabs[0].(map[string]interface{}); ok {
					if tabRenderer, ok := tab["tabRenderer"].(map[string]interface{}); ok {
						if content, ok := tabRenderer["content"].(map[string]interface{}); ok {
							if richGridRenderer, ok := content["richGridRenderer"].(map[string]interface{}); ok {
								if contents, ok := richGridRenderer["contents"].([]interface{}); ok {
									for _, item := range contents {
										if richSectionRenderer, ok := item.(map[string]interface{})["richSectionRenderer"]; ok {
											if content, ok := richSectionRenderer.(map[string]interface{})["content"]; ok {
												if richShelfRenderer, ok := content.(map[string]interface{})["richShelfRenderer"]; ok {
													if title, ok := richShelfRenderer.(map[string]interface{})["title"].(map[string]interface{}); ok {
														if runs, ok := title["runs"].([]interface{}); ok && len(runs) > 0 {
															if run, ok := runs[0].(map[string]interface{}); ok {
																if text, ok := run["text"].(string); ok && text == "Shorts" {
																	if shelfContents, ok := richShelfRenderer.(map[string]interface{})["contents"].([]interface{}); ok {
																		for _, shelfItem := range shelfContents {
																			if richItemRenderer, ok := shelfItem.(map[string]interface{})["richItemRenderer"]; ok {
																				if content, ok := richItemRenderer.(map[string]interface{})["content"]; ok {
																					if reelItemRenderer, ok := content.(map[string]interface{})["reelItemRenderer"]; ok {
																						shortVideo := ExtractShortVideoData(reelItemRenderer.(map[string]interface{}))
																						shortVideos = append(shortVideos, shortVideo)
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
				}
			}
		}
	}
	return shortVideos, nil
}

func ExtractVideoData(videoRenderer map[string]interface{}) map[string]interface{} {
	video := make(map[string]interface{})
	if id, ok := videoRenderer["videoId"].(string); ok {
		video["id"] = id
	}
	video["type"] = "video"
	if thumbnail, ok := videoRenderer["thumbnail"].(map[string]interface{}); ok {
		video["thumbnail"] = thumbnail
	}
	if title, ok := videoRenderer["title"].(map[string]interface{}); ok {
		if runs, ok := title["runs"].([]interface{}); ok && len(runs) > 0 {
			if run, ok := runs[0].(map[string]interface{}); ok {
				if text, ok := run["text"].(string); ok {
					video["title"] = text
				}
			}
		}
	}
	if ownerText, ok := videoRenderer["ownerText"].(map[string]interface{}); ok {
		if runs, ok := ownerText["runs"].([]interface{}); ok && len(runs) > 0 {
			if run, ok := runs[0].(map[string]interface{}); ok {
				if text, ok := run["text"].(string); ok {
					video["channelTitle"] = text
				}
			}
		}
	}
	if lengthText, ok := videoRenderer["lengthText"].(map[string]interface{}); ok {
		if simpleText, ok := lengthText["simpleText"].(string); ok {
			video["length"] = simpleText
		}
	}
	video["isLive"] = false
	if badges, ok := videoRenderer["badges"].([]interface{}); ok {
		for _, badge := range badges {
			if badgeMap, ok := badge.(map[string]interface{}); ok {
				if metadataBadgeRenderer, ok := badgeMap["metadataBadgeRenderer"].(map[string]interface{}); ok {
					if style, ok := metadataBadgeRenderer["style"].(string); ok && style == "BADGE_STYLE_TYPE_LIVE_NOW" {
						video["isLive"] = true
						break
					}
				}
			}
		}
	}
	return video
}

func ExtractShortVideoData(reelItemRenderer map[string]interface{}) map[string]interface{} {
	shortVideo := make(map[string]interface{})
	shortVideo["id"] = reelItemRenderer["videoId"]
	shortVideo["type"] = "reel"
	if thumbnail, ok := reelItemRenderer["thumbnail"].(map[string]interface{}); ok {
		if thumbnails, ok := thumbnail["thumbnails"].([]interface{}); ok && len(thumbnails) > 0 {
			shortVideo["thumbnail"] = thumbnails[0]
		}
	}
	if headline, ok := reelItemRenderer["headline"].(map[string]interface{}); ok {
		shortVideo["title"] = headline["simpleText"]
	}
	if inlinePlaybackEndpoint, ok := reelItemRenderer["inlinePlaybackEndpoint"].(map[string]interface{}); ok {
		shortVideo["inlinePlaybackEndpoint"] = inlinePlaybackEndpoint
	} else {
		shortVideo["inlinePlaybackEndpoint"] = map[string]interface{}{}
	}
	return shortVideo
}

func ExtractSuggestions(initdata map[string]interface{}) []map[string]interface{} {
	var suggestions []map[string]interface{}
	if contents, ok := initdata["contents"].(map[string]interface{}); ok {
		if twoColumnWatchNextResults, ok := contents["twoColumnWatchNextResults"].(map[string]interface{}); ok {
			if secondaryResults, ok := twoColumnWatchNextResults["secondaryResults"].(map[string]interface{}); ok {
				if secondaryResultsRenderer, ok := secondaryResults["secondaryResults"].(map[string]interface{}); ok {
					if results, ok := secondaryResultsRenderer["results"].([]interface{}); ok {
						for _, result := range results {
							if compactVideoRenderer, ok := result.(map[string]interface{})["compactVideoRenderer"]; ok {
								suggestion := ExtractCompactVideoRenderer(compactVideoRenderer.(map[string]interface{}))
								suggestions = append(suggestions, suggestion)
							}
						}
					}
				}
			}
		}
	}
	return suggestions
}

func ExtractCompactVideoRenderer(compactVideoRenderer map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	result["id"] = compactVideoRenderer["videoId"]
	result["type"] = "video"
	if thumbnail, ok := compactVideoRenderer["thumbnail"].(map[string]interface{}); ok {
		result["thumbnail"] = thumbnail["thumbnails"]
	}
	result["title"] = ExtractTextFromRuns(compactVideoRenderer["title"])
	result["channelTitle"] = ExtractTextFromRuns(compactVideoRenderer["shortBylineText"])
	result["shortBylineText"] = ExtractTextFromRuns(compactVideoRenderer["shortBylineText"])
	if lengthText, ok := compactVideoRenderer["lengthText"].(map[string]interface{}); ok {
		result["length"] = lengthText["simpleText"]
	}
	result["isLive"] = false
	if badges, ok := compactVideoRenderer["badges"].([]interface{}); ok {
		for _, badge := range badges {
			if badgeMap, ok := badge.(map[string]interface{}); ok {
				if metadataBadgeRenderer, ok := badgeMap["metadataBadgeRenderer"].(map[string]interface{}); ok {
					if style, ok := metadataBadgeRenderer["style"].(string); ok && style == "BADGE_STYLE_TYPE_LIVE_NOW" {
						result["isLive"] = true
						break
					}
				}
			}
		}
	}
	return result
}
