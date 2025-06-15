package services_test

import (
	"fmt"
	"go-fiber-unittest/domain/entities"
	"go-fiber-unittest/src/mock/repositories"
	"go-fiber-unittest/src/services"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockIBlogRepository struct {
	mock.Mock
}

var (
	mockBlogService = new(repositories.MockIBlogRepository)
	MockUserService = new(repositories.MockIUserRepository)
	mockS3Provider  = new(repositories.MockS3Provider)
	BlogSV          = services.NewBlogsService(mockBlogService, MockUserService, mockS3Provider)
)

func RefreshBlogSV() {
	mockBlogService = new(repositories.MockIBlogRepository)
	MockUserService = new(repositories.MockIUserRepository)
	mockS3Provider = new(repositories.MockS3Provider)
	BlogSV = services.NewBlogsService(mockBlogService, MockUserService, mockS3Provider)
}

func TestGetAll(t *testing.T) {
	t.Run("TestGetAll", func(t *testing.T) {
		RefreshBlogSV()

		mockBlogService.On("FindAll", 1, 10).Return([]entities.BlogModel{}, int64(1), fmt.Errorf("error"))

		_, err := BlogSV.GetAll("1", "10")
		assert.Equal(t, fmt.Errorf("error"), err)

	})
	t.Run("success", func(t *testing.T) {
		RefreshBlogSV()

		mockBlogService.On("FindAll", 1, 10).Return([]entities.BlogModel{}, int64(1), nil)

		_, err := BlogSV.GetAll("1", "10")
		assert.Equal(t, nil, err)

	})
}

func TestGetOne(t *testing.T) {
	t.Run("TestGetOne", func(t *testing.T) {
		RefreshBlogSV()

		mockBlogService.On("FindOne", "1").Return(entities.BlogModel{}, fmt.Errorf("error"))

		_, err := BlogSV.GetOne("1")
		assert.Equal(t, fmt.Errorf("error"), err)

	})
	t.Run("success", func(t *testing.T) {
		RefreshBlogSV()

		mockBlogService.On("FindOne", "1").Return(entities.BlogModel{}, nil)

		_, err := BlogSV.GetOne("1")
		assert.Equal(t, nil, err)

	})
}
func TestInsert(t *testing.T) {
	t.Run("TestFindByID", func(t *testing.T) {
		RefreshBlogSV()

		MockUserService.On("FindByID", mock.Anything).Return(&entities.UserDataFormat{}, fmt.Errorf("User not found"))

		err := BlogSV.Insert(entities.BlogModel{})
		assert.Equal(t, fmt.Errorf("User not found"), err)

	})
	t.Run("TestCheck_permission", func(t *testing.T) {
		RefreshBlogSV()

		MockUserService.On("FindByID", mock.Anything).Return(&entities.UserDataFormat{Role: "user"}, nil)
		err := BlogSV.Insert(entities.BlogModel{})
		assert.Equal(t, fmt.Errorf("you don't have permission"), err)

	})
	t.Run("TestGetNextBlogID_error", func(t *testing.T) {
		RefreshBlogSV()

		MockUserService.On("FindByID", mock.Anything).Return(&entities.UserDataFormat{Role: "admin"}, nil)
		mockBlogService.On("GetNextBlogID").Return(1, fmt.Errorf("error"))
		err := BlogSV.Insert(entities.BlogModel{})
		assert.Equal(t, fmt.Errorf("error"), err)

	})
	t.Run("TestInsert_error", func(t *testing.T) {
		RefreshBlogSV()

		MockUserService.On("FindByID", mock.Anything).Return(&entities.UserDataFormat{Role: "admin"}, nil)
		mockBlogService.On("GetNextBlogID").Return(1, nil)
		mockBlogService.On("Insert", mock.Anything).Return(fmt.Errorf("error"))
		err := BlogSV.Insert(entities.BlogModel{})
		assert.Equal(t, fmt.Errorf("error"), err)

	})
	t.Run("TestInsert_success", func(t *testing.T) {
		RefreshBlogSV()

		MockUserService.On("FindByID", mock.Anything).Return(&entities.UserDataFormat{Role: "admin"}, nil)
		mockBlogService.On("GetNextBlogID").Return(1, nil)
		mockBlogService.On("Insert", mock.Anything).Return(nil)
		err := BlogSV.Insert(entities.BlogModel{})
		assert.Equal(t, nil, err)

	})
}
func TestUpdate(t *testing.T) {
	t.Run("TestFindByID", func(t *testing.T) {
		RefreshBlogSV()

		MockUserService.On("FindByID", mock.Anything).Return(&entities.UserDataFormat{}, fmt.Errorf("User not found"))

		err := BlogSV.Update(entities.BlogModel{}, "1", "1")
		assert.Equal(t, fmt.Errorf("User not found"), err)

	})
	t.Run("TestCheck_permission", func(t *testing.T) {
		RefreshBlogSV()

		MockUserService.On("FindByID", mock.Anything).Return(&entities.UserDataFormat{Role: "user"}, nil)
		err := BlogSV.Update(entities.BlogModel{}, "1", "1")
		assert.Equal(t, fmt.Errorf("you don't have permission"), err)

	})
	t.Run("TestFindOne_error", func(t *testing.T) {
		RefreshBlogSV()

		MockUserService.On("FindByID", mock.Anything).Return(&entities.UserDataFormat{Role: "admin"}, nil)
		mockBlogService.On("FindOne", mock.Anything).Return(entities.BlogModel{}, fmt.Errorf("failed to find blog "))
		err := BlogSV.Update(entities.BlogModel{}, "1", "1")
		assert.Equal(t, fmt.Errorf("failed to find blog "), err)

	})
	t.Run("TestUpdate_error", func(t *testing.T) {
		RefreshBlogSV()

		MockUserService.On("FindByID", mock.Anything).Return(&entities.UserDataFormat{Role: "admin"}, nil)
		mockBlogService.On("FindOne", mock.Anything).Return(entities.BlogModel{}, nil)
		mockBlogService.On("Update", mock.Anything, mock.Anything).Return(fmt.Errorf("failed to update blog "))
		err := BlogSV.Update(entities.BlogModel{}, "1", "1")

		assert.Equal(t, fmt.Errorf("failed to update blog "), err)

	})
	t.Run("TestUpdate_success", func(t *testing.T) {
		RefreshBlogSV()

		MockUserService.On("FindByID", mock.Anything).Return(&entities.UserDataFormat{Role: "admin"}, nil)
		mockBlogService.On("FindOne", mock.Anything).Return(entities.BlogModel{}, nil)
		mockBlogService.On("Update", mock.Anything, mock.Anything).Return(nil)
		err := BlogSV.Update(entities.BlogModel{}, "1", "1")

		assert.Equal(t, nil, err)

	})
}
func TestDelete(t *testing.T) {
	t.Run("TestFindByID", func(t *testing.T) {
		RefreshBlogSV()

		MockUserService.On("FindByID", mock.Anything).Return(&entities.UserDataFormat{}, fmt.Errorf("User not found"))

		err := BlogSV.Delete("1", "1")
		assert.Equal(t, fmt.Errorf("User not found"), err)

	})
	t.Run("TestCheck_permission", func(t *testing.T) {
		RefreshBlogSV()

		MockUserService.On("FindByID", mock.Anything).Return(&entities.UserDataFormat{Role: "user"}, nil)
		err := BlogSV.Delete("1", "1")
		assert.Equal(t, fmt.Errorf("you don't have permission"), err)

	})
	t.Run("TestDeleteKeyNameImage_error", func(t *testing.T) {
		RefreshBlogSV()

		MockUserService.On("FindByID", mock.Anything).Return(&entities.UserDataFormat{Role: "admin"}, nil)
		mockS3Provider.On("DeleteKeyNameImage", mock.Anything).Return(fmt.Errorf("failed to delete image from S3"))
		err := BlogSV.Delete("1", "1")
		assert.Equal(t, fmt.Errorf("failed to delete image from S3"), err)

	})
	t.Run("TestDelete_error", func(t *testing.T) {
		RefreshBlogSV()

		MockUserService.On("FindByID", mock.Anything).Return(&entities.UserDataFormat{Role: "admin"}, nil)
		mockS3Provider.On("DeleteKeyNameImage", mock.Anything).Return(nil)
		mockBlogService.On("Delete", mock.Anything).Return(fmt.Errorf("failed to delete blog"))
		err := BlogSV.Delete("1", "1")
		assert.Equal(t, fmt.Errorf("failed to delete blog"), err)

	})
	t.Run("TestDelete_success", func(t *testing.T) {
		RefreshBlogSV()

		MockUserService.On("FindByID", mock.Anything).Return(&entities.UserDataFormat{Role: "admin"}, nil)
		mockS3Provider.On("DeleteKeyNameImage", mock.Anything).Return(nil)
		mockBlogService.On("Delete", mock.Anything).Return(nil)
		err := BlogSV.Delete("1", "1")
		assert.Equal(t, nil, err)

	})

}
func TestUploadImage(t *testing.T) {
	t.Run("TestFindByID", func(t *testing.T) {
		RefreshBlogSV()

		MockUserService.On("FindByID", mock.Anything).Return(&entities.UserDataFormat{}, fmt.Errorf("User not found"))

		_, err := BlogSV.UploadImage("1", "1", []byte{})
		assert.Equal(t, fmt.Errorf("User not found"), err)

	})
	t.Run("TestCheck_permission", func(t *testing.T) {
		RefreshBlogSV()

		MockUserService.On("FindByID", mock.Anything).Return(&entities.UserDataFormat{Role: "user"}, nil)
		_, err := BlogSV.UploadImage("1", "1", []byte{})
		assert.Equal(t, fmt.Errorf("you don't have permission"), err)

	})

	t.Run("TestUploadS3FromString_error", func(t *testing.T) {
		RefreshBlogSV()

		MockUserService.On("FindByID", mock.Anything).Return(&entities.UserDataFormat{Role: "admin"}, nil)
		mockS3Provider.On("CreateKeyNameImageBlog", mock.Anything, mock.Anything).Return("key-name", "image/webp")
		mockS3Provider.On("UploadS3FromString", mock.Anything, mock.Anything, mock.Anything).Return("", fmt.Errorf("error"))
		_, err := BlogSV.UploadImage("1", "1", []byte{})
		assert.Equal(t, fmt.Errorf("error"), err)

	})
	t.Run("TestUploadImage_error", func(t *testing.T) {
		RefreshBlogSV()

		MockUserService.On("FindByID", mock.Anything).Return(&entities.UserDataFormat{Role: "admin"}, nil)
		mockS3Provider.On("CreateKeyNameImageBlog", mock.Anything, mock.Anything).Return("key-name", "image/webp")
		mockS3Provider.On("UploadS3FromString", mock.Anything, mock.Anything, mock.Anything).Return("", nil)
		mockBlogService.On("UploadImage", mock.Anything, mock.Anything, mock.Anything).Return(fmt.Errorf("failed to delete blog from repository"))
		_, err := BlogSV.UploadImage("1", "1", []byte{})
		assert.Equal(t, fmt.Errorf("failed to delete blog from repository"), err)

	})
	t.Run("TestUploadImage_success", func(t *testing.T) {
		RefreshBlogSV()

		MockUserService.On("FindByID", mock.Anything).Return(&entities.UserDataFormat{Role: "admin"}, nil)
		mockS3Provider.On("CreateKeyNameImageBlog", mock.Anything, mock.Anything).Return("key-name", "image/webp")
		mockS3Provider.On("UploadS3FromString", mock.Anything, mock.Anything, mock.Anything).Return("", nil)
		mockBlogService.On("UploadImage", mock.Anything, mock.Anything, mock.Anything).Return(nil)
		_, err := BlogSV.UploadImage("1", "1", []byte{})
		assert.Equal(t, nil, err)

	})

}
