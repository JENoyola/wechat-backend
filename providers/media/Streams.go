package media

// storeVideo stores a new video to the provider
func storeVideo(groupID, fileTitle string, videoContent []byte) (VideoPlayback, error) {

	var res VideoPlayback

	// create library
	library, err := createVideoLibrary(groupID)
	if err != nil {
		return res, err
	}

	// update library
	err = updateLibrary(library.Id)
	if err != nil {
		return res, err
	}
	// create video file
	videoID, err := createVideoFile(fileTitle, library.ApiKey, library.Id)
	if err != nil {
		return res, err
	}
	// upload video file
	err = uploadVideoContent(library.Id, videoID, library.ApiKey, videoContent)
	if err != nil {
		return res, err
	}

	res, err = getVideoPlayData(library.Id, videoID, library.ApiKey)
	if err != nil {
		return res, err
	}

	return res, nil
}

// deleteLibraryData Deletes the whole library and its contents
func deleteLibraryData(LibraryID int, API_KEY string) error {
	return deleteLibrary(LibraryID, API_KEY)
}
