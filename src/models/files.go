package models

import (
	"context"
	config "dfs/src/config"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)


type FileMetadata struct {
	Id             primitive.ObjectID     `bson:"_id" json:"_id,omitempty"`
	FileName       string                 `bson:"fileName" json:"fileName"`
	FileSize       uint64                 `bson:"fileSize" json:"fileSize"`
	NumberOfChunks uint64                 `bson:"numberOfChunks" json:"numberOfChunks"`
	Replicas       [][]string             `bson:"replicas" json:"replicas"`
	ClientId       string                 `bson:"clientId" json:"clientId"`
	SharedUser     []primitive.ObjectID   `bson:"sharedUser" json:"sharedUser"`
    FileIdentifier string                 `bson:"fileIdentifier" json:"fileIdentifier"`
}

func (file *FileMetadata) CreateFile() error {
    log.Println("======== Getting mongo client to insert ==============")
    client, err := config.GetMongoClient()
    log.Println("======== Got mongo client to insert ==============")
    if err != nil {
        return err
    }
    collection := client.Database("dfs").Collection("files")
    fi, errors := collection.InsertOne(context.Background(), file)
    log.Println("======== Inserted file ==============")
    log.Println(fi)
    if errors != nil {
        return errors
    }
    return nil
}


func GetAllFiles() ([]FileMetadata, error) {
    var files []FileMetadata
    client, err := config.GetMongoClient()
    if err != nil {
        return files, err
    }
    collection := client.Database("dfs").Collection("files")
    cursor, err := collection.Find(context.Background(), bson.M{})
    if err != nil {
        return files, err
    }
    defer cursor.Close(context.Background())
    for cursor.Next(context.Background()) {
        var file FileMetadata
        cursor.Decode(&file)
        files = append(files, file)
    }
    return files, nil

}

func GetFileById(id string) (FileMetadata, error) {
    
    log.Println("======== Getting mongo client to get file ==============")
    var file FileMetadata
    client, err := config.GetMongoClient()
    log.Println("======== Got mongo client to get file ==============")
    if err != nil {
        return file, err
    }
    log.Println("======== Getting collection to get file ==============")
    log.Println(id)
    if err != nil {
        log.Println("======== Error getting file ==============" + err.Error())
        return file, err
    }
    collection := client.Database("dfs").Collection("files")
    err = collection.FindOne(context.Background(), bson.M{"fileIdentifier": id}).Decode(&file)
    if err != nil {
        log.Println("======== Error getting file ==============" + err.Error())
        return file, err
    }
    log.Println("======== Got file ==============")
    log.Println("=================" + file.FileName + "=================")
    return file, nil
}

func (file *FileMetadata) UpdateFile() error {
    log.Println("======== Getting mongo client to update ==============")
    client, err := config.GetMongoClient()
    log.Println("======== Got mongo client to update ==============")
    if err != nil {
        return err
    }

    collection := client.Database("dfs").Collection("files")
    g, errors := collection.ReplaceOne(context.Background(), bson.M{"fileIdentifier": file.FileIdentifier}, file)
    log.Println("======== Updated file ==============", g)
    if errors != nil {
        return errors
    }
    return nil
}

func (file *FileMetadata) Delete() error {
    client, err := config.GetMongoClient()
    if err != nil {
        return err
    }
    collection := client.Database("dfs").Collection("files")
    _, err = collection.DeleteOne(context.Background(), bson.M{"fileIdentifier": file.FileIdentifier})
    if err != nil {
        return err
    }
    return nil
}