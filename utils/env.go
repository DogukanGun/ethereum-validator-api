package utils

import (
	"github.com/joho/godotenv"
	"log"
)

// InitializeENV loads environment variables from the specified .env file.
// It returns true if the file was successfully loaded, false otherwise.
//
// Parameters:
//   - envFileName: The path to the .env file to be loaded
//
// Returns:
//   - bool: true if environment variables were loaded successfully, false if there was an error
func InitializeENV(envFileName string) bool {
	// First check if env file exists & load the file if it exists
	dotenvErr := godotenv.Load(envFileName)

	if dotenvErr == nil {
		log.Println(".env file located and loaded.")
		return true
	} else {
		log.Println("Failed to load .env file: ", dotenvErr)
		return false
	}
}
