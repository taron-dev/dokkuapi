package postgres

import (
	"fmt"
	"os"
)

// ServiceExists return is service already exists
func ServiceExists(serviceName string) bool {
	postgresRoot := os.Getenv("POSTGRES_ROOT")
	servicePath := fmt.Sprintf("%v/%v", postgresRoot, serviceName)

	return directoryExists(servicePath)
}

// directoryExists returns if a path exists and is a directory
func directoryExists(filePath string) bool {
	fi, err := os.Stat(filePath)
	if err != nil {
		return false
	}

	return fi.IsDir()
}
