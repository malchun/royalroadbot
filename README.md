# RoyalRoadBot

A simple web service that scrapes and displays the top 10 popular books from RoyalRoad.com.

## Quick Start Guide

### Prerequisites
- Go 1.20 or higher
- Docker and Docker Compose (optional, for containerized deployment)

### Local Development

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/royalroadbot.git
   cd royalroadbot
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Run the application:
   ```bash
   go run main.go
   ```

4. Access the web service at [http://localhost:8090](http://localhost:8090)

### Docker Deployment

1. Build and start the container:
   ```bash
   docker-compose up --build
   ```

2. Access the web service at [http://localhost:8090](http://localhost:8090)

## Project Structure

- `main.go`: Contains the core application logic
- `Dockerfile`: Instructions for building the Docker container
- `docker-compose.yaml`: Configuration for Docker Compose deployment

## Code Overview

The application performs the following tasks:
1. Scrapes the "active-popular" fiction list from RoyalRoad.com
2. Extracts the top 10 book titles and links
3. Presents them as a simple HTML list via a web server

## Areas for Improvement

### 1. Error Handling
- Add more robust error handling, especially for network failures
- Implement retries for web scraping
- Add logging for debugging purposes

### 2. Code Organization
- Consider breaking the application into packages:
  - `scraper` for the web scraping logic
  - `server` for the HTTP server implementation
  - `models` for data structures

### 3. Performance
- Add caching to prevent scraping RoyalRoad on every request
- Implement proper rate limiting to be respectful to the target website

### 4. User Experience
- Add CSS styling to improve the presentation
- Consider adding more details about each book (cover images, ratings, etc.)
- Implement pagination or filtering options

### 5. Testing
- Add unit tests for the scraping logic
- Add integration tests for the HTTP endpoints

### 6. Configuration
- Move hardcoded values to environment variables or a config file
- Make the port configurable
- Allow setting the number of books to display

### 7. Documentation
- Add godoc comments to functions and types
- Create API documentation if you expand the service

## Best Practices

1. **Rate Limiting**: Be respectful when scraping websites. Consider adding delays between requests or implementing a proper rate limiter.

2. **Terms of Service**: Ensure you're complying with RoyalRoad's terms of service regarding scraping.

3. **Error Handling**: Never ignore errors. Log them appropriately and return meaningful error messages.

4. **Code Reviews**: Request code reviews for any significant changes.

5. **Git Workflow**: Use feature branches and submit PRs for changes instead of committing directly to main.

## Common Issues and Solutions

### The scraper isn't finding any books
- Check if RoyalRoad's HTML structure has changed
- Use browser developer tools to identify updated selectors

### Docker container exits immediately
- Check logs with `docker-compose logs`
- Ensure the Go application is properly building

### Port conflicts
- If port 8090 is already in use, modify the port in both `main.go` and `docker-compose.yaml`

## Future Enhancements

1. Add a database to store book information
2. Create an API endpoint to return books in JSON format
3. Implement user authentication to allow saving favorite books
4. Set up scheduled scraping to maintain up-to-date information
5. Add metrics and monitoring

## Resources

- [Colly Documentation](http://go-colly.org/docs/)
- [Go HTTP Package Documentation](https://pkg.go.dev/net/http)
- [Docker Compose Documentation](https://docs.docker.com/compose/)
