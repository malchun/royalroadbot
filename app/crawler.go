package main

import (
	"log"

	"github.com/gocolly/colly/v2"
)

// fetchPopularBooks scrapes the RoyalRoad website for popular books
// and returns the top 10 books with their titles and links
func fetchPopularBooks() ([]Book, error) {
	// Map of kinds and their corresponding URLs
	crawlURLs := map[string]string{
		"popular": "https://www.royalroad.com/fictions/active-popular",
	}
	return fetchBooks(crawlURLs["popular"])
}

func fetchBooks(crawlUrl string) ([]Book, error) {
	c := colly.NewCollector()

	var books []Book

	c.OnHTML(".fiction-list-item", func(e *colly.HTMLElement) {
		title := e.ChildText(".fiction-title")
		link := e.ChildAttr(".fiction-title a", "href")
		if title != "" && link != "" {
			books = append(books, Book{
				Title: title,
				Link:  "https://www.royalroad.com" + link,
			})
		}
	})

	err := c.Visit(crawlUrl)
	if err != nil {
		return nil, err
	}

	if len(books) > 10 {
		books = books[:10]
	}

	// Popular books are now kept in memory only (not saved to database)
	log.Printf("Fetched %d popular books", len(books))

	return books, nil
}
