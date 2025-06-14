package gateways

import (
	"bn-crud-ads/domain/entities"
	"bn-crud-ads/src/middlewares"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func (h HTTPGateway) GetAdsData(ctx *fiber.Ctx) error {

	data, err := h.AdsService.GetAds()
	if err != nil {
		return ctx.Status(fiber.StatusForbidden).JSON(entities.ResponseModel{Message: "cannot get all ads data"})
	}

	return ctx.Status(fiber.StatusOK).JSON(entities.ResponseModel{Message: "get ads success", Data: data})
}

func (h HTTPGateway) SetAdsRedis(ctx *fiber.Ctx) error {

	status := h.AdsService.SetAdsRedis()
	if !status {
		return ctx.Status(fiber.StatusForbidden).JSON(entities.ResponseModel{Message: "set redis unsuccess"})
	}
	return ctx.Status(fiber.StatusOK).JSON(entities.ResponseModel{Message: "success"})
}

func (h HTTPGateway) GetAdsNoneToken(ctx *fiber.Ctx) error {
	data, err := h.AdsService.GetAds()
	if err != nil {
		return ctx.Status(fiber.StatusForbidden).JSON(entities.ResponseModel{Message: "cannot get all ads data"})
	}
	return ctx.Status(fiber.StatusOK).JSON(entities.ResponseModel{Message: "success", Data: data})
}

func (h HTTPGateway) UpdateAds(ctx *fiber.Ctx) error {
	adsModel := new(entities.AdsModel)

	if err := ctx.BodyParser(&adsModel); err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	data, err := h.AdsService.UpdateAds(*adsModel)
	if err != nil {
		return ctx.Status(fiber.StatusForbidden).JSON(entities.ResponseModel{Message: "cannot update ads data"})
	}
	return ctx.Status(fiber.StatusOK).JSON(entities.ResponseMessage{Message: data})
}

func (h HTTPGateway) GetMarketplaceSound(ctx *fiber.Ctx) error {
	data, err := h.AdsService.GetMLPSound()
	if err != nil {
		return ctx.Status(fiber.StatusForbidden).JSON(entities.ResponseModel{Message: "get ads unsuccess", Status: 403})
	}
	if len(data) > 0 {
		return ctx.Status(fiber.StatusOK).JSON(entities.ResponseModel{Message: "get ads success", Data: data})
	} else {
		return ctx.Status(fiber.StatusOK).JSON(entities.ResponseModel{Message: "empty ads", Data: data})
	}
}
func (h HTTPGateway) GetAdsRedis(ctx *fiber.Ctx) error {

	data, err := h.AdsService.GetAdsRedis()

	if err != nil {
		return ctx.Status(fiber.StatusForbidden).JSON(entities.ResponseModel{Message: "cannot get redis"})
	}
	return ctx.Status(fiber.StatusOK).JSON(entities.ResponseModel{Message: "get redis success", Data: data})
}
func (h HTTPGateway) GetAdsAlertMessage(ctx *fiber.Ctx) error {

	data, err := h.AdsService.GetAdsAlertMessage()

	if err != nil {
		return ctx.Status(fiber.StatusForbidden).JSON(entities.ResponseModel{Message: "cannot get adsAlertmessage"})
	}
	return ctx.Status(fiber.StatusOK).JSON(entities.ResponseModel{Message: "get adsAlertmessage success", Data: data})
}

func (h HTTPGateway) Insert_ads(ctx *fiber.Ctx) error {
	tokendata, err := middlewares.DecodeJWTToken(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(entities.ResponseModel{Message: "Unauthorization Token."})
	}
	// รับค่าจาก data_form
	nameAds := ctx.FormValue("name_ads")
	speaker := ctx.FormValue("speaker")
	description := ctx.FormValue("description")
	isPreviewStr := ctx.FormValue("is_preview")
	isPreview := false

	if isPreviewStr == "true" || isPreviewStr == "1" {
		isPreview = true
	}

	// ตรวจสอบฟิลด์ที่จำเป็น ยกเว้น name_ads
	if speaker == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(entities.ResponseModel{Message: "speaker is required"})
	}

	if description == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(entities.ResponseModel{Message: "description is required"})
	}

	var fileBytes []byte
	var fileExtension string

	if imageFile, err := ctx.FormFile("image_file"); err == nil && imageFile != nil {
		// ดึงชื่อไฟล์และนามสกุลไฟล์
		filename := imageFile.Filename
		fileExtension = strings.ToLower(filepath.Ext(filename)) // แปลงเป็นตัวพิมพ์เล็กเพื่อตรวจสอบง่ายขึ้น

		// ตรวจสอบประเภทไฟล์จากนามสกุล
		if fileExtension != ".png" && fileExtension != ".jpg" && fileExtension != ".jpeg" && fileExtension != ".svg" {
			return ctx.Status(fiber.StatusBadRequest).JSON(entities.ResponseModel{
				Message: "Unsupported file type. Only SVG, PNG, and JPG/JPEG are allowed",
			})
		}

		// อ่านเนื้อหาไฟล์
		fileContent, err := imageFile.Open()
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(entities.ResponseModel{Message: "Failed to open uploaded file"})
		}
		defer fileContent.Close()

		fileBytes, err = io.ReadAll(fileContent)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(entities.ResponseModel{Message: "Failed to read uploaded file"})
		}
	} else {
		return ctx.Status(fiber.StatusBadRequest).JSON(entities.ResponseModel{Message: "image_file is required"})
	}

	// เรียกใช้บริการ AdsService พร้อมส่งนามสกุลไฟล์
	image, audio, err := h.AdsService.Insert_ads(description, nameAds, speaker, fileBytes, fileExtension, *tokendata.Token, isPreview)
	if err != nil {
		errorMessage := fmt.Sprintf("cannot insert ads data: %v", err)
		return ctx.Status(fiber.StatusForbidden).JSON(entities.ResponseModel{
			Message: errorMessage,
			Data:    nil, // ไม่มีข้อมูลเพิ่มเติมให้ส่งกลับ
		})
	}

	// กรณีสำเร็จ
	return ctx.Status(fiber.StatusOK).JSON(entities.ResponseModel{
		Message: "Ads inserted successfully",
		Data: map[string]interface{}{
			"image": image,
			"audio": audio,
		},
	})

}
func (h HTTPGateway) Delete_ads(ctx *fiber.Ctx) error {
	_, err := middlewares.DecodeJWTToken(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(entities.ResponseModel{Message: "Unauthorization Token."})
	}
	params := ctx.Queries()
	ads_id := params["ads_id"]
	if ads_id == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(entities.ResponseMessage{Message: "ads_id is required"})
	}

	err = h.AdsService.Delete_ads(ads_id)
	if err != nil {
		errorMessage := fmt.Sprintf("cannot Delete ads data: %v", err)
		return ctx.Status(fiber.StatusForbidden).JSON(entities.ResponseModel{Message: errorMessage})
	}
	return ctx.Status(fiber.StatusOK).JSON(entities.ResponseMessage{Message: "success"})
}
func (h HTTPGateway) Update(ctx *fiber.Ctx) error {
	_, err := middlewares.DecodeJWTToken(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(entities.ResponseModel{Message: "Unauthorization Token."})
	}
	params := ctx.Queries()
	ads_id := params["ads_id"]
	if ads_id == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(entities.ResponseMessage{Message: "ads_id is required"})
	}
	nameAds := ctx.FormValue("name_ads")
	speaker := ctx.FormValue("speaker")
	description := ctx.FormValue("description")

	// ตรวจสอบฟิลด์ที่จำเป็น ยกเว้น name_ads
	if speaker == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(entities.ResponseModel{Message: "speaker is required"})
	}

	if description == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(entities.ResponseModel{Message: "description is required"})
	}

	var fileBytes []byte
	var fileExtension string

	if imageFile, err := ctx.FormFile("image_file"); err == nil && imageFile != nil {
		// ดึงชื่อไฟล์และนามสกุลไฟล์
		filename := imageFile.Filename
		fileExtension = strings.ToLower(filepath.Ext(filename)) // แปลงเป็นตัวพิมพ์เล็กเพื่อตรวจสอบง่ายขึ้น

		// ตรวจสอบประเภทไฟล์จากนามสกุล
		if fileExtension != ".png" && fileExtension != ".jpg" && fileExtension != ".jpeg" && fileExtension != ".svg" {
			return ctx.Status(fiber.StatusBadRequest).JSON(entities.ResponseModel{
				Message: "Unsupported file type. Only SVG, PNG, and JPG/JPEG are allowed",
			})
		}

		// อ่านเนื้อหาไฟล์
		fileContent, err := imageFile.Open()
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(entities.ResponseModel{Message: "Failed to open uploaded file"})
		}
		defer fileContent.Close()

		fileBytes, err = io.ReadAll(fileContent)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(entities.ResponseModel{Message: "Failed to read uploaded file"})
		}
	} else {
		return ctx.Status(fiber.StatusBadRequest).JSON(entities.ResponseModel{Message: "image_file is required"})
	}
	image, audio, err := h.AdsService.Update_ads(ads_id, description, nameAds, speaker, fileBytes, fileExtension)
	if err != nil {
		errorMessage := fmt.Sprintf("cannot update ads data: %v", err)
		return ctx.Status(fiber.StatusForbidden).JSON(entities.ResponseModel{Message: errorMessage})
	}
	responseData := map[string]interface{}{}
	if image != "" {
		responseData["image"] = image
	}
	if audio != "" {
		responseData["audio"] = audio
	}

	// คืนค่า JSON Response
	return ctx.Status(fiber.StatusOK).JSON(map[string]interface{}{
		"message": "success",
		"data":    responseData,
	})
}

func (h HTTPGateway) Getall_ads(ctx *fiber.Ctx) error {

	data, err := h.AdsService.Getall_ads()
	if err != nil {
		errorMessage := fmt.Sprintf("cannot get all ads data: %v", err)
		return ctx.Status(fiber.StatusForbidden).JSON(entities.ResponseModel{Message: errorMessage})
	}

	return ctx.Status(fiber.StatusOK).JSON(entities.ResponseModel{Message: "get ads success", Data: data})
}

func (h HTTPGateway) Find_one_ads(ctx *fiber.Ctx) error {
	params := ctx.Queries()
	ads_id := params["ads_id"]
	if ads_id == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(entities.ResponseMessage{Message: "ads_id is required"})
	}
	data, err := h.AdsService.FindOne(ads_id)
	if err != nil {
		errorMessage := fmt.Sprintf("cannot get one ads data: %v", err)
		return ctx.Status(fiber.StatusForbidden).JSON(entities.ResponseModel{Message: errorMessage})
	}

	return ctx.Status(fiber.StatusOK).JSON(entities.ResponseModel{Message: "get ads success", Data: data})
}
func (h HTTPGateway) Find_one_random_ads(ctx *fiber.Ctx) error {
	data, err := h.AdsService.Find_one_random_ads()
	if err != nil {
		errorMessage := fmt.Sprintf("cannot get one ads data: %v", err)
		return ctx.Status(fiber.StatusForbidden).JSON(entities.ResponseModel{Message: errorMessage})
	}

	return ctx.Status(fiber.StatusOK).JSON(entities.ResponseModel{Message: "get ads success", Data: data})
}
