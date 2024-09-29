package media

// ImageRequest basic model that represents the request the main server sends to media provider
type ImageRequest struct {
	UserID      string   `json:"user_id"`
	TargetID    string   `json:"target_id"`
	ContentID   string   `json:"content_id"`
	MediaSource [][]byte `json:"media_content"`
}

// VideoRequest basic model that represent the request the main server sends to the media provider
type VideoRequest struct {
	UserID      string `json:"user_id"`
	TargetID    string `json:"target_id"`
	ContentID   string `json:"content_id"`
	MediaSource []byte `json:"media_content"`
}

// VideoResponse basic model that represents the response of the video operations
type VideoResponse struct {
	ContentID   string `json:"content_id"`
	MediaSource string `json:"media_content"`
}

// ImageResponse basic model that represents the response of the image operations
type ImageResponse struct {
	ContentID   string   `json:"content_id"`
	MediaSource []string `json:"media_source"`
	Thumbnails  []string `json:"thumbnails"`
}

type LibraryResponse struct {
	Id     int    `json:"Id"`
	ApiKey string `json:"ApiKey"`
}

type VideoPlayback struct {
	GUID      string `json:"guid"`
	Thumbnail string `json:"thumbnailUrl"`
	Src       string `json:"videoPlaylistUrl"`
}
