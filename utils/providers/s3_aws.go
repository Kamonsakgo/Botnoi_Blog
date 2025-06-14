package providers

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"golang.org/x/exp/rand"
)

type S3Provider struct {
	Svc *s3.S3
}

type IS3provider interface {
	UploadS3FromString(fileName []byte, keyName string, contentType string) (string, error)
	HashString(userID string) (string, error)
	CreateKeyNameImageAds(fileName string, typeFile string) (keyName string, contentType string, error error)
	DeleteKeyNameImage(url string) error
	CreateKeyNameImageAudio(pathName string, typeFile string) (string, string, error)
	UploadToS3(keyName string, contentBytes []byte, fileType string) (string, error)
	DeleteFromS3(key string) error
	// CreateKeyNameImage(data entities.SpeakerModel, imageType string, fileName string) (string, string)
}

func NewS3Provider() IS3provider {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(os.Getenv("AWS_REGION")),
		Credentials: credentials.NewStaticCredentials(os.Getenv("AWS_ACCESS_KEY_ID"), os.Getenv("AWS_SECRET_ACCESS_KEY"), ""),
	})
	if err != nil {
		fmt.Println("failed to connect aws", err)
	}
	svc := s3.New(sess)
	return &S3Provider{
		Svc: svc,
	}
}
func (sp *S3Provider) UploadToS3(keyName string, contentBytes []byte, fileType string) (string, error) {
	// โหลด bucket name และ region จาก environment variable
	bucket := os.Getenv("S3_PICTURE_BUCKET_NAME")
	region := os.Getenv("AWS_REGION")
	if bucket == "" || region == "" {
		return "", fmt.Errorf("missing required environment variables: S3_PICTURE_BUCKET_NAME or AWS_REGION")
	}

	// กำหนด content type
	contentType := "audio/" + fileType

	// อัพโหลดไฟล์ไปที่ S3
	_, err := sp.Svc.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(keyName),
		Body:        bytes.NewReader(contentBytes),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		fmt.Printf("Failed to upload to S3: %v\n", err)
		return "", err
	}

	// สร้าง public URL
	voiceURL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", bucket, region, keyName)
	return voiceURL, nil
}
func (s *S3Provider) UploadS3FromString(fileName []byte, keyName string, contentType string) (string, error) {
	_, err := s.Svc.PutObject(&s3.PutObjectInput{
		Body:         bytes.NewReader(fileName),
		Bucket:       aws.String(os.Getenv("S3_PICTURE_BUCKET_NAME")),
		Key:          aws.String(keyName),
		ContentType:  aws.String(contentType),
		Metadata:     map[string]*string{"Content-Disposition": aws.String("inline")},
		ACL:          aws.String("public-read"),
		CacheControl: aws.String("no-cache"),
	})
	if err != nil {
		return "", err
	}

	req, _ := s.Svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(os.Getenv("S3_PICTURE_BUCKET_NAME")),
		Key:    aws.String(keyName),
	})

	_, err = req.Presign(100 * 365 * 24 * time.Hour) // 100 years expiration
	if err != nil {
		return "", err
	}

	fullURL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", os.Getenv("S3_PICTURE_BUCKET_NAME"), os.Getenv("AWS_REGION"), keyName)

	return fullURL, nil
}

func (s *S3Provider) HashString(userID string) (string, error) {
	byteID := []byte(userID)
	hash := sha256.New()
	hash.Write(byteID)
	hashed := hash.Sum(nil)
	hexDigest := hex.EncodeToString(hashed)
	return hexDigest, nil
}
func (s *S3Provider) GenerateRandomNumber() string {

	rand.Seed(uint64(time.Now().UnixNano()))
	randomNumber := rand.Intn(9000000) + 1000000
	return fmt.Sprintf("%d", randomNumber)
}
func (s *S3Provider) CreateKeyNameImageAds(pathName string, typeFile string) (string, string, error) {
	// Create the key name for the image ad
	keyName := fmt.Sprintf("Ads/%v.%v", pathName, typeFile)

	// Determine the content type based on the file extension
	var contentType string
	switch strings.ToLower(typeFile) { // Convert typeFile to lowercase
	case "jpg", "jpeg":
		contentType = "image/jpeg"
	case "png":
		contentType = "image/png"
	case "svg":
		contentType = "image/svg+xml"
	default:
		// Return an error for unsupported file types
		return "", "", fmt.Errorf("unsupported file extension: %s", typeFile)
	}

	// Return the key name and content type
	return keyName, contentType, nil
}
func (s *S3Provider) CreateKeyNameImageAudio(pathName string, typeFile string) (string, string, error) {
	// Create the key name for the image ad
	keyName := fmt.Sprintf("Audio/%v.%v", pathName, typeFile)

	// Determine the content type based on the file extension
	var contentType string
	switch strings.ToLower(typeFile) { // Convert typeFile to lowercase
	case "jpg", "jpeg":
		contentType = "image/jpeg"
	case "png":
		contentType = "image/png"
	case "svg":
		contentType = "image/svg+xml"
	default:
		// Return an error for unsupported file types
		return "", "", fmt.Errorf("unsupported file extension: %s", typeFile)
	}

	// Return the key name and content type
	return keyName, contentType, nil
}
func (s *S3Provider) DeleteKeyNameImage(url string) error {

	baseURL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/", os.Getenv("S3_PICTURE_BUCKET_NAME"), os.Getenv("AWS_REGION"))
	keyname := strings.TrimPrefix(url, baseURL)
	fullKey := "Ads/" + keyname

	resp, err := s.Svc.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(os.Getenv("S3_PICTURE_BUCKET_NAME")),
		Prefix: aws.String(fullKey),
	})
	if err != nil {
		return fmt.Errorf("failed to list objects in S3: %v", err)
	}

	if len(resp.Contents) > 0 {
		var objectsToDelete []*s3.ObjectIdentifier
		for _, object := range resp.Contents {
			objectsToDelete = append(objectsToDelete, &s3.ObjectIdentifier{
				Key: object.Key,
			})
		}
		_, err := s.Svc.DeleteObjects(&s3.DeleteObjectsInput{
			Bucket: aws.String(os.Getenv("S3_PICTURE_BUCKET_NAME")),
			Delete: &s3.Delete{
				Objects: objectsToDelete,
			},
		})
		if err != nil {
			return fmt.Errorf("failed to delete objects: %v", err)
		}
	}

	return nil
}
func (sp *S3Provider) DeleteFromS3(key string) error {
	// ลบวัตถุ (Object) จาก S3
	_, err := sp.Svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(os.Getenv("S3_PICTURE_BUCKET_NAME")), // ดึงชื่อ Bucket จาก Environment Variable
		Key:    aws.String(key),                                 // ระบุ Key ของไฟล์ที่ต้องการลบ
	})
	if err != nil {
		return fmt.Errorf("failed to delete object from S3: %w", err)
	}

	// ยืนยันการลบไฟล์
	err = sp.Svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(os.Getenv("S3_PICTURE_BUCKET_NAME")),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("failed to confirm deletion of object from S3: %w", err)
	}

	return nil
}
