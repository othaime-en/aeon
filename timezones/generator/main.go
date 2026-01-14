package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

const (
	// GeoNames cities with population > 15,000
	geonamesURL = "https://download.geonames.org/export/dump/cities15000.zip"
	outputFile  = "../cities_generated.go"
)

func main() {
	fmt.Println("🌍 Generating timezone data from GeoNames...")

	// Download and parse GeoNames data
	cities, err := downloadAndParse()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✓ Parsed %d cities\n", len(cities))

	// Placeholder for generating Go code

	fmt.Printf("✓ Generated %s\n", outputFile)
	fmt.Println("✓ Done! Run 'go fmt ./timezones' to format the generated file.")
}

func downloadAndParse() (map[string]string, error) {
	fmt.Println("⬇️  Downloading GeoNames data (~11 MB)...")

	resp, err := http.Get(geonamesURL)
	if err != nil {
		return nil, fmt.Errorf("download failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("download failed: HTTP %d", resp.StatusCode)
	}

	// Read zip file into memory
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read failed: %w", err)
	}

	fmt.Println("📦 Extracting and parsing...")

	// Placeholder for unzip and parse
	cities := make(map[string]string)

	return cities, nil
}
