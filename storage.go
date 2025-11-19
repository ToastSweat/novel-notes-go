package main

import (
	"encoding/json"
	"os"
)

func LoadLibrary(filename string) (Library, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		// If the file doesn't exist, return an empty Library (first run)
		if os.IsNotExist(err) {
			empty := Library{
				Bookcases:    []Bookcase{},
				TotalScore:   0,
				LastRollOver: "",
			}
			return empty, nil
		}
		// Some other error (permissions, etc.)
		return Library{}, err
	}

	var lib Library
	if err := json.Unmarshal(data, &lib); err != nil {
		return Library{}, err
	}

	return lib, nil
}

func SaveLibrary(lib Library, filename string) error {
	data, err := json.MarshalIndent(lib, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		return err
	}

	return nil
}
