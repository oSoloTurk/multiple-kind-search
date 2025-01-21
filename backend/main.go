package main

import (
	"log"
	"os"
)

func main() {
	// ... existing code ...

	// ... after initializing repo
	if err := repo.CreateIndices(); err != nil {
		log.Printf("Error creating indices: %v", err)
	}
} 