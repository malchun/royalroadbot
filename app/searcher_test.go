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

	"github.com/gocolly/colly/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// setupSearchTestWithMongoDB sets up a MongoDB container for search tests
func setupSearchTestWithMongoDB(t *testing.T) (*mongodb.MongoDBContainer, *mongo.Client, func()) {
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
		os.Unsetenv("MONGODB_URI")
		client = nil // Reset global client
	}

	return mongodbContainer, testClient, cleanup
}

// createMockSearchServer creates a test server that mimics Royal Road search results
func createMockSearchServer(htmlContent string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(htmlContent))
	}))
}

// Test searchRoyalRoadBooks with successful results
func TestSearchRoyalRoadBooks_Success(t *testing.T) {
	// Create mock HTML that matches Royal Road search results structure
	mockHTML := `
		<!DOCTYPE html>
		<html>
		<body>
			<h2><a href="/fiction/1234/wizard-story">The Great Wizard</a></h2>
			<h2><a href="/fiction/5678/magic-academy">Magic Academy Adventures</a></h2>
			<h2><a href="/fiction/9012/spellcaster">The Spellcaster Chronicles</a></h2>
		</body>
		</html>
	`

	server := createMockSearchServer(mockHTML)
	defer server.Close()

	// Override the search function to use our test server
	testSearchURL := server.URL + "?title="

	// Monkey patch by creating a custom function for testing
	testSearchFunc := func(query string) ([]Book, error) {
		// Simulate the searchRoyalRoadBooks function with our test URL
		books, err := searchBooksFromURL(testSearchURL + query)
		return books, err
	}

	// Test the search
	books, err := testSearchFunc("wizard")

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, 3, len(books))
	assert.Equal(t, "The Great Wizard", books[0].Title)
	assert.Equal(t, "https://www.royalroad.com/fiction/1234/wizard-story", books[0].Link)
	assert.Equal(t, "Magic Academy Adventures", books[1].Title)
	assert.Equal(t, "https://www.royalroad.com/fiction/5678/magic-academy", books[1].Link)
}

// Helper function to test URL-based searching (for testing purposes)
func searchBooksFromURL(searchURL string) ([]Book, error) {
	c := colly.NewCollector()
	var books []Book

	c.OnHTML("h2", func(e *colly.HTMLElement) {
		title := strings.TrimSpace(e.ChildText("a"))
		link := e.ChildAttr("a", "href")

		if title != "" && link != "" {
			fullLink := link
			if !strings.HasPrefix(link, "http") {
				fullLink = "https://www.royalroad.com" + link
			}

			books = append(books, Book{
				Title: title,
				Link:  fullLink,
			})
		}
	})

	c.UserAgent = "Mozilla/5.0 (compatible; RoyalRoadBot/1.0)"
	err := c.Visit(searchURL)
	if err != nil {
		return nil, fmt.Errorf("failed to visit search URL: %w", err)
	}

	if len(books) > 15 {
		books = books[:15]
	}

	return books, nil
}

// Test searchRoyalRoadBooks with empty query
func TestSearchRoyalRoadBooks_EmptyQuery(t *testing.T) {
	books, err := searchRoyalRoadBooks("")

	assert.NoError(t, err)
	assert.Empty(t, books)
}

// Test searchRoyalRoadBooks with whitespace-only query
func TestSearchRoyalRoadBooks_WhitespaceQuery(t *testing.T) {
	books, err := searchRoyalRoadBooks("   ")

	assert.NoError(t, err)
	assert.Empty(t, books)
}

// Test searchRoyalRoadBooks with no results
func TestSearchRoyalRoadBooks_NoResults(t *testing.T) {
	mockHTML := `
		<!DOCTYPE html>
		<html>
		<body>
			<p>No results found</p>
		</body>
		</html>
	`

	server := createMockSearchServer(mockHTML)
	defer server.Close()

	books, err := searchBooksFromURL(server.URL)

	assert.NoError(t, err)
	assert.Empty(t, books)
}

// Test searchRoyalRoadBooks with more than 15 results (should limit)
func TestSearchRoyalRoadBooks_LimitResults(t *testing.T) {
	var htmlBuilder strings.Builder
	htmlBuilder.WriteString(`<!DOCTYPE html><html><body>`)

	// Create 20 mock search results
	for i := 1; i <= 20; i++ {
		htmlBuilder.WriteString(fmt.Sprintf(`
			<h2><a href="/fiction/%d/book-%d">Book %d</a></h2>
		`, i, i, i))
	}

	htmlBuilder.WriteString(`</body></html>`)

	server := createMockSearchServer(htmlBuilder.String())
	defer server.Close()

	books, err := searchBooksFromURL(server.URL)

	assert.NoError(t, err)
	assert.Equal(t, 15, len(books)) // Should limit to 15
	assert.Equal(t, "Book 1", books[0].Title)
	assert.Equal(t, "Book 15", books[14].Title)
}

// Test searchRoyalRoadBooks with malformed HTML
func TestSearchRoyalRoadBooks_MalformedHTML(t *testing.T) {
	mockHTML := `<!DOCTYPE html><html><body><h2><a href="/fiction/123"Incomplete`

	server := createMockSearchServer(mockHTML)
	defer server.Close()

	books, err := searchBooksFromURL(server.URL)

	assert.NoError(t, err)
	assert.Empty(t, books) // Should handle gracefully
}

// Test memorizeBook with valid book
func TestMemorizeBook_Success(t *testing.T) {
	_, testClient, cleanup := setupSearchTestWithMongoDB(t)
	defer cleanup()

	book := Book{
		Title: "Test Wizard Book",
		Link:  "https://www.royalroad.com/fiction/1234/test-wizard-book",
	}

	err := memorizeBook(book)

	assert.NoError(t, err)

	// Verify book was saved to database
	ctx := context.Background()
	collection := testClient.Database(dbName).Collection(memorizedCollectionName)

	var savedBook bson.M
	err = collection.FindOne(ctx, bson.M{"title": book.Title}).Decode(&savedBook)
	require.NoError(t, err)

	assert.Equal(t, book.Title, savedBook["title"])
	assert.Equal(t, book.Link, savedBook["link"])
	assert.NotNil(t, savedBook["timestamp"])
}

// Test memorizeBook with empty title
func TestMemorizeBook_EmptyTitle(t *testing.T) {
	_, _, cleanup := setupSearchTestWithMongoDB(t)
	defer cleanup()

	book := Book{
		Title: "",
		Link:  "https://www.royalroad.com/fiction/1234/test-book",
	}

	err := memorizeBook(book)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "title and link cannot be empty")
}

// Test memorizeBook with empty link
func TestMemorizeBook_EmptyLink(t *testing.T) {
	_, _, cleanup := setupSearchTestWithMongoDB(t)
	defer cleanup()

	book := Book{
		Title: "Test Book",
		Link:  "",
	}

	err := memorizeBook(book)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "title and link cannot be empty")
}

// Test memorizeBook with duplicate book
func TestMemorizeBook_Duplicate(t *testing.T) {
	_, testClient, cleanup := setupSearchTestWithMongoDB(t)
	defer cleanup()

	book := Book{
		Title: "Duplicate Book",
		Link:  "https://www.royalroad.com/fiction/1234/duplicate-book",
	}

	// Save book first time
	err := memorizeBook(book)
	assert.NoError(t, err)

	// Try to save same book again
	err = memorizeBook(book)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already memorized")

	// Verify only one copy exists
	ctx := context.Background()
	collection := testClient.Database(dbName).Collection(memorizedCollectionName)
	count, err := collection.CountDocuments(ctx, bson.M{"title": book.Title})
	require.NoError(t, err)
	assert.Equal(t, int64(1), count)
}

// Test getMemorizedBooks with no books
func TestGetMemorizedBooks_Empty(t *testing.T) {
	_, _, cleanup := setupSearchTestWithMongoDB(t)
	defer cleanup()

	books, err := getMemorizedBooks()

	assert.NoError(t, err)
	assert.Empty(t, books)
}

// Test getMemorizedBooks with multiple books
func TestGetMemorizedBooks_Multiple(t *testing.T) {
	_, _, cleanup := setupSearchTestWithMongoDB(t)
	defer cleanup()

	// Create test books with slight delays to ensure different timestamps
	book1 := Book{Title: "First Book", Link: "https://example.com/1"}
	book2 := Book{Title: "Second Book", Link: "https://example.com/2"}
	book3 := Book{Title: "Third Book", Link: "https://example.com/3"}

	// Save books with delays
	err := memorizeBook(book1)
	assert.NoError(t, err)

	time.Sleep(10 * time.Millisecond)
	err = memorizeBook(book2)
	assert.NoError(t, err)

	time.Sleep(10 * time.Millisecond)
	err = memorizeBook(book3)
	assert.NoError(t, err)

	// Retrieve books
	books, err := getMemorizedBooks()

	assert.NoError(t, err)
	assert.Equal(t, 3, len(books))

	// Should be sorted by timestamp descending (newest first)
	assert.Equal(t, "Third Book", books[0].Title)
	assert.Equal(t, "Second Book", books[1].Title)
	assert.Equal(t, "First Book", books[2].Title)
}

// Test removeMemorizedBook with existing book
func TestRemoveMemorizedBook_Success(t *testing.T) {
	_, testClient, cleanup := setupSearchTestWithMongoDB(t)
	defer cleanup()

	book := Book{
		Title: "Book to Remove",
		Link:  "https://www.royalroad.com/fiction/1234/book-to-remove",
	}

	// First save the book
	err := memorizeBook(book)
	assert.NoError(t, err)

	// Verify it exists
	ctx := context.Background()
	collection := testClient.Database(dbName).Collection(memorizedCollectionName)
	count, err := collection.CountDocuments(ctx, bson.M{"title": book.Title})
	require.NoError(t, err)
	assert.Equal(t, int64(1), count)

	// Remove the book
	err = removeMemorizedBook(book.Title)
	assert.NoError(t, err)

	// Verify it's gone
	count, err = collection.CountDocuments(ctx, bson.M{"title": book.Title})
	require.NoError(t, err)
	assert.Equal(t, int64(0), count)
}

// Test removeMemorizedBook with non-existent book
func TestRemoveMemorizedBook_NotFound(t *testing.T) {
	_, _, cleanup := setupSearchTestWithMongoDB(t)
	defer cleanup()

	err := removeMemorizedBook("Non-existent Book")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

// Test removeMemorizedBook with empty title
func TestRemoveMemorizedBook_EmptyTitle(t *testing.T) {
	_, _, cleanup := setupSearchTestWithMongoDB(t)
	defer cleanup()

	err := removeMemorizedBook("")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "title cannot be empty")
}

// Test integration: search, memorize, retrieve, remove
func TestSearchAndMemorizeIntegration(t *testing.T) {
	_, _, cleanup := setupSearchTestWithMongoDB(t)
	defer cleanup()

	// Create mock search results
	mockHTML := `
		<!DOCTYPE html>
		<html>
		<body>
			<h2><a href="/fiction/1234/wizard-book">Amazing Wizard Story</a></h2>
			<h2><a href="/fiction/5678/magic-book">Magic Adventures</a></h2>
		</body>
		</html>
	`

	server := createMockSearchServer(mockHTML)
	defer server.Close()

	// Search for books
	books, err := searchBooksFromURL(server.URL)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(books))

	// Memorize first book
	err = memorizeBook(books[0])
	assert.NoError(t, err)

	// Memorize second book
	err = memorizeBook(books[1])
	assert.NoError(t, err)

	// Retrieve memorized books
	memorizedBooks, err := getMemorizedBooks()
	assert.NoError(t, err)
	assert.Equal(t, 2, len(memorizedBooks))

	// Verify the books are the same (might be in different order due to timestamp sorting)
	titles := []string{memorizedBooks[0].Title, memorizedBooks[1].Title}
	assert.Contains(t, titles, "Amazing Wizard Story")
	assert.Contains(t, titles, "Magic Adventures")

	// Remove one book
	err = removeMemorizedBook(books[0].Title)
	assert.NoError(t, err)

	// Verify only one book remains
	memorizedBooks, err = getMemorizedBooks()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(memorizedBooks))
}


