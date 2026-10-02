package main

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// 
type Server struct {
	db *mongo.Database
}

type InferenceSession struct {
    ID          primitive.ObjectID `bson:"_id,omitempty"`
    Facts       InferenceFacts     `bson:"facts"`
    Active      bool               `bson:"active"`
}

type InferenceFacts map[string]interface{}