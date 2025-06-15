package services

import (
	//"fmt"

	"fmt"
	"go-fiber-unittest/domain/entities"
	"go-fiber-unittest/domain/repositories"
	"go-fiber-unittest/src/utils/providers"
	"strconv"
)

type blogsService struct {
	BlogsRepository repositories.IBlogsRepository
	S3Provider      providers.IS3provider
	UsersRepository repositories.IUsersRepository
}

type IBlogsService interface {
	Insert(data entities.BlogModel) error
	Update(data entities.BlogModel, userID string, blog_id string) error
	Delete(userID string, blog_id string) error
	GetAll(page string, limit string) (data entities.BlogResponseModel, err error)
	GetOne(blog_id string) (data entities.BlogModel, err error)
	UploadImage(userID string, blog_id string, imageBytes []byte) (URL string, err error)
}

func NewBlogsService(repo0 repositories.IBlogsRepository, repo1 repositories.IUsersRepository, s3 providers.IS3provider) IBlogsService {
	return &blogsService{
		BlogsRepository: repo0,
		UsersRepository: repo1,
		S3Provider:      s3,
	}
}

func (sv blogsService) Insert(data entities.BlogModel) error {
	admin, err := sv.UsersRepository.FindByID(data.UserID)
	if err != nil {
		return err
	}

	if admin.Role != "superadmin" && admin.Role != "admin" {
		return fmt.Errorf("you don't have permission")
	}
	blogID, err := sv.BlogsRepository.GetNextBlogID()
	if err != nil {
		return err
	}
	data.BlogID = strconv.Itoa(blogID)
	err = sv.BlogsRepository.Insert(data)
	if err != nil {
		return err
	}

	return nil

}
func (sv blogsService) Update(data entities.BlogModel, blog_id string, userID string) error {
	admin, err := sv.UsersRepository.FindByID(userID)
	if err != nil {
		return err
	}
	if admin.Role != "superadmin" && admin.Role != "admin" {
		return fmt.Errorf("you don't have permission")
	}
	_, err = sv.BlogsRepository.FindOne(blog_id)
	if err != nil {
		return fmt.Errorf("failed to find blog ")
	}

	err = sv.BlogsRepository.Update(data, blog_id)
	if err != nil {
		return err
	}

	return nil

}

func (s *blogsService) GetAll(page string, limit string) (data entities.BlogResponseModel, err error) {
	pageInt, _ := strconv.Atoi(page)
	limitInt, _ := strconv.Atoi(limit)
	blogs, totalCount, err := s.BlogsRepository.FindAll(pageInt, limitInt)
	if err != nil {
		return data, err
	}

	// Calculate the total number of pages
	totalPages := (totalCount + int64(limitInt) - 1) / int64(limitInt) // Round up division

	// Return the paginated data
	data = entities.BlogResponseModel{
		Blogs:      blogs,
		Page:       pageInt,
		TotalPages: int(totalPages),
		TotalCount: totalCount,
	}
	return data, nil
}

func (s *blogsService) GetOne(blog_id string) (data entities.BlogModel, err error) {
	blogs, err := s.BlogsRepository.FindOne(blog_id)
	if err != nil {
		return data, err
	}
	return blogs, nil
}
func (sv blogsService) Delete(userID string, blog_id string) error {
	admin, err := sv.UsersRepository.FindByID(userID)
	if err != nil {
		return fmt.Errorf("User not found")
	}

	if admin.Role != "superadmin" && admin.Role != "admin" {
		return fmt.Errorf("you don't have permission")
	}

	pathName := "blog_id:" + blog_id
	err = sv.S3Provider.DeleteKeyNameImage(pathName)
	if err != nil {
		return fmt.Errorf("failed to delete image from S3")
	}
	err = sv.BlogsRepository.Delete(blog_id)
	if err != nil {
		return fmt.Errorf("failed to delete blog")
	}
	return nil
}
func (sv blogsService) UploadImage(userID string, blog_id string, imageBytes []byte) (URL string, err error) {
	admin, err := sv.UsersRepository.FindByID(userID)
	if err != nil {
		return "", fmt.Errorf("User not found")
	}

	if admin.Role != "superadmin" && admin.Role != "admin" {
		return "", fmt.Errorf("you don't have permission")
	}
	pathName := "blog_id:" + blog_id
	key, contentType := sv.S3Provider.CreateKeyNameImageBlog(pathName, "webp")
	URL, err = sv.S3Provider.UploadS3FromString(imageBytes, key, contentType)
	if err != nil {
		return "", err
	}
	name := URL[len(URL)-12:]
	err = sv.BlogsRepository.UploadImage(blog_id, name, URL)
	if err != nil {
		return "", fmt.Errorf("failed to UploadImage from repository")
	}

	return URL, nil
}
