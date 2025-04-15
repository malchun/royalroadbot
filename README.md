# RoyalRoadBot

A web service that scrapes, stores, and displays the top 10 popular books from RoyalRoad.com.

## Quick Start Guide

### Prerequisites
- Go 1.23 or higher
- Docker and Docker Compose (optional, for containerized deployment)
- MongoDB (included in docker-compose configuration)

### Local Development

1. Clone the repository:
```bash
git clone https://github.com/yourusername/royalroadbot.git
cd royalroadbot
```

2. Build and run the application:
```bash
just build
just run
```

3. Access the web service at [http://localhost:8090](http://localhost:8090)

### Docker Deployment

The project includes several Docker Compose configurations:

1. Standard deployment (with MongoDB):
```bash
docker-compose up --build
```

2. Test environment:
```bash
docker-compose -f docker-compose-test.yaml up --build
```

3. MongoDB only (for local development):
```bash
docker-compose -f docker-compose-mongo.yaml up
```

Access the web service at [http://localhost:8090](http://localhost:8090)

## Project Structure

- `app/`
  - `main.go`: Core application logic and web server
  - `main_page.go`: HTML template rendering
  - `database.go`: MongoDB integration and data persistence
  - `database_test.go`: Database operation tests
- `Dockerfile`: Instructions for building the Docker container
- `docker-compose.yaml`: Main Docker Compose configuration
- `docker-compose-test.yaml`: Test environment configuration
- `docker-compose-mongo.yaml`: MongoDB-only configuration
- `justfile`: Task automation commands
- `go.mod`, `go.sum`: Go dependencies

## Code Overview

The application performs the following tasks:
1. Scrapes the "active-popular" fiction list from RoyalRoad.com using Colly
2. Extracts the top 10 book titles and links
3. Stores the book data in MongoDB for persistence
4. Presents them as a simple HTML list via a web server

## Main Dependencies

- [Colly](http://go-colly.org/docs/): Web scraping framework
- [MongoDB Go Driver](https://pkg.go.dev/go.mongodb.org/mongo-driver): Database operations
- [Testify](https://github.com/stretchr/testify): Testing framework

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
- Optimize database queries

### 4. User Experience
- Add CSS styling to improve the presentation
- Consider adding more details about each book (cover images, ratings, etc.)
- Implement pagination or filtering options

### 5. Testing
- Expand unit test coverage
- Add integration tests for the HTTP endpoints
- Add end-to-end tests with Docker Compose test environment

### 6. Configuration
- Move hardcoded values to environment variables or a config file
- Make the port configurable
- Allow setting the number of books to display
- Add MongoDB configuration options

### 7. Documentation
- Add godoc comments to functions and types
- Create API documentation if you expand the service
- Document database schema and operations

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
- Check logs for any scraping errors

### Docker container exits immediately
- Check logs with `docker-compose logs`
- Ensure the Go application is properly building
- Verify MongoDB connection settings

### Port conflicts
- If port 8090 is already in use, modify the port in both `main.go` and `docker-compose.yaml`
- Check if MongoDB port (27017) is available when running locally

### Database connection issues
- Verify MongoDB is running and accessible
- Check connection string configuration
- Ensure proper network configuration in Docker Compose

## Future Enhancements

1. Create an API endpoint to return books in JSON format
2. Implement user authentication to allow saving favorite books
3. Set up scheduled scraping to maintain up-to-date information
4. Add metrics and monitoring
5. Implement book category filtering
6. Add a search function for stored books

## Resources

- [Colly Documentation](http://go-colly.org/docs/)
- [Go HTTP Package Documentation](https://pkg.go.dev/net/http)
- [MongoDB Go Driver Documentation](https://pkg.go.dev/go.mongodb.org/mongo-driver)
- [Docker Compose Documentation](https://docs.docker.com/compose/)