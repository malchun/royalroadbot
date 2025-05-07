package main

import (
	"context"
	"testing"
	"time"

	"os"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// setupTestDatabase starts a MongoDB container and returns a cleanup function
func setupTestDatabase(t *testing.T) (func(), string) {
	// Create MongoDB container request
	ctx := context.Background()
	mongodbContainer, err := mongodb.RunContainer(ctx,
		testcontainers.WithImage("mongo:latest"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("Waiting for connections").WithStartupTimeout(time.Second*30),
		),
	)
	require.NoError(t, err, "Failed to start MongoDB container")

	// Get connection URI
	connectionURI, err := mongodbContainer.ConnectionString(ctx)
	require.NoError(t, err, "Failed to get MongoDB connection string")

	// Set environment variable for MongoDB URI
	os.Setenv("MONGODB_URI", connectionURI)

	// Return cleanup function
	cleanup := func() {
		// Reset the global client to nil to ensure tests don't interfere with each other
		client = nil
		// Terminate the container
		if err := mongodbContainer.Terminate(ctx); err != nil {
			t.Fatalf("Failed to terminate MongoDB container: %v", err)
		}
	}

	return cleanup, connectionURI
}

// TestSaveAndGetBooks tests the save and retrieve functionality with a real MongoDB instance
func TestSaveAndGetBooks(t *testing.T) {
	// Set up test database
	cleanup, _ := setupTestDatabase(t)
	defer cleanup()

	// Create test data
	books := []Book{
		{Title: "Test Book 1", Link: "https://example.com/book1"},
		{Title: "Test Book 2", Link: "https://example.com/book2"},
	}

	// Save books using the real implementation
	// This should use the actual database connection
	err := saveBooksWithMetadata(books)
	assert.NoError(t, err)

	// Retrieve books using the real implementation
	retrievedBooks, err := getBooksWithMetadata()
	assert.NoError(t, err)

	// Verify the books were saved correctly
	assert.Equal(t, len(books), len(retrievedBooks))
	for i, book := range books {
		assert.Equal(t, book.Title, retrievedBooks[i].Title)
		assert.Equal(t, book.Link, retrievedBooks[i].Link)
	}
}

// TestDatabaseConnection tests the connection to the MongoDB database
func TestDatabaseConnection(t *testing.T) {
	// Set up test database
	cleanup, connectionURI := setupTestDatabase(t)
	defer cleanup()

	// Test direct connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	clientOptions := options.Client().ApplyURI(connectionURI)
	testClient, err := mongo.Connect(ctx, clientOptions)
	require.NoError(t, err, "Failed to connect to MongoDB")
	defer testClient.Disconnect(ctx)

	// Verify connection with ping
	err = testClient.Ping(ctx, nil)
	require.NoError(t, err, "Failed to ping MongoDB")
}
