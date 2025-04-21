package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Book struct is defined in main_page.go
type testBook struct {
	Title string
	Link  string
}

// Store mock books for testing
var mockBooks []testBook

// Declarations for test mocking
var originalConnectFunc func()
var originalSaveFunc func([]Book) error

// Mock implementations for testing
func init() {
	// Save original functions
	originalConnectFunc = ConnectDBFunc
	originalSaveFunc = saveBooksWithMetadataFunc

	// Set up mock implementations
	ConnectDBFunc = func() {
		// Mock implementation - do nothing
		fmt.Println("Mock: Connected to MongoDB")
	}

	saveBooksWithMetadataFunc = func(books []Book) error {
		// Convert Book to testBook and save to our in-memory storage
		mockBooks = make([]testBook, len(books))
		for i, b := range books {
			mockBooks[i] = testBook{
				Title: b.Title,
				Link:  b.Link,
			}
		}
		fmt.Println("Mock: Saved", len(books), "books to database")
		return nil
	}
}

// Mock function to get books from our in-memory storage
func getTestBooks() ([]testBook, error) {
	return mockBooks, nil
}

func TestSaveAndGetBooks(t *testing.T) {
	// Create test data
	books := []Book{
		{Title: "Test Book 1", Link: "https://example.com/book1"},
		{Title: "Test Book 2", Link: "https://example.com/book2"},
	}

	// Reset mock books
	mockBooks = nil

	// Save books using the saveBooksWithMetadata function
	// This will use our mock implementation
	err := saveBooksWithMetadata(books)
	assert.NoError(t, err)

	// Retrieve books using our test function
	retrievedBooks, err := getTestBooks()
	assert.NoError(t, err)

	// Verify the books were saved to our mock storage
	assert.Equal(t, len(books), len(retrievedBooks))
	for i, book := range books {
		assert.Equal(t, book.Title, retrievedBooks[i].Title)
		assert.Equal(t, book.Link, retrievedBooks[i].Link)
	}
}
