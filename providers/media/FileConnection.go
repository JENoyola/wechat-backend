package media

// InsertFile Inserts a new file to the provider and returs the url
func (m *Media) InsertFile(content []byte, filename string) (string, error) {
	return storeFile(content, filename)
}
