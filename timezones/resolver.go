package timezones

import (
	"fmt"
	"strings"
	"time"
)

// Resolve takes a city name/alias and returns the IANA timezone location
func Resolve(name string) (*time.Location, error) {
	// Normalize input
	normalized := strings.ToLower(strings.TrimSpace(name))
	normalized = strings.ReplaceAll(normalized, "_", " ")

	// Check manual aliases first (NYC -> new york)
	if canonical, ok := ManualAliases[normalized]; ok {
		normalized = canonical
	}

	// Check generated cities map
	if tz, ok := GeneratedCities[normalized]; ok {
		return time.LoadLocation(tz)
	}

	// Fallback: try IANA timezone variations
	variations := []string{
		name,
		strings.ReplaceAll(name, " ", "_"),
		"America/" + strings.ReplaceAll(name, " ", "_"),
		"Europe/" + strings.ReplaceAll(name, " ", "_"),
		"Asia/" + strings.ReplaceAll(name, " ", "_"),
		"Africa/" + strings.ReplaceAll(name, " ", "_"),
		"Australia/" + strings.ReplaceAll(name, " ", "_"),
		"Pacific/" + strings.ReplaceAll(name, " ", "_"),
	}

	for _, v := range variations {
		if loc, err := time.LoadLocation(v); err == nil {
			return loc, nil
		}
	}

	return nil, fmt.Errorf("unknown location: %s", name)
}