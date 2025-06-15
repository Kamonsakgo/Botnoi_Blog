package repositories

import "github.com/stretchr/testify/mock"

type MockS3Provider struct {
	mock.Mock
}

func (m *MockS3Provider) UploadS3FromString(fileName []byte, keyName string, contentType string) (string, error) {
	args := m.Called(fileName, keyName, contentType)
	return args.String(0), args.Error(1)
}

func (m *MockS3Provider) HashString(userID string) (string, error) {
	args := m.Called(userID)
	return args.String(0), args.Error(1)
}

func (m *MockS3Provider) CreateKeyNameImageBlog(fileName string, typeFile string) (string, string) {
	args := m.Called(fileName, typeFile)
	return args.String(0), args.String(1)
}

func (m *MockS3Provider) DeleteKeyNameImage(url string) error {
	args := m.Called(url)
	return args.Error(0)
}

func (m *MockS3Provider) CreateKeyNameImageHL(typeFile string) (string, string) {
	args := m.Called(typeFile)
	return args.String(0), args.String(1)
}
func (m *MockS3Provider) DeleteKeyNameImageHL(keyName string) error {
	args := m.Called(keyName)
	return args.Error(0)
}
