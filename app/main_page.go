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