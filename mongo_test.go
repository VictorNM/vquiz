package main

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"os"
	"testing"
)

var db *mongo.Database

func TestInsertQuestion(t *testing.T) {
	assert := assert.New(t)

	err := insertQuestion(context.Background(), db, Question{
		Content: "hello",
		Answer:  "world",
	})

	assert.NoError(err)

	questions, err := getQuestions(context.TODO(), db)

	assert.NoError(err)
	assert.Len(questions, 1)
	assert.Equal("hello", questions[0].Content)
}

func TestMain(m *testing.M) {
	client, err := connect(context.Background(), os.Getenv("MONGO_URL"))
	if err != nil {
		log.Fatalf("Connect mongo failed: %v", err)
	}

	db = client.Database("test")

	code := m.Run()

	if err := db.Drop(context.Background()); err != nil {
		log.Fatalf("Drop database failed: %v", err)
	}

	os.Exit(code)
}
