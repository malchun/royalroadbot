package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Setup function for tests with MongoDB testcontainer
func setupTestWithMongoDB(t *testing.T) (*mongodb.MongoDBContainer, *mongo.Client, func()) {

	// Create MongoDB container
	ctx := context.Background()
	mongodbContainer, err := mongodb.RunContainer(ctx,
		testcontainers.WithImage("mongo:latest"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("Waiting for connections").
				WithOccurrence(1).
				WithStartupTimeout(time.Second*30),
		),
	)
	require.NoError(t, err)

	// Get connection URI
	mongodbURI, err := mongodbContainer.ConnectionString(ctx)
	require.NoError(t, err)

	// Override environment variable
	os.Setenv("MONGODB_URI", mongodbURI)

	// Create a new client directly for tests
	clientOptions := options.Client().ApplyURI(mongodbURI)
	testClient, err := mongo.Connect(ctx, clientOptions)
	require.NoError(t, err)

	// Clear the test database
	err = testClient.Database(dbName).Drop(ctx)
	require.NoError(t, err)

	// Return a cleanup function
	cleanup := func() {
		testClient.Disconnect(ctx)
		mongodbContainer.Terminate(ctx)

		// Clear environment variable
		os.Unsetenv("MONGODB_URI")
	}

	return mongodbContainer, testClient, cleanup
}

// Helper function to verify that NO books are saved in MongoDB for popular books
func verifyNoBooksInMongoDB(t *testing.T, testClient *mongo.Client) {
	ctx := context.Background()
	collection := testClient.Database(dbName).Collection(collectionName)

	count, err := collection.CountDocuments(ctx, bson.M{})
	require.NoError(t, err)
	assert.Equal(t, int64(0), count, "Popular books should not be saved to database")
}

// Positive Tests

// Test successfully fetching books (happy path)
func TestFetchBooks_Success(t *testing.T) {
	_, testClient, cleanup := setupTestWithMongoDB(t)
	defer cleanup()

	// Create a test server with mock HTML content
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`
			<!DOCTYPE html>
			<html>
			<body>
				<div class="fiction-list-item">
					<h2 class="fiction-title"><a href="/fiction/1234">Test Book 1</a></h2>
				</div>
				<div class="fiction-list-item">
					<h2 class="fiction-title"><a href="/fiction/5678">Test Book 2</a></h2>
				</div>
			</body>
			</html>
		`))
	}))
	defer server.Close()

	// Call the function we're testing
	books, err := fetchBooks(server.URL)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, 2, len(books))
	assert.Equal(t, "Test Book 1", books[0].Title)
	assert.Equal(t, "https://www.royalroad.com/fiction/1234", books[0].Link)
	assert.Equal(t, "Test Book 2", books[1].Title)
	assert.Equal(t, "https://www.royalroad.com/fiction/5678", books[1].Link)

	// Verify popular books are NOT saved to MongoDB (memory-only now)
	verifyNoBooksInMongoDB(t, testClient)
}

// Test that we only return a maximum of 10 books
func TestFetchBooks_MaximumTenBooks(t *testing.T) {
	_, testClient, cleanup := setupTestWithMongoDB(t)
	defer cleanup()

	// Create HTML with 15 books
	var htmlBuilder strings.Builder
	htmlBuilder.WriteString(`<!DOCTYPE html><html><body>`)

	for i := 1; i <= 15; i++ {
		htmlBuilder.WriteString(fmt.Sprintf(`
			<div class="fiction-list-item">
				<h2 class="fiction-title"><a href="/fiction/%d">Test Book %d</a></h2>
			</div>
		`, i, i))
	}

	htmlBuilder.WriteString(`</body></html>`)

	// Create a test server with our HTML
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(htmlBuilder.String()))
	}))
	defer server.Close()

	// Call the function we're testing
	books, err := fetchBooks(server.URL)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, 10, len(books)) // Should limit to 10 books

	// Check the first and last books
	assert.Equal(t, "Test Book 1", books[0].Title)
	assert.Equal(t, "Test Book 10", books[9].Title)

	// Verify popular books are NOT saved to MongoDB (memory-only now)
	verifyNoBooksInMongoDB(t, testClient)
}

// Test handling empty response still works
func TestFetchBooks_EmptyResponse(t *testing.T) {
	_, testClient, cleanup := setupTestWithMongoDB(t)
	defer cleanup()

	// Create a test server with empty HTML
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<!DOCTYPE html><html><body></body></html>`))
	}))
	defer server.Close()

	// Call the function we're testing
	books, err := fetchBooks(server.URL)

	// Assertions
	assert.NoError(t, err) // No error should be returned
	assert.Empty(t, books) // Should return empty slice

	// Verify no books were saved in MongoDB
	verifyNoBooksInMongoDB(t, testClient)
}

// Test that we handle and filter items with missing data correctly
func TestFetchBooks_MissingData(t *testing.T) {
	_, testClient, cleanup := setupTestWithMongoDB(t)
	defer cleanup()

	// Create a test server with some items missing titles or links
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`
			<!DOCTYPE html>
			<html>
			<body>
				<div class="fiction-list-item">
					<h2 class="fiction-title"><a href="/fiction/1234">Test Book 1</a></h2>
				</div>
				<div class="fiction-list-item">
					<h2 class="fiction-title"><a href="">Missing Link</a></h2>
				</div>
				<div class="fiction-list-item">
					<h2 class="fiction-title"><a href="/fiction/5678"></a></h2>
				</div>
				<div class="fiction-list-item">
					<h2 class="fiction-title"><a href="/fiction/9012">Test Book 4</a></h2>
				</div>
			</body>
			</html>
		`))
	}))
	defer server.Close()

	// Call the function we're testing
	books, err := fetchBooks(server.URL)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, 2, len(books)) // Should only have 2 valid books
	assert.Equal(t, "Test Book 1", books[0].Title)
	assert.Equal(t, "Test Book 4", books[1].Title)

	// Verify popular books are NOT saved to MongoDB (memory-only now)
	verifyNoBooksInMongoDB(t, testClient)
}

// Negative Tests

// Test handling invalid URLs
func TestFetchBooks_InvalidURL(t *testing.T) {
	_, testClient, cleanup := setupTestWithMongoDB(t)
	defer cleanup()

	// Call with an invalid URL
	books, err := fetchBooks("not-a-valid-url")

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, books)

	// Verify no books were saved in MongoDB
	verifyNoBooksInMongoDB(t, testClient)
}

// Test handling server errors
func TestFetchBooks_ServerError(t *testing.T) {
	_, testClient, cleanup := setupTestWithMongoDB(t)
	defer cleanup()

	// Create a test server that returns an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	// Call the function we're testing
	books, err := fetchBooks(server.URL)

	// The colly library doesn't always return an error for HTTP status codes,
	// so we may just get an empty slice
	if err != nil {
		assert.Error(t, err)
		assert.Nil(t, books)
	} else {
		assert.Empty(t, books)
	}

	// Verify no books were saved in MongoDB
	verifyNoBooksInMongoDB(t, testClient)
}

// Test that the function doesn't crash with malformed HTML
func TestFetchBooks_MalformedHTML(t *testing.T) {
	_, testClient, cleanup := setupTestWithMongoDB(t)
	defer cleanup()

	// Create a test server with malformed HTML
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<!DOCTYPE html><html><body><div class="fiction-list-item"><h2 class="fiction-title">Broken HTML`))
	}))
	defer server.Close()

	// Call the function we're testing
	books, err := fetchBooks(server.URL)

	// Assertions - should handle malformed HTML gracefully
	assert.NoError(t, err)
	assert.Empty(t, books) // Either empty or nil books

	// Verify no books were saved in MongoDB
	verifyNoBooksInMongoDB(t, testClient)
}
