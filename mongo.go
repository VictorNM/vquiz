package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func connect(ctx context.Context, uri string) (*mongo.Client, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	return client, nil
}

func insertQuestion(ctx context.Context, db *mongo.Database, q Question) error {
	q.ID = primitive.NewObjectID()

	_, err := db.Collection("questions").InsertOne(ctx, q)
	if err != nil {
		return fmt.Errorf("insert question error: %v", err)
	}

	return nil
}

func getQuestions(ctx context.Context, db *mongo.Database) ([]Question, error) {
	cur, err := db.Collection("questions").Find(ctx, bson.D{})
	if err != nil {
		return nil, fmt.Errorf("find questions failed: %v", err)
	}
	defer cur.Close(ctx)

	var res []Question

	for cur.Next(ctx) {
		var q Question
		err := cur.Decode(&q)

		if err != nil {
			return nil, fmt.Errorf("decode questions failed: %v", err)
		}
		res = append(res, q)
	}

	return res, nil
}
