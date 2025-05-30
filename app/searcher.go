package main

import (
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/gocolly/colly/v2"
)

// searchRoyalRoadBooks searches Royal Road for books matching the given query
// and returns up to 15 books with their titles and links
func searchRoyalRoadBooks(query string) ([]Book, error) {
	if strings.TrimSpace(query) == "" {
		return []Book{}, nil
	}

	// Construct the search URL
	searchURL := fmt.Sprintf("https://www.royalroad.com/fictions/search?title=%s", url.QueryEscape(query))
	
	c := colly.NewCollector()
	var books []Book

	// Target the h2 elements that contain book titles and links
	c.OnHTML("h2", func(e *colly.HTMLElement) {
		// Look for anchor tags within h2 elements
		title := strings.TrimSpace(e.ChildText("a"))
		link := e.ChildAttr("a", "href")
		
		if title != "" && link != "" {
			// Ensure we have the full URL
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

	// Set user agent to be respectful
	c.UserAgent = "Mozilla/5.0 (compatible; RoyalRoadBot/1.0)"

	err := c.Visit(searchURL)
	if err != nil {
		return nil, fmt.Errorf("failed to visit search URL: %w", err)
	}

	// Limit results to 15 books
	if len(books) > 15 {
		books = books[:15]
	}

	return books, nil
}

// memorizeBook saves a specific book to the database in the memorizedBooks collection
func memorizeBook(book Book) error {
	if strings.TrimSpace(book.Title) == "" || strings.TrimSpace(book.Link) == "" {
		return fmt.Errorf("book title and link cannot be empty")
	}

	err := saveBookToMemory(book)
	if err != nil {
		log.Printf("Failed to memorize book '%s': %v", book.Title, err)
		return fmt.Errorf("failed to save book to memory: %w", err)
	}

	log.Printf("Successfully memorized book: %s", book.Title)
	return nil
}

// removeMemorizedBook removes a book from the memorized collection by title
func removeMemorizedBook(title string) error {
	if strings.TrimSpace(title) == "" {
		return fmt.Errorf("book title cannot be empty")
	}

	err := removeBookFromMemory(title)
	if err != nil {
		log.Printf("Failed to remove memorized book '%s': %v", title, err)
		return fmt.Errorf("failed to remove book from memory: %w", err)
	}

	log.Printf("Successfully removed memorized book: %s", title)
	return nil
}

// getMemorizedBooks retrieves all memorized books from the database
func getMemorizedBooks() ([]Book, error) {
	books, err := loadMemorizedBooks()
	if err != nil {
		log.Printf("Failed to load memorized books: %v", err)
		return nil, fmt.Errorf("failed to load memorized books: %w", err)
	}

	return books, nil
}