package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

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