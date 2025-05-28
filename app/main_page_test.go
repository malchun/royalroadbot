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

	// Verify HTMX is included
	assert.Contains(t, html, "https://unpkg.com/htmx.org")
	
	// Verify HTMX attributes are present for search functionality
	assert.Contains(t, html, "hx-post=\"/search\"")
	assert.Contains(t, html, "hx-trigger=\"input changed delay:500ms, search\"")
	assert.Contains(t, html, "hx-target=\"#book-results\"")
	
	// Verify refresh button with HTMX attributes
	assert.Contains(t, html, "hx-get=\"/refresh\"")
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

func TestRenderBookList(t *testing.T) {
	// Test books data
	books := []Book{
		{Title: "Test Book 1", Link: "https://example.com/book1"},
		{Title: "Test Book 2", Link: "https://example.com/book2"},
	}

	// Get the template
	tmpl, err := renderBookList(books)

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
	
	// Verify it only contains the book list section, not the entire page
	assert.Contains(t, html, "<ul class=\"book-list\">")
	assert.NotContains(t, html, "<!DOCTYPE html>")
	assert.NotContains(t, html, "<title>")
}

func TestEmptyBookList(t *testing.T) {
	// Empty books slice
	var books []Book
	
	// Get the template
	tmpl, err := renderBookList(books)
	
	// Verify no error occurred
	assert.NoError(t, err)
	
	// Create a buffer to hold the template output
	var buffer strings.Builder
	
	// Execute the template with empty data
	err = tmpl.Execute(&buffer, books)
	
	// Verify no error occurred during execution
	assert.NoError(t, err)
	
	// Get the rendered HTML
	html := buffer.String()
	
	// Verify the HTML contains the "no results" message
	assert.Contains(t, html, "No books found matching your search.")
	assert.Contains(t, html, "class=\"no-results\"")
}
