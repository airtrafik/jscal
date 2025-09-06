package jscal

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// LocalDateTime represents a date-time value without timezone information.
// This type is defined in RFC 8984 Section 1.4.4.
// It's a type alias of time.Time but marshals/unmarshals without timezone.
type LocalDateTime time.Time

// NewLocalDateTime creates a LocalDateTime from a time.Time value,
// preserving the actual time but ignoring timezone for formatting.
func NewLocalDateTime(t time.Time) *LocalDateTime {
	ldt := LocalDateTime(t)
	return &ldt
}

// Time converts the LocalDateTime to a time.Time value.
func (ldt *LocalDateTime) Time() time.Time {
	if ldt == nil {
		return time.Time{}
	}
	return time.Time(*ldt)
}

// String returns the LocalDateTime in RFC 3339 format without timezone.
func (ldt LocalDateTime) String() string {
	t := time.Time(ldt)
	// Use RFC3339Nano to preserve nanoseconds
	s := t.Format(time.RFC3339Nano)
	// Remove the timezone suffix (Z or +00:00)
	if strings.HasSuffix(s, "Z") {
		return s[:len(s)-1]
	}
	if idx := strings.LastIndex(s, "+"); idx > 0 {
		return s[:idx]
	}
	if idx := strings.LastIndex(s, "-"); idx > 10 { // Avoid matching the date separators
		return s[:idx]
	}
	return s
}

// MarshalJSON implements json.Marshaler.
func (ldt LocalDateTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(ldt.String())
}

// UnmarshalJSON implements json.Unmarshaler.
func (ldt *LocalDateTime) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	parsed, err := ParseLocalDateTime(s)
	if err != nil {
		return err
	}
	*ldt = *parsed
	return nil
}

// ParseLocalDateTime parses a string in RFC 3339 format without timezone.
func ParseLocalDateTime(s string) (*LocalDateTime, error) {
	// Add 'Z' if no timezone is specified to make it valid RFC3339
	if !strings.Contains(s, "Z") && !strings.Contains(s, "+") && !strings.Contains(s, "-") {
		// Check if there's a T to ensure it's a datetime format
		if strings.Contains(s, "T") {
			s = s + "Z"
		}
	}

	// Try parsing with various formats
	formats := []string{
		time.RFC3339,
		time.RFC3339Nano,
		"2006-01-02T15:04:05",
		"2006-01-02T15:04:05.999999999",
	}

	var t time.Time
	var err error
	for _, format := range formats {
		t, err = time.Parse(format, s)
		if err == nil {
			break
		}
	}

	if err != nil {
		return nil, fmt.Errorf("invalid LocalDateTime format: %s", s)
	}

	return NewLocalDateTime(t), nil
}

// Equal returns true if the two LocalDateTime values represent the same moment.
// This comparison ignores timezone differences.
func (ldt *LocalDateTime) Equal(other *LocalDateTime) bool {
	if ldt == nil || other == nil {
		return ldt == other
	}

	t1 := time.Time(*ldt)
	t2 := time.Time(*other)

	// Compare the actual time values (ignoring location)
	return t1.Year() == t2.Year() &&
		t1.Month() == t2.Month() &&
		t1.Day() == t2.Day() &&
		t1.Hour() == t2.Hour() &&
		t1.Minute() == t2.Minute() &&
		t1.Second() == t2.Second() &&
		t1.Nanosecond() == t2.Nanosecond()
}

// IsZero returns true if the LocalDateTime is the zero value.
func (ldt *LocalDateTime) IsZero() bool {
	if ldt == nil {
		return true
	}
	return time.Time(*ldt).IsZero()
}

// Before returns true if ldt is before other.
func (ldt *LocalDateTime) Before(other *LocalDateTime) bool {
	if ldt == nil || other == nil {
		return false
	}
	t1 := time.Time(*ldt)
	t2 := time.Time(*other)

	// Create times in the same location for comparison
	utc1 := time.Date(t1.Year(), t1.Month(), t1.Day(), t1.Hour(), t1.Minute(), t1.Second(), t1.Nanosecond(), time.UTC)
	utc2 := time.Date(t2.Year(), t2.Month(), t2.Day(), t2.Hour(), t2.Minute(), t2.Second(), t2.Nanosecond(), time.UTC)

	return utc1.Before(utc2)
}

// After returns true if ldt is after other.
func (ldt *LocalDateTime) After(other *LocalDateTime) bool {
	if ldt == nil || other == nil {
		return false
	}
	t1 := time.Time(*ldt)
	t2 := time.Time(*other)

	// Create times in the same location for comparison
	utc1 := time.Date(t1.Year(), t1.Month(), t1.Day(), t1.Hour(), t1.Minute(), t1.Second(), t1.Nanosecond(), time.UTC)
	utc2 := time.Date(t2.Year(), t2.Month(), t2.Day(), t2.Hour(), t2.Minute(), t2.Second(), t2.Nanosecond(), time.UTC)

	return utc1.After(utc2)
}

// Add returns the LocalDateTime with the given duration added.
func (ldt LocalDateTime) Add(d time.Duration) LocalDateTime {
	t := time.Time(ldt)
	return LocalDateTime(t.Add(d))
}

// Sub returns the duration between two LocalDateTime values.
func (ldt LocalDateTime) Sub(other LocalDateTime) time.Duration {
	t1 := time.Time(ldt)
	t2 := time.Time(other)
	return t1.Sub(t2)
}

// Format formats the LocalDateTime using the given layout.
func (ldt LocalDateTime) Format(layout string) string {
	t := time.Time(ldt)
	return t.Format(layout)
}

// Year returns the year of the LocalDateTime.
func (ldt LocalDateTime) Year() int {
	return time.Time(ldt).Year()
}

// Month returns the month of the LocalDateTime.
func (ldt LocalDateTime) Month() time.Month {
	return time.Time(ldt).Month()
}

// Day returns the day of the LocalDateTime.
func (ldt LocalDateTime) Day() int {
	return time.Time(ldt).Day()
}

// Hour returns the hour of the LocalDateTime.
func (ldt LocalDateTime) Hour() int {
	return time.Time(ldt).Hour()
}

// Minute returns the minute of the LocalDateTime.
func (ldt LocalDateTime) Minute() int {
	return time.Time(ldt).Minute()
}

// Second returns the second of the LocalDateTime.
func (ldt LocalDateTime) Second() int {
	return time.Time(ldt).Second()
}

// Nanosecond returns the nanosecond of the LocalDateTime.
func (ldt LocalDateTime) Nanosecond() int {
	return time.Time(ldt).Nanosecond()
}
