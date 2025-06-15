package services

import (
	//"fmt"

	"fmt"
	"go-fiber-unittest/domain/entities"
	"go-fiber-unittest/domain/repositories"
	"go-fiber-unittest/src/utils/providers"
	"log"
	"strconv"
	"strings"
)

type HighlightsService struct {
	HighlightsRepository repositories.IHighlightsRepository
	S3Provider           providers.IS3provider
	UsersRepository      repositories.IUsersRepository
}

type IHighlightsService interface {
	Insert(data entities.HighlightModel, imageBytes []byte) error
	Delete(highlightID string, userID string) error
	Update(data entities.HighlightModel, imageBytes []byte, UserID string) error
	GetAll(page string, limit string) (data entities.HighlightResponseModel, err error)
	GetOne(highlightID string) (data entities.HighlightModel, err error)
}

func NewHighlightsService(repo0 repositories.IHighlightsRepository, repo1 repositories.IUsersRepository, s3 providers.IS3provider) IHighlightsService {
	return &HighlightsService{
		HighlightsRepository: repo0,
		UsersRepository:      repo1,
		S3Provider:           s3,
	}
}

func (sv HighlightsService) Insert(data entities.HighlightModel, imageBytes []byte) error {

	admin, err := sv.UsersRepository.FindByID(data.UserID)
	if err != nil {
		return err
	}
	if admin.Role != "superadmin" && admin.Role != "admin" {
		return fmt.Errorf("you don't have permission")
	}
	existing, err := sv.HighlightsRepository.FindOne(data.HighlightID)
	if existing != nil {
		return fmt.Errorf("highlight_id '%s' already exists", data.HighlightID)
	}
	if err != nil && err.Error() != fmt.Sprintf("highlight with id '%s' not found", data.HighlightID) {
		return fmt.Errorf("error checking highlight_id: %w", err)
	}
	key, contentType := sv.S3Provider.CreateKeyNameImageHL("webp")
	URL, err := sv.S3Provider.UploadS3FromString(imageBytes, key, contentType)
	if err != nil {
		return err
	}
	data.ImageURL = URL
	err = sv.HighlightsRepository.Insert(data)
	if err != nil {
		return err
	}

	return nil
}
func (sv HighlightsService) Update(data entities.HighlightModel, imageBytes []byte, UserID string) error {
	// Step 1: Validate user permissions
	admin, err := sv.UsersRepository.FindByID(UserID)
	if err != nil {
		return fmt.Errorf("error fetching user: %w", err)
	}
	if admin.Role != "superadmin" && admin.Role != "admin" {
		return fmt.Errorf("you don't have permission to perform this action")
	}

	// Step 2: Check if highlight_id exists
	existing, err := sv.HighlightsRepository.FindOne(data.HighlightID)
	if existing == nil {
		return fmt.Errorf("highlight_id '%s' does not exist", data.HighlightID)
	}
	if err != nil {
		return fmt.Errorf("error checking highlight_id: %w", err)
	}

	// Step 3: Process and upload image (if provided)
	if len(imageBytes) > 0 {
		// Extract the filename from the existing URL before overwriting
		if existing.ImageURL != "" {
			parts := strings.Split(existing.ImageURL, "/")
			oldKey := parts[len(parts)-1] // Get the last part of the URL
			if err := sv.S3Provider.DeleteKeyNameImageHL(oldKey); err != nil {
				return fmt.Errorf("error deleting old image from S3: %w", err)
			}
		}

		key, contentType := sv.S3Provider.CreateKeyNameImageHL("webp")
		URL, err := sv.S3Provider.UploadS3FromString(imageBytes, key, contentType)
		if err != nil {
			return fmt.Errorf("error uploading image to S3: %w", err)
		}
		data.ImageURL = URL
	} else {
		data.ImageURL = existing.ImageURL // Keep existing image if no new image provided
	}

	// Step 4: Update highlight data into repository
	err = sv.HighlightsRepository.Update(data, data.HighlightID)
	if err != nil {
		return fmt.Errorf("error updating highlight: %w", err)
	}

	// Step 5: Log success and return
	log.Printf("Highlight with ID '%s' updated successfully", data.HighlightID)
	return nil
}

func (sv HighlightsService) Delete(highlightID string, userID string) error {
	// Step 1: Validate user permissions
	admin, err := sv.UsersRepository.FindByID(userID)
	if err != nil {
		return fmt.Errorf("error fetching user: %w", err)
	}
	if admin.Role != "superadmin" && admin.Role != "admin" {
		return fmt.Errorf("you don't have permission to perform this action")
	}

	// Step 2: Check if highlight_id exists
	existing, err := sv.HighlightsRepository.FindOne(highlightID)
	if existing == nil {
		return fmt.Errorf("highlight_id '%s' does not exist", highlightID)
	}
	if err != nil {
		return fmt.Errorf("error checking highlight_id: %w", err)
	}
	if existing.ImageURL != "" {
		parts := strings.Split(existing.ImageURL, "/")
		oldKey := parts[len(parts)-1] // Get the last part of the URL
		if err := sv.S3Provider.DeleteKeyNameImageHL(oldKey); err != nil {
			return fmt.Errorf("error deleting old image from S3: %w", err)
		}
	}
	// Step 3: Delete highlight data from repository
	err = sv.HighlightsRepository.Delete(highlightID)
	if err != nil {
		return fmt.Errorf("error deleting highlight: %w", err)
	}

	// Step 4: Log success and return
	log.Printf("Highlight with ID '%s' deleted successfully", highlightID)
	return nil
}
func (s *HighlightsService) GetAll(page string, limit string) (data entities.HighlightResponseModel, err error) {
	pageInt, _ := strconv.Atoi(page)
	limitInt, _ := strconv.Atoi(limit)
	highlights, totalCount, err := s.HighlightsRepository.FindAll(pageInt, limitInt)
	if err != nil {
		return data, err
	}

	// Calculate the total number of pages
	totalPages := (totalCount + int64(limitInt) - 1) / int64(limitInt) // Round up division

	// Return the paginated data
	data = entities.HighlightResponseModel{
		Highlights: highlights,
		Page:       pageInt,
		TotalPages: int(totalPages),
		TotalCount: totalCount,
	}
	return data, nil
}

func (s *HighlightsService) GetOne(highlight_id string) (data entities.HighlightModel, err error) {
	highlights, err := s.HighlightsRepository.FindOne(highlight_id)
	if err != nil {
		return data, err
	}
	return *highlights, nil
}
