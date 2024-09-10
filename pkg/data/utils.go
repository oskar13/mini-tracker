package data

import (
	"fmt"
	"log"
	"os"
	"strings"
)

// Reads out password string from a file specified, file must only contain that password and nothing else, all whitespace will be trimmed
func ReadPassword(filePath string) string {
	data, err := os.ReadFile(filePath)
	fmt.Println("Reading:")
	fmt.Println(filePath)
	if err != nil {
		log.Fatalf("Failed to read password file: %v", err)
	}
	fmt.Println(string(data))

	var dataStr = string(data)
	return strings.TrimSpace(dataStr)
}
