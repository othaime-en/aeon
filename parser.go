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
// - Relative times: "tomorrow 10am", "in 2 hours", "yesterday 3pm"
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

	// Try relative time parsing first
	if result, err := parseRelativeTime(input, refTime); err == nil {
		return result, nil
	}

	// Fall back to simple time parsing (original behavior)
	if result, err := parseSimpleTime(input, refTime); err == nil {
		return result, nil
	}

	return ParsedTime{}, fmt.Errorf("could not parse time: %s", input)
}

// parseRelativeTime handles relative time expressions
func parseRelativeTime(input string, refTime time.Time) (ParsedTime, error) {
	// "in X hours/minutes/days"
	inPattern := regexp.MustCompile(`^in\s+(\d+)\s+(hour|hours|minute|minutes|min|mins|day|days)$`)
	if matches := inPattern.FindStringSubmatch(input); matches != nil {
		amount, _ := strconv.Atoi(matches[1])
		unit := matches[2]

		var duration time.Duration
		switch {
		case strings.HasPrefix(unit, "hour"):
			duration = time.Duration(amount) * time.Hour
		case strings.HasPrefix(unit, "min"):
			duration = time.Duration(amount) * time.Minute
		case strings.HasPrefix(unit, "day"):
			duration = time.Duration(amount) * 24 * time.Hour
		}

		result := refTime.Add(duration)
		return ParsedTime{Time: result, Original: input}, nil
	}

	// "tomorrow [time]" or just "tomorrow"
	if strings.HasPrefix(input, "tomorrow") {
		tomorrow := refTime.AddDate(0, 0, 1)
		remaining := strings.TrimSpace(strings.TrimPrefix(input, "tomorrow"))

		if remaining == "" {
			// Just "tomorrow" - use 9am as default
			result := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 9, 0, 0, 0, refTime.Location())
			return ParsedTime{Time: result, Original: input}, nil
		}

		// "tomorrow 3pm"
		timeResult, err := parseSimpleTime(remaining, tomorrow)
		if err != nil {
			return ParsedTime{}, err
		}
		return ParsedTime{Time: timeResult.Time, Original: input}, nil
	}

	// "yesterday [time]" or just "yesterday"
	if strings.HasPrefix(input, "yesterday") {
		yesterday := refTime.AddDate(0, 0, -1)
		remaining := strings.TrimSpace(strings.TrimPrefix(input, "yesterday"))

		if remaining == "" {
			result := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 9, 0, 0, 0, refTime.Location())
			return ParsedTime{Time: result, Original: input}, nil
		}

		timeResult, err := parseSimpleTime(remaining, yesterday)
		if err != nil {
			return ParsedTime{}, err
		}
		return ParsedTime{Time: timeResult.Time, Original: input}, nil
	}

	// "next monday/tuesday/etc [time]"
	weekdays := map[string]time.Weekday{
		"sunday": time.Sunday, "sun": time.Sunday,
		"monday": time.Monday, "mon": time.Monday,
		"tuesday": time.Tuesday, "tue": time.Tuesday, "tues": time.Tuesday,
		"wednesday": time.Wednesday, "wed": time.Wednesday,
		"thursday": time.Thursday, "thu": time.Thursday, "thur": time.Thursday, "thurs": time.Thursday,
		"friday": time.Friday, "fri": time.Friday,
		"saturday": time.Saturday, "sat": time.Saturday,
	}

	nextPattern := regexp.MustCompile(`^next\s+(\w+)(?:\s+(.+))?$`)
	if matches := nextPattern.FindStringSubmatch(input); matches != nil {
		dayName := matches[1]
		timeStr := matches[2]

		if targetWeekday, ok := weekdays[dayName]; ok {
			// Find next occurrence of this weekday
			currentWeekday := refTime.Weekday()
			daysUntil := int(targetWeekday - currentWeekday)
			if daysUntil <= 0 {
				daysUntil += 7
			}

			targetDate := refTime.AddDate(0, 0, daysUntil)

			if timeStr == "" {
				// Default to 9am
				result := time.Date(targetDate.Year(), targetDate.Month(), targetDate.Day(), 9, 0, 0, 0, refTime.Location())
				return ParsedTime{Time: result, Original: input}, nil
			}

			// Parse the time part
			timeResult, err := parseSimpleTime(timeStr, targetDate)
			if err != nil {
				return ParsedTime{}, err
			}
			return ParsedTime{Time: timeResult.Time, Original: input}, nil
		}
	}

	return ParsedTime{}, fmt.Errorf("not a relative time")
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
