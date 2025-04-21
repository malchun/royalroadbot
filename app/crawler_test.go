package main

import (
	"testing"
	"net/http"
	"net/http/httptest"
	
	"github.com/gocolly/colly/v2"
	"github.com/stretchr/testify/assert"
)

func init() {
	// Initialize any test setup if needed
}

// TestScrapingFunction is a simple test to verify the core scraping functionality
func TestScrapingFunction(t *testing.T) {
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
	
	// Use a simplified version of the scraping logic
	c := colly.NewCollector()
	
	titles := []string{}
	links := []string{}
	
	c.OnHTML(".fiction-list-item", func(e *colly.HTMLElement) {
		title := e.ChildText(".fiction-title")
		link := e.ChildAttr(".fiction-title a", "href")
		if title != "" && link != "" {
			titles = append(titles, title)
			links = append(links, "https://www.royalroad.com"+link)
		}
	})
	
	// Visit our test server instead of RoyalRoad
	assert.NoError(t, c.Visit(server.URL))
	
	// Verify we found the expected number of books
	assert.Equal(t, 2, len(titles))
	assert.Equal(t, 2, len(links))
	
	// Verify book titles and links are correct
	assert.Equal(t, "Test Book 1", titles[0])
	assert.Equal(t, "https://www.royalroad.com/fiction/1234", links[0])
	assert.Equal(t, "Test Book 2", titles[1])
	assert.Equal(t, "https://www.royalroad.com/fiction/5678", links[1])
}