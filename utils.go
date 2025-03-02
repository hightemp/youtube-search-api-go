package youtubesearchapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
)

// YoutubeInitData определена в types.go
// httpClient определен в client.go

func GetYoutubeInitData(url string) (*YoutubeInitData, error) {
	resp, err := HttpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("ошибка при выполнении запроса: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ошибка при чтении ответа: %v", err)
	}

	bodyStr := string(body)

	initData := &YoutubeInitData{}

	// Извлечение ytInitialData
	ytInitDataRegex := regexp.MustCompile(`var ytInitialData = (.+?);</script>`)
	matches := ytInitDataRegex.FindStringSubmatch(bodyStr)
	if len(matches) > 1 {
		err = json.Unmarshal([]byte(matches[1]), &initData.Initdata)
		if err != nil {
			return nil, fmt.Errorf("ошибка при разборе ytInitialData: %v", err)
		}
	} else {
		return nil, fmt.Errorf("не удалось найти ytInitialData")
	}

	// Извлечение innertubeApiKey
	apiKeyRegex := regexp.MustCompile(`"innertubeApiKey":"(.+?)"`)
	matches = apiKeyRegex.FindStringSubmatch(bodyStr)
	if len(matches) > 1 {
		initData.APIToken = matches[1]
	}

	// Извлечение INNERTUBE_CONTEXT
	contextRegex := regexp.MustCompile(`INNERTUBE_CONTEXT":(.+?)}},`)
	matches = contextRegex.FindStringSubmatch(bodyStr)
	if len(matches) > 1 {
		err = json.Unmarshal([]byte(matches[1]+`}}`), &initData.Context)
		if err != nil {
			return nil, fmt.Errorf("ошибка при разборе INNERTUBE_CONTEXT: %v", err)
		}
	}

	return initData, nil
}

func ExtractTextFromRuns(data interface{}) string {
	if runs, ok := data.(map[string]interface{})["runs"].([]interface{}); ok {
		var texts []string
		for _, run := range runs {
			if runMap, ok := run.(map[string]interface{}); ok {
				if text, ok := runMap["text"].(string); ok {
					texts = append(texts, text)
				}
			}
		}
		return strings.Join(texts, "")
	}
	return ""
}

func IsLiveVideo(viewCount interface{}) bool {
	if viewCountRenderer, ok := viewCount.(map[string]interface{}); ok {
		if _, ok := viewCountRenderer["isLive"]; ok {
			return true
		}
	}
	return false
}

func ExtractContinuationToken(initdata map[string]interface{}) string {
	if contents, ok := initdata["contents"].(map[string]interface{}); ok {
		if twoColumnSearchResultsRenderer, ok := contents["twoColumnSearchResultsRenderer"].(map[string]interface{}); ok {
			if primaryContents, ok := twoColumnSearchResultsRenderer["primaryContents"].(map[string]interface{}); ok {
				if sectionListRenderer, ok := primaryContents["sectionListRenderer"].(map[string]interface{}); ok {
					if contentsList, ok := sectionListRenderer["contents"].([]interface{}); ok {
						for _, content := range contentsList {
							if contentMap, ok := content.(map[string]interface{}); ok {
								if continuationItemRenderer, ok := contentMap["continuationItemRenderer"].(map[string]interface{}); ok {
									if continuationEndpoint, ok := continuationItemRenderer["continuationEndpoint"].(map[string]interface{}); ok {
										if continuationCommand, ok := continuationEndpoint["continuationCommand"].(map[string]interface{}); ok {
											if token, ok := continuationCommand["token"].(string); ok {
												return token
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
	return ""
}
