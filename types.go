package jscal

import (
	"fmt"
	"strings"
	"time"
)

// Common types for JSCalendar implementation

// Participant represents an event participant
type Participant struct {
	Type                 *string           `json:"@type,omitempty"` // Always "Participant"
	Name                 *string           `json:"name,omitempty"`
	Email                *string           `json:"email,omitempty"`
	SendTo               map[string]string `json:"sendTo,omitempty"` // Contact methods (e.g., {"imip": "mailto:..."})
	Kind                 *string           `json:"kind,omitempty"`   // individual, group, resource, location
	Roles                map[string]bool   `json:"roles,omitempty"`  // owner, attendee, optional, informational, chair, contact
	LocationId           *string           `json:"locationId,omitempty"`
	Language             *string           `json:"language,omitempty"`
	ParticipationStatus  *string           `json:"participationStatus,omitempty"` // needs-action, accepted, declined, tentative, delegated
	ParticipationComment *string           `json:"participationComment,omitempty"`
	ExpectReply          *bool             `json:"expectReply,omitempty"`
	ScheduleAgent        *string           `json:"scheduleAgent,omitempty"` // server, client, none
	ScheduleForceSend    *bool             `json:"scheduleForceSend,omitempty"`
	ScheduleSequence     *int              `json:"scheduleSequence,omitempty"`
	ScheduleStatus       []string          `json:"scheduleStatus,omitempty"`
	ScheduleUpdated      *time.Time        `json:"scheduleUpdated,omitempty"`
	SentBy               *string           `json:"sentBy,omitempty"` // Email address of sender
	InvitedBy            *string           `json:"invitedBy,omitempty"`
	DelegatedTo          map[string]bool   `json:"delegatedTo,omitempty"`
	DelegatedFrom        map[string]bool   `json:"delegatedFrom,omitempty"`
	MemberOf             map[string]bool   `json:"memberOf,omitempty"`
	Links                map[string]*Link  `json:"links,omitempty"`
}

// Location represents a physical or virtual location
type Location struct {
	Type          *string          `json:"@type,omitempty"`
	Name          *string          `json:"name,omitempty"`
	Description   *string          `json:"description,omitempty"`
	LocationTypes map[string]bool  `json:"locationTypes,omitempty"`
	RelativeTo    *string          `json:"relativeTo,omitempty"`
	TimeZone      *string          `json:"timeZone,omitempty"`
	Coordinates   *string          `json:"coordinates,omitempty"` // geo: URI
	Links         map[string]*Link `json:"links,omitempty"`
	Rel           *string          `json:"rel,omitempty"` // start, end
	Title         *string          `json:"title,omitempty"`
}

// VirtualLocation represents a virtual meeting location
type VirtualLocation struct {
	Type        string          `json:"@type"`
	Name        *string         `json:"name,omitempty"`
	Description *string         `json:"description,omitempty"`
	URI         string          `json:"uri"`
	Features    map[string]bool `json:"features,omitempty"` // RFC 8984: String[Boolean] - predefined: audio, chat, feed, moderator, phone, screen, video (extensible)
}

// Link represents a URI with metadata
type Link struct {
	Type        *string `json:"@type,omitempty"`
	Href        string  `json:"href"`
	Cid         *string `json:"cid,omitempty"`
	ContentType *string `json:"contentType,omitempty"`
	Size        *int    `json:"size,omitempty"`
	Rel         *string `json:"rel,omitempty"`
	Display     *string `json:"display,omitempty"` // badge, graphic, fullsize, thumbnail
	Title       *string `json:"title,omitempty"`
}

// RecurrenceRule represents recurrence rules
type RecurrenceRule struct {
	Type           string         `json:"@type"`
	Frequency      string         `json:"frequency"` // yearly, monthly, weekly, daily, hourly, minutely, secondly
	Interval       *int           `json:"interval,omitempty"`
	RScale         *string        `json:"rscale,omitempty"`
	Skip           *string        `json:"skip,omitempty"`           // omit, backward, forward
	FirstDayOfWeek *int           `json:"firstDayOfWeek,omitempty"` // 0=Monday, 1=Tuesday, etc.
	ByDay          []NDay         `json:"byDay,omitempty"`
	ByMonthDay     []int          `json:"byMonthDay,omitempty"`
	ByMonth        []string       `json:"byMonth,omitempty"`
	ByYearDay      []int          `json:"byYearDay,omitempty"`
	ByWeekNo       []int          `json:"byWeekNo,omitempty"`
	ByHour         []int          `json:"byHour,omitempty"`
	ByMinute       []int          `json:"byMinute,omitempty"`
	BySecond       []int          `json:"bySecond,omitempty"`
	BySetPos       []int          `json:"bySetPos,omitempty"`
	Count          *int           `json:"count,omitempty"`
	Until          *LocalDateTime `json:"until,omitempty"`
}

// NDay represents a day of the week with optional nth occurrence
type NDay struct {
	Day         string `json:"day"` // mo, tu, we, th, fr, sa, su
	NthOfPeriod *int   `json:"nthOfPeriod,omitempty"`
}

// Alert represents a notification/reminder
type Alert struct {
	Type         string               `json:"@type"`
	Trigger      *OffsetTrigger       `json:"trigger,omitempty"`
	Acknowledged *time.Time           `json:"acknowledged,omitempty"`
	RelatedTo    map[string]*Relation `json:"relatedTo,omitempty"`
	Action       *string              `json:"action,omitempty"` // display, email
}

// OffsetTrigger represents when an alert should fire
type OffsetTrigger struct {
	Type       string  `json:"@type"`
	Offset     string  `json:"offset"`               // ISO 8601 duration
	RelativeTo *string `json:"relativeTo,omitempty"` // start, end
}

// Relation represents relationships to other objects
type Relation struct {
	Type     string          `json:"@type"`
	Relation map[string]bool `json:"relation,omitempty"` // first, next, child, parent
}

// UTCDateTime represents a UTC date-time
type UTCDateTime struct {
	Type    string `json:"@type"`
	UTCTime string `json:"utcTime"` // ISO 8601 format: YYYY-MM-DDTHH:MM:SSZ
}

// Helper functions

// NewParticipant creates a new participant with basic info
func NewParticipant(name, email string) *Participant {
	return &Participant{
		Name:  &name,
		Email: &email,
		Roles: map[string]bool{"attendee": true},
	}
}

// NewLocation creates a new location with basic info
func NewLocation(name string) *Location {
	return &Location{
		Type: String("Location"),
		Name: &name,
	}
}

// NewVirtualLocation creates a new virtual location
func NewVirtualLocation(name, uri string) *VirtualLocation {
	return &VirtualLocation{
		Type: "VirtualLocation",
		Name: &name,
		URI:  uri,
	}
}

// NewLink creates a new link
func NewLink(href string) *Link {
	return &Link{
		Href: href,
	}
}

// NewRecurrenceRule creates a new RecurrenceRule with the given frequency
func NewRecurrenceRule(frequency string) *RecurrenceRule {
	return &RecurrenceRule{
		Type:      "RecurrenceRule",
		Frequency: frequency,
	}
}

// String returns a pointer to the string value
func String(s string) *string {
	return &s
}

// Int returns a pointer to the int value
func Int(i int) *int {
	return &i
}

// Bool returns a pointer to the bool value
func Bool(b bool) *bool {
	return &b
}

func TimePtr(t time.Time) *time.Time {
	return &t
}

// FormatDayOfWeek converts a day name to JSCalendar format
func FormatDayOfWeek(day string) string {
	switch strings.ToUpper(day) {
	case "MONDAY", "MO":
		return "mo"
	case "TUESDAY", "TU":
		return "tu"
	case "WEDNESDAY", "WE":
		return "we"
	case "THURSDAY", "TH":
		return "th"
	case "FRIDAY", "FR":
		return "fr"
	case "SATURDAY", "SA":
		return "sa"
	case "SUNDAY", "SU":
		return "su"
	default:
		return strings.ToLower(day)
	}
}

// ParseNDay parses an RRULE BYDAY value into NDay
func ParseNDay(value string) (*NDay, error) {
	if len(value) < 2 {
		return nil, fmt.Errorf("invalid day value: %s", value)
	}

	// Check if there's a numeric prefix
	var nthOfPeriod *int
	dayPart := value

	if len(value) > 2 {
		// Extract numeric and day parts
		numPart := value[:len(value)-2]
		dayPart = value[len(value)-2:]

		if len(numPart) > 0 {
			// Validate that numPart contains only integer characters
			// Allow optional negative sign at start
			validChars := true
			for i, r := range numPart {
				if i == 0 && r == '-' {
					continue // Allow negative sign at start
				}
				if r < '0' || r > '9' {
					validChars = false
					break
				}
			}

			if !validChars {
				return nil, fmt.Errorf("invalid occurrence format: numeric prefix must be integer, got %s", numPart)
			}

			// Parse the validated integer
			var num int
			if n, err := fmt.Sscanf(numPart, "%d", &num); n == 1 && err == nil {
				nthOfPeriod = &num
			} else {
				return nil, fmt.Errorf("invalid occurrence format: failed to parse %s as integer", numPart)
			}
		}
	}

	// Validate the day part - must be exactly 2 characters
	if len(dayPart) != 2 {
		return nil, fmt.Errorf("invalid day format: must be 2-character day code, got %s", dayPart)
	}

	// Validate the day part
	formatted := FormatDayOfWeek(dayPart)
	validDays := map[string]bool{
		"mo": true, "tu": true, "we": true, "th": true,
		"fr": true, "sa": true, "su": true,
	}
	if !validDays[formatted] {
		return nil, fmt.Errorf("invalid day: %s", dayPart)
	}

	return &NDay{
		Day:         formatted,
		NthOfPeriod: nthOfPeriod,
	}, nil
}
