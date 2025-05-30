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

func TestRenderTabbedMain(t *testing.T) {
	// Test books data
	books := []Book{
		{Title: "Test Book 1", Link: "https://example.com/book1"},
		{Title: "Test Book 2", Link: "https://example.com/book2"},
	}

	// Get the template
	tmpl, err := renderTabbedMain(books)

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

	// Verify essential HTML structure
	assert.Contains(t, html, "<!DOCTYPE html>")
	assert.Contains(t, html, "<title>Royal Road - Books Explorer</title>")
	assert.Contains(t, html, "Royal Road - Books Explorer")

	// Verify HTMX is included
	assert.Contains(t, html, "https://unpkg.com/htmx.org")

	// Verify tabbed interface structure
	assert.Contains(t, html, "class=\"tabs\"")
	assert.Contains(t, html, "class=\"tab-buttons\"")
	assert.Contains(t, html, "class=\"tab-content\"")
	assert.Contains(t, html, "data-tab=\"popular\"")
	assert.Contains(t, html, "data-tab=\"search\"")
	assert.Contains(t, html, "Popular Books")
	assert.Contains(t, html, "Search & Memorize")

	// Verify theme toggle functionality
	assert.Contains(t, html, "theme-toggle")
	assert.Contains(t, html, "toggleTheme()")

	// Verify search functionality in search tab
	assert.Contains(t, html, "id=\"searchQuery\"")
	assert.Contains(t, html, "hx-post=\"/search-books\"")
	assert.Contains(t, html, "id=\"search-results\"")
	assert.Contains(t, html, "id=\"memorized-results\"")

	// Verify popular books tab content
	assert.Contains(t, html, "id=\"popular-tab\"")
	assert.Contains(t, html, "Test Book 1")
	assert.Contains(t, html, "Test Book 2")
}

func TestRenderSearchResults(t *testing.T) {
	// Test books data
	books := []Book{
		{Title: "Search Result 1", Link: "https://example.com/search1"},
		{Title: "Search Result 2", Link: "https://example.com/search2"},
	}

	// Get the template
	tmpl, err := renderSearchResults(books)

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
	assert.Contains(t, html, "Search Result 1")
	assert.Contains(t, html, "https://example.com/search1")
	assert.Contains(t, html, "Search Result 2")
	assert.Contains(t, html, "https://example.com/search2")

	// Verify memorize buttons are present
	assert.Contains(t, html, "memorize-btn")
	assert.Contains(t, html, "hx-post=\"/memorize-book\"")

	// Verify it contains the book list structure for search results
	assert.Contains(t, html, "class=\"book-list\"")
	assert.Contains(t, html, "class=\"book-item\"")
	assert.Contains(t, html, "class=\"book-actions\"")

	// Verify it's a partial template (no full HTML structure)
	assert.NotContains(t, html, "<!DOCTYPE html>")
	assert.NotContains(t, html, "<title>")
}

func TestRenderSearchResultsEmpty(t *testing.T) {
	// Empty books slice
	var books []Book

	// Get the template
	tmpl, err := renderSearchResults(books)

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
	assert.Contains(t, html, "No books found")
	assert.Contains(t, html, "class=\"no-results\"")
}

func TestRenderMemorizedBooks(t *testing.T) {
	// Test books data
	books := []Book{
		{Title: "Memorized Book 1", Link: "https://example.com/memo1"},
		{Title: "Memorized Book 2", Link: "https://example.com/memo2"},
	}

	// Get the template
	tmpl, err := renderMemorizedBooks(books)

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
	assert.Contains(t, html, "Memorized Book 1")
	assert.Contains(t, html, "https://example.com/memo1")
	assert.Contains(t, html, "Memorized Book 2")
	assert.Contains(t, html, "https://example.com/memo2")

	// Verify remove buttons are present
	assert.Contains(t, html, "remove-btn")
	assert.Contains(t, html, "hx-post=\"/remove-memorized-book\"")

	// Verify memorized books header
	assert.Contains(t, html, "memorized-header")
	assert.Contains(t, html, "2 memorized books")
	assert.Contains(t, html, "memorized-count")

	// Verify it contains the book list structure
	assert.Contains(t, html, "class=\"book-list\"")
	assert.Contains(t, html, "class=\"book-item\"")
	assert.Contains(t, html, "class=\"book-actions\"")

	// Verify it's a partial template (no full HTML structure)
	assert.NotContains(t, html, "<!DOCTYPE html>")
	assert.NotContains(t, html, "<title>")
}

func TestRenderMemorizedBooksEmpty(t *testing.T) {
	// Empty books slice
	var books []Book

	// Get the template
	tmpl, err := renderMemorizedBooks(books)

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

	// Verify the HTML contains the "no memorized books" message
	assert.Contains(t, html, "No memorized books yet")
	assert.Contains(t, html, "class=\"no-results\"")
	assert.Contains(t, html, "memorized-header")
	assert.Contains(t, html, "0 memorized books")
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
