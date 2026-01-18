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
