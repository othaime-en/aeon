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

	loc, err := time.LoadLocation(normalized)
	if err != nil {
		return nil, fmt.Errorf("unknown location: %s", name)
	}
	return loc, nil
}