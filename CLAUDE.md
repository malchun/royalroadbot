# CLAUDE.md - AI Assistant Guide for RoyalRoadBot

## Project Overview

**RoyalRoadBot** is a Go web application that scrapes, stores, and displays the top 10 popular books from RoyalRoad.com. It features a modern web interface with real-time search capabilities and data persistence through MongoDB.

### Key Information
- **Language**: Go 1.24
- **License**: MIT
- **Author**: malchun
- **Main Port**: 8090
- **Database**: MongoDB

## Architecture & Structure

### Core Components

```
royalroadbot/
├── app/                    # Main application code
│   ├── main.go            # HTTP server and route handlers
│   ├── model.go           # Data structures (Book)
│   ├── crawler.go         # Web scraping logic (Colly)
│   ├── database.go        # MongoDB operations
│   ├── main_page.go       # Template rendering with embedded filesystem
│   ├── templates/         # HTML templates directory
│   │   ├── main.html     # Main page template with theme support
│   │   └── book_list.html # Partial template for HTMX updates
│   └── *_test.go          # Comprehensive test suite
├── bin/                   # Compiled binaries
├── docker-compose*.yaml   # Multiple Docker configurations
├── Dockerfile             # Multi-stage container build
├── justfile              # Task automation (like Makefile)
├── go.mod/go.sum         # Go dependencies
└── README.md             # User documentation
```

### Technology Stack

**Backend:**
- Go 1.24 with standard library HTTP server
- Colly v2.1.0 for web scraping
- MongoDB Go Driver v1.17.3 for data persistence
- Testcontainers for integration testing

**Frontend:**
- Server-side HTML templates with embedded filesystem
- HTMX 1.9.6 for dynamic interactions
- Modern CSS with responsive design and theme support
- Minimal JavaScript for theme switching
- Dark/Light theme toggle with localStorage persistence

**Infrastructure:**
- Docker with multi-stage builds
- Docker Compose for orchestration
- MongoDB for data storage
- Mongo Express for database administration

## Key Features

### Web Scraping
- Scrapes RoyalRoad.com's "active-popular" fiction list
- Extracts book titles and links using CSS selectors
- Limits results to top 10 books
- Handles network errors gracefully

### Web Interface
- Real-time search with 500ms debounce
- Manual refresh capability
- Responsive design for all devices
- HTMX-powered partial page updates
- Direct links to books on RoyalRoad.com
- **Dark/Light theme toggle** with persistent user preference
- **CSS Custom Properties** for clean theming system
- **Template modularity** with separate files for main page and partials

### Data Management
- MongoDB persistence with automatic connection management
- In-memory caching for performance
- Metadata storage with timestamps
- Connection pooling and proper cleanup

## Development Workflow

### Quick Start Commands (using Just)

```bash
# Local development with MongoDB in Docker
just run-dev-mongo     # Start MongoDB + Mongo Express
just run-dev-local     # Build and run app locally

# Full containerized deployment
just build-docker      # Build all containers
just run              # Start full stack

# Testing
just test-local       # Run tests locally
just test-docker      # Run tests in containers

# Utilities
just clean            # Clean up containers/images
just logs             # View container logs
```

### Docker Configurations

1. **docker-compose.yaml** - Main application stack
2. **docker-compose-mongo.yaml** - MongoDB only
3. **docker-compose-dev.yaml** - Adds Mongo Express web client
4. **docker-compose-test.yaml** - Test environment

### Environment Variables

- `MONGODB_URI` - MongoDB connection string
  - Local dev: `mongodb://admin:password@127.0.0.1:27017`
  - Docker: Handled by docker-compose networking

## Code Structure Analysis

### Main Application (`app/main.go`)
- HTTP server on port 8090
- Three main routes:
  - `/` - Main page with book list
  - `/search` - HTMX search endpoint
  - `/refresh` - Manual data refresh
- Global caching with `cachedBooks` variable
- Error handling with HTTP status codes

### Data Model (`app/model.go`)
```go
type Book struct {
    Title string
    Link  string
}
```

### Web Scraping (`app/crawler.go`)
- Uses Colly framework
- Targets `.fiction-list-item` CSS selector
- Extracts `.fiction-title` and link attributes
- Automatically saves to MongoDB after scraping
- Robust error handling

### Database Layer (`app/database.go`)
- MongoDB connection management
- BSON serialization
- Test environment detection
- Connection cleanup and resource management
- Collections: `royalRoadBooks.books`

### Frontend (`app/main_page.go` + `app/templates/`)
- **Embedded filesystem** using Go's `//go:embed` directive
- **Modular templates**: `main.html` for full page, `book_list.html` for partials
- **Theme support**: CSS custom properties for light/dark modes
- HTMX integration for dynamic updates
- Search indicators and loading states
- **Template separation**: HTML moved from inline strings to external files

## Testing Strategy

### Test Coverage
- **Unit Tests**: All major functions tested
- **Integration Tests**: Database operations with testcontainers
- **HTTP Tests**: Handler testing with httptest
- **Template Tests**: HTML rendering verification

### Test Files
- `crawler_test.go` - Web scraping with mock servers
- `database_test.go` - MongoDB operations
- `main_page_test.go` - Template rendering
- `main_test.go` - HTTP handlers

### Testing Infrastructure
- Testcontainers for real MongoDB instances
- HTTP test servers for scraping tests
- Comprehensive assertions with testify
- Automatic cleanup in test teardown

## AI Assistant Guidelines

### When Working on This Project

1. **File Organization**: All application code is in `/app` directory
2. **Single Package**: Everything uses `package main`
3. **Testing**: Always run tests when making changes
4. **Docker**: Prefer containerized development when possible
5. **Dependencies**: Check `go.mod` before adding new packages

### Common Tasks

**Adding New Features:**
- Add handler to `main.go`
- Update templates in `main_page.go` if needed
- Add corresponding tests
- Update justfile if new commands needed

**Database Changes:**
- Modify structures in `model.go`
- Update database operations in `database.go`
- Add migration logic if needed
- Update tests accordingly

**Frontend Updates:**
- Modify templates in `app/templates/` directory
- Update `main_page.go` for template loading changes
- Maintain HTMX compatibility
- Keep responsive design principles
- Test theme switching functionality
- Test across different screen sizes

### Development Best Practices

1. **Always use the justfile** for common tasks
2. **Run tests locally first** before Docker deployment
3. **Check MongoDB Express** at http://localhost:8081 for data inspection
4. **Use testcontainers** for integration tests requiring database
5. **Follow Go conventions** for naming and structure

### Debugging Tips

**Application Issues:**
- Check logs with `just logs`
- Verify MongoDB connection
- Test scraping endpoints manually

**Database Issues:**
- Use Mongo Express web interface
- Check connection strings
- Verify container networking

**Frontend Issues:**
- Check browser developer tools
- Verify HTMX requests in network tab
- Test template rendering separately

## External Dependencies

### Go Modules (Key Dependencies)
```
github.com/gocolly/colly/v2 v2.1.0         # Web scraping
go.mongodb.org/mongo-driver v1.17.3        # MongoDB client
github.com/stretchr/testify v1.10.0        # Testing framework
github.com/testcontainers/testcontainers-go # Integration testing
```

### Runtime Dependencies
- MongoDB 6.0+ (via Docker)
- Go 1.23+ (currently using 1.24)
- Docker and Docker Compose

## Security Considerations

1. **No hardcoded credentials** - Uses environment variables
2. **Input validation** - Search queries are sanitized
3. **Rate limiting** - Respectful scraping behavior
4. **Error handling** - No sensitive information leakage
5. **Container security** - Multi-stage builds minimize attack surface

## Performance Characteristics

- **Memory usage**: Minimal (single binary + cached books)
- **Response time**: Sub-100ms for cached data
- **Scraping time**: 1-3 seconds depending on network
- **Database ops**: Millisecond-level for small datasets
- **Concurrent users**: Limited by Go's goroutine model

## Future Enhancement Areas

1. **API Development**: RESTful endpoints for programmatic access
2. **Caching Strategy**: Redis for distributed caching
3. **Monitoring**: Prometheus metrics and health checks
4. **Authentication**: User accounts and personalization
5. **Book Details**: Enhanced metadata scraping
6. **Pagination**: Support for larger datasets
7. **Rate Limiting**: Sophisticated request throttling
8. **Enhanced Theming**: Custom color schemes, high contrast mode
9. **Accessibility**: Keyboard shortcuts, screen reader improvements

## Troubleshooting Common Issues

### Build Failures
- Check Go version compatibility
- Verify all dependencies in go.mod
- Clear module cache: `go clean -modcache`

### Container Issues
- Port conflicts: Check if 8090/27017/8081 are available
- Memory issues: Ensure Docker has sufficient resources
- Network issues: Verify docker-compose networking

### Scraping Failures
- RoyalRoad structure changes: Update CSS selectors
- Network timeouts: Check connectivity
- Rate limiting: Implement delays between requests

## Recent Major Changes

### Template Refactoring & Theme Support
- **Template Extraction**: HTML templates moved from inline strings to separate files
- **Embedded Filesystem**: Templates compiled into binary using `//go:embed`
- **Theme System**: Dark/Light mode toggle with CSS custom properties
- **Persistence**: Theme preference saved in localStorage
- **Modularity**: Separation of main page and partial templates for better maintainability

### Technical Implementation Details
- **File Structure**: Templates now in `app/templates/` directory
- **Loading Mechanism**: `template.ParseFS()` instead of `template.Parse()`
- **Theme Variables**: CSS custom properties for consistent theming
- **JavaScript Integration**: Minimal JS for theme switching functionality
- **Backward Compatibility**: All existing functionality preserved

This guide provides comprehensive information for AI assistants working on the RoyalRoadBot project. Always refer to the README.md for user-facing documentation and this file for development guidance.