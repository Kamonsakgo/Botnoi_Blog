package repositories

import (
	"bn-crud-ads/domain/datasources"
	"bn-crud-ads/domain/entities"
	"context"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
)

type IWrongMessageRepository interface {
	InsertWrongMessage(data entities.WrongMessage) error
}

type wrongMessageRepository struct {
	Context    context.Context
	Collection *mongo.Collection
}

func NewWrongMessageRepository(db *datasources.MongoDB) IWrongMessageRepository {
	return &wrongMessageRepository{
		Context:    db.Context,
		Collection: db.MongoDB.Database(os.Getenv("DATABASE_NAME")).Collection("wrong_message"),
	}
}

func (repo *wrongMessageRepository) InsertWrongMessage(data entities.WrongMessage) error {
	_, err := repo.Collection.InsertOne(repo.Context, data)
	return err
}
