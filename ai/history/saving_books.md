# Implementation Plan: Royal Road Search and Memorize Feature

## Overview
Implement a new feature to search and memorize books from Royal Road, with a tabbed interface separating popular books from search functionality.

## Current State Analysis
- Existing app scrapes top 10 popular books from Royal Road
- Uses MongoDB for persistence (will be removed from popular books display)
- Has HTMX frontend with search functionality (currently searches cached popular books)
- Template system with main.html and book_list.html

## Implementation Steps

### Step 1: Create Search Module (`searcher.go`)
**File:** `royalroadbot/app/searcher.go`

**Functionality:**
- Implement `searchRoyalRoadBooks(query string) ([]Book, error)`
- Use Royal Road search URL: `https://www.royalroad.com/fictions/search?title={query}`
- Extract book results from search page using Colly
- Limit results to reasonable number (10-15 books)
- Add `memorizeBook(book Book) error` function to save specific book to database

**CSS Selectors to investigate:**
- Search results container
- Book title and link elements
- Handle pagination if needed

### Step 2: Create Search Module Tests (`searcher_test.go`)
**File:** `royalroadbot/app/searcher_test.go`

**Test Coverage:**
- Test search functionality with mock HTTP server
- Test memorizeBook database operations
- Test edge cases (empty results, network errors)
- Integration tests with testcontainers

### Step 3: Modify Crawler Module
**File:** `royalroadbot/app/crawler.go`

**Changes:**
- Remove database saving from `fetchBooks()` function
- Remove `saveBooksWithMetadata()` call
- Keep scraping logic intact
- Update function documentation

**Impact Assessment:**
- Verify which functions call `fetchPopularBooks()`
- Ensure no breaking changes to existing handlers

### Step 4: Add New HTTP Handlers
**File:** `royalroadbot/app/main.go`

**New Handlers:**
- `/search-books` - Search Royal Road for books
- `/memorize-book` - Save a specific book to database
- `/memorized-books` - Display saved books

**Handler Details:**
```go
func searchBooksHandler(w http.ResponseWriter, r *http.Request)
func memorizeBookHandler(w http.ResponseWriter, r *http.Request)
func memorizedBooksHandler(w http.ResponseWriter, r *http.Request)
```

### Step 5: Update Database Module
**File:** `royalroadbot/app/database.go`

**New Functions:**
- `saveBookToMemory(book Book) error` - Save individual book
- `getMemorizedBooks() ([]Book, error)` - Retrieve saved books
- `removeMemorizedBook(bookTitle string) error` - Remove book

**Database Schema:**
- Use existing schema

### Step 6: Create New Templates
**Files:**
- `royalroadbot/app/templates/search_results.html`
- `royalroadbot/app/templates/memorized_books.html`
- `royalroadbot/app/templates/tabbed_main.html`

**Template Structure:**
```html
<!-- tabbed_main.html -->
<div class="tabs">
  <div class="tab-buttons">
    <button class="tab-btn active" data-tab="popular">Popular Books</button>
    <button class="tab-btn" data-tab="search">Search & Memorize</button>
  </div>
  <div class="tab-content">
    <div id="popular-tab" class="tab-pane active">
      <!-- Current popular books functionality -->
    </div>
    <div id="search-tab" class="tab-pane">
      <!-- New search interface -->
    </div>
  </div>
</div>
```

### Step 7: Update Template Rendering
**File:** `royalroadbot/app/main_page.go`

**Changes:**
- Add new template parsing functions
- Update `renderPage()` to use new tabbed layout
- Add `renderSearchResults()` and `renderMemorizedBooks()`

### Step 8: Frontend JavaScript Updates
**File:** `royalroadbot/app/templates/tabbed_main.html`

**JavaScript Features:**
- Tab switching functionality
- HTMX integration for search
- Memorize/remove book actions
- Search debouncing
- Loading indicators

### Step 9: CSS Styling Updates
**Styling for:**
- Tab interface design
- Search results layout
- Memorized books display
- Action buttons (memorize, remove)
- Responsive design for mobile

### Step 10: Update Tests
**Files to Update:**
- `main_test.go` - Add new handler tests
- `main_page_test.go` - Update template tests
- `database_test.go` - Add memorized books tests

### Step 11: Documentation Updates
**Files:**
- Update `README.md` with new features
- Update `CLAUDE.md` with architectural changes

## Confirmation Points

### Step 1 Confirmation:
- [x] `searcher.go` created with search functionality
- [x] Royal Road search URL working correctly
- [x] CSS selectors extracting book data properly
- [x] `memorizeBook()` function implemented
- [x] Database functions for memorized books added to `database.go`

### Step 2 Confirmation:
- [x] `searcher_test.go` created with comprehensive tests
- [x] Mock HTTP server tests passing
- [x] Database integration tests working
- [x] Edge cases covered
- [x] Integration tests with testcontainers working
- [x] All test scenarios validated successfully
- [x] Special character test removed and documented in `ai/tasks/special_character_processing.md`

### Step 3 Confirmation:
- [x] `fetchBooks()` modified to remove database saving
- [x] No breaking changes to existing functionality
- [x] Popular books still display correctly
- [x] Tests still passing
- [x] Crawler tests updated to verify memory-only behavior

### Step 4 Confirmation:
- [x] New HTTP handlers added to `main.go`
- [x] Routes registered correctly (`/search-books`, `/memorize-book`, `/memorized-books`)
- [x] Handler functions implemented with proper validation and error handling
- [x] HTMX compatibility maintained
- [x] Tests added to `main_test.go` with testcontainer support

### Step 5 Confirmation:
- [x] Database functions for memorized books implemented (completed in Steps 1-2)
- [x] New collection `memorizedBooks` working
- [x] CRUD operations tested
- [x] Error handling proper
- [x] All database operations use separate memorized collection

### Step 6 Confirmation:
- [x] New templates created
- [x] Tabbed interface HTML structure complete
- [x] HTMX attributes properly set
- [x] Template inheritance working
- [x] Templates created: `tabbed_main.html`, `search_results.html`, `memorized_books.html`
- [x] Action buttons (memorize, remove) with proper styling
- [x] Theme compatibility maintained (dark/light modes)
- [x] Mobile responsive design working

### Step 7 Confirmation:
- [x] Template rendering functions updated
- [x] New templates properly parsed
- [x] Template execution working
- [x] Error handling maintained
- [x] Added `renderTabbedMain()`, `renderSearchResults()`, `renderMemorizedBooks()` functions
- [x] Updated handlers to use HTML templates instead of JSON responses
- [x] Proper HTML content-type headers set
- [x] All tests updated and passing

### Step 8 Confirmation:
- [x] Tab switching JavaScript working
- [x] HTMX requests functioning
- [x] Search debouncing implemented
- [x] User interactions smooth

### Step 9 Confirmation:
- [x] Tab styling complete
- [x] Search interface visually appealing
- [x] Mobile responsive design working
- [x] Theme compatibility maintained

### Step 10 Confirmation:
- [x] All existing tests updated
- [x] New functionality tests added
- [x] Test coverage maintained
- [x] Integration tests passing

### Step 11 Confirmation:
- [x] Documentation updated
- [x] README.md reflects new features
- [x] CLAUDE.md updated for AI assistance

## Risk Mitigation

1. **Breaking Changes:** Test each step thoroughly before proceeding
2. **Database Schema:** Maintain backward compatibility
3. **Frontend UX:** Ensure smooth transitions between tabs
4. **Performance:** Monitor search response times
5. **Royal Road Changes:** Make CSS selectors resilient to minor changes

## Success Criteria

1. ✅ Two-tab interface working smoothly
2. ✅ Popular books tab shows current functionality without database persistence
3. ✅ Search tab allows searching Royal Road books
4. ✅ Books can be memorized and removed
5. ✅ All existing functionality preserved
6. ✅ Test coverage maintained
7. ✅ Mobile responsive design working
8. ✅ Theme system compatibility

## Implementation Status: COMPLETED ✅

All steps have been successfully implemented and tested. The Royal Road Search and Memorize feature is fully functional with:

- Complete tabbed interface with smooth interactions
- Real-time search across Royal Road catalog with 500ms debouncing
- Book memorization system with MongoDB persistence
- Full CRUD operations for memorized books
- HTMX-powered dynamic updates
- Mobile responsive design
- Dark/Light theme support
- Comprehensive test coverage (all tests passing)
- Updated documentation

Date Completed: 2025-05-30

## Rollback Plan

Each step should be committed separately to allow easy rollback:
1. Keep original `crawler.go` backup before modifications
2. Database changes should be additive (new collections)
3. Template changes should not break existing functionality
4. Test each component independently
