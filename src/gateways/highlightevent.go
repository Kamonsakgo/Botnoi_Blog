package gateways

import (
	"go-fiber-unittest/domain/entities"
	"go-fiber-unittest/src/middlewares"
	"io"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

func (h HTTPGateway) InsertHL_event(ctx *fiber.Ctx) error {
	// Decode JWT token to validate the user
	tokenData, err := middlewares.DecodeJWTToken(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(entities.ResponseModel{Message: "Unauthorized token."})
	}

	userID := tokenData.UserID

	// Required fields validation
	requiredFields := map[string]string{
		"title":          "title is required",
		"location":       "location is required",
		"date":           "date is required",
		"category":       "category is required",
		"highlight_id":   "highlight_id is required",
		"speaker":        "speaker is required",
		"location_event": "location_event is required",
	}

	for field, errMsg := range requiredFields {
		if ctx.FormValue(field) == "" {
			return ctx.Status(fiber.StatusBadRequest).JSON(entities.ResponseModel{Message: errMsg})
		}
	}

	// Handle image file upload if provided
	var fileBytes []byte
	if imageFile, err := ctx.FormFile("imagefile"); err == nil && imageFile != nil {
		fileContent, err := imageFile.Open()
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(entities.ResponseModel{Message: "Failed to open uploaded file."})
		}
		defer fileContent.Close()

		fileBytes, err = io.ReadAll(fileContent)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(entities.ResponseModel{Message: "Failed to read uploaded file."})
		}
	}

	// Prepare the body data
	bodydata := entities.HighlightModel{
		Title:         ctx.FormValue("title"),
		UserID:        userID,
		HighlightID:   ctx.FormValue("highlight_id"),
		Category:      strings.Split(ctx.FormValue("category"), ","),
		Location:      ctx.FormValue("location"),
		LocationEvent: ctx.FormValue("location_event"),
		Date:          ctx.FormValue("date"),
		Speaker:       ctx.FormValue("speaker"),
		CreatedAt:     time.Now().UTC().Add(7 * time.Hour),
		ImageURL:      "", // ImageURL can be updated later if needed
	}

	// Insert the new highlight event
	err = h.HighlightService.Insert(bodydata, fileBytes)
	if err != nil {
		return ctx.Status(fiber.StatusForbidden).JSON(entities.ResponseModel{Message: err.Error()})
	}

	// Return success response
	return ctx.Status(fiber.StatusOK).JSON(entities.ResponseModel{Message: "Success"})
}
func (h HTTPGateway) UpdateHL_event(ctx *fiber.Ctx) error {
	// Decode JWT token to validate the user
	tokenData, err := middlewares.DecodeJWTToken(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(entities.ResponseModel{Message: "Unauthorized token."})
	}

	params := ctx.Queries()
	highlightID := params["highlight_id"]
	if highlightID == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(entities.ResponseModel{Message: "highlight_id is required"})
	}
	requiredFields := map[string]string{
		"title":          "title is required",
		"location":       "location is required",
		"date":           "date is required",
		"category":       "category is required",
		"highlight_id":   "highlight_id is required",
		"speaker":        "speaker is required",
		"location_event": "location_event is required",
	}

	for field, errMsg := range requiredFields {
		if ctx.FormValue(field) == "" {
			return ctx.Status(fiber.StatusBadRequest).JSON(entities.ResponseModel{Message: errMsg})
		}
	}
	// Prepare updated data
	updatedData := entities.HighlightModel{
		Title:         ctx.FormValue("title"),
		Category:      strings.Split(ctx.FormValue("category"), ","),
		Location:      ctx.FormValue("location"),
		LocationEvent: ctx.FormValue("location_event"),
		Date:          ctx.FormValue("date"),
		Speaker:       ctx.FormValue("speaker"),
		LastUpdateAt:  time.Now().UTC().Add(7 * time.Hour),
		HighlightID:   highlightID,
	}

	// Handle image file upload if provided
	var fileBytes []byte
	if imageFile, err := ctx.FormFile("imagefile"); err == nil && imageFile != nil {
		fileContent, err := imageFile.Open()
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(entities.ResponseModel{Message: "Failed to open uploaded file."})
		}
		defer fileContent.Close()

		fileBytes, err = io.ReadAll(fileContent)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(entities.ResponseModel{Message: "Failed to read uploaded file."})
		}
	}

	// Call update service
	err = h.HighlightService.Update(updatedData, fileBytes, tokenData.UserID)
	if err != nil {
		return ctx.Status(fiber.StatusForbidden).JSON(entities.ResponseModel{Message: err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(entities.ResponseModel{Message: "Update successful"})
}
func (h HTTPGateway) DeleteHL_event(ctx *fiber.Ctx) error {
	// Decode JWT token to validate the user
	tokenData, err := middlewares.DecodeJWTToken(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(entities.ResponseModel{Message: "Unauthorized token."})
	}

	// Extract highlight_id from request
	params := ctx.Queries()
	highlightID := params["highlight_id"]
	if highlightID == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(entities.ResponseModel{Message: "highlight_id is required"})
	}

	// Call delete service
	err = h.HighlightService.Delete(highlightID, tokenData.UserID)
	if err != nil {
		return ctx.Status(fiber.StatusForbidden).JSON(entities.ResponseModel{Message: err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(entities.ResponseModel{Message: "Delete successful"})
}
func (h *HTTPGateway) GetAllHL_event(ctx *fiber.Ctx) error {
	params := ctx.Queries()

	page := params["page"]
	if page == "" {
		page = "1"
	}
	limit := params["limit"]
	if limit == "" {
		limit = "10"
	}
	blogs, err := h.HighlightService.GetAll(page, limit)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(entities.ResponseModel{Message: "Failed to get all blog"})
	}
	return ctx.Status(fiber.StatusOK).JSON(entities.ResponseModel{
		Message: "success",
		Data:    blogs,
	})
}

func (h *HTTPGateway) GetOneHL_event(ctx *fiber.Ctx) error {
	params := ctx.Queries()
	highlightID := params["highlight_id"]
	if highlightID == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(entities.ResponseModel{Message: "highlightID is required"})
	}

	// Call GetOne and store the result in a separate variable
	highlight, err := h.HighlightService.GetOne(highlightID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(entities.ResponseModel{Message: "highlight not found"})
	}

	return ctx.Status(fiber.StatusOK).JSON(entities.ResponseModel{
		Message: "success",
		Data:    highlight, // Return the full highlight object or its details
	})
}
