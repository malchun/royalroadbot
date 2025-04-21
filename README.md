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

For local development with MongoDB client:
```bash
just dev-mongo     # Start MongoDB and Mongo Express in containers
just build         # Build the Go application locally
just dev-run       # Build and run with local MongoDB (combines the above two steps)
```

For containerized deployment:
```bash
just docker-build  # Build Docker images
just run           # Start all services via Docker Compose
```

3. Testing the application:
```bash
just test          # Run tests locally
just test-docker   # Run tests in Docker container
```

4. Access the services:
- Web application: [http://localhost:8090](http://localhost:8090)
- MongoDB client (Mongo Express): [http://localhost:8081](http://localhost:8081) (when using dev-mongo)

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

4. Development environment with MongoDB + Mongo Express:
```bash
docker-compose -f docker-compose-dev.yaml up
```

Access the web service at [http://localhost:8090](http://localhost:8090)

For MongoDB Express web client, access [http://localhost:8081](http://localhost:8081)

## Project Structure

- `app/`
  - `main.go`: Web server setup and request handling
  - `crawler.go`: Web scraping functionality for RoyalRoad.com
  - `main_page.go`: HTML template rendering and Book data structure
  - `database.go`: MongoDB integration and data persistence
  - `database_test.go`: Database operation tests
  - `crawler_test.go`: Web scraper tests
  - `main_page_test.go`: Template rendering tests
- `Dockerfile`: Instructions for building the Docker container
- `docker-compose.yaml`: Main Docker Compose configuration
- `docker-compose-test.yaml`: Test environment configuration
- `docker-compose-mongo.yaml`: MongoDB-only configuration
- `docker-compose-dev.yaml`: Development environment with MongoDB and Mongo Express
- `justfile`: Task automation commands:
  - `build`: Build the Go application locally
  - `docker-build`: Build Docker containers
  - `dev-mongo`: Start MongoDB with Mongo Express for development
  - `dev-run`: Build and run with local MongoDB
  - `run`: Run all services with Docker Compose
  - `test`: Run tests locally
  - `test-docker`: Run tests with Docker
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
- Further improve the application structure by creating dedicated packages:
  - `models` for data structures
  - `api` for any future API endpoints

### 3. Performance
- Add caching to prevent scraping RoyalRoad on every request
- Implement proper rate limiting to be respectful to the target website
- Optimize database queries

### 4. User Experience
- Add CSS styling to improve the presentation
- Consider adding more details about each book (cover images, ratings, etc.)
- Implement pagination or filtering options

### 5. Testing
- Add integration tests for the HTTP endpoints
- Add end-to-end tests with Docker Compose test environment
- Fix remaining test dependency issues

### 6. Configuration
- Move hardcoded values to environment variables or a config file
- Make the port configurable
- Allow setting the number of books to display

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