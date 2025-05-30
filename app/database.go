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
	dbName                    = "royalRoadBooks"
	collectionName            = "books"
	memorizedCollectionName   = "memorizedBooks"
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

// isInTestEnvironment checks if we're running in a test environment
func isInTestEnvironment() bool {
	for _, arg := range os.Args {
		if arg == "-test.v" || arg == "-test.run" {
			return true
		}
	}
	return false
}

func saveBooksWithMetadata(books []Book) error {
	ConnectDB()
	// Only disconnect if not in a test environment
	if !isInTestEnvironment() {
		defer client.Disconnect(context.TODO())
	}

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

// saveBookToMemory saves a single book to the memorized books collection
func saveBookToMemory(book Book) error {
	ConnectDB()
	// Only disconnect if not in a test environment
	if !isInTestEnvironment() {
		defer client.Disconnect(context.TODO())
	}

	collection := client.Database(dbName).Collection(memorizedCollectionName)
	
	// Check if book already exists to avoid duplicates
	existing := collection.FindOne(context.TODO(), bson.M{"title": book.Title})
	if existing.Err() == nil {
		return fmt.Errorf("book '%s' is already memorized", book.Title)
	}
	
	_, err := collection.InsertOne(context.TODO(), bson.M{
		"title":     book.Title,
		"link":      book.Link,
		"timestamp": time.Now(),
	})
	if err != nil {
		return fmt.Errorf("failed to insert memorized book: %v", err)
	}
	return nil
}

// loadMemorizedBooks retrieves all memorized books from the database
func loadMemorizedBooks() ([]Book, error) {
	ConnectDB()
	// Only disconnect if not in a test environment
	if !isInTestEnvironment() {
		defer client.Disconnect(context.TODO())
	}

	var books []Book
	collection := client.Database(dbName).Collection(memorizedCollectionName)
	
	// Sort by timestamp descending (newest first)
	opts := options.Find().SetSort(bson.D{{"timestamp", -1}})
	cursor, err := collection.Find(context.TODO(), bson.M{}, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to find memorized books: %v", err)
	}
	
	// Create a slice to hold the raw documents
	var rawBooks []bson.M
	if err = cursor.All(context.TODO(), &rawBooks); err != nil {
		return nil, fmt.Errorf("error decoding memorized books: %v", err)
	}
	
	// Convert to Book structs
	for _, rawBook := range rawBooks {
		if title, ok := rawBook["title"].(string); ok {
			if link, ok := rawBook["link"].(string); ok {
				books = append(books, Book{
					Title: title,
					Link:  link,
				})
			}
		}
	}
	
	return books, nil
}

// removeBookFromMemory removes a book from the memorized collection by title
func removeBookFromMemory(title string) error {
	ConnectDB()
	// Only disconnect if not in a test environment
	if !isInTestEnvironment() {
		defer client.Disconnect(context.TODO())
	}

	collection := client.Database(dbName).Collection(memorizedCollectionName)
	result, err := collection.DeleteOne(context.TODO(), bson.M{"title": title})
	if err != nil {
		return fmt.Errorf("failed to delete memorized book: %v", err)
	}
	
	if result.DeletedCount == 0 {
		return fmt.Errorf("book '%s' not found in memorized collection", title)
	}
	
	return nil
}
