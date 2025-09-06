package jscal

import (
	"fmt"
	"strings"
	"time"
)

// parseISO8601Duration parses an ISO 8601 duration string to time.Duration
func parseISO8601Duration(duration string) (time.Duration, error) {
	// Parser for ISO 8601 durations like PT1H, P1D, PT30M, P1Y2M3DT4H5M6S
	// Supports: Years (Y), Months (M), Weeks (W), Days (D), Hours (H), Minutes (M), Seconds (S)
	// Supports fractional values for all units
	// Supports negative durations (e.g., -PT15M)

	if duration == "" {
		return 0, fmt.Errorf("invalid ISO 8601 duration: empty string")
	}
	
	// Check for negative duration
	negative := false
	if strings.HasPrefix(duration, "-") {
		negative = true
		duration = duration[1:]
	}
	
	if !strings.HasPrefix(duration, "P") {
		return 0, fmt.Errorf("invalid ISO 8601 duration: must start with P")
	}

	duration = duration[1:] // Remove P
	if duration == "" {
		return 0, nil // "P" alone is valid (0 duration)
	}

	var result time.Duration

	// Check for time portion
	timeIndex := strings.Index(duration, "T")
	var datePart, timePart string

	if timeIndex >= 0 {
		datePart = duration[:timeIndex]
		timePart = duration[timeIndex+1:]
	} else {
		datePart = duration
	}

	// Parse date part (PYMWD)
	if datePart != "" {
		remaining := datePart

		// Parse Years (approximate: 365 days)
		if idx := strings.Index(remaining, "Y"); idx >= 0 {
			years := 0.0
			if n, err := fmt.Sscanf(remaining[:idx], "%f", &years); n == 1 && err == nil {
				result += time.Duration(years * 365 * 24 * float64(time.Hour))
			}
			remaining = remaining[idx+1:]
		}

		// Parse Months (approximate: 30 days)
		if idx := strings.Index(remaining, "M"); idx >= 0 {
			months := 0.0
			if n, err := fmt.Sscanf(remaining[:idx], "%f", &months); n == 1 && err == nil {
				result += time.Duration(months * 30 * 24 * float64(time.Hour))
			}
			remaining = remaining[idx+1:]
		}

		// Parse Weeks
		if idx := strings.Index(remaining, "W"); idx >= 0 {
			weeks := 0.0
			if n, err := fmt.Sscanf(remaining[:idx], "%f", &weeks); n == 1 && err == nil {
				result += time.Duration(weeks * 7 * 24 * float64(time.Hour))
			}
			remaining = remaining[idx+1:]
		}

		// Parse Days
		if idx := strings.Index(remaining, "D"); idx >= 0 {
			days := 0.0
			if n, err := fmt.Sscanf(remaining[:idx], "%f", &days); n == 1 && err == nil {
				result += time.Duration(days * 24 * float64(time.Hour))
			}
		}
	}

	// Parse time part (HMS)
	if timePart != "" {
		remaining := timePart

		// Parse Hours
		if idx := strings.Index(remaining, "H"); idx >= 0 {
			hours := 0.0
			if n, err := fmt.Sscanf(remaining[:idx], "%f", &hours); n == 1 && err == nil {
				result += time.Duration(hours * float64(time.Hour))
			}
			remaining = remaining[idx+1:]
		}

		// Parse Minutes
		if idx := strings.Index(remaining, "M"); idx >= 0 {
			minutes := 0.0
			if n, err := fmt.Sscanf(remaining[:idx], "%f", &minutes); n == 1 && err == nil {
				result += time.Duration(minutes * float64(time.Minute))
			}
			remaining = remaining[idx+1:]
		}

		// Parse Seconds
		if idx := strings.Index(remaining, "S"); idx >= 0 {
			seconds := 0.0
			if n, err := fmt.Sscanf(remaining[:idx], "%f", &seconds); n == 1 && err == nil {
				result += time.Duration(seconds * float64(time.Second))
			}
		}
	}
	
	// Apply negative flag if set
	if negative {
		result = -result
	}

	return result, nil
}