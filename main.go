package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gocolly/colly/v2"
)

type Book struct {
	Title string
	Link  string
}

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

	return books, nil
}

func booksHandler(w http.ResponseWriter, r *http.Request) {
	books, err := fetchPopularBooks()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch books %s", err), http.StatusInternalServerError)
		return
	}

	var sb strings.Builder
	sb.WriteString("<html><body><h1>Top 10 Popular Books on Royal Road</h1><ul>")
	for _, book := range books {
		sb.WriteString(fmt.Sprintf("<li><a href=\"%s\">%s</a></li>", book.Link, book.Title))
	}
	sb.WriteString("</ul></body></html>")

	fmt.Fprint(w, sb.String())
}

func main() {
	http.HandleFunc("/", booksHandler)
	fmt.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}
