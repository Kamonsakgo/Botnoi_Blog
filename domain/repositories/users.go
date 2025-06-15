package repositories

import (
	"context"
	"fmt"
	. "go-fiber-unittest/domain/datasources"
	"go-fiber-unittest/domain/entities"
	"os"

	fiberlog "github.com/gofiber/fiber/v2/log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type usersRepository struct {
	Context    context.Context
	Collection *mongo.Collection
}

type IUsersRepository interface {
	InsertNewUser(data *entities.NewUserBody) bool
	FindAll() ([]entities.UserDataFormat, error)
	FindByID(userID string) (*entities.UserDataFormat, error)
}

func NewUsersRepository(db *MongoDB) IUsersRepository {
	return &usersRepository{
		Context:    db.Context,
		Collection: db.MongoDB.Database(os.Getenv("DATABASE_NAME")).Collection("users"),
	}
}

func (repo usersRepository) InsertNewUser(data *entities.NewUserBody) bool {
	if _, err := repo.Collection.InsertOne(repo.Context, data); err != nil {
		fiberlog.Errorf("Users -> InsertNewUser: %s \n", err)
		return false
	}
	return true
}

func (repo usersRepository) FindAll() ([]entities.UserDataFormat, error) {
	options := options.Find()
	filter := bson.M{}
	cursor, err := repo.Collection.Find(repo.Context, filter, options)
	if err != nil {
		fiberlog.Errorf("Users -> FindAll: %s \n", err)
		return nil, err
	}
	defer cursor.Close(repo.Context)
	pack := make([]entities.UserDataFormat, 0)
	for cursor.Next(repo.Context) {
		var item entities.UserDataFormat

		err := cursor.Decode(&item)
		if err != nil {
			continue
		}

		pack = append(pack, item)
	}
	return pack, nil
}
func (repo usersRepository) FindByID(userID string) (*entities.UserDataFormat, error) {

	result := entities.UserDataFormat{}
	user := repo.Collection.FindOne(repo.Context, bson.M{"user_id": userID}).Decode(&result)

	if user == mongo.ErrNoDocuments {
		return nil, fmt.Errorf("User not found")
	}
	return &result, nil
}
