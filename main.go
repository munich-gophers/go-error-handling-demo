// error-handling-demo/main.go
package main

import (
	"database/sql" // To simulate a common error type like sql.ErrNoRows
	"errors"
	"fmt"
	"os"
)

// --- Sentinel Errors ---
var ErrResourceNotFound = errors.New("resource not found")

// --- Custom Error Type ---
type ConfigError struct {
	FileName string
	Op       string
	Err      error // Underlying error
}

func (ce *ConfigError) Error() string {
	return fmt.Sprintf("config error during '%s' for file '%s': %v", ce.Op, ce.FileName, ce.Err)
}

// Unwrap allows errors.Is and errors.As to work on the wrapped error.
func (ce *ConfigError) Unwrap() error {
	return ce.Err
}

// --- Simulated Functions that Can Fail ---

// Simulates trying to load a configuration file.
func loadAppConfig(filePath string) error {
	// Simulate a low-level file system error
	if filePath == "missing.json" {
		// Pretend os.Open failed
		originalErr := os.ErrNotExist // A common error from the os package
		return &ConfigError{
			FileName: filePath,
			Op:       "open file",
			Err:      fmt.Errorf("failed to open: %w", originalErr), // Wrap the OS error
		}
	}
	// Simulate another type of configuration issue
	if filePath == "invalid.json" {
		originalErr := errors.New("invalid JSON structure")
		return &ConfigError{
			FileName: filePath,
			Op:       "parse content",
			Err:      originalErr, // Wrap a generic error
		}
	}
	fmt.Printf("Successfully loaded config: %s\n", filePath)
	return nil
}

// Simulates fetching data from a database.
func fetchData(queryID string) (string, error) {
	if queryID == "123" {
		return "Sample Data", nil
	}
	if queryID == "notfound_db" {
		// Simulate a database error where no rows are returned
		return "", fmt.Errorf("database query for ID '%s' failed: %w", queryID, sql.ErrNoRows)
	}
	if queryID == "custom_resource_err" {
		return "", fmt.Errorf("could not retrieve resource data: %w", ErrResourceNotFound)
	}
	return "", fmt.Errorf("unknown query ID: %s", queryID)
}

// Higher-level function that calls other functions.
func processRequest(configPath string, dataQueryID string) error {
	fmt.Printf("\n--- Processing request with config '%s' and data query '%s' ---\n", configPath, dataQueryID)

	err := loadAppConfig(configPath)
	if err != nil {
		// Handle error from loadAppConfig
		var ce *ConfigError
		if errors.As(err, &ce) {
			fmt.Printf("Detailed Config Error: Operation '%s' on file '%s' failed.\n", ce.Op, ce.FileName)
			// We can also check the underlying error if needed
			if errors.Is(ce.Unwrap(), os.ErrNotExist) {
				fmt.Println("  Underlying cause: File does not exist. Attempting fallback...")
				// return fallback() // Or some other recovery mechanism
			}
		} else {
			fmt.Println("An unexpected error occurred while loading config.")
		}
		// Regardless of type, we return the wrapped error to preserve full context
		return fmt.Errorf("failed during config stage: %w", err)
	}

	data, err := fetchData(dataQueryID)
	if err != nil {
		// Handle error from fetchData
		if errors.Is(err, sql.ErrNoRows) {
			fmt.Println("Data Fetch Error: The requested data was not found in the database (sql.ErrNoRows).")
			// Potentially return a user-friendly error or default data
		} else if errors.Is(err, ErrResourceNotFound) {
			fmt.Println("Data Fetch Error: A specific custom resource was not found (ErrResourceNotFound).")
		} else {
			fmt.Println("An unknown error occurred while fetching data.")
		}
		return fmt.Errorf("failed during data fetching stage: %w", err)
	}

	fmt.Printf("Successfully processed request. Data: %s\n", data)
	return nil
}

func main() {
	// Scenario 1: Config file missing
	err := processRequest("missing.json", "123")
	if err != nil {
		fmt.Printf("Main Error Handler: %v\n", err) // Full wrapped error chain
	}

	// Scenario 2: Invalid config content
	err = processRequest("invalid.json", "123")
	if err != nil {
		fmt.Printf("Main Error Handler: %v\n", err)
	}

	// Scenario 3: Data not found in DB (sql.ErrNoRows)
	err = processRequest("valid.json", "notfound_db")
	if err != nil {
		fmt.Printf("Main Error Handler: %v\n", err)
	}

	// Scenario 4: Custom resource not found (ErrResourceNotFound)
	err = processRequest("valid.json", "custom_resource_err")
	if err != nil {
		fmt.Printf("Main Error Handler: %v\n", err)
	}

	// Scenario 5: Success
	err = processRequest("valid.json", "123")
	if err != nil {
		fmt.Printf("Main Error Handler: %v\n", err)
	}
}
