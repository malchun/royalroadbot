# RoyalRoadBot

A web service that scrapes, stores, and displays the top 10 popular books from RoyalRoad.com.

## Quick Start Guide

### Prerequisites
- Go 1.23 or higher (currently using Go 1.24)
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
just run-dev-mongo    # Start MongoDB and Mongo Express in containers
just build           # Build the Go application locally
just run-dev-local   # Build and run with local MongoDB
```

For containerized deployment:
```bash
just build-docker    # Build Docker images
just run            # Start all services via Docker Compose
```

3. Testing the application:
```bash
just test-local     # Run tests locally
just test-docker    # Run tests in Docker container
```

4. Access the services:
- Web application: [http://localhost:8090](http://localhost:8090)
- MongoDB client (Mongo Express): [http://localhost:8081](http://localhost:8081) (when using run-dev-mongo)

### Docker Deployment

The project includes several Docker Compose configurations:

1. Standard deployment (with MongoDB):
```bash
just run  # or docker-compose up --build
```

2. Test environment:
```bash
just test-docker  # or docker-compose -f docker-compose-test.yaml up --build
```

3. MongoDB only (for local development):
```bash
just run-mongo  # or docker-compose -f docker-compose-mongo.yaml up
```

4. Development environment with MongoDB + Mongo Express:
```bash
just run-dev-mongo  # or docker-compose -f docker-compose-mongo.yaml -f docker-compose-dev.yaml up
```

Access the web service at [http://localhost:8090](http://localhost:8090)

For MongoDB Express web client, access [http://localhost:8081](http://localhost:8081)

## Project Structure

- `app/`
  - `main.go`: Web server setup and request handling
  - `model.go`: Book data structure definition
  - `crawler.go`: Web scraping functionality for RoyalRoad.com
  - `main_page.go`: HTML template rendering for the front-end
  - `database.go`: MongoDB integration and data persistence
  - `templates/`: HTML templates directory
    - `main.html`: Main page template with theme support
    - `book_list.html`: Partial template for HTMX updates
  - `database_test.go`: Database operation tests
  - `crawler_test.go`: Web scraper tests
  - `main_page_test.go`: Template rendering tests
- `Dockerfile`: Instructions for building the Docker container
- `Dockerfile.test`: Instructions for building the test container
- `docker-compose.yaml`: Main Docker Compose configuration
- `docker-compose-test.yaml`: Test environment configuration
- `docker-compose-mongo.yaml`: MongoDB-only configuration
- `docker-compose-dev.yaml`: Development environment with MongoDB and Mongo Express
- `justfile`: Task automation commands for building, running, and testing
- `go.mod`, `go.sum`: Go dependencies

## Current Functionality

The application currently performs the following tasks:
1. Scrapes the "active-popular" fiction list from RoyalRoad.com using Colly
2. Extracts the top 10 book titles and links
3. Stores the book data in MongoDB for persistence
4. Presents books as a styled HTML list via a web server with a search function

### Web Interface Features:
- Clean, responsive UI with modern styling
- **Dark/Light theme toggle** with persistent user preference
- Client-side search functionality for filtering books
- HTMX-powered real-time search with debouncing
- Direct links to the books on RoyalRoad.com
- **Modular template system** with embedded filesystem

## Main Dependencies

- **Go 1.24**: Latest stable version of Go
- **[Colly v2.1.0](http://go-colly.org/docs/)**: Web scraping framework
- **[MongoDB Go Driver v1.17.3](https://pkg.go.dev/go.mongodb.org/mongo-driver)**: Database operations
- **[Testify v1.10.0](https://github.com/stretchr/testify)**: Testing framework
- **Docker & Docker Compose**: Containerization and service orchestration
- **Just**: Task runner for command automation

## Development Commands

The project includes a `justfile` with many helpful commands:

- `just build` - Build the Go application locally
- `just build-docker` - Build Docker containers
- `just run-dev-mongo` - Start MongoDB with Mongo Express for development
- `just run-dev-local` - Build and run with local MongoDB
- `just run` - Run all services with Docker Compose
- `just run-mongo` - Run just the MongoDB service
- `just run-detached` - Run containers in detached mode
- `just stop` - Stop running containers
- `just stop-all` - Stop all containers (including dev environments)
- `just logs` - Show logs from running containers
- `just clean` - Remove containers, images, and volumes
- `just restart` - Rebuild and restart all containers
- `just test-docker` - Run tests in Docker environment
- `just test-local` - Run tests locally

## Areas for Improvement

### 1. Error Handling
- Add more robust error handling, especially for network failures
- Implement retries for web scraping
- Add structured logging for monitoring and debugging

### 2. Code Organization
- ✅ **Improved template organization** - Templates extracted to separate files
- Further improve the application structure by creating dedicated packages:
  - `models` for data structures
  - `api` for REST endpoints
  - `storage` for database operations

### 3. Performance
- Add caching to prevent scraping RoyalRoad on every request
- Implement proper rate limiting to be respectful to the target website
- Optimize database queries and add indexes

### 4. User Experience
- ✅ **Dark/Light theme support** - Toggle with persistent preferences
- Add more details about each book (cover images, ratings, synopsis)
- Implement pagination for larger datasets
- Add sorting options (by popularity, rating, etc.)

### 5. Testing
- Add integration tests for the HTTP endpoints
- Add end-to-end tests with Docker Compose test environment
- Increase test coverage

### 6. Configuration
- Move hardcoded values to environment variables or a config file
- Make the port configurable
- Allow setting the number of books to display

### 7. Documentation
- Add godoc comments to functions and types
- Create API documentation for future endpoints
- Document database schema and operations

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

## Recent Enhancements

### Theme Support & Template Refactoring (Latest)
- **Dark/Light Theme Toggle**: Users can switch between themes with a button
- **Theme Persistence**: User preference saved in localStorage
- **CSS Custom Properties**: Clean variable-based theming system
- **Template Extraction**: HTML templates moved to separate files in `app/templates/`
- **Embedded Filesystem**: Templates compiled into binary using Go's embed directive
- **Improved Modularity**: Separate templates for main page and HTMX partials

## Future Enhancements

1. Create a REST API for programmatic access to book data
2. Implement user accounts to save favorite books
3. Set up scheduled scraping to maintain up-to-date information
4. Add metrics and monitoring (Prometheus/Grafana)
5. Implement book category filtering and tags
6. Add server-side pagination and advanced search options
7. Collect and display more book metadata (ratings, chapters, etc.)
8. Add more theme options (custom colors, high contrast mode)
9. Implement keyboard shortcuts for theme switching

## Resources

- [Colly Documentation](http://go-colly.org/docs/)
- [Go HTTP Package Documentation](https://pkg.go.dev/net/http)
- [MongoDB Go Driver Documentation](https://pkg.go.dev/go.mongodb.org/mongo-driver)
- [Docker Compose Documentation](https://docs.docker.com/compose/)