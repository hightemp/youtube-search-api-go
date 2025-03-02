package youtubesearchapi

import (
	"net/http"
	"net/url"
	"os"
	"time"
)

var HttpClient *http.Client

func init() {
	proxyURL, _ := url.Parse(os.Getenv("HTTP_PROXY"))
	HttpClient = &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		},
	}
}
