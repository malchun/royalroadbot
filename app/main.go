package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
)

var (
	cachedBooks []Book
	booksMutex  sync.RWMutex
)

func booksHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	
	// If no cached books, fetch them
	booksMutex.RLock()
	needsFetch := len(cachedBooks) == 0
	booksMutex.RUnlock()
	
	if needsFetch {
		booksMutex.Lock()
		// Double-check after acquiring write lock
		if len(cachedBooks) == 0 {
			cachedBooks, err = fetchPopularBooks()
			if err != nil {
				booksMutex.Unlock()
				http.Error(w, fmt.Sprintf("Failed to fetch books: %s", err), http.StatusInternalServerError)
				return
			}
		}
		booksMutex.Unlock()
	}

	booksMutex.RLock()
	booksCopy := make([]Book, len(cachedBooks))
	copy(booksCopy, cachedBooks)
	booksMutex.RUnlock()

	tmpl, err := renderPage(booksCopy)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse template: %s", err), http.StatusInternalServerError)
		return
	}

	// Execute the template with the books data
	err = tmpl.Execute(w, booksCopy)
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
	booksMutex.RLock()
	var filteredBooks []Book
	if searchQuery != "" {
		for _, book := range cachedBooks {
			if strings.Contains(strings.ToLower(book.Title), searchQuery) {
				filteredBooks = append(filteredBooks, book)
			}
		}
	} else {
		filteredBooks = make([]Book, len(cachedBooks))
		copy(filteredBooks, cachedBooks)
	}
	booksMutex.RUnlock()
	
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
	newBooks, err := fetchPopularBooks()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch books: %s", err), http.StatusInternalServerError)
		return
	}
	
	booksMutex.Lock()
	cachedBooks = newBooks
	booksCopy := make([]Book, len(cachedBooks))
	copy(booksCopy, cachedBooks)
	booksMutex.Unlock()
	
	// Render just the book list part
	tmpl, err := renderBookList(booksCopy)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse template: %s", err), http.StatusInternalServerError)
		return
	}
	
	err = tmpl.Execute(w, booksCopy)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to execute template: %s", err), http.StatusInternalServerError)
		return
	}
}

func main() {
	// Initialize books on startup
	var err error
	initialBooks, err := fetchPopularBooks()
	if err != nil {
		log.Printf("Warning: Failed to pre-fetch books: %s", err)
	} else {
		booksMutex.Lock()
		cachedBooks = initialBooks
		booksMutex.Unlock()
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
