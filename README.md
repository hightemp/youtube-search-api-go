# youtube-search-api-go

A Go package for interacting with YouTube's search functionality and retrieving video, channel, and playlist data.
Based on https://github.com/damonwonghv/youtube-search-api.

## Features

- Search for videos, channels, and playlists
- Get video details
- Get channel information
- Get playlist data
- Retrieve short video data
- Get suggested videos

## Usage

Install

```bash
go get -u github.com/yourusername/youtube-search-api-go
```

Import the package in your Go code:

```go
import "github.com/yourusername/youtube-search-api-go/pkg/youtubesearchapi"
```

Use the provided functions to interact with YouTube data. For example:

```go
result, err := youtubesearchapi.GetData("golang tutorial", true, 5, nil)
if err != nil {
    // Handle error
}
// Process the result
```

## License

This project is licensed under the MIT License.