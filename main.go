package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go <input_directory> <output_file>")
		return
	}

	dir := os.Args[1]
	outputFile := os.Args[2]

	outputFilePath, err := filepath.Abs(outputFile)
	if err != nil {
		fmt.Printf("Failed to determine absolute path for output file: %s\n", err)
		return
	}

	files, err := filepath.Glob(filepath.Join(dir, "*.csv"))
	if err != nil {
		fmt.Printf("Failed to read directory: %s\n", err)
		return
	}

	if len(files) == 0 {
		fmt.Println("No CSV files found in the specified directory.")
		return
	}

	outFile, err := os.Create(outputFilePath)
	if err != nil {
		fmt.Printf("Failed to create output file: %s\n", err)
		return
	}
	defer outFile.Close()

	w := csv.NewWriter(outFile)
	defer w.Flush()

	h := false

	for _, file := range files {
		f, err := os.Open(file)
		if err != nil {
			fmt.Printf("Failed to open file %s: %s\n", file, err)
			continue
		}
		defer f.Close()

		r := csv.NewReader(f)
		records, err := r.ReadAll()
		if err != nil {
			fmt.Printf("Failed to read file %s: %s\n", file, err)
			continue
		}

		if !h {
			w.Write(records[0])
			h = true
		}

		for _, record := range records[1:] {
			w.Write(record)
		}
	}

	fmt.Printf("Merged %d CSV files into %s\n", len(files), outputFile)
}
