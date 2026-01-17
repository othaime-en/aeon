package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// ParsedTime represents a parsed time with its components
type ParsedTime struct {
	Time     time.Time
	Original string
}

// parseTimeWithContext parses a time string with support for:
// - Natural language: "noon", "midnight", "now"
// - Traditional: "3pm", "15:04"
func parseTimeWithContext(input string, refTime time.Time) (ParsedTime, error) {
	input = strings.ToLower(strings.TrimSpace(input))

	if input == "" {
		return ParsedTime{}, fmt.Errorf("empty time string")
	}

	// Natural language shortcuts
	switch input {
	case "now":
		return ParsedTime{Time: refTime, Original: input}, nil
	case "noon":
		return ParsedTime{
			Time:     time.Date(refTime.Year(), refTime.Month(), refTime.Day(), 12, 0, 0, 0, refTime.Location()),
			Original: input,
		}, nil
	case "midnight":
		return ParsedTime{
			Time:     time.Date(refTime.Year(), refTime.Month(), refTime.Day(), 0, 0, 0, 0, refTime.Location()),
			Original: input,
		}, nil
	}

	// Fall back to simple time parsing (original behavior)
	if result, err := parseSimpleTime(input, refTime); err == nil {
		return result, nil
	}

	return ParsedTime{}, fmt.Errorf("could not parse time: %s", input)
}

// parseSimpleTime handles basic time formats (original parseTime logic)
func parseSimpleTime(input string, baseDate time.Time) (ParsedTime, error) {
	input = strings.ToLower(strings.TrimSpace(input))

	// Handle "noon" and "midnight" if passed as simple time
	switch input {
	case "noon":
		result := time.Date(baseDate.Year(), baseDate.Month(), baseDate.Day(), 12, 0, 0, 0, baseDate.Location())
		return ParsedTime{Time: result, Original: input}, nil
	case "midnight":
		result := time.Date(baseDate.Year(), baseDate.Month(), baseDate.Day(), 0, 0, 0, 0, baseDate.Location())
		return ParsedTime{Time: result, Original: input}, nil
	}

	// Try parsing with various time formats
	timeFormats := []string{
		"3pm",
		"3:04pm",
		"15:04",
		"3PM",
		"3:04PM",
		"15",
		"3",
	}

	for _, format := range timeFormats {
		parsed, err := time.Parse(format, input)
		if err == nil {
			// Combine parsed time with base date
			result := time.Date(
				baseDate.Year(), baseDate.Month(), baseDate.Day(),
				parsed.Hour(), parsed.Minute(), parsed.Second(),
				0, baseDate.Location(),
			)
			return ParsedTime{Time: result, Original: input}, nil
		}
	}

	// Try manual parsing for formats like "3pm", "10am"
	ampmPattern := regexp.MustCompile(`^(\d{1,2})(?::(\d{2}))?\s*(am|pm)?$`)
	if matches := ampmPattern.FindStringSubmatch(input); matches != nil {
		hour, _ := strconv.Atoi(matches[1])
		minute := 0
		if matches[2] != "" {
			minute, _ = strconv.Atoi(matches[2])
		}

		// Handle AM/PM
		if matches[3] == "pm" && hour < 12 {
			hour += 12
		} else if matches[3] == "am" && hour == 12 {
			hour = 0
		}

		if hour >= 0 && hour < 24 && minute >= 0 && minute < 60 {
			result := time.Date(
				baseDate.Year(), baseDate.Month(), baseDate.Day(),
				hour, minute, 0, 0, baseDate.Location(),
			)
			return ParsedTime{Time: result, Original: input}, nil
		}
	}

	return ParsedTime{}, fmt.Errorf("invalid time format")
}
