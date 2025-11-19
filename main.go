package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("Novel Notes", "Version 1.0")

	// 1. Load library from file (or get an empty one on first run)
	lib, err := LoadLibrary("novel_notes.json")
	if err != nil {
		fmt.Println("Error loading library:", err)
		return
	}

	// 2. If no bookcases yet, build some sample data (first run only)
	if len(lib.Bookcases) == 0 {
		fmt.Println("No bookcases found, creating sample data...")
		initializeSampleData(&lib)
	}

	// 3. Auto-rollover once per day
	AutoRollover(&lib)

	// 4. Read command-line arguments (skip program name) and handle commands
	args := os.Args[1:]
	handleCommand(&lib, args)

	// 5. Save back to file (even if nothing changed, it's fine)
	if err := SaveLibrary(lib, "novel_notes.json"); err != nil {
		fmt.Println("Error saving library:", err)
	} else {
		fmt.Println("Library saved to novel_notes.json")
	}
}
