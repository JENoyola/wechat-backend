package media

// ENDPOINT FOR VIDEO

func (m *Media) StoreVideo(groupID, filename string, content []byte) (VideoPlayback, error) {
	return storeVideo(groupID, filename, content)
}

func (m *Media) DeleteLibraryContent(libraryID int, API_KEY string) error {
	return deleteLibraryData(libraryID, API_KEY)
}
