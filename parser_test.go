package main

import (
	"testing"
	"time"
)

func TestParseTimeWithContext(t *testing.T) {
	// Reference time: Friday, January 16, 2026 at 2:30 PM UTC
	refTime := time.Date(2026, 1, 16, 14, 30, 0, 0, time.UTC)

	tests := []struct {
		name        string
		input       string
		refTime     time.Time
		wantHour    int
		wantMinute  int
		wantDay     int
		wantMonth   time.Month
		wantYear    int
		shouldError bool
	}{
		// Natural language
		{
			name:       "now",
			input:      "now",
			refTime:    refTime,
			wantHour:   14,
			wantMinute: 30,
			wantDay:    16,
			wantMonth:  time.January,
			wantYear:   2026,
		},
		{
			name:       "noon",
			input:      "noon",
			refTime:    refTime,
			wantHour:   12,
			wantMinute: 0,
			wantDay:    16,
			wantMonth:  time.January,
			wantYear:   2026,
		},
		{
			name:       "midnight",
			input:      "midnight",
			refTime:    refTime,
			wantHour:   0,
			wantMinute: 0,
			wantDay:    16,
			wantMonth:  time.January,
			wantYear:   2026,
		},

		// Simple times
		{
			name:       "3pm",
			input:      "3pm",
			refTime:    refTime,
			wantHour:   15,
			wantMinute: 0,
			wantDay:    16,
			wantMonth:  time.January,
			wantYear:   2026,
		},
		{
			name:       "10am",
			input:      "10am",
			refTime:    refTime,
			wantHour:   10,
			wantMinute: 0,
			wantDay:    16,
			wantMonth:  time.January,
			wantYear:   2026,
		},
		{
			name:       "3:30pm",
			input:      "3:30pm",
			refTime:    refTime,
			wantHour:   15,
			wantMinute: 30,
			wantDay:    16,
			wantMonth:  time.January,
			wantYear:   2026,
		},
		{
			name:       "15:04",
			input:      "15:04",
			refTime:    refTime,
			wantHour:   15,
			wantMinute: 4,
			wantDay:    16,
			wantMonth:  time.January,
			wantYear:   2026,
		},

		// Relative - in X hours/minutes/days
		{
			name:       "in 2 hours",
			input:      "in 2 hours",
			refTime:    refTime,
			wantHour:   16,
			wantMinute: 30,
			wantDay:    16,
			wantMonth:  time.January,
			wantYear:   2026,
		},
		{
			name:       "in 30 minutes",
			input:      "in 30 minutes",
			refTime:    refTime,
			wantHour:   15,
			wantMinute: 0,
			wantDay:    16,
			wantMonth:  time.January,
			wantYear:   2026,
		},
		{
			name:       "in 1 day",
			input:      "in 1 day",
			refTime:    refTime,
			wantHour:   14,
			wantMinute: 30,
			wantDay:    17,
			wantMonth:  time.January,
			wantYear:   2026,
		},
		{
			name:       "in 3 days",
			input:      "in 3 days",
			refTime:    refTime,
			wantHour:   14,
			wantMinute: 30,
			wantDay:    19,
			wantMonth:  time.January,
			wantYear:   2026,
		},

		// Tomorrow/Yesterday
		{
			name:       "tomorrow",
			input:      "tomorrow",
			refTime:    refTime,
			wantHour:   9,
			wantMinute: 0,
			wantDay:    17,
			wantMonth:  time.January,
			wantYear:   2026,
		},
		{
			name:       "tomorrow 3pm",
			input:      "tomorrow 3pm",
			refTime:    refTime,
			wantHour:   15,
			wantMinute: 0,
			wantDay:    17,
			wantMonth:  time.January,
			wantYear:   2026,
		},
		{
			name:       "tomorrow noon",
			input:      "tomorrow noon",
			refTime:    refTime,
			wantHour:   12,
			wantMinute: 0,
			wantDay:    17,
			wantMonth:  time.January,
			wantYear:   2026,
		},
		{
			name:       "yesterday 10am",
			input:      "yesterday 10am",
			refTime:    refTime,
			wantHour:   10,
			wantMinute: 0,
			wantDay:    15,
			wantMonth:  time.January,
			wantYear:   2026,
		},

		// Next weekday
		{
			name:       "next monday",
			input:      "next monday",
			refTime:    refTime, // Friday Jan 16
			wantHour:   9,
			wantMinute: 0,
			wantDay:    19, // Monday Jan 19
			wantMonth:  time.January,
			wantYear:   2026,
		},
		{
			name:       "next monday 3pm",
			input:      "next monday 3pm",
			refTime:    refTime,
			wantHour:   15,
			wantMinute: 0,
			wantDay:    19,
			wantMonth:  time.January,
			wantYear:   2026,
		},
		{
			name:       "next friday",
			input:      "next friday",
			refTime:    refTime, // Friday Jan 16
			wantHour:   9,
			wantMinute: 0,
			wantDay:    23, // Next Friday Jan 23
			wantMonth:  time.January,
			wantYear:   2026,
		},
		{
			name:       "next tuesday 10:30am",
			input:      "next tuesday 10:30am",
			refTime:    refTime,
			wantHour:   10,
			wantMinute: 30,
			wantDay:    20,
			wantMonth:  time.January,
			wantYear:   2026,
		},

		// Dates with times
		{
			name:       "2026-01-20 3pm",
			input:      "2026-01-20 3pm",
			refTime:    refTime,
			wantHour:   15,
			wantMinute: 0,
			wantDay:    20,
			wantMonth:  time.January,
			wantYear:   2026,
		},
		{
			name:       "2026-02-14 noon",
			input:      "2026-02-14 noon",
			refTime:    refTime,
			wantHour:   12,
			wantMinute: 0,
			wantDay:    14,
			wantMonth:  time.February,
			wantYear:   2026,
		},
		{
			name:       "jan 20 3pm",
			input:      "jan 20 3pm",
			refTime:    refTime,
			wantHour:   15,
			wantMinute: 0,
			wantDay:    20,
			wantMonth:  time.January,
			wantYear:   2026,
		},
		{
			name:       "february 14 10:30am",
			input:      "february 14 10:30am",
			refTime:    refTime,
			wantHour:   10,
			wantMinute: 30,
			wantDay:    14,
			wantMonth:  time.February,
			wantYear:   2026,
		},

		// Error cases
		{
			name:        "empty string",
			input:       "",
			refTime:     refTime,
			shouldError: true,
		},
		{
			name:        "invalid format",
			input:       "asdfghjkl",
			refTime:     refTime,
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseTimeWithContext(tt.input, tt.refTime)

			if tt.shouldError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if result.Time.Hour() != tt.wantHour {
				t.Errorf("hour = %d, want %d", result.Time.Hour(), tt.wantHour)
			}
			if result.Time.Minute() != tt.wantMinute {
				t.Errorf("minute = %d, want %d", result.Time.Minute(), tt.wantMinute)
			}
			if result.Time.Day() != tt.wantDay {
				t.Errorf("day = %d, want %d", result.Time.Day(), tt.wantDay)
			}
			if result.Time.Month() != tt.wantMonth {
				t.Errorf("month = %v, want %v", result.Time.Month(), tt.wantMonth)
			}
			if result.Time.Year() != tt.wantYear {
				t.Errorf("year = %d, want %d", result.Time.Year(), tt.wantYear)
			}
		})
	}
}

func TestParseRelativeTime(t *testing.T) {
	refTime := time.Date(2026, 1, 16, 14, 30, 0, 0, time.UTC)

	tests := []struct {
		name      string
		input     string
		wantDay   int
		wantHour  int
		wantError bool
	}{
		{
			name:     "in 1 hour",
			input:    "in 1 hour",
			wantDay:  16,
			wantHour: 15,
		},
		{
			name:     "in 5 hours",
			input:    "in 5 hours",
			wantDay:  16,
			wantHour: 19,
		},
		{
			name:     "tomorrow default",
			input:    "tomorrow",
			wantDay:  17,
			wantHour: 9,
		},
		{
			name:     "yesterday default",
			input:    "yesterday",
			wantDay:  15,
			wantHour: 9,
		},
		{
			name:      "not relative",
			input:     "3pm",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseRelativeTime(tt.input, refTime)

			if tt.wantError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if result.Time.Day() != tt.wantDay {
				t.Errorf("day = %d, want %d", result.Time.Day(), tt.wantDay)
			}
			if result.Time.Hour() != tt.wantHour {
				t.Errorf("hour = %d, want %d", result.Time.Hour(), tt.wantHour)
			}
		})
	}
}

func TestParseDateWithTime(t *testing.T) {
	refTime := time.Date(2026, 1, 16, 14, 30, 0, 0, time.UTC)

	tests := []struct {
		name      string
		input     string
		wantYear  int
		wantMonth time.Month
		wantDay   int
		wantHour  int
		wantError bool
	}{
		{
			name:      "ISO format",
			input:     "2026-03-15 3pm",
			wantYear:  2026,
			wantMonth: time.March,
			wantDay:   15,
			wantHour:  15,
		},
		{
			name:      "month name",
			input:     "march 15 3pm",
			wantYear:  2026,
			wantMonth: time.March,
			wantDay:   15,
			wantHour:  15,
		},
		{
			name:      "numeric date",
			input:     "3/15 3pm",
			wantYear:  2026,
			wantMonth: time.March,
			wantDay:   15,
			wantHour:  15,
		},
		{
			name:      "not a date",
			input:     "tomorrow 3pm",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseDateWithTime(tt.input, refTime)

			if tt.wantError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if result.Time.Year() != tt.wantYear {
				t.Errorf("year = %d, want %d", result.Time.Year(), tt.wantYear)
			}
			if result.Time.Month() != tt.wantMonth {
				t.Errorf("month = %v, want %v", result.Time.Month(), tt.wantMonth)
			}
			if result.Time.Day() != tt.wantDay {
				t.Errorf("day = %d, want %d", result.Time.Day(), tt.wantDay)
			}
			if result.Time.Hour() != tt.wantHour {
				t.Errorf("hour = %d, want %d", result.Time.Hour(), tt.wantHour)
			}
		})
	}
}

func TestParseSimpleTime(t *testing.T) {
	refTime := time.Date(2026, 1, 16, 14, 30, 0, 0, time.UTC)

	tests := []struct {
		name       string
		input      string
		wantHour   int
		wantMinute int
		wantError  bool
	}{
		{
			name:       "3pm",
			input:      "3pm",
			wantHour:   15,
			wantMinute: 0,
		},
		{
			name:       "10:30am",
			input:      "10:30am",
			wantHour:   10,
			wantMinute: 30,
		},
		{
			name:       "15:04",
			input:      "15:04",
			wantHour:   15,
			wantMinute: 4,
		},
		{
			name:       "noon",
			input:      "noon",
			wantHour:   12,
			wantMinute: 0,
		},
		{
			name:       "midnight",
			input:      "midnight",
			wantHour:   0,
			wantMinute: 0,
		},
		{
			name:      "invalid",
			input:     "invalid",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseSimpleTime(tt.input, refTime)

			if tt.wantError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if result.Time.Hour() != tt.wantHour {
				t.Errorf("hour = %d, want %d", result.Time.Hour(), tt.wantHour)
			}
			if result.Time.Minute() != tt.wantMinute {
				t.Errorf("minute = %d, want %d", result.Time.Minute(), tt.wantMinute)
			}
		})
	}
}
