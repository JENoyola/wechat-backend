package media

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
)

var (
	ErrNoInserted = errors.New("content not inserted")
)

// DIRECT REQUEST TO IMAGES AND FILES

// storeProfileImages creates a request to store an avatar to storage
func storeProfileImages(file []byte, filename string) (string, error) {

	alog := StartLogger()

	url := fmt.Sprintf("%s/%s/%s", os.Getenv("BASE_URL"), os.Getenv("PROFILES_PATH"), filename)

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(file))
	if err != nil {
		alog.ErrorLog(err.Error())
		return "", err
	}

	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("AccessKey", os.Getenv("PROFILES_AUTH"))
	req.Header.Set("accept", "application/json")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		alog.ErrorLog(err.Error())
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		return "", ErrNoInserted
	}

	return fmt.Sprintf("%s/%s", os.Getenv("PROFILES_URL"), filename), nil
}

// storeImagesMedia creates a request to store an image to storage
func storeImagesMedia(files [][]byte, filename []string) (ImageResponse, error) {

	alog := StartLogger()

	var response ImageResponse

	for i, file := range files {

		url := fmt.Sprintf("%s/%s/%s", os.Getenv("BASE_URL"), os.Getenv("CONTENT_PATH"), filename[i])

		req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(file))
		if err != nil {
			alog.ErrorLog(err.Error())
			return response, err
		}

		req.Header.Set("Content-Type", "application/octet-stream")
		req.Header.Set("AccessKey", os.Getenv("CONTENT_AUTH"))
		req.Header.Set("accept", "application/json")

		client := http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			alog.ErrorLog(err.Error())
			return response, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			return response, ErrNoInserted
		}

		response.MediaSource = append(response.MediaSource, fmt.Sprintf("%s/%s", os.Getenv("CONTENT_URL"), filename[i]))
	}

	return response, nil

}

// storeFile creates a request to store a file to storage
func storeFile(file []byte, filename string) (string, error) {

	alog := StartLogger()

	url := fmt.Sprintf("%s/%s/%s", os.Getenv("BASE_URL"), os.Getenv("FILES_PATH"), filename)

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(file))
	if err != nil {
		alog.ErrorLog(err.Error())
		return "", err
	}

	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("AccessKey", os.Getenv("FILES_AUTH"))
	req.Header.Set("accept", "application/json")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		alog.ErrorLog(err.Error())
		return "", err
	}
	defer resp.Body.Close()

	var response struct {
		HttpCode int
	}

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		alog.ErrorLog(err.Error())
		return "", err
	}

	return fmt.Sprintf("%s/%s", os.Getenv("FILES_URL"), filename), nil
}
