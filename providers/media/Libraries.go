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
	ErrFailedCreatingLibrary       = errors.New("failed creating library")
	ErrFailedUpdatingLibrary       = errors.New("failed updating library")
	ErrFailedDeletingLibrary       = errors.New("failed deleting library")
	ErrFailedCreatingVideoFile     = errors.New("failed creating video file")
	ErrFailedUploadingVideoContent = errors.New("failed uploading video content")
	ErrFailedGettingVideoData      = errors.New("failed getting video data")
)

// createVideoLibrary
func createVideoLibrary(libraryName string) (LibraryResponse, error) {

	alog := StartLogger()

	var payload struct {
		Name string `json:"Name"`
	}

	var response LibraryResponse

	payload.Name = libraryName

	body, err := json.Marshal(payload)
	if err != nil {
		alog.ErrorLog(err.Error())
		return response, err
	}

	url := os.Getenv("BASE_LIBRARY_URL")

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		alog.ErrorLog(err.Error())
		return response, err
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("AccessKey", os.Getenv("MEDIA_PKEY"))

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		alog.ErrorLog(err.Error())
		return response, err
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		alog.ErrorLog(err.Error())
		return response, err
	}

	if res.StatusCode != http.StatusCreated {
		return response, ErrFailedCreatingLibrary
	}

	return response, nil
}

func updateLibrary(libraryID int) error {

	alog := StartLogger()

	var payload struct {
		BlockNoneReferrer                 bool   `json:"BlockNoneReferrer"`
		EnableMP4Fallback                 bool   `json:"EnableMP4Fallback"`
		KeepOriginalFiles                 bool   `json:"KeepOriginalFiles"`
		AllowDirectPlay                   bool   `json:"AllowDirectPlay"`
		EnableTranscribing                bool   `json:"EnableTranscribing"`
		EnableTranscribingTitleGeneration bool   `json:"EnableTranscribingTitleGeneration"`
		EnabledResolutions                string `json:"EnabledResolutions"`
	}
	payload.BlockNoneReferrer = false
	payload.EnableMP4Fallback = false
	payload.KeepOriginalFiles = false
	payload.AllowDirectPlay = true
	payload.EnableTranscribing = false
	payload.EnableTranscribingTitleGeneration = false
	payload.EnabledResolutions = "480p, 1080p"

	body, err := json.Marshal(payload)
	if err != nil {
		alog.ErrorLog(err.Error())
		return err
	}

	url := fmt.Sprintf("%s/%d", os.Getenv("BASE_LIBRARY_URL"), libraryID)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		alog.ErrorLog(err.Error())
		return err
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("AccessKey", os.Getenv("MEDIA_PKEY"))

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		alog.ErrorLog(err.Error())
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return ErrFailedUpdatingLibrary
	}

	return nil
}

func deleteLibrary(libraryID int, API_KEY string) error {

	alog := StartLogger()

	url := fmt.Sprintf("%s/%d", os.Getenv("BASE_LIBRARY_URL"), libraryID)

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		alog.ErrorLog()
		return err
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("AccessKey", API_KEY)

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		alog.ErrorLog(err.Error())
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusNoContent {
		alog.ErrorLog(ErrFailedDeletingLibrary.Error())
		return ErrFailedDeletingLibrary
	}

	return nil
}

func createVideoFile(title, API_KEY string, LibraryID int) (string, error) {

	//
	alog := StartLogger()

	var payload struct {
		Title string `json:"title"`
	}

	payload.Title = title

	body, err := json.Marshal(payload)
	if err != nil {
		alog.ErrorLog(err.Error())
		return "", err
	}

	url := fmt.Sprintf("%s/%d/videos", os.Getenv("BASE_VIDEO_URL"), LibraryID)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		alog.ErrorLog(err.Error())
		return "", err
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("AccessKey", API_KEY)

	client := http.Client{}

	res, err := client.Do(req)
	if err != nil {
		alog.ErrorLog(err.Error())
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", ErrFailedCreatingVideoFile
	}

	var response struct {
		GID string `json:"guid"`
	}

	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		alog.ErrorLog(err.Error())
		return "", err
	}

	return response.GID, nil
}

func uploadVideoContent(LibraryID int, VideoID string, API_KEY string, content []byte) error {

	alog := StartLogger()

	url := fmt.Sprintf("%s/%d/videos/%s", os.Getenv("BASE_VIDEO_URL"), LibraryID, VideoID)

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(content))
	if err != nil {
		alog.ErrorLog(err.Error())
		return err
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("AccessKey", API_KEY)

	client := http.Client{}

	res, err := client.Do(req)
	if err != nil {
		alog.ErrorLog(err.Error())
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return ErrFailedUploadingVideoContent
	}

	return nil

}

func getVideoPlayData(LibraryID int, VideoID string, API_KEY string) (VideoPlayback, error) {

	alog := StartLogger()

	var result VideoPlayback

	url := fmt.Sprintf("%s/%d/videos/%s/play", os.Getenv("BASE_VIDEO_URL"), LibraryID, VideoID)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		alog.ErrorLog(err.Error())
		return result, err
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("AccessKey", API_KEY)

	client := http.Client{}

	res, err := client.Do(req)
	if err != nil {
		alog.ErrorLog(err.Error())
		return result, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return result, ErrFailedGettingVideoData
	}

	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		return result, err
	}

	result.GUID = fmt.Sprintf("%d$%s", LibraryID, VideoID)

	return result, nil
}
