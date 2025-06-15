package providers

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"go-fiber-unittest/domain/entities"
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
	CreateKeyNameImageBlog(fileName string, typeFile string) (keyName string, contentType string)
	DeleteKeyNameImage(url string) error
	CreateKeyNameImageHL(typeFile string) (keyName string, contentType string)
	DeleteKeyNameImageHL(keyName string) error
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

func (s *S3Provider) UploadS3FromString(fileName []byte, keyName string, contentType string) (string, error) {
	_, err := s.Svc.PutObject(&s3.PutObjectInput{
		Body:         bytes.NewReader(fileName),
		Bucket:       aws.String(os.Getenv("S3_PICTURE_BUCKET_NAME")),
		Key:          aws.String(keyName),
		ContentType:  aws.String("image/svg+xml"),
		Metadata:     map[string]*string{"Content-Disposition": aws.String("inline")}, // เปลี่ยนเป็น inline
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
func (s *S3Provider) CreateKeyNameImageBlog(pathName string, typeFile string) (keyName string, contentType string) {
	randomBlogID := s.GenerateRandomNumber()
	keyName = fmt.Sprintf("blogs/generals/%v/%v.%v", pathName, randomBlogID, typeFile)
	contentType = fmt.Sprintf("image/%v", typeFile)
	return keyName, contentType
}

func CreateKeyNameImage(data entities.BlogModel, imageType string, fileName string) (string, string) {

	blogid := strings.ToLower(data.BlogID)

	keyName := fmt.Sprintf("picture/%v/face_%v.%v", blogid, blogid, imageType)
	if fileName != "" {
		keyName = fmt.Sprintf("picture/%v/face_%v.%v", blogid, blogid, imageType)
	}
	contentType := fmt.Sprintf("image/%v", imageType)
	return keyName, contentType
}

func (s *S3Provider) DeleteKeyNameImage(url string) error {

	baseURL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/", os.Getenv("S3_PICTURE_BUCKET_NAME"), os.Getenv("AWS_REGION"))
	keyname := strings.TrimPrefix(url, baseURL)
	fullKey := "blogs/generals/" + keyname

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
func (s *S3Provider) CreateKeyNameImageHL(typeFile string) (keyName string, contentType string) {
	randomHLID := s.GenerateRandomNumber()
	keyName = fmt.Sprintf("highlight/%v.%v", randomHLID, typeFile)
	contentType = fmt.Sprintf("image/%v", typeFile)
	return keyName, contentType
}
func (s *S3Provider) DeleteKeyNameImageHL(keyName string) error {
	// สร้าง full key ที่จะใช้ในการลบ
	fullKey := "highlight/" + keyName

	// ตรวจสอบว่ามี object ใน S3 ที่ตรงกับ fullKey หรือไม่
	resp, err := s.Svc.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(os.Getenv("S3_PICTURE_BUCKET_NAME")),
		Prefix: aws.String(fullKey),
	})
	if err != nil {
		return fmt.Errorf("failed to list objects in S3: %v", err)
	}

	// ถ้ามี objects ที่ตรงกับ fullKey ให้ลบออก
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
