package database

import (
	"context"
	"dfs/internal/config"
	"dfs/internal/interfaces"
	"dfs/internal/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type FileRepository struct {
	client *mongo.Client
}

func (repository *FileRepository) Create(file *models.FileMetadata) (any, error) {
	collection := repository.client.Database("dfs").Collection("files")
	_, errors := collection.InsertOne(context.Background(), file)
	if errors != nil {
		return nil, errors
	}
	return nil, nil
}

func (repository *FileRepository) GetAll() (*[]*models.FileMetadata, error) {
	collection := repository.client.Database("dfs").Collection("files")
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	var files []*models.FileMetadata
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var file models.FileMetadata
		cursor.Decode(&file)
		files = append(files, &file)
	}
	return &files, nil
}

func (repository *FileRepository) GetByID(id string) (*models.FileMetadata, error) {
	var file models.FileMetadata
	collection := repository.client.Database("dfs").Collection("files")
	err := collection.FindOne(context.Background(), bson.M{"fileIdentifier": id}).Decode(&file)
	if err != nil {
		return nil, err
	}
	return &file, nil
}

func NewRepository(env *config.Environment) interfaces.FileRepositoryInterface {
	clientOptions := options.Client().ApplyURI(env.MongoURI)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoClient, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		panic(err)
	}
	err = mongoClient.Ping(ctx, nil)
	if err != nil {
		panic(err)
	}

	return &FileRepository{
		client: mongoClient,
	}

}
