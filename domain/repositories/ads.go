package repositories

import (
	"bn-crud-ads/domain/datasources"
	"bn-crud-ads/domain/entities"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	fiberlog "github.com/gofiber/fiber/v2/log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type adsRepository struct {
	Context    context.Context
	Collection *mongo.Collection
	AlertColl  *mongo.Collection
}

type IAdsRepository interface {
	FindAds(level int) (*bson.M, error)
	Update(data entities.AdsModel) (string, error)
	FindMkpSound() ([]bson.M, error)
	FindAdsAlertMessage() (*entities.AdsAlertMessage, error)
	FindAds_ID() (int, error)
	Insert_ads(data *entities.AdsDataFormat) error
	Delete_ads(ads_id string) error
	Update_ads(data entities.AdsDataFormat, ads_id string) error
	FindOne(ads_id string) (data entities.AdsDataFormat, err error)
	Getall_ads() ([]entities.AdsDataFormat, error)
	GetNextAdsID() (int, error)
}

func (repo adsRepository) GetNextAdsID() (int, error) {
	const maxRetries = 3
	var retryCount = 0
	for {
		result := repo.Collection.FindOneAndUpdate(
			repo.Context,
			bson.M{"_id": "Ads_id"},
			bson.M{"$inc": bson.M{"seq": 1}},
			options.FindOneAndUpdate().SetUpsert(true),
		)

		if result.Err() != nil {
			if retryCount < maxRetries {
				retryCount++
				time.Sleep(time.Second * time.Duration(retryCount))
				fmt.Printf("Retrying... Attempt %d\n", retryCount)
				continue
			}
			return 0, result.Err()
		}

		var counterDoc struct{ Seq int }
		err := result.Decode(&counterDoc)
		if err != nil {
			if retryCount < maxRetries {
				retryCount++
				time.Sleep(time.Second * time.Duration(retryCount))
				fmt.Printf("Retrying... Attempt %d\n", retryCount)
				continue
			}
			return 0, err
		}
		return counterDoc.Seq, nil
	}
}
func NewAdsRepository(db *datasources.MongoDB) IAdsRepository {
	return &adsRepository{
		Context:    db.Context,
		Collection: db.MongoDB.Database(os.Getenv("DATABASE_NAME")).Collection("ads"),
		AlertColl:  db.MongoDB.Database(os.Getenv("DATABASE_NAME")).Collection("alert_message"),
	}
}

func (repo adsRepository) FindAds(level int) (*bson.M, error) {

	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: bson.D{{Key: "level", Value: level}, {Key: "status", Value: true}}}},
		bson.D{{Key: "$sample", Value: bson.D{{Key: "size", Value: 1}}}},
	}

	cursor, err := repo.Collection.Aggregate(repo.Context, pipeline)
	if err != nil {
		log.Fatal(err)
	}

	defer cursor.Close(repo.Context)

	pack := make([]*bson.M, 0)
	for cursor.Next(repo.Context) {
		var item *bson.M
		err := cursor.Decode(&item)
		if err != nil {
			continue
		}

		pack = append(pack, item)
	}
	return pack[0], nil

}

func (repo adsRepository) Update(data entities.AdsModel) (string, error) {
	for _, updateData := range data.AdsPlay {

		filter := bson.M{"id": updateData.ID}
		update := bson.M{"$inc": bson.M{"play": updateData.Play}}

		_, err := repo.Collection.UpdateMany(repo.Context, filter, update)
		if err != nil {
			panic(err)
		}
	}

	return "success update ads", nil
}

func (repo adsRepository) FindMkpSound() ([]bson.M, error) {
	projection := bson.M{"_id": 0, "play": 0}
	var opt = options.Find().SetProjection(projection)
	filter := bson.M{"url": bson.M{"$exists": true, "$ne": nil}}
	cursor, err := repo.Collection.Find(repo.Context, filter, opt)
	if err != nil {
		fiberlog.Errorf("Ads -> FindMkpSound: %s \n", err)
		return nil, err
	}

	defer cursor.Close(repo.Context)

	pack := make([]bson.M, 0)
	for cursor.Next(repo.Context) {
		var item bson.M
		err := cursor.Decode(&item)
		if err != nil {
			continue
		}

		pack = append(pack, item)
	}
	return pack, nil
}

func (repo adsRepository) FindAdsAlertMessage() (*entities.AdsAlertMessage, error) {

	// Retrieves the first matching document
	projection := bson.M{"ads": true, "_id": false}
	opts := options.FindOne().SetProjection(projection)
	var result *entities.AdsAlertMessage
	err := repo.AlertColl.FindOne(repo.Context, bson.M{}, opts).Decode(&result)
	// Prints a message if no documents are matched or if any
	// other errors occur during the operation
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, err
		}
		panic(err)
	}
	return result, nil
}
func (repo adsRepository) FindAds_ID() (int, error) {
	// Aggregation pipeline: แปลง id เป็นตัวเลขและหาค่าสูงสุด
	pipeline := mongo.Pipeline{
		{
			{"$project", bson.D{
				{"idInt", bson.D{{"$toInt", "$id"}}}, // แปลง id เป็น int
				{"_id", 0},
			}},
		},
		{
			{"$sort", bson.D{{"idInt", -1}}}, // เรียงลำดับจากมากไปน้อย
		},
		{
			{"$limit", 1}, // เลือกเอกสารแรกสุด
		},
	}

	cursor, err := repo.Collection.Aggregate(repo.Context, pipeline)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(repo.Context)

	// อ่านเอกสารที่ได้
	var results []struct {
		IDInt int `bson:"idInt"`
	}
	if err := cursor.All(repo.Context, &results); err != nil {
		return 0, err
	}

	if len(results) == 0 {
		return 0, fmt.Errorf("no documents found")
	}

	return results[0].IDInt, nil
}
func (repo adsRepository) Insert_ads(data *entities.AdsDataFormat) error {
	// เพิ่มข้อมูลลงในคอลเลกชัน
	_, err := repo.Collection.InsertOne(repo.Context, data)
	if err != nil {
		return fmt.Errorf("failed to insert Ads: %v", err)
	}
	return nil
}
func (repo adsRepository) Delete_ads(ads_id string) error {
	filter := bson.M{"id": ads_id}
	_, err := repo.Collection.DeleteOne(repo.Context, filter)
	if err != nil {
		return fmt.Errorf("failed to delete Ads")
	}
	return nil
}
func (repo *adsRepository) Update_ads(data entities.AdsDataFormat, ads_id string) error {
	filter := bson.M{"id": ads_id}
	update := bson.M{
		"$set": bson.M{
			"description": data.Description,
			"url":         data.Url,
			"image_url":   data.Image_url,
		},
	}
	_, err := repo.Collection.UpdateOne(repo.Context, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update blog ")
	}

	return nil
}
func (repo adsRepository) FindOne(ads_id string) (data entities.AdsDataFormat, err error) {
	filter := bson.M{"id": ads_id}
	if err := repo.Collection.FindOne(repo.Context, filter).Decode(&data); err != nil {
		return data, fmt.Errorf("failed to find ads")
	}
	return data, nil
}
func (repo adsRepository) Getall_ads() ([]entities.AdsDataFormat, error) {
	var data []entities.AdsDataFormat

	findOptions := options.Find().SetSkip(1)
	cursor, err := repo.Collection.Find(repo.Context, bson.M{}, findOptions)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch all ads: %w", err)
	}
	defer func() {
		if err := cursor.Close(repo.Context); err != nil {
			fmt.Printf("Failed to close cursor: %v\n", err)
		}
	}()

	// ดึงข้อมูลและ decode ทีละ record
	for cursor.Next(repo.Context) {
		var ad entities.AdsDataFormat
		if err := cursor.Decode(&ad); err != nil {
			return nil, fmt.Errorf("failed to decode ad: %w", err)
		}
		data = append(data, ad)
	}

	// ตรวจสอบข้อผิดพลาดในการ iteration
	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor iteration error: %w", err)
	}

	return data, nil
}
