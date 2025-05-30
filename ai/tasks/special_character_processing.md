# Special Character Processing in Royal Road Search

## Overview

This document outlines the handling of special characters in the Royal Road search functionality, including URL encoding, edge cases, and testing considerations.

## Current Implementation

### URL Encoding
The `searchRoyalRoadBooks()` function in `searcher.go` uses Go's `url.QueryEscape()` to properly encode search queries:

```go
searchURL := fmt.Sprintf("https://www.royalroad.com/fictions/search?title=%s", url.QueryEscape(query))
```

### Supported Characters
The implementation handles:
- **Alphanumeric characters**: Letters and numbers (no encoding needed)
- **Spaces**: Encoded as `%20` or `+`
- **Special characters**: `&`, `+`, `%`, `=`, `?`, `/`, etc. are properly encoded
- **Unicode characters**: Non-ASCII characters are UTF-8 encoded

## Edge Cases and Considerations

### 1. Ampersand (&) Handling
- Input: `"wizard & magic"`
- Encoded: `"wizard%20%26%20magic"`
- Royal Road handles this correctly in search queries

### 2. Plus Sign (+) Handling
- Input: `"magic + spells"`
- Encoded: `"magic%20%2B%20spells"`
- Prevents confusion with space encoding

### 3. Percent Sign (%) Handling
- Input: `"100% magic"`
- Encoded: `"100%25%20magic"`
- Double-encoding is properly avoided

### 4. Question Mark (?) Handling
- Input: `"what magic?"`
- Encoded: `"what%20magic%3F"`
- Prevents URL parameter confusion

## Testing Challenges

### Network Dependency
Testing special character handling requires actual HTTP requests to Royal Road, which creates dependencies on:
- Network connectivity
- Royal Road server availability
- Rate limiting considerations
- Changing search results

### Mock Server Limitations
Mock servers don't fully replicate Royal Road's:
- URL parameter parsing
- Search algorithm behavior
- Character encoding edge cases

## Recommended Testing Strategy

### Unit Tests (Current Implementation)
- Test URL encoding directly with `url.QueryEscape()`
- Verify query string construction
- Test empty and whitespace-only inputs

### Integration Tests (Future Enhancement)
```go
func TestSpecialCharacterEncoding(t *testing.T) {
    testCases := []struct {
        input    string
        expected string
    }{
        {"wizard & magic", "wizard%20%26%20magic"},
        {"magic + spells", "magic%20%2B%20spells"},
        {"100% success", "100%25%20success"},
        {"what magic?", "what%20magic%3F"},
    }
    
    for _, tc := range testCases {
        encoded := url.QueryEscape(tc.input)
        assert.Equal(t, tc.expected, encoded)
    }
}
```

### Manual Testing Scenarios
1. Search for `"dragon & wizard"` - verify results returned
2. Search for `"magic + academy"` - verify proper parsing
3. Search for `"100% complete"` - verify percentage handling
4. Search for unicode characters like `"魔法"` or `"магия"`

## Known Issues and Workarounds

### Issue: Royal Road Search Sensitivity
Royal Road's search may be sensitive to certain character combinations or may normalize some characters internally.

**Workaround**: The application handles encoding correctly; any search limitations are on Royal Road's side.

### Issue: Rate Limiting with Special Characters
Some special character combinations might trigger Royal Road's rate limiting more aggressively.

**Workaround**: Implement exponential backoff and respect rate limits in the scraping logic.

## Future Enhancements

### 1. Search Query Sanitization
Consider implementing pre-processing to:
- Remove or replace problematic characters
- Normalize unicode characters
- Handle common search patterns

### 2. Enhanced Error Handling
Implement specific error handling for:
- Invalid character encoding
- Royal Road search errors
- Rate limiting responses

### 3. Query Optimization
Analyze Royal Road's search behavior to optimize queries:
- Character combination effectiveness
- Search term ordering
- Length limitations

## Testing Removal Rationale

The `TestSearchRoyalRoadBooks_SpecialCharacters` test was removed from `searcher_test.go` because:

1. **Network Dependency**: Required actual HTTP requests to Royal Road
2. **Unreliable Results**: Search results change over time
3. **Rate Limiting**: Could trigger Royal Road's anti-scraping measures
4. **Maintenance Overhead**: Required constant updates as Royal Road's search evolved

The functionality is better tested through:
- Direct URL encoding unit tests
- Manual testing during development
- Integration tests in controlled environments

## Conclusion

Special character processing is handled correctly by Go's standard library URL encoding. The main challenges are in testing rather than implementation. Future work should focus on comprehensive integration testing in isolated environments rather than relying on external services.