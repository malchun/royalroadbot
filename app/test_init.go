package main

import (
	"fmt"
)

// setupTestMocks configures all the mock functions for tests
func setupTestMocks() {
	// Override saveBooksWithMetadata to do nothing in tests
	saveBooksWithMetadataFunc = func(books []Book) error {
		// Don't actually save to database in tests
		return nil
	}

	// Override database connection function
	ConnectDBFunc = func() {
		// Do nothing - don't try to connect to MongoDB in tests
		// Just print a message for debugging
		fmt.Println("Test mode: Not connecting to MongoDB")
	}
}