package youtubesearchapi

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	code := m.Run()

	os.Exit(code)
}

func TestGetData(t *testing.T) {
	result, err := GetData("golang tutorial", true, 5, nil)
	if err != nil {
		t.Logf("GetData failed: %v", err)
		t.Skip("Skipping test due to network error")
	}
	if len(result.Items) == 0 {
		t.Error("GetData returned no items")
	}
	if result.NextPage == nil {
		t.Error("GetData returned nil NextPage")
	}
}

func TestGetVideoDetails(t *testing.T) {
	result, err := GetVideoDetails("dQw4w9WgXcQ")
	if err != nil {
		t.Logf("GetVideoDetails failed: %v", err)
		t.Skip("Skipping test due to network error")
	}
	if result.ID == "" {
		t.Error("GetVideoDetails returned empty ID")
	}
	if result.Title == "" {
		t.Error("GetVideoDetails returned empty Title")
	}
	if result.Channel == "" {
		t.Error("GetVideoDetails returned empty Channel")
	}
	if result.Description == "" {
		t.Error("GetVideoDetails returned empty Description")
	}
	if result.Thumbnail == nil {
		t.Error("GetVideoDetails returned nil Thumbnail")
	}
	if len(result.Keywords) == 0 {
		t.Error("GetVideoDetails returned no Keywords")
	}
	if len(result.Suggestion) == 0 {
		t.Error("GetVideoDetails returned no Suggestions")
	}
}

func TestGetPlaylistData(t *testing.T) {
	result, err := GetPlaylistData("PL0Zuz27SZ-6Oi6xNtL_fwCrwpuqylMsgT", 5)
	if err != nil {
		t.Logf("GetPlaylistData failed: %v", err)
		t.Skip("Skipping test due to network error")
	}
	if len(result.Items) == 0 {
		t.Error("GetPlaylistData returned no items")
	}
	if result.Metadata == nil {
		t.Error("GetPlaylistData returned nil metadata")
	}
}

func TestGetSuggestData(t *testing.T) {
	result, err := GetSuggestData(5)
	if err != nil {
		t.Logf("GetSuggestData failed: %v", err)
		t.Skip("Skipping test due to network error")
	}
	if len(result.Items) == 0 {
		t.Error("GetSuggestData returned no items")
	}
}

func TestGetChannelById(t *testing.T) {
	result, err := GetChannelById("UCCezIgC97PvUuR4_gbFUs5g")
	if err != nil {
		t.Logf("GetChannelById failed: %v", err)
		t.Skip("Skipping test due to network error")
	}
	if result == nil {
		t.Error("GetChannelById returned nil result")
	}
	if len(result) == 0 {
		t.Error("GetChannelById returned no data")
	}
}

func TestGetShortVideo(t *testing.T) {
	result, err := GetShortVideo()
	if err != nil {
		t.Logf("GetShortVideo failed: %v", err)
		t.Skip("Skipping test due to network error")
	}
	if len(result) == 0 {
		t.Error("GetShortVideo returned no data")
	}
	for _, video := range result {
		if video.Type != "short" {
			t.Errorf("Expected video type 'short', got '%s'", video.Type)
		}
		if video.ID == "" {
			t.Error("GetShortVideo returned a video with empty ID")
		}
		if video.Title == "" {
			t.Error("GetShortVideo returned a video with empty Title")
		}
	}
}
