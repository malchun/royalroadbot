package main

import (
	"html/template"
)

type Book struct {
	Title string
	Link  string
}

func renderPage(books []Book) (*template.Template, error) {
	// Parse the HTML template
	tmpl, err := template.New("bookTemplate").Parse(`
<!DOCTYPE html>
<html lang="en">
<head>
				<meta charset="UTF-8">
				<meta name="viewport" content="width=device-width, initial-scale=1.0">
				<title>Royal Road - Popular Books</title>
				<style>
								body {
												font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
												line-height: 1.6;
												color: #333;
												max-width: 800px;
												margin: 0 auto;
												padding: 20px;
												background-color: #f5f5f5;
								}
								h1 {
												color: #2c3e50;
												text-align: center;
												margin-bottom: 30px;
												border-bottom: 2px solid #3498db;
												padding-bottom: 10px;
								}
								.search-container {
												margin-bottom: 20px;
												text-align: center;
								}
								#searchInput {
												padding: 8px 15px;
												width: 70%;
												border: 1px solid #ddd;
												border-radius: 4px;
												font-size: 16px;
								}
								.book-list {
												list-style-type: none;
												padding: 0;
								}
								.book-item {
												background-color: white;
												margin-bottom: 10px;
												padding: 15px;
												border-radius: 5px;
												box-shadow: 0 2px 5px rgba(0,0,0,0.1);
												transition: transform 0.2s;
								}
								.book-item:hover {
												transform: translateY(-3px);
												box-shadow: 0 5px 15px rgba(0,0,0,0.1);
								}
								.book-item a {
												color: #3498db;
												text-decoration: none;
												font-weight: bold;
												font-size: 18px;
								}
								.book-item a:hover {
												text-decoration: underline;
								}
								.no-results {
												text-align: center;
												font-style: italic;
												color: #7f8c8d;
												display: none;
								}
								footer {
												margin-top: 30px;
												text-align: center;
												font-size: 14px;
												color: #7f8c8d;
								}
				</style>
</head>
<body>
				<h1>Top 10 Popular Books on Royal Road</h1>

				<div class="search-container">
								<input type="text" id="searchInput" placeholder="Search for books..." oninput="searchBooks()">
				</div>

				<ul class="book-list" id="bookList">
								{{range .}}
								<li class="book-item">
												<a href="{{.Link}}" target="_blank">{{.Title}}</a>
								</li>
								{{end}}
				</ul>

				<div class="no-results" id="noResults">
								No books found matching your search.
				</div>

				<footer>
								Data scraped from Royal Road's Active Popular Fiction List
				</footer>

				<script>
								function searchBooks() {
												const input = document.getElementById('searchInput');
												const filter = input.value.toUpperCase();
												const bookList = document.getElementById('bookList');
												const books = bookList.getElementsByTagName('li');
												const noResults = document.getElementById('noResults');

												let resultsFound = false;

												for (let i = 0; i < books.length; i++) {
																const bookTitle = books[i].getElementsByTagName('a')[0];
																const txtValue = bookTitle.textContent || bookTitle.innerText;

																if (txtValue.toUpperCase().indexOf(filter) > -1) {
																				books[i].style.display = "";
																				resultsFound = true;
																} else {
																				books[i].style.display = "none";
																}
												}

												if (resultsFound) {
																noResults.style.display = "none";
												} else {
																noResults.style.display = "block";
												}
								}
				</script>
</body>
</html>`)

	if err != nil {
		return nil, err
	}

	return tmpl, nil
}
