package jscal

import (
	"testing"
	"time"
)

func TestParseISO8601Duration(t *testing.T) {
	tests := []struct {
		name     string
		duration string
		want     time.Duration
		wantErr  bool
	}{
		// Basic durations
		{"1 hour", "PT1H", time.Hour, false},
		{"30 minutes", "PT30M", 30 * time.Minute, false},
		{"45 seconds", "PT45S", 45 * time.Second, false},
		{"1.5 hours", "PT1.5H", 90 * time.Minute, false},
		{"2.5 minutes", "PT2.5M", 150 * time.Second, false},

		// Combined durations
		{"1h 30m", "PT1H30M", 90 * time.Minute, false},
		{"2h 15m 30s", "PT2H15M30S", 2*time.Hour + 15*time.Minute + 30*time.Second, false},
		{"1h 0m 45s", "PT1H0M45S", time.Hour + 45*time.Second, false},

		// Date components (converted to hours)
		{"1 day", "P1D", 24 * time.Hour, false},
		{"1 week", "P1W", 7 * 24 * time.Hour, false},
		{"2 weeks", "P2W", 14 * 24 * time.Hour, false},
		{"1 month", "P1M", 30 * 24 * time.Hour, false},
		{"1 year", "P1Y", 365 * 24 * time.Hour, false},

		// Combined date and time
		{"1 day 2 hours", "P1DT2H", 26 * time.Hour, false},
		{"1 week 3 days", "P1W3D", 10 * 24 * time.Hour, false},
		{"1 year 2 months 3 days", "P1Y2M3DT0H", (365 + 60 + 3) * 24 * time.Hour, false},

		// Edge cases
		{"0 duration", "PT0S", 0, false},
		{"only P", "P", 0, false},
		{"only PT", "PT", 0, false},
		{"negative duration", "-PT1H", -time.Hour, false},
		{"invalid format", "1H30M", 0, true},
		{"missing P", "T1H", 0, true},
		{"text", "one hour", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseISO8601Duration(tt.duration)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseISO8601Duration() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("parseISO8601Duration() = %v, want %v", got, tt.want)
			}
		})
	}
}
