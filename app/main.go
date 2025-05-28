package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

var cachedBooks []Book

func booksHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	
	// If no cached books, fetch them
	if len(cachedBooks) == 0 {
		cachedBooks, err = fetchPopularBooks()
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to fetch books %s", err), http.StatusInternalServerError)
			return
		}
	}

	tmpl, err := renderPage(cachedBooks)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse template: %s", err), http.StatusInternalServerError)
		return
	}

	// Execute the template with the books data
	err = tmpl.Execute(w, cachedBooks)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to execute template: %s", err), http.StatusInternalServerError)
		return
	}
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}
	
	searchQuery := strings.ToLower(r.FormValue("search"))
	
	// Filter the cached books based on the search query
	var filteredBooks []Book
	if searchQuery != "" {
		for _, book := range cachedBooks {
			if strings.Contains(strings.ToLower(book.Title), searchQuery) {
				filteredBooks = append(filteredBooks, book)
			}
		}
	} else {
		filteredBooks = cachedBooks
	}
	
	// Render just the book list part
	tmpl, err := renderBookList(filteredBooks)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse template: %s", err), http.StatusInternalServerError)
		return
	}
	
	err = tmpl.Execute(w, filteredBooks)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to execute template: %s", err), http.StatusInternalServerError)
		return
	}
}

func refreshHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	
	// Refetch books from the source
	cachedBooks, err = fetchPopularBooks()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch books %s", err), http.StatusInternalServerError)
		return
	}
	
	// Render just the book list part
	tmpl, err := renderBookList(cachedBooks)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse template: %s", err), http.StatusInternalServerError)
		return
	}
	
	err = tmpl.Execute(w, cachedBooks)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to execute template: %s", err), http.StatusInternalServerError)
		return
	}
}

func main() {
	// Initialize books on startup
	var err error
	cachedBooks, err = fetchPopularBooks()
	if err != nil {
		log.Printf("Warning: Failed to pre-fetch books: %s", err)
	}
	
	// Register routes
	http.HandleFunc("/", booksHandler)
	http.HandleFunc("/search", searchHandler)
	http.HandleFunc("/refresh", refreshHandler)
	
	fmt.Println("Starting server on :8090")
	if err := http.ListenAndServe(":8090", nil); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}
