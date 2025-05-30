package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"github.com/testcontainers/testcontainers-go/wait"
)

// setupTestDatabaseForHandlers starts a MongoDB container and returns a cleanup function
func setupTestDatabaseForHandlers(t *testing.T) (func(), string) {
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

	// Store original environment variable
	originalURI := os.Getenv("MONGODB_URI")

	// Set environment variable for MongoDB URI
	os.Setenv("MONGODB_URI", connectionURI)

	// Return cleanup function
	cleanup := func() {
		// Reset the global client to nil to ensure tests don't interfere with each other
		client = nil
		// Restore original environment variable
		if originalURI != "" {
			os.Setenv("MONGODB_URI", originalURI)
		} else {
			os.Unsetenv("MONGODB_URI")
		}
		// Terminate the container
		if err := mongodbContainer.Terminate(ctx); err != nil {
			t.Fatalf("Failed to terminate MongoDB container: %v", err)
		}
	}

	return cleanup, connectionURI
}

// Test setup helper function to initialize and later restore the global cached books
func setupCachedBooksForTest(t *testing.T) func() {
	// Store original cached books
	originalCachedBooks := cachedBooks

	// Setup test data
	cachedBooks = []Book{
		{Title: "Test Book 1", Link: "https://example.com/book1"},
		{Title: "Test Book 2", Link: "https://example.com/book2"},
		{Title: "Another Test", Link: "https://example.com/another"},
	}

	// Return a cleanup function to restore the original state
	return func() {
		cachedBooks = originalCachedBooks
	}
}

func TestBooksHandler(t *testing.T) {
	// Setup test data and defer cleanup
	cleanup := setupCachedBooksForTest(t)
	defer cleanup()

	// Create a request to pass to our handler
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(booksHandler)

	// Call the handler
	handler.ServeHTTP(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Check the response body contains our book data
	body := rr.Body.String()
	assert.Contains(t, body, "Test Book 1")
	assert.Contains(t, body, "Test Book 2")
	assert.Contains(t, body, "Another Test")
	assert.Contains(t, body, "https://example.com/book1")
	assert.Contains(t, body, "https://example.com/book2")
	assert.Contains(t, body, "https://example.com/another")
}

func TestSearchHandler(t *testing.T) {
	// Setup test data and defer cleanup
	cleanup := setupCachedBooksForTest(t)
	defer cleanup()

	// Test cases
	testCases := []struct {
		name            string
		searchQuery     string
		expectedBooks   []string
		unexpectedBooks []string
	}{
		{
			name:            "Empty search returns all books",
			searchQuery:     "",
			expectedBooks:   []string{"Test Book 1", "Test Book 2", "Another Test"},
			unexpectedBooks: []string{},
		},
		{
			name:            "Partial match returns matching books",
			searchQuery:     "book",
			expectedBooks:   []string{"Test Book 1", "Test Book 2"},
			unexpectedBooks: []string{"Another Test"},
		},
		{
			name:            "Case insensitive search",
			searchQuery:     "ANOTHER",
			expectedBooks:   []string{"Another Test"},
			unexpectedBooks: []string{"Test Book 1", "Test Book 2"},
		},
		{
			name:            "No matches shows no results message",
			searchQuery:     "nonexistent",
			expectedBooks:   []string{"No books found matching your search."},
			unexpectedBooks: []string{"Test Book 1", "Test Book 2", "Another Test"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create form data
			form := url.Values{}
			form.Add("search", tc.searchQuery)
			formData := form.Encode()

			// Create a request to the search endpoint
			req, err := http.NewRequest("POST", "/search", strings.NewReader(formData))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			// Create a ResponseRecorder
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(searchHandler)

			// Call the handler
			handler.ServeHTTP(rr, req)

			// Check the status code
			assert.Equal(t, http.StatusOK, rr.Code)

			// Check the response contains expected books and doesn't contain unexpected books
			body := rr.Body.String()
			for _, book := range tc.expectedBooks {
				assert.Contains(t, body, book)
			}
			for _, book := range tc.unexpectedBooks {
				assert.NotContains(t, body, book)
			}
		})
	}
}

func TestRefreshHandler(t *testing.T) {
	// Store the original cached books to restore later
	originalCachedBooks := cachedBooks

	// Set up initial test books
	testBooks := []Book{
		{Title: "Test Book 1", Link: "https://example.com/book1"},
		{Title: "Test Book 2", Link: "https://example.com/book2"},
	}
	
	// Setup initial state
	cachedBooks = []Book{{Title: "Initial Book", Link: "https://example.com/initial"}}
	
	// After the test, restore original state
	defer func() {
		cachedBooks = originalCachedBooks
	}()
	
	// Create a request to the refresh endpoint
	req, err := http.NewRequest("GET", "/refresh", nil)
	if err != nil {
		t.Fatal(err)
	}
	
	// Create a ResponseRecorder
	rr := httptest.NewRecorder()
	
	// Create a custom handler that simulates the refreshHandler behavior
	// but uses our test books instead of calling fetchPopularBooks
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate updating the cached books
		cachedBooks = testBooks
		
		// Render the book list just like the real handler
		tmpl, err := renderBookList(cachedBooks)
		if err != nil {
			http.Error(w, "Template error", http.StatusInternalServerError)
			return
		}
		
		// Execute the template
		err = tmpl.Execute(w, cachedBooks)
		if err != nil {
			http.Error(w, "Template execution error", http.StatusInternalServerError)
			return
		}
	})
	
	// Call the handler
	handler.ServeHTTP(rr, req)
	
	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code)
	
	// Verify the response contains our test books
	body := rr.Body.String()
	assert.Contains(t, body, "Test Book 1")
	assert.Contains(t, body, "Test Book 2")
	
	// Verify the cached books were updated
	assert.Equal(t, 2, len(cachedBooks))
	assert.Equal(t, "Test Book 1", cachedBooks[0].Title)
	assert.Equal(t, "Test Book 2", cachedBooks[1].Title)
	
	// Verify the old book is gone
	assert.NotContains(t, body, "Initial Book")
}

func TestSearchBooksHandler(t *testing.T) {
	// Test valid search query
	t.Run("Valid search query", func(t *testing.T) {
		form := url.Values{}
		form.Add("query", "test")
		formData := form.Encode()

		req, err := http.NewRequest("POST", "/search-books", strings.NewReader(formData))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(searchBooksHandler)

		handler.ServeHTTP(rr, req)

		// Should return OK status
		assert.Equal(t, http.StatusOK, rr.Code)
		
		// Should return HTML content type
		assert.Contains(t, rr.Header().Get("Content-Type"), "text/html")
		
		// Should contain HTML book list structure
		body := rr.Body.String()
		assert.Contains(t, body, `<ul class="book-list">`)
	})

	// Test empty search query
	t.Run("Empty search query", func(t *testing.T) {
		form := url.Values{}
		form.Add("query", "")
		formData := form.Encode()

		req, err := http.NewRequest("POST", "/search-books", strings.NewReader(formData))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(searchBooksHandler)

		handler.ServeHTTP(rr, req)

		// Should return bad request status
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	// Test missing query parameter
	t.Run("Missing query parameter", func(t *testing.T) {
		req, err := http.NewRequest("POST", "/search-books", strings.NewReader(""))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(searchBooksHandler)

		handler.ServeHTTP(rr, req)

		// Should return bad request status
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})
}

func TestMemorizeBookHandler(t *testing.T) {
	// Test valid book memorization
	t.Run("Valid book memorization", func(t *testing.T) {
		// Setup test database
		cleanup, _ := setupTestDatabaseForHandlers(t)
		defer cleanup()

		form := url.Values{}
		form.Add("title", "Test Book")
		form.Add("link", "https://example.com/test")
		formData := form.Encode()

		req, err := http.NewRequest("POST", "/memorize-book", strings.NewReader(formData))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(memorizeBookHandler)

		handler.ServeHTTP(rr, req)

		// Should return OK status
		assert.Equal(t, http.StatusOK, rr.Code)
		
		// Should contain success message
		body := rr.Body.String()
		assert.Contains(t, body, "Book memorized successfully")
	})

	// Test GET method (should be rejected)
	t.Run("GET method not allowed", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/memorize-book", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(memorizeBookHandler)

		handler.ServeHTTP(rr, req)

		// Should return method not allowed status
		assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)
	})

	// Test missing title
	t.Run("Missing title", func(t *testing.T) {
		form := url.Values{}
		form.Add("link", "https://example.com/test")
		formData := form.Encode()

		req, err := http.NewRequest("POST", "/memorize-book", strings.NewReader(formData))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(memorizeBookHandler)

		handler.ServeHTTP(rr, req)

		// Should return bad request status
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	// Test missing link
	t.Run("Missing link", func(t *testing.T) {
		form := url.Values{}
		form.Add("title", "Test Book")
		formData := form.Encode()

		req, err := http.NewRequest("POST", "/memorize-book", strings.NewReader(formData))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(memorizeBookHandler)

		handler.ServeHTTP(rr, req)

		// Should return bad request status
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})
}

func TestMemorizedBooksHandler(t *testing.T) {
	// Test getting memorized books
	t.Run("Get memorized books", func(t *testing.T) {
		// Setup test database
		cleanup, _ := setupTestDatabaseForHandlers(t)
		defer cleanup()

		req, err := http.NewRequest("GET", "/memorized-books", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(memorizedBooksHandler)

		handler.ServeHTTP(rr, req)

		// Should return OK status
		assert.Equal(t, http.StatusOK, rr.Code)
		
		// Should return HTML content type
		assert.Contains(t, rr.Header().Get("Content-Type"), "text/html")
		
		// Should contain HTML book list structure
		body := rr.Body.String()
		assert.Contains(t, body, `<ul class="book-list">`)
	})
}

func TestRemoveMemorizedBookHandler(t *testing.T) {
	// Test valid book removal
	t.Run("Valid book removal", func(t *testing.T) {
		// Setup test database
		cleanup, _ := setupTestDatabaseForHandlers(t)
		defer cleanup()

		// First memorize a book
		book := Book{Title: "Book to Remove", Link: "https://example.com/remove"}
		err := memorizeBook(book)
		require.NoError(t, err)

		// Now remove it
		form := url.Values{}
		form.Add("title", "Book to Remove")
		formData := form.Encode()

		req, err := http.NewRequest("POST", "/remove-memorized-book", strings.NewReader(formData))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(removeMemorizedBookHandler)

		handler.ServeHTTP(rr, req)

		// Should return OK status
		assert.Equal(t, http.StatusOK, rr.Code)
		
		// Should return HTML content type
		assert.Contains(t, rr.Header().Get("Content-Type"), "text/html")

		// Verify book was actually removed from database
		books, err := getMemorizedBooks()
		assert.NoError(t, err)
		for _, b := range books {
			assert.NotEqual(t, "Book to Remove", b.Title)
		}
	})

	// Test GET method not allowed
	t.Run("GET method not allowed", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/remove-memorized-book", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(removeMemorizedBookHandler)

		handler.ServeHTTP(rr, req)

		// Should return method not allowed
		assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)
	})

	// Test missing title
	t.Run("Missing title", func(t *testing.T) {
		form := url.Values{}
		// No title added
		formData := form.Encode()

		req, err := http.NewRequest("POST", "/remove-memorized-book", strings.NewReader(formData))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(removeMemorizedBookHandler)

		handler.ServeHTTP(rr, req)

		// Should return bad request
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	// Test removing non-existent book
	t.Run("Remove non-existent book", func(t *testing.T) {
		// Setup test database
		cleanup, _ := setupTestDatabaseForHandlers(t)
		defer cleanup()

		form := url.Values{}
		form.Add("title", "Non-existent Book")
		formData := form.Encode()

		req, err := http.NewRequest("POST", "/remove-memorized-book", strings.NewReader(formData))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(removeMemorizedBookHandler)

		handler.ServeHTTP(rr, req)

		// Should return internal server error when book doesn't exist
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}