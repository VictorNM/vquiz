package main

import "go.mongodb.org/mongo-driver/bson/primitive"

type Question struct {
	ID      primitive.ObjectID `bson:"_id"`
	Content string             `bson:"content"`
	Answer  string             `bson:"answer"`
}
