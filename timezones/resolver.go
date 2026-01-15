package timezones

import (
	"fmt"
	"time"
)

// Resolve takes a city name/alias and returns the IANA timezone location
func Resolve(name string) (*time.Location, error) {
	loc, err := time.LoadLocation(name)
	if err != nil {
		return nil, fmt.Errorf("unknown location: %s", name)
	}
	return loc, nil
}
