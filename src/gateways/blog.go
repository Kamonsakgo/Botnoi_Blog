package gateways

import (
	"go-fiber-unittest/domain/entities"
	"go-fiber-unittest/src/middlewares"
	"io"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

func (h HTTPGateway) Insertblog(ctx *fiber.Ctx) error {
	tokenData, err := middlewares.DecodeJWTToken(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(entities.ResponseModel{Message: "Unauthorization Token."})
	}
	userID := tokenData.UserID
	//token := tokenData.Token

	if title := ctx.FormValue("title"); title == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(entities.ResponseModel{Message: "title is required"})
	}
	if content := ctx.FormValue("content"); content == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(entities.ResponseModel{Message: "content is required"})
	}
	if category := ctx.FormValue("category"); category == "" || (category != "event" && category != "article") {
		return ctx.Status(fiber.StatusBadRequest).JSON(entities.ResponseModel{
			Message: "category must be 'event' or 'article'",
		})
	}
	if tag := ctx.FormValue("tag"); tag == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(entities.ResponseModel{Message: "tag is required"})
	}
	if typeblog := ctx.FormValue("type"); typeblog == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(entities.ResponseModel{Message: "type is required"})
	}
	if location := ctx.FormValue("location"); location == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(entities.ResponseModel{Message: "location is required"})
	}

	bodydata := entities.BlogModel{
		Title:     ctx.FormValue("title"),
		UserID:    userID,
		Content:   ctx.FormValue("content"),
		Category:  ctx.FormValue("category"),
		Tag:       strings.Split(ctx.FormValue("tag"), ","),
		Type:      strings.Split(ctx.FormValue("type"), ","),
		Location:  ctx.FormValue("location"),
		CreatedAt: time.Now().UTC().Add(7 * time.Hour),
		ImageURL:  nil,
	}
	status := h.BlogService.Insert(bodydata)
	if status != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(entities.ResponseModel{Message: "Failed to insert new blog"})
	}
	return ctx.Status(fiber.StatusCreated).JSON(entities.ResponseModel{
		Message: "Blog created successfully",
	})
}
func (h HTTPGateway) Updateblog(ctx *fiber.Ctx) error {
	tokenData, err := middlewares.DecodeJWTToken(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(entities.ResponseModel{Message: "Unauthorization Token."})
	}
	userID := tokenData.UserID
	//token := tokenData.Token
	params := ctx.Queries()
	blog_id := params["blog_id"]
	if blog_id == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(entities.ResponseMessage{Message: "blog_id is required"})
	}
	if title := ctx.FormValue("title"); title == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(entities.ResponseModel{Message: "title is required"})
	}
	if content := ctx.FormValue("content"); content == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(entities.ResponseModel{Message: "content is required"})
	}
	if category := ctx.FormValue("category"); category == "" || (category != "event" && category != "article") {
		return ctx.Status(fiber.StatusBadRequest).JSON(entities.ResponseModel{
			Message: "category must be 'event' or 'article'",
		})
	}
	if tag := ctx.FormValue("tag"); tag == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(entities.ResponseModel{Message: "tag is required"})
	}
	if typeblog := ctx.FormValue("type"); typeblog == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(entities.ResponseModel{Message: "type is required"})
	}
	if location := ctx.FormValue("location"); location == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(entities.ResponseModel{Message: "location is required"})
	}
	bodydata := entities.BlogModel{
		Title:        ctx.FormValue("title"),
		Content:      ctx.FormValue("content"),
		Category:     ctx.FormValue("category"),
		Tag:          strings.Split(ctx.FormValue("tag"), ","),
		Type:         strings.Split(ctx.FormValue("type"), ","),
		Location:     ctx.FormValue("location"),
		HL_id:        ctx.FormValue("highlight_id"),
		LastUpdateAt: time.Now().UTC().Add(7 * time.Hour),
	}
	status := h.BlogService.Update(bodydata, blog_id, userID)
	if status != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(entities.ResponseModel{Message: "Failed to update  blog"})
	}
	return ctx.Status(fiber.StatusCreated).JSON(entities.ResponseModel{
		Message: "Blog update successfully",
	})
}

func (h *HTTPGateway) GetAllblog(ctx *fiber.Ctx) error {
	params := ctx.Queries()

	page := params["page"]
	if page == "" {
		page = "1"
	}
	limit := params["limit"]
	if limit == "" {
		limit = "10"
	}
	blogs, err := h.BlogService.GetAll(page, limit)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(entities.ResponseModel{Message: "Failed to get all blog"})
	}
	return ctx.Status(fiber.StatusOK).JSON(entities.ResponseModel{
		Message: "success",
		Data:    blogs,
	})
}

func (h *HTTPGateway) GetOneBlog(ctx *fiber.Ctx) error {
	params := ctx.Queries()
	blog_id := params["blog_id"]
	if blog_id == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(entities.ResponseModel{Message: "blog_id is required"})
	}
	blog, err := h.BlogService.GetOne(blog_id)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(entities.ResponseModel{Message: "blog not found"})
	}
	return ctx.Status(fiber.StatusOK).JSON(entities.ResponseModel{
		Message: "success",
		Data:    blog,
	})
}

func (h HTTPGateway) Deleteblog(ctx *fiber.Ctx) error {
	tokenData, err := middlewares.DecodeJWTToken(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(entities.ResponseModel{Message: "Unauthorization Token."})
	}
	userID := tokenData.UserID
	params := ctx.Queries()
	blog_id := params["blog_id"]
	if blog_id == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(entities.ResponseMessage{Message: "blog_id is required"})
	}

	status := h.BlogService.Delete(userID, blog_id)
	if status != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(entities.ResponseModel{Message: "Failed to Delete blog"})
	}
	return ctx.Status(fiber.StatusCreated).JSON(entities.ResponseModel{
		Message: "Blog Delete successfully",
	})
}
func (h HTTPGateway) UploadImage(ctx *fiber.Ctx) error {
	tokenData, err := middlewares.DecodeJWTToken(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(entities.ResponseModel{Message: "Unauthorization Token."})
	}
	userID := tokenData.UserID
	params := ctx.Queries()
	blog_id := params["blog_id"]
	if blog_id == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(entities.ResponseMessage{Message: "blog_id is required"})
	}
	var fileBytes []byte
	if imageFile, err := ctx.FormFile("imagefile"); err == nil && imageFile != nil {
		fileContent, err := imageFile.Open()
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(entities.ResponseModel{Message: "Failed to open uploaded file"})
		}
		defer fileContent.Close()

		fileBytes, err = io.ReadAll(fileContent)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(entities.ResponseModel{Message: "Failed to read uploaded file"})
		}
	}

	URL, status := h.BlogService.UploadImage(userID, blog_id, fileBytes)
	if status != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(entities.ResponseModel{Message: "Failed to UploadImage blog"})
	}
	return ctx.Status(fiber.StatusCreated).JSON(entities.ResponseModel{
		Message: "Blog UploadImage successfully",
		Data:    URL,
	})
}
