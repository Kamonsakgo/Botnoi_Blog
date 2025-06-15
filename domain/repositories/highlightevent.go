package repositories

import (
	"context"
	"fmt"
	"go-fiber-unittest/domain/datasources"
	"go-fiber-unittest/domain/entities"
	"os"
	"time"

	//"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type HighlightsRepository struct {
	Context    context.Context
	Collection *mongo.Collection
}

type IHighlightsRepository interface {
	Insert(data entities.HighlightModel) error
	FindOne(highlightID string) (*entities.HighlightModel, error)
	Delete(highlight_id string) error
	Update(data entities.HighlightModel, highlight_id string) error
	FindAll(page int, limit int) ([]entities.HighlightModel, int64, error)
}

func NewHighlightsRepository(db *datasources.MongoDB) IHighlightsRepository {
	return &HighlightsRepository{
		Context:    db.Context,
		Collection: db.MongoDB.Database(os.Getenv("DATABASE_NAME")).Collection("HL_events"),
	}
}

func (repo HighlightsRepository) Insert(data entities.HighlightModel) error {
	if _, err := repo.Collection.InsertOne(repo.Context, data); err != nil {
		return err
	}

	return nil
}
func (repo HighlightsRepository) FindOne(highlightID string) (*entities.HighlightModel, error) {
	// Find one highlight by highlight_id
	filter := bson.M{"highlight_id": highlightID}
	var result entities.HighlightModel
	err := repo.Collection.FindOne(repo.Context, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("highlight with id '%s' not found", highlightID)
		}
		return nil, err
	}

	return &result, nil
}
func (repo HighlightsRepository) Delete(highlight_id string) error {
	filter := bson.M{"highlight_id": highlight_id}
	_, err := repo.Collection.DeleteOne(repo.Context, filter)
	if err != nil {
		return fmt.Errorf("failed to delete highlight")
	}

	return nil
}
func (repo HighlightsRepository) Update(data entities.HighlightModel, highlight_id string) error {
	filter := bson.M{"highlight_id": highlight_id}
	update := bson.M{
		"$set": bson.M{
			"title":          data.Title,
			"category":       data.Category,
			"location":       data.Location,
			"last_update_at": data.LastUpdateAt,
			"location_event": data.LocationEvent,
			"speaker":        data.Speaker,
			"date":           data.Date,
			"image_url":      data.ImageURL,
		},
	}
	_, err := repo.Collection.UpdateOne(repo.Context, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update highlight ")
	}

	return nil
}
func (repo *HighlightsRepository) FindAll(page int, limit int) ([]entities.HighlightModel, int64, error) {
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
	var highlights []entities.HighlightModel
	for cursor.Next(ctx) {
		var highlight entities.HighlightModel
		if err := cursor.Decode(&highlight); err != nil {
			// Log the error and skip to the next document.
			// You might want to use a logging library here.
			continue
		}
		highlights = append(highlights, highlight)
	}

	// Check for any errors encountered during iteration.
	if err := cursor.Err(); err != nil {
		return nil, 0, err
	}

	// Return the found blogs and total count
	return highlights, total, nil
}
