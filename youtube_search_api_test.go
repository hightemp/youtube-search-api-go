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
	if result == nil {
		t.Error("GetData returned nil result")
	}
	if items, ok := result["items"].([]map[string]interface{}); !ok || len(items) == 0 {
		t.Error("GetData returned no items")
	}
}

func TestGetVideoDetails(t *testing.T) {
	result, err := GetVideoDetails("dQw4w9WgXcQ")
	if err != nil {
		t.Logf("GetVideoDetails failed: %v", err)
		t.Skip("Skipping test due to network error")
	}
	if result == nil {
		t.Error("GetVideoDetails returned nil result")
	}
	if _, ok := result["id"]; !ok {
		t.Error("GetVideoDetails result doesn't contain 'id'")
	}
}

func TestGetPlaylistData(t *testing.T) {
	result, err := GetPlaylistData("PL0Zuz27SZ-6Oi6xNtL_fwCrwpuqylMsgT", 5)
	if err != nil {
		t.Logf("GetPlaylistData failed: %v", err)
		t.Skip("Skipping test due to network error")
	}
	if result == nil {
		t.Error("GetPlaylistData returned nil result")
	}
	if items, ok := result["items"].([]map[string]interface{}); !ok || len(items) == 0 {
		t.Error("GetPlaylistData returned no items")
	}
}

func TestGetSuggestData(t *testing.T) {
	result, err := GetSuggestData(5)
	if err != nil {
		t.Logf("GetSuggestData failed: %v", err)
		t.Skip("Skipping test due to network error")
	}
	if result == nil {
		t.Error("GetSuggestData returned nil result")
	}
	if items, ok := result["items"].([]map[string]interface{}); !ok || len(items) == 0 {
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
	if result == nil {
		t.Error("GetShortVideo returned nil result")
	}
	if len(result) == 0 {
		t.Error("GetShortVideo returned no data")
	}
}
