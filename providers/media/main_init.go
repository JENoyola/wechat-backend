package media

type MediaMock struct {
	StoreVideoMockFunc        func(string, string, []byte) (VideoPlayback, error)
	InsertFileMockFunc        func([]byte, string) (string, error)
	InsertUserAvatarMockFunc  func([]byte, string) (string, error)
	InsertGroupAvatarMockFunc func([]byte, string) (string, error)
	InsetImagesMockFunc       func([][]byte, []string) (ImageResponse, error)
}

func (m *MediaMock) StoreVideo(groupID, filename string, content []byte) (VideoPlayback, error) {
	if m.StoreVideoMockFunc != nil {
		return m.StoreVideoMockFunc(groupID, filename, content)
	}
	return VideoPlayback{}, nil
}

func (m *MediaMock) InsertFile(content []byte, filename string) (string, error) {
	if m.InsertFileMockFunc != nil {
		return m.InsertFileMockFunc(content, filename)
	}
	return "", nil
}

func (m *MediaMock) InsertUserAvatar(content []byte, filename string) (string, error) {
	if m.InsertUserAvatarMockFunc != nil {
		return m.InsertUserAvatarMockFunc(content, filename)
	}
	return "", nil
}

func (m *MediaMock) InsertGroupAvatar(content []byte, filename string) (string, error) {
	if m.InsertGroupAvatarMockFunc != nil {
		return m.InsertGroupAvatar(content, filename)
	}
	return "", nil
}

func (m *MediaMock) InsetImages(images [][]byte, filenames []string) (ImageResponse, error) {
	if m.InsetImagesMockFunc != nil {
		return m.InsetImagesMockFunc(images, filenames)
	}
	return ImageResponse{}, nil
}
