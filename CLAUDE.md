# CLAUDE.md - AI Assistant Guide for RoyalRoadBot

## Project Overview

**RoyalRoadBot** is a Go web application that scrapes popular books from RoyalRoad.com and provides comprehensive search functionality with book memorization features. It features a modern tabbed web interface with real-time search capabilities across the entire RoyalRoad catalog and persistent book collection management through MongoDB.

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
│   ├── crawler.go         # Web scraping logic for popular books (Colly)
│   ├── searcher.go        # Royal Road search and book memorization
│   ├── database.go        # MongoDB operations
│   ├── main_page.go       # Template rendering with embedded filesystem
│   ├── templates/         # HTML templates directory
│   │   ├── main.html     # Legacy main page template
│   │   ├── tabbed_main.html # New tabbed interface with full functionality
│   │   ├── book_list.html # Partial template for popular books HTMX updates
│   │   ├── search_results.html # Partial template for search results
│   │   └── memorized_books.html # Partial template for memorized books
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

### Popular Books Tab
- Scrapes RoyalRoad.com's "active-popular" fiction list
- Extracts book titles and links using CSS selectors
- Limits results to top 10 books
- No database persistence (memory-only for popular books)
- Client-side search filtering
- Direct memorization from popular list

### Search & Memorize Tab
- **Real-time search** across RoyalRoad's entire catalog
- **Book memorization** - Save favorite books to personal collection
- **Memorized books management** - View and remove saved books
- **Search debouncing** with 500ms delay for optimal performance
- **HTMX-powered interactions** for smooth user experience

### Web Interface
- **Tabbed interface** - Clean separation of popular books and search functionality
- **Theme system** - Dark/Light mode toggle with persistent user preference
- **Responsive design** for all devices and screen sizes
- **HTMX-powered partial page updates** without full page reloads
- **CSS Custom Properties** for clean theming system
- **Template modularity** with separate files for different interface sections
- **Action buttons** for memorizing and removing books
- **Loading indicators** and status messages for user feedback

### Data Management
- **MongoDB persistence** for memorized books with proper metadata
- **Separate collections** - Popular books (memory-only) vs memorized books (persistent)
- **Connection management** with automatic cleanup and pooling
- **CRUD operations** for memorized book collection
- **Duplicate prevention** for memorized books

## Development Workflow

### Quick Start Commands (using Just)

```bash
# Local development with MongoDB in Docker
just run-dev-mongo     # Start MongoDB + Mongo Express
just run-dev-local     # Build and run app locally

# Full containerized deployment
just rebuild-all      # Build all containers
just run              # Start full stack

# Testing
just test-local       # Run tests locally

# Utilities
just clean            # Clean up containers/images
just logs             # View container logs
```

### Docker Configurations

1. **docker-compose.yaml** - Main application stack
2. **docker-compose-mongo.yaml** - MongoDB only
3. **docker-compose-dev.yaml** - Adds Mongo Express web client

### Environment Variables

- `MONGODB_URI` - MongoDB connection string
  - Local dev: `mongodb://admin:password@127.0.0.1:27017`
  - Docker: Handled by docker-compose networking

## Code Structure Analysis

### Main Application (`app/main.go`)
- HTTP server on port 8090
- Seven main routes:
  - `/` - Main tabbed interface page
  - `/search` - Legacy HTMX search endpoint for popular books
  - `/refresh` - Manual refresh of popular books
  - `/search-books` - Royal Road search endpoint
  - `/memorize-book` - Save book to memorized collection
  - `/memorized-books` - Get memorized books list
  - `/remove-memorized-book` - Remove book from memorized collection
- Global caching with `cachedBooks` variable for popular books
- Comprehensive error handling with proper HTTP status codes
- HTMX-compatible HTML responses for dynamic updates

### Data Model (`app/model.go`)
```go
type Book struct {
    Title string
    Link  string
}
```

### Popular Books Scraping (`app/crawler.go`)
- Uses Colly framework for RoyalRoad popular books
- Targets `.fiction-list-item` CSS selector
- Extracts `.fiction-title` and link attributes
- **Memory-only storage** (no database persistence for popular books)
- Robust error handling with graceful fallbacks

### Search Functionality (`app/searcher.go`)
- **Royal Road search integration** using their search API
- **Book memorization** with duplicate prevention
- **Memorized book retrieval** and management
- **Book removal** from memorized collection
- Comprehensive error handling and logging
- Input validation and sanitization

### Database Layer (`app/database.go`)
- MongoDB connection management with automatic detection
- BSON serialization for Book structures
- Test environment detection for proper cleanup
- Connection pooling and resource management
- **Dual collection strategy**:
  - `royalRoadBooks.books` - Legacy collection (not used in current implementation)
  - `royalRoadBooks.memorizedBooks` - Active collection for user's memorized books
- **CRUD operations** for memorized books with proper error handling
- **Duplicate prevention** and existence checking

### Frontend (`app/main_page.go` + `app/templates/`)
- **Embedded filesystem** using Go's `//go:embed` directive
- **Comprehensive template system**:
  - `tabbed_main.html` - Main tabbed interface with full functionality
  - `search_results.html` - Partial for Royal Road search results
  - `memorized_books.html` - Partial for memorized books display
  - `book_list.html` - Partial for popular books HTMX updates
  - `main.html` - Legacy template (maintained for compatibility)
- **Advanced theme support**: CSS custom properties for light/dark modes
- **HTMX integration** for seamless dynamic updates
- **Interactive elements**: Tab switching, search debouncing, action buttons
- **Status indicators**: Loading states, success/error messages
- **Mobile responsiveness**: Optimized layouts for all screen sizes
- **Template separation**: Clean separation of concerns with modular design

## Testing Strategy

### Test Coverage
- **Unit Tests**: All major functions comprehensively tested
- **Integration Tests**: Database operations with testcontainers
- **HTTP Handler Tests**: Complete handler testing with httptest
- **Template Tests**: HTML rendering verification for all templates
- **Search Tests**: Royal Road search functionality with mock servers
- **Memorization Tests**: Book saving, retrieval, and removal operations

### Test Files
- `crawler_test.go` - Popular books scraping with mock servers
- `searcher_test.go` - Search functionality and book memorization
- `database_test.go` - MongoDB operations including memorized books
- `main_page_test.go` - Template rendering for all templates
- `main_test.go` - HTTP handlers including new search and memorization endpoints

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
- Add handler to `main.go` (follow existing pattern)
- Update appropriate templates in `app/templates/` directory
- Add corresponding tests in relevant `*_test.go` files
- Update justfile if new commands needed
- Maintain HTMX compatibility for dynamic interactions

**Database Changes:**
- Modify structures in `model.go` if needed
- Update database operations in `database.go` (consider separate collections)
- Add migration logic if schema changes required
- Update tests accordingly (use testcontainers for integration tests)
- Follow existing memorized books pattern for new collections

**Frontend Updates:**
- Modify templates in `app/templates/` directory (use embedded filesystem)
- Update `main_page.go` for new template functions
- Maintain HTMX compatibility for partial updates
- Follow responsive design principles
- Test theme switching functionality across all new interfaces
- Test across different screen sizes and devices
- Consider tab structure when adding new functionality

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

### Search & Memorize Feature Implementation
- **Tabbed Interface**: Complete redesign with separate Popular Books and Search & Memorize tabs
- **Royal Road Search Integration**: Real-time search across entire RoyalRoad catalog
- **Book Memorization System**: Save and manage favorite books with MongoDB persistence
- **HTMX-Powered Interactions**: Smooth, dynamic user experience without page reloads
- **Advanced Template System**: Modular templates for different interface sections
- **Comprehensive Testing**: Full test coverage including new handlers, templates, and database operations

### Template Refactoring & Theme Support
- **Template Extraction**: HTML templates moved from inline strings to separate files
- **Embedded Filesystem**: Templates compiled into binary using `//go:embed`
- **Theme System**: Dark/Light mode toggle with CSS custom properties
- **Persistence**: Theme preference saved in localStorage
- **Enhanced Modularity**: Multiple template files for different UI sections

### Technical Implementation Details
- **Architecture**: Dual-tab system with separate functionality per tab
- **File Structure**: Expanded templates in `app/templates/` directory
- **Database Strategy**: Separate collections for popular vs memorized books
- **Loading Mechanism**: `template.ParseFS()` with multiple template functions
- **JavaScript Enhancement**: Tab switching, search debouncing, HTMX event handling
- **Mobile Responsiveness**: Optimized layouts for all screen sizes
- **Backward Compatibility**: All existing functionality preserved and enhanced

This guide provides comprehensive information for AI assistants working on the RoyalRoadBot project. Always refer to the README.md for user-facing documentation and this file for development guidance.
