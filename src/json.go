package jsonS3

import (
	"encoding/json"
	"log"
	"os"
)

func CreateJSON[T any](data []T) []byte {
	result, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Fatalf("failed to create JSON: %v\n", err)
	}
	return result
}

func SaveLocal(filepath string, data []byte) error {
	file, err := os.Create(filepath)
	if err != nil {
		log.Fatalf("failed to create local file: %v", err)
		return err
	}
	defer file.Close()

	_, err = file.Write(data)

	if err != nil {
		log.Printf("failed to save locally: %v\n", err)
	} else {
		log.Printf("💾 Saved to local: %s\n", filepath)
	}
	return err
}
