package main

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/tormoder/fit"
)

func main() {
	// Define command line flags
	filesFlag := flag.String("files", "", "Space-delimited set of filenames or a directory containing .fit files")
	outputFlag := flag.String("output", ".", "Directory to output files to, defaults to current directory")

	flag.Parse()

	if *filesFlag == "" {
		log.Fatal("The -files flag is required")
	}

	// Get a list of .fit files
	fileList := getFitFiles(*filesFlag)

	// Ensure the output directory exists
	err := os.MkdirAll(*outputFlag, os.ModePerm)
	if err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	// Process each .fit file
	for _, file := range fileList {
		processFitFile(file, *outputFlag)
	}
}

// getFitFiles returns a list of .fit files based on the provided input (filesFlag)
func getFitFiles(filesFlag string) []string {
	var files []string

	// Check if the input is a directory
	if info, err := os.Stat(filesFlag); err == nil && info.IsDir() {
		files, err = filepath.Glob(filepath.Join(filesFlag, "*.[fF][iI][tT]"))
		if err != nil {
			log.Fatalf("Failed to list .fit files: %v", err)
		}
	} else {
		// Split the space-delimited string into file paths
		files = strings.Split(filesFlag, " ")
	}

	return files
}

// processFitFile processes the .fit file to remove GPS info and save it with a random name
func processFitFile(filePath, outputDir string) {
	// Open the .fit file
	fitFile, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Failed to open file %s: %v", filePath, err)
	}
	defer fitFile.Close()

	// Decode the .fit file
	fitData, err := fit.Decode(fitFile)
	if err != nil {
		log.Fatalf("Failed to decode .fit file %s: %v", filePath, err)
	}

	// Get the actual activity
	// Remove GPS data from the file
	removeGPSData(fitData)

	// Generate a random filename
	newFileName := generateRandomFilename() + ".fit"
	newFilePath := filepath.Join(outputDir, newFileName)

	// Save the modified .fit file
	newFile, err := os.Create(newFilePath)
	if err != nil {
		log.Fatalf("Failed to create output file %s: %v", newFilePath, err)
	}
	defer newFile.Close()

	err = fit.Encode(newFile, fitData, binary.LittleEndian)
	if err != nil {
		log.Fatalf("Failed to encode .fit file %s: %v", newFilePath, err)
	}

	fmt.Printf("Processed %s and saved as %s\n", filePath, newFilePath)
}

// removeGPSData strips the GPS information from the .fit data
func removeGPSData(fitData *fit.File) {
	activity, err := fitData.Activity()
	if err != nil {
		log.Fatalf("Failed to decode activity %v", err)
		return
	}

	for _, lap := range activity.Laps {
		lap.StartPositionLat = fit.NewLatitudeInvalid()
		lap.StartPositionLong = fit.NewLongitudeInvalid()
		lap.EndPositionLat = fit.NewLatitudeInvalid()
		lap.EndPositionLong = fit.NewLongitudeInvalid()
	}

	for _, sesh := range activity.Sessions {
		sesh.StartPositionLat = fit.NewLatitudeInvalid()
		sesh.StartPositionLong = fit.NewLongitudeInvalid()
		sesh.EndPositionLat = fit.NewLatitudeInvalid()
		sesh.EndPositionLong = fit.NewLongitudeInvalid()

		sesh.NecLat = fit.NewLatitudeInvalid()
		sesh.NecLong = fit.NewLongitudeInvalid()
		sesh.SwcLat = fit.NewLatitudeInvalid()
		sesh.SwcLong = fit.NewLongitudeInvalid()
	}

	for _, record := range activity.Records {
		record.PositionLat = fit.NewLatitudeInvalid()
		record.PositionLong = fit.NewLongitudeInvalid()
	}
}

// generateRandomFilename generates a random 8-character alphanumeric string
func generateRandomFilename() string {
	randBytes := make([]byte, 4) // 4 bytes will result in 8 hex characters
	_, err := rand.Read(randBytes)
	if err != nil {
		log.Fatalf("Failed to generate random filename: %v", err)
	}
	return hex.EncodeToString(randBytes)
}
