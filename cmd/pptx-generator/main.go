package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"haoma/internal/application/services"
)

func main() {
	// Define command-line flags
	inputFile := flag.String("input", "", "Path to input Marp markdown file")
	outputFile := flag.String("output", "", "Path to output PowerPoint file (.pptx)")
	
	flag.Parse()

	// Validate inputs
	if *inputFile == "" {
		fmt.Println("Usage: pptx-generator -input <input.md> -output <output.pptx>")
		fmt.Println("\nExample:")
		fmt.Println("  pptx-generator -input data/yggdrasil-slides.md -output yggdrasil.pptx")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Set default output file if not provided
	if *outputFile == "" {
		baseName := filepath.Base(*inputFile)
		ext := filepath.Ext(baseName)
		*outputFile = baseName[:len(baseName)-len(ext)] + ".pptx"
	}

	// Check if input file exists
	if _, err := os.Stat(*inputFile); os.IsNotExist(err) {
		log.Fatalf("Input file does not exist: %s", *inputFile)
	}

	// Create presentation generator
	generator := services.NewPresentationGenerator()

	// Generate presentation
	fmt.Printf("Generating PowerPoint presentation...\n")
	fmt.Printf("  Input:  %s\n", *inputFile)
	fmt.Printf("  Output: %s\n", *outputFile)

	if err := generator.GenerateFromMarp(*inputFile, *outputFile); err != nil {
		log.Fatalf("Failed to generate presentation: %v", err)
	}

	fmt.Printf("\n✓ Successfully generated: %s\n", *outputFile)
}
