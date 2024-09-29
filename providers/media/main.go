package media

import (
	"log"

	"github.com/joho/godotenv"
)

type MediaHUB interface {
	StoreVideo(groupID, filename string, content []byte) (VideoPlayback, error)
	InsertFile(content []byte, filename string) (string, error)
	InsertUserAvatar(content []byte, filename string) (string, error)
	InsertGroupAvatar(content []byte, filename string) (string, error)
	InsetImages(images [][]byte, filenames []string) (ImageResponse, error)
}

type Media struct{}

func NewMediaService() (*Media, error) {

	err := godotenv.Load("services.env")
	if err != nil {
		log.Fatal(err.Error())
	}

	return &Media{}, nil
}
