package repositories

import (
	"context"
	"fmt"
	"go-fiber-unittest/domain/datasources"
	"go-fiber-unittest/domain/entities"
	"os"

	//"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BlogsRepository struct {
	Context    context.Context
	Collection *mongo.Collection
}

type IBlogsRepository interface {
	Insert(data entities.BlogModel) error
	GetNextBlogID() (int, error)
	Update(data entities.BlogModel, blog_id string) error
	FindAll(page int, limit int) ([]entities.BlogModel, int64, error)
	Delete(blogid string) error
	FindOne(blogid string) (data entities.BlogModel, err error)
	UploadImage(blogid string, name string, url string) error
}

func NewBlogsRepository(db *datasources.MongoDB) IBlogsRepository {
	return &BlogsRepository{
		Context:    db.Context,
		Collection: db.MongoDB.Database(os.Getenv("DATABASE_NAME")).Collection("blogs"),
	}
}

func (repo BlogsRepository) GetNextBlogID() (int, error) {
	const maxRetries = 3
	var retryCount = 0
	for {
		result := repo.Collection.FindOneAndUpdate(
			repo.Context,
			bson.M{"_id": "blog_id"},
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

func (repo BlogsRepository) Insert(data entities.BlogModel) error {

	if _, err := repo.Collection.InsertOne(repo.Context, data); err != nil {
		return err
	}

	return nil
}
func (repo *BlogsRepository) Update(data entities.BlogModel, blog_id string) error {
	filter := bson.M{"blog_id": blog_id}

	update := bson.M{
		"$set": bson.M{
			"title":          data.Title,
			"content":        data.Content,
			"category":       data.Category,
			"tag":            data.Tag,
			"type":           data.Type,
			"location":       data.Location,
			"last_update_at": data.LastUpdateAt,
			"HL_id":          data.HL_id,
		},
	}
	_, err := repo.Collection.UpdateOne(repo.Context, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update blog ")
	}

	return nil
}
func (repo BlogsRepository) FindOne(blogid string) (data entities.BlogModel, err error) {
	filter := bson.M{"blog_id": blogid}
	if err := repo.Collection.FindOne(repo.Context, filter).Decode(&data); err != nil {
		return data, fmt.Errorf("failed to find blog")
	}

	return data, nil
}

func (repo BlogsRepository) Delete(blogid string) error {
	filter := bson.M{"blog_id": blogid}
	_, err := repo.Collection.DeleteOne(repo.Context, filter)
	if err != nil {
		return fmt.Errorf("failed to delete blog")
	}

	return nil
}

func (repo *BlogsRepository) FindAll(page int, limit int) ([]entities.BlogModel, int64, error) {
	// Create a context with a timeout to avoid hanging operations.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Calculate skip and limit values for MongoDB pagination
	skip := int64((page - 1) * limit)
	limit64 := int64(limit)

	// Define find options with skip and limit for pagination.
	findOptions := options.Find()
	findOptions.SetSkip(skip)
	findOptions.SetLimit(limit64)

	// Retrieve the total count of documents in the collection
	total, err := repo.Collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, err
	}

	// Execute the Find operation with the defined options.
	cursor, err := repo.Collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	// Iterate through the cursor and decode each document into a BlogModel.
	var blogs []entities.BlogModel
	for cursor.Next(ctx) {
		var blog entities.BlogModel
		if err := cursor.Decode(&blog); err != nil {
			// Log the error and skip to the next document.
			// You might want to use a logging library here.
			continue
		}
		blogs = append(blogs, blog)
	}

	// Check for any errors encountered during iteration.
	if err := cursor.Err(); err != nil {
		return nil, 0, err
	}

	// Return the found blogs and total count
	return blogs, total, nil
}

func (repo BlogsRepository) UploadImage(blogid string, name string, url string) error {
	newImage := bson.M{
		"name": name,
		"url":  url,
	}
	filter := bson.M{"blog_id": blogid}

	update := bson.M{
		"$push": bson.M{"image_url": newImage},
	}

	_, err := repo.Collection.UpdateOne(repo.Context, filter, update)
	if err != nil {
		return fmt.Errorf("failed to upload image")
	}

	return nil
}
