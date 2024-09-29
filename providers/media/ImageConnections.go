package media

// ENPOINT FOR IMAGES

// InsertUserAvatar Inserts a new User avatar to provider and returns the url
func (m *Media) InsertUserAvatar(content []byte, filename string) (string, error) {
	return storeProfileImages(content, filename)
}

// InsertGroupAvatar Insert a new Group avatar to provider and return the url
func (m *Media) InsertGroupAvatar(content []byte, filename string) (string, error) {
	return storeProfileImages(content, filename)
}

func (m *Media) InsetImages(images [][]byte, filenames []string) (ImageResponse, error) {
	return storeImagesMedia(images, filenames)
}
