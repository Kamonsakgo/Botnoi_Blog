package repositories

import (
	"go-fiber-unittest/domain/entities"

	"github.com/stretchr/testify/mock"
)

type MockIBlogRepository struct {
	mock.Mock
}

func (m *MockIBlogRepository) FindAll(page int, limit int) ([]entities.BlogModel, int64, error) {
	args := m.Called(page, limit)
	return args.Get(0).([]entities.BlogModel), args.Get(1).(int64), args.Error(2)
}

func (m *MockIBlogRepository) FindOne(blogid string) (data entities.BlogModel, err error) {
	args := m.Called(blogid)
	return args.Get(0).(entities.BlogModel), args.Error(1)
}

func (m *MockIBlogRepository) UploadImage(blogid string, name string, url string) error {
	args := m.Called(blogid, name, url)
	return args.Error(0)
}

func (m *MockIBlogRepository) Delete(blogid string) error {
	args := m.Called(blogid)
	return args.Error(0)
}

func (m *MockIBlogRepository) Update(data entities.BlogModel, blog_id string) error {
	args := m.Called(data, blog_id)
	return args.Error(0)
}

func (m *MockIBlogRepository) GetNextBlogID() (int, error) {
	args := m.Called()
	return args.Int(0), args.Error(1)
}

func (m *MockIBlogRepository) Insert(data entities.BlogModel) error {
	args := m.Called(data)
	return args.Error(0)
}
