package main

import (
	"embed"
	"html/template"
)

//go:embed templates/*
var templateFS embed.FS

func renderPage(books []Book) (*template.Template, error) {
	// Parse the main HTML template from embedded filesystem
	tmpl, err := template.ParseFS(templateFS, "templates/main.html")
	if err != nil {
		return nil, err
	}

	return tmpl, nil
}

// renderBookList renders just the book list for HTMX partial updates
func renderBookList(books []Book) (*template.Template, error) {
	// Parse the partial book list template from embedded filesystem
	tmpl, err := template.ParseFS(templateFS, "templates/book_list.html")
	if err != nil {
		return nil, err
	}

	return tmpl, nil
}

// renderTabbedMain renders the main tabbed interface
func renderTabbedMain(books []Book) (*template.Template, error) {
	// Parse the tabbed main template from embedded filesystem
	tmpl, err := template.ParseFS(templateFS, "templates/tabbed_main.html")
	if err != nil {
		return nil, err
	}

	return tmpl, nil
}

// renderSearchResults renders search results for HTMX partial updates
func renderSearchResults(books []Book) (*template.Template, error) {
	// Parse the search results template from embedded filesystem
	tmpl, err := template.ParseFS(templateFS, "templates/search_results.html")
	if err != nil {
		return nil, err
	}

	return tmpl, nil
}

// renderMemorizedBooks renders memorized books for HTMX partial updates
func renderMemorizedBooks(books []Book) (*template.Template, error) {
	// Parse the memorized books template from embedded filesystem
	tmpl, err := template.ParseFS(templateFS, "templates/memorized_books.html")
	if err != nil {
		return nil, err
	}

	return tmpl, nil
}