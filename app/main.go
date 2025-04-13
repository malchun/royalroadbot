package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gocolly/colly/v2"
)

func fetchPopularBooks() ([]Book, error) {
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

	err := c.Visit("https://www.royalroad.com/fictions/active-popular")
	if err != nil {
		return nil, err
	}

	if len(books) > 10 {
		books = books[:10]
	}

	// Save fetched books to MongoDB
	err = saveBooksWithMetadata(books)
	if err != nil {
		log.Printf("Failed to save books: %v", err)
	} else {
		log.Println("Books saved successfully!")
	}

	return books, nil
}

func booksHandler(w http.ResponseWriter, r *http.Request) {
	books, err := fetchPopularBooks()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch books %s", err), http.StatusInternalServerError)
		return
	}

	tmpl, err := renderPage(books)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse template: %s", err), http.StatusInternalServerError)
		return
	}

	// Execute the template with the books data
	err = tmpl.Execute(w, books)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to execute template: %s", err), http.StatusInternalServerError)
		return
	}
}

func main() {
	http.HandleFunc("/", booksHandler)
	fmt.Println("Starting server on :8090")
	if err := http.ListenAndServe(":8090", nil); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}
