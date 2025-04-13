package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	dbName         = "royalRoadBooks"
	collectionName = "books"
)

var client *mongo.Client

func ConnectDB() {
	mongoURI := os.Getenv("MONGODB_URI")

	var err error
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")
}

// saveBooksWithMetadata saves the list of books with their metadata to MongoDB database
func saveBooksWithMetadata(books []Book) error {
	ConnectDB()
	defer client.Disconnect(context.TODO())

	collection := client.Database(dbName).Collection(collectionName)
	for _, book := range books {
		_, err := collection.InsertOne(context.TODO(), bson.M{
			"title": book.Title,
			"link":  book.Link,
		})
		if err != nil {
			return fmt.Errorf("failed to insert book: %v", err)
		}
	}
	return nil
}

// getBooksWithMetadata retrieves the list of books with their metadata from MongoDB database
func getBooksWithMetadata() ([]Book, error) {
	ConnectDB()
	defer client.Disconnect(context.TODO())

	var books []Book
	collection := client.Database(dbName).Collection(collectionName)
	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to find books: %v", err)
	}
	if err = cursor.All(context.TODO(), &books); err != nil {
		return nil, fmt.Errorf("error decoding into struct: %v", err)
	}
	return books, nil
}
