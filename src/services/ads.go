package services

import (
	"bn-crud-ads/domain/entities"
	"bn-crud-ads/domain/repositories"
	"bn-crud-ads/utils/providers"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type AdsService struct {
	AdsRepository    repositories.IAdsRepository
	RedisRepository  repositories.IRedisConnectionRepository
	S3Provider       providers.IS3provider
	wrongMessageRepo repositories.IWrongMessageRepository
}

type IAdsService interface {
	GetAds() (*bson.M, error)
	UpdateAds(data entities.AdsModel) (string, error)
	GetMLPSound() ([]bson.M, error)
	GetAdsRedis() (*entities.AdsAlertMessage, error)
	SetAdsRedis() bool
	GetAdsAlertMessage() (*entities.AdsAlertMessage, error)
	Insert_ads(description string, nameAds string, speaker string, fileBytes []byte, filename string, tokendata string, isPreview bool) (Url_Image string, Url_AUdio string, error error)
	Delete_ads(ads_id string) error
	Update_ads(ads_id string, description string, nameAds string, speaker string, fileBytes []byte, fileExtension string) (Url_Image string, Url_AUdio string, error error)
	Getall_ads() ([]entities.AdsDataFormat, error)
	FindOne(adsID string) (entities.AdsDataFormat, error)
	Find_one_random_ads() (entities.AdsDataFormat, error)
}

func NewAdsService(repo0 repositories.IAdsRepository, cache0 repositories.IRedisConnectionRepository, s3 providers.IS3provider, wrongMessageRepo repositories.IWrongMessageRepository) IAdsService {
	return &AdsService{
		AdsRepository:    repo0,
		RedisRepository:  cache0,
		S3Provider:       s3,
		wrongMessageRepo: wrongMessageRepo,
	}
}

func (sv AdsService) GetAds() (*bson.M, error) {
	percent, err1 := sv.GetAdsRedis()
	var data []int
	var probabilities []float64
	if err1 != nil {
		return nil, err1
	}
	if percent != nil {
		data = percent.Ads.Level
		probabilities = percent.Ads.Percent
	} else {
		mongo, err := sv.GetAdsAlertMessage()
		if err != nil {
			return nil, err
		}
		data = mongo.Ads.Level
		probabilities = mongo.Ads.Percent

	}

	randomized_data := random_choice(data, probabilities)
	adData, err := sv.AdsRepository.FindAds(randomized_data)
	if err != nil {
		return nil, err
	}
	delete(*adData, "level")
	delete(*adData, "_id")

	return adData, nil

}

func (sv AdsService) UpdateAds(data entities.AdsModel) (string, error) {
	adData, err := sv.AdsRepository.Update(data)
	if err != nil {
		return "", err
	}

	return adData, nil

}

func (sv AdsService) GetMLPSound() ([]bson.M, error) {
	adData, err := sv.AdsRepository.FindMkpSound()
	if err != nil {
		return nil, err
	}

	return adData, nil
}
func (sv AdsService) GetAdsRedis() (*entities.AdsAlertMessage, error) {
	adData, err := sv.RedisRepository.GetAdsAlertMessage()

	return adData, err

}
func (sv AdsService) GetAdsAlertMessage() (*entities.AdsAlertMessage, error) {
	adData, err := sv.AdsRepository.FindAdsAlertMessage()
	if err != nil {
		return nil, err
	}
	return adData, err

}
func (sv AdsService) SetAdsRedis() bool {
	adData, err := sv.AdsRepository.FindAdsAlertMessage()
	if err != nil {
		log.Println(err)
		return false
	}
	adDataJson, _ := json.Marshal(adData)
	setStatus := sv.RedisRepository.SetAdsData(adDataJson)

	return setStatus

}

func random_choice(data []int, weights []float64) int {
	// Calculate cumulative weights
	cumulativeWeights := make([]float64, len(weights))
	cumulativeWeights[0] = weights[0]
	for i := 1; i < len(weights); i++ {
		cumulativeWeights[i] = cumulativeWeights[i-1] + weights[i]
	}

	// Generate a random number between 0 and the sum of weights
	randomNumber := rand.Float64() * cumulativeWeights[len(cumulativeWeights)-1]

	// Find the index where the random number falls in the cumulative weights
	index := -1
	for i := 0; i < len(cumulativeWeights); i++ {
		if randomNumber <= cumulativeWeights[i] {
			index = i
			break
		}
	}

	// Retrieve the corresponding element from the data
	randomChoice := data[index]
	return randomChoice
}
func bsonMToMap(bm bson.M) map[string]interface{} {
	result := make(map[string]interface{})

	for key, value := range bm {
		switch v := value.(type) {
		case bson.M:
			// If the value is another bson.M, recursively convert it to a map
			result[key] = bsonMToMap(v)
		default:
			// Otherwise, simply assign the value to the result map
			result[key] = v
		}
	}

	return result
}
func generateRandomString(length int) string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
func (sv AdsService) Insert_ads(description string, nameAds string, speaker string, fileBytes []byte, filename string, tokendata string, isPreview bool) (Url_Image string, Url_Audio string, error error) {
	// สร้างชื่อโฆษณาหากไม่ระบุ
	if nameAds == "" {
		nameAds = generateRandomString(5)
	}

	// สร้าง pathName และ key สำหรับ S3
	pathName := "speaker_id_" + speaker + "/" + nameAds
	key, contentType, _ := sv.S3Provider.CreateKeyNameImageAds(pathName, strings.TrimPrefix(filename, "."))

	// อัปโหลดไฟล์ภาพไปยัง S3
	url, err := sv.S3Provider.UploadS3FromString(fileBytes, key, contentType)
	if err != nil {
		return "", "", fmt.Errorf("failed to upload image: %w", err)
	}

	// สร้างเสียงและอัปโหลดไปยัง S3
	audioContent, err := sv.generateVoice(description, speaker)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate voice: %w", err)
	}
	fileName := fmt.Sprintf("audio/speaker_id_%s/%s.%s", speaker, nameAds, "mp3")
	voiceURL, err := sv.S3Provider.UploadToS3(fileName, audioContent, "mp3")
	if err != nil {
		return "", "", fmt.Errorf("failed to upload audio: %w", err)
	}

	if isPreview {
		// หากเป็น Preview ให้ตั้งเวลาลบไฟล์
		go func() {
			time.Sleep(5 * time.Minute)

			// ลบไฟล์ภาพ
			err := sv.S3Provider.DeleteFromS3(key)
			if err != nil {
				fmt.Printf("Failed to delete preview image from S3: %v\n", err)
			}

			// ลบไฟล์เสียง
			err = sv.S3Provider.DeleteFromS3(fileName)
			if err != nil {
				fmt.Printf("Failed to delete preview audio from S3: %v\n", err)
			}
		}()
	} else {
		// หากไม่ใช่ Preview ให้บันทึกข้อมูลลงฐานข้อมูล
		Adsid, err := sv.AdsRepository.GetNextAdsID()
		if err != nil {
			return "", "", fmt.Errorf("failed to retrieve ID: %w", err)
		}
		newIDStr := strconv.Itoa(Adsid)

		// เตรียมข้อมูลสำหรับบันทึกในฐานข้อมูล
		data := &entities.AdsDataFormat{
			Id:          newIDStr, // ใช้ ID ที่เพิ่มค่าแล้ว
			Image_url:   url,
			Play:        0,
			Url:         voiceURL,
			Description: description,
			Level:       4, // ระดับที่กำหนด
			Language:    "th",
			Status:      true,
			Nameads:     nameAds,
		}

		// บันทึกข้อมูลลงฐานข้อมูล
		err = sv.AdsRepository.Insert_ads(data)
		if err != nil {
			return "", "", fmt.Errorf("failed to insert ad: %w", err)
		}
	}

	// คืนค่า URL ของภาพและเสียง
	return url, voiceURL, nil
}

func (sv AdsService) Delete_ads(ads_id string) error {
	// ค้นหาโฆษณาในฐานข้อมูล
	ad, err := sv.AdsRepository.FindOne(ads_id)
	if err != nil {
		return fmt.Errorf("failed to find ad: %w", err)
	}

	// ตรวจสอบว่า URL มีค่าหรือไม่ก่อนลบ
	if ad.Url != "" {
		audioKey, err := extractKeyFromURL(ad.Url)
		if err != nil {
			return fmt.Errorf("failed to extract audio key from URL: %w", err)
		}

		// ลบไฟล์เสียงจาก S3
		err = sv.S3Provider.DeleteFromS3(audioKey)
		if err != nil {
			return fmt.Errorf("failed to delete audio from S3: %w", err)
		}
	}

	if ad.Image_url != "" {
		imageKey, err := extractKeyFromURL(ad.Image_url)
		if err != nil {
			return fmt.Errorf("failed to extract image key from URL: %w", err)
		}

		// ลบไฟล์ภาพจาก S3
		err = sv.S3Provider.DeleteFromS3(imageKey)
		if err != nil {
			return fmt.Errorf("failed to delete image from S3: %w", err)
		}
	}

	// ลบข้อมูลโฆษณาจากฐานข้อมูล
	err = sv.AdsRepository.Delete_ads(ads_id)
	if err != nil {
		return fmt.Errorf("failed to delete ad from database: %w", err)
	}

	return nil
}

// ฟังก์ชันสำหรับแยก Key จาก URL
func extractKeyFromURL(inputURL string) (string, error) {
	// ตรวจสอบว่า URL มี hostname และ path
	if !strings.Contains(inputURL, "amazonaws.com/") {
		return "", fmt.Errorf("invalid S3 URL")
	}

	// แยกเอาเฉพาะส่วน Path หลังจาก "amazonaws.com/"
	parts := strings.Split(inputURL, "amazonaws.com/")
	if len(parts) < 2 {
		return "", fmt.Errorf("failed to extract key from URL")
	}

	// คืนค่า Key ที่เหลือหลัง "amazonaws.com/"
	return parts[1], nil
}

func (sv AdsService) Update_ads(ads_id string, description string, nameAds string, speaker string, fileBytes []byte, fileExtension string) (Url_Image string, Url_Audio string, error error) {

	// ค้นหาโฆษณาใน repository
	ad, err := sv.AdsRepository.FindOne(ads_id)
	if err != nil {
		return "", "", fmt.Errorf("failed to find ad: %w", err)
	}
	if nameAds == "" {
		nameAds = generateRandomString(5)
	}
	// ตรวจสอบว่า description มีการเปลี่ยนแปลงหรือไม่
	var voiceURL string
	if ad.Description != description {
		// สร้าง voice content ใหม่เฉพาะเมื่อ description เปลี่ยนแปลง
		audioContent, err := sv.generateVoice(description, speaker)
		if err != nil {
			return "", "", fmt.Errorf("failed to generate voice: %w", err)
		}

		// อัปโหลดไฟล์เสียงไปยัง S3
		fileName := fmt.Sprintf("audio/speaker_id_%s/%s.%s", speaker, nameAds, "mp3")
		voiceURL, err = sv.S3Provider.UploadToS3(fileName, audioContent, "mp3")
		if err != nil {
			return "", "", fmt.Errorf("failed to upload audio file: %w", err)
		}
	} else {
		// ใช้ voice URL เดิมหาก description ไม่เปลี่ยนแปลง
		voiceURL = ad.Url
	}

	// ตรวจสอบ fileBytes และอัปโหลดรูปภาพหากจำเป็น
	var imageURL string
	if len(fileBytes) > 0 {
		// ดำเนินการอัปโหลดรูปภาพหาก fileBytes ไม่ใช่ค่าว่าง
		pathName := "speaker_id_" + speaker + "/" + nameAds
		key, contentType, _ := sv.S3Provider.CreateKeyNameImageAds(pathName, strings.TrimPrefix(fileExtension, "."))
		imageURL, err = sv.S3Provider.UploadS3FromString(fileBytes, key, contentType)
		if err != nil {
			return "", "", fmt.Errorf("failed to upload image file: %w", err)
		}
	} else {
		// ใช้ image URL เดิมหากไม่มีการอัปโหลดไฟล์ใหม่
		imageURL = ad.Image_url
	}

	// สร้างโครงสร้างข้อมูลโฆษณาใหม่
	data := &entities.AdsDataFormat{
		Url:         voiceURL,
		Description: description,
		Image_url:   imageURL,
		Nameads:     nameAds,
	}

	// อัปเดตโฆษณาใน repository
	err = sv.AdsRepository.Update_ads(*data, ads_id)
	if err != nil {
		return "", "", fmt.Errorf("failed to update ad: %w", err)
	}

	// ส่งกลับ URL ของ voice และ image
	return voiceURL, imageURL, nil
}

func (sv AdsService) Getall_ads() ([]entities.AdsDataFormat, error) {
	// Call the repository to get all ads

	data, err := sv.AdsRepository.Getall_ads()
	if err != nil {
		return nil, fmt.Errorf("failed to get all ads: %w", err)
	}

	// Return the data and nil error
	return data, nil
}
func (sv AdsService) FindOne(adsID string) (entities.AdsDataFormat, error) {
	// Call the repository to get a single ad by ID
	data, err := sv.AdsRepository.FindOne(adsID)
	if err != nil {
		return entities.AdsDataFormat{}, fmt.Errorf("failed to get ad with ID %s: %w", adsID, err)
	}

	// Return the data and nil error
	return data, nil
}
func (sv AdsService) Find_one_random_ads() (entities.AdsDataFormat, error) {
	// เรียก repository เพื่อดึงข้อมูลโฆษณาทั้งหมด
	allAds, err := sv.AdsRepository.Getall_ads()
	if err != nil {
		return entities.AdsDataFormat{}, fmt.Errorf("ไม่สามารถดึงข้อมูลโฆษณาได้: %w", err)
	}

	// ตรวจสอบว่ามีโฆษณาหรือไม่
	if len(allAds) == 0 {
		return entities.AdsDataFormat{}, fmt.Errorf("ไม่มีโฆษณาให้ใช้งาน")
	}

	// สร้างตัวเลขสุ่ม
	rand.Seed(time.Now().UnixNano())      // ตั้งค่า seed สำหรับตัวเลขสุ่ม
	randomIndex := rand.Intn(len(allAds)) // สุ่ม index

	// คืนค่าข้อมูลโฆษณาแบบสุ่ม
	return allAds[randomIndex], nil
}
func (sv AdsService) generateVoice(description string, speaker string) ([]byte, error) {
	// แปลง speaker เป็น int พร้อมตรวจสอบข้อผิดพลาด
	speakerint, err := strconv.Atoi(speaker)
	if err != nil {
		return nil, fmt.Errorf("invalid speaker value: %v", err)
	}

	// สร้าง TTSRequest payload
	utilsMsgData := providers.TTSRequest{
		Text:           description,
		Speaker:        speakerint,
		Format:         "mp3",
		SentencePause:  0.5,
		ParagraphPause: 0.3,
		Volume:         "1",
		Speed:          "1",
		Language:       "th",
	}

	// เรียกใช้ API
	audioContent, ttsStatusCode, err := providers.CallTTSAPI(utilsMsgData)
	if err != nil {
		switch ttsStatusCode {
		case 429:
			return nil, errors.New("queue is full")
		case 400:
			if os.Getenv("LANGID") == "active" {
				// Re-perform language detection if needed
				// Add logic here if required
			}
			return nil, errors.New("invalid language for this speaker")
		case 403:
			// จัดการข้อความ error
			errorMessage := err.Error()
			var detail string
			if strings.Contains(errorMessage, "{\"detail\":") {
				idx := strings.Index(errorMessage, "{\"detail\":")
				jsonStr := errorMessage[idx:]
				var errorResp struct {
					Detail string `json:"detail"`
				}
				if json.Unmarshal([]byte(jsonStr), &errorResp) == nil {
					detail = errorResp.Detail
				} else {
					detail = "text error"
				}
			} else {
				detail = "text error"
			}

			// บันทึกข้อความที่ผิดพลาดลงใน repository
			wrongMessage := entities.WrongMessage{
				Language: "th",
				Message:  description,
				Speaker:  fmt.Sprintf("%d", speakerint), // ใช้ fmt.Sprintf แทน string()
				Datetime: entities.DateTimeBangkok(),
			}
			if err := sv.wrongMessageRepo.InsertWrongMessage(wrongMessage); err != nil {
				// Log repository error
				fmt.Printf("Failed to insert wrong message: %v\n", err)
			}

			return nil, errors.New(detail)
		default:
			return nil, fmt.Errorf("TTS API call failed: %v", err)
		}
	}

	return audioContent, nil
}
