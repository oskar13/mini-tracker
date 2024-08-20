package data

import (
	"fmt"
	"log"
	"os"
	"strings"
)

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
