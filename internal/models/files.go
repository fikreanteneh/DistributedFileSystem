package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FileMetadata struct {
	Id             primitive.ObjectID   `bson:"_id" json:"_id,omitempty"`
	FileName       string               `bson:"fileName" json:"fileName"`
	FileSize       uint64               `bson:"fileSize" json:"fileSize"`
	NumberOfChunks uint64               `bson:"numberOfChunks" json:"numberOfChunks"`
	Replicas       [][]string           `bson:"replicas" json:"replicas"`
	ClientId       string               `bson:"clientId" json:"clientId"`
	SharedUser     []primitive.ObjectID `bson:"sharedUser" json:"sharedUser"`
	FileIdentifier string               `bson:"fileIdentifier" json:"fileIdentifier"`
}
