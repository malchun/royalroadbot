package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRenderPage(t *testing.T) {
	// Test books data
	books := []Book{
		{Title: "Test Book 1", Link: "https://example.com/book1"},
		{Title: "Test Book 2", Link: "https://example.com/book2"},
	}

	// Get the template
	tmpl, err := renderPage(books)

	// Verify no error occurred
	assert.NoError(t, err)

	// Create a buffer to hold the template output
	var buffer strings.Builder

	// Execute the template with our test data
	err = tmpl.Execute(&buffer, books)

	// Verify no error occurred during execution
	assert.NoError(t, err)

	// Get the rendered HTML
	html := buffer.String()

	// Verify the HTML contains our test book titles and links
	assert.Contains(t, html, "Test Book 1")
	assert.Contains(t, html, "https://example.com/book1")
	assert.Contains(t, html, "Test Book 2")
	assert.Contains(t, html, "https://example.com/book2")

	// Verify essential HTML structure
	assert.Contains(t, html, "<!DOCTYPE html>")
	assert.Contains(t, html, "<title>Royal Road - Popular Books</title>")
	assert.Contains(t, html, "<h1>Top 10 Popular Books on Royal Road</h1>")

	// Verify search functionality is included
	assert.Contains(t, html, "searchBooks()")
	assert.Contains(t, html, "<input type=\"text\" id=\"searchInput\"")
}

func TestBookStructure(t *testing.T) {
	// Create a book
	book := Book{
		Title: "Test Book",
		Link:  "https://example.com/book",
	}

	// Verify book fields
	assert.Equal(t, "Test Book", book.Title)
	assert.Equal(t, "https://example.com/book", book.Link)
}
