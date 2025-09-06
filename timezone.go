package jscal

import "time"

// TimeZone represents a custom timezone definition according to RFC 8984 Section 4.7.2
type TimeZone struct {
	// Type identifier
	Type *string `json:"@type,omitempty"` // Should be "TimeZone"

	// Timezone identifier (e.g., "America/New_York")
	TzId string `json:"tzId"`

	// Last modified time of the timezone definition
	Updated *time.Time `json:"updated,omitempty"`

	// URL to the timezone definition
	URL *string `json:"url,omitempty"`

	// Validity period start
	ValidUntil *time.Time `json:"validUntil,omitempty"`

	// Aliases for this timezone
	Aliases []string `json:"aliases,omitempty"`

	// Standard offset from UTC in the format ±HH:MM
	StandardOffset *string `json:"standardOffset,omitempty"`

	// Daylight saving offset from UTC in the format ±HH:MM
	DaylightOffset *string `json:"daylightOffset,omitempty"`

	// Standard timezone components (for complex definitions)
	Standard []TimeZoneRule `json:"standard,omitempty"`

	// Daylight timezone components (for complex definitions)
	Daylight []TimeZoneRule `json:"daylight,omitempty"`
}

// TimeZoneRule represents a timezone transition rule
type TimeZoneRule struct {
	// Start date-time for this rule
	Start *LocalDateTime `json:"start,omitempty"`

	// Offset from UTC in format ±HH:MM
	OffsetFrom string `json:"offsetFrom"`
	OffsetTo   string `json:"offsetTo"`

	// Recurrence rule for when this transition occurs
	RecurrenceRules []RecurrenceRule `json:"recurrenceRules,omitempty"`

	// Timezone names
	Names map[string]string `json:"names,omitempty"`

	// Comments about this rule
	Comments []string `json:"comments,omitempty"`
}

// NewTimeZone creates a new TimeZone with the given ID
func NewTimeZone(tzId string) *TimeZone {
	tzType := "TimeZone"
	return &TimeZone{
		Type: &tzType,
		TzId: tzId,
	}
}
