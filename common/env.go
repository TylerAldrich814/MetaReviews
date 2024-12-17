package common

import (
	"log"
	"syscall"
)

// Helper function for querying a packages local .env file.
// If the enviornment variable is not found, the fallback is
// returned instead.
//
// - In order for this helper function to work, you must include
//   '_ "github.com/joho/godotenv/autoload"' in the files imports.
func EnvString(key, fallback string) string {
  if val, ok := syscall.Getenv(key); ok {
    return val
  }
  log.Printf("Failed to find env \"%s\"", key)
  return fallback
}
