package main

import (
	"fmt"
	"log"
	"net/http"
)

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
