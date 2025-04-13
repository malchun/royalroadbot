package main

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	testMongoURI       = "mongodb://admin:password@mongo:27017" // Use a different URI for testing if needed
	testDBName         = "royalRoadBooksTest"
	testCollectionName = "booksTest"
)

var testClient *mongo.Client

func TestMain(m *testing.M) {
	// Setup MongoDB client for tests
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	clientOptions := options.Client().ApplyURI(testMongoURI)
	var err error
	testClient, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = testClient.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB for testing!")

	// Run tests
	m.Run()

	// Teardown: Disconnect the client
	testClient.Disconnect(context.TODO())
}

func TestSaveAndGetBooks(t *testing.T) {
	books := []Book{
		{Title: "Test Book 1", Link: "https://example.com/book1"},
		{Title: "Test Book 2", Link: "https://example.com/book2"},
	}

	// Save books to database
	err := saveBooksWithMetadata(books)
	assert.NoError(t, err)

	// Retrieve books from database
	retrievedBooks, err := getBooksWithMetadata()
	assert.NoError(t, err)
	assert.Equal(t, len(books), len(retrievedBooks))
	for i, book := range books {
		assert.Equal(t, book.Title, retrievedBooks[i].Title)
		assert.Equal(t, book.Link, retrievedBooks[i].Link)
	}
}
