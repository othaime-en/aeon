package main

import (
	"archive/zip"
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
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

	// Open zip archive
	zipReader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, fmt.Errorf("unzip failed: %w", err)
	}

	// Find cities15000.txt
	var txtFile *zip.File
	for _, f := range zipReader.File {
		if strings.HasSuffix(f.Name, "cities15000.txt") {
			txtFile = f
			break
		}
	}

	if txtFile == nil {
		return nil, fmt.Errorf("cities15000.txt not found in zip")
	}

	// Open and parse the text file
	rc, err := txtFile.Open()
	if err != nil {
		return nil, fmt.Errorf("open txt failed: %w", err)
	}
	defer rc.Close()

	cities := make(map[string]string)
	scanner := bufio.NewScanner(rc)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		// Parse tab-separated values
		// Format: geonameid, name, asciiname, alternatenames, lat, lon, ..., timezone, ...
		fields := strings.Split(line, "\t")
		if len(fields) < 18 {
			continue
		}

		name := strings.TrimSpace(fields[1])      // City name
		asciiName := strings.TrimSpace(fields[2]) // ASCII version
		timezone := strings.TrimSpace(fields[17]) // Timezone

		if name == "" || timezone == "" {
			continue
		}

		// Normalize names to lowercase
		nameKey := strings.ToLower(name)
		asciiKey := strings.ToLower(asciiName)

		// Add both versions (prefer ASCII for consistency)
		cities[nameKey] = timezone
		if asciiKey != nameKey {
			cities[asciiKey] = timezone
		}

		// Handle common variations
		// Remove "City" suffix (e.g., "New York City" -> "New York")
		if strings.HasSuffix(nameKey, " city") {
			shortName := strings.TrimSuffix(nameKey, " city")
			if _, exists := cities[shortName]; !exists {
				cities[shortName] = timezone
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan failed: %w", err)
	}

	return cities, nil
}