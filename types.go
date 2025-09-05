package jscal

import (
	"fmt"
	"strings"
	"time"
)

// Common types for JSCalendar implementation

// Participant represents an event participant
type Participant struct {
	Name               *string            `json:"name,omitempty"`
	Email              *string            `json:"email,omitempty"`
	Kind               *string            `json:"kind,omitempty"` // individual, group, resource, location
	Roles              map[string]bool    `json:"roles,omitempty"` // owner, attendee, optional, informational, chair, contact
	LocationId         *string            `json:"locationId,omitempty"`
	Language           *string            `json:"language,omitempty"`
	ParticipationStatus *string           `json:"participationStatus,omitempty"` // needs-action, accepted, declined, tentative, delegated
	ParticipationComment *string          `json:"participationComment,omitempty"`
	ExpectReply        *bool              `json:"expectReply,omitempty"`
	ScheduleAgent      *string            `json:"scheduleAgent,omitempty"` // server, client, none
	ScheduleForceSend  *bool              `json:"scheduleForceSend,omitempty"`
	ScheduleSequence   *int               `json:"scheduleSequence,omitempty"`
	ScheduleStatus     []string           `json:"scheduleStatus,omitempty"`
	ScheduleUpdated    *time.Time         `json:"scheduleUpdated,omitempty"`
	InvitedBy          *string            `json:"invitedBy,omitempty"`
	DelegatedTo        map[string]bool    `json:"delegatedTo,omitempty"`
	DelegatedFrom      map[string]bool    `json:"delegatedFrom,omitempty"`
	MemberOf           map[string]bool    `json:"memberOf,omitempty"`
	Links              map[string]*Link   `json:"links,omitempty"`
	Progress           *string            `json:"progress,omitempty"` // needs-action, in-process, completed, failed, cancelled
	ProgressUpdated    *time.Time         `json:"progressUpdated,omitempty"`
	PercentComplete    *int               `json:"percentComplete,omitempty"`
}

// Location represents a physical or virtual location
type Location struct {
	Type        *string            `json:"@type,omitempty"`
	Name        *string            `json:"name,omitempty"`
	Description *string            `json:"description,omitempty"`
	LocationTypes map[string]bool  `json:"locationTypes,omitempty"`
	RelativeTo  *string            `json:"relativeTo,omitempty"`
	TimeZone    *string            `json:"timeZone,omitempty"`
	Coordinates *string            `json:"coordinates,omitempty"` // geo: URI
	Links       map[string]*Link   `json:"links,omitempty"`
}

// VirtualLocation represents a virtual meeting location
type VirtualLocation struct {
	Type        string             `json:"@type"`
	Name        *string            `json:"name,omitempty"`
	Description *string            `json:"description,omitempty"`
	URI         string             `json:"uri"`
	Features    []string           `json:"features,omitempty"` // audio, video, chat, screen, phone
}

// Link represents a URI with metadata
type Link struct {
	Type         *string            `json:"@type,omitempty"`
	Href         string             `json:"href"`
	Cid          *string            `json:"cid,omitempty"`
	ContentType  *string            `json:"contentType,omitempty"`
	Size         *int               `json:"size,omitempty"`
	Rel          *string            `json:"rel,omitempty"`
	Display      *string            `json:"display,omitempty"` // badge, graphic, fullsize, thumbnail
	Title        *string            `json:"title,omitempty"`
}

// RecurrenceRule represents recurrence rules
type RecurrenceRule struct {
	Type        string      `json:"@type"`
	Frequency   string      `json:"frequency"` // yearly, monthly, weekly, daily, hourly, minutely, secondly
	Interval    *int        `json:"interval,omitempty"`
	RScale      *string     `json:"rscale,omitempty"`
	Skip        *string     `json:"skip,omitempty"` // omit, backward, forward
	FirstDayOfWeek *int     `json:"firstDayOfWeek,omitempty"` // 0=Monday, 1=Tuesday, etc.
	ByDay       []NDay      `json:"byDay,omitempty"`
	ByMonthDay  []int       `json:"byMonthDay,omitempty"`
	ByMonth     []string    `json:"byMonth,omitempty"`
	ByYearDay   []int       `json:"byYearDay,omitempty"`
	ByWeekNo    []int       `json:"byWeekNo,omitempty"`
	ByHour      []int       `json:"byHour,omitempty"`
	ByMinute    []int       `json:"byMinute,omitempty"`
	BySecond    []int       `json:"bySecond,omitempty"`
	BySetPos    []int       `json:"bySetPos,omitempty"`
	Count       *int        `json:"count,omitempty"`
	Until       *time.Time  `json:"until,omitempty"`
}

// NDay represents a day of the week with optional nth occurrence
type NDay struct {
	Day string `json:"day"` // mo, tu, we, th, fr, sa, su
	NthOfPeriod *int `json:"nthOfPeriod,omitempty"`
}

// Alert represents a notification/reminder
type Alert struct {
	Type        string              `json:"@type"`
	Trigger     *OffsetTrigger      `json:"trigger,omitempty"`
	Acknowledged *time.Time         `json:"acknowledged,omitempty"`
	RelatedTo   map[string]*Relation `json:"relatedTo,omitempty"`
	Action      *string             `json:"action,omitempty"` // display, email
}

// OffsetTrigger represents when an alert should fire
type OffsetTrigger struct {
	Type       string  `json:"@type"`
	Offset     string  `json:"offset"` // ISO 8601 duration
	RelativeTo *string `json:"relativeTo,omitempty"` // start, end
}

// Relation represents relationships to other objects
type Relation struct {
	Type string `json:"@type"`
	Relation map[string]bool `json:"relation,omitempty"` // first, next, child, parent
}

// TimeZone represents timezone information
type TimeZone struct {
	Type        string     `json:"@type"`
	TzId        string     `json:"tzId"`
	Updated     *time.Time `json:"updated,omitempty"`
	URL         *string    `json:"url,omitempty"`
	ValidUntil  *time.Time `json:"validUntil,omitempty"`
	Aliases     map[string]bool `json:"aliases,omitempty"`
	Standard    []TimeZoneRule `json:"standard,omitempty"`
	Daylight    []TimeZoneRule `json:"daylight,omitempty"`
}

// TimeZoneRule represents timezone rules
type TimeZoneRule struct {
	Type           string          `json:"@type"`
	Start          time.Time       `json:"start"`
	OffsetFrom     string          `json:"offsetFrom"`
	OffsetTo       string          `json:"offsetTo"`
	RecurrenceRules []RecurrenceRule `json:"recurrenceRules,omitempty"`
	RecurrenceOverrides map[string]*TimeZoneRule `json:"recurrenceOverrides,omitempty"`
	Names          map[string]string `json:"names,omitempty"`
	Comments       []string        `json:"comments,omitempty"`
}

// LocalDateTime represents a date-time without timezone
type LocalDateTime struct {
	Type     string `json:"@type"`
	DateTime string `json:"dateTime"` // ISO 8601 format: YYYY-MM-DDTHH:MM:SS
}

// UTCDateTime represents a UTC date-time
type UTCDateTime struct {
	Type     string `json:"@type"`
	UTCTime  string `json:"utcTime"` // ISO 8601 format: YYYY-MM-DDTHH:MM:SSZ
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
		// Try to parse the numeric part
		numPart := value[:len(value)-2]
		dayPart = value[len(value)-2:]
		
		if num := 0; len(numPart) > 0 {
			if _, err := fmt.Sscanf(numPart, "%d", &num); err == nil {
				nthOfPeriod = &num
			}
		}
	}
	
	return &NDay{
		Day:         FormatDayOfWeek(dayPart),
		NthOfPeriod: nthOfPeriod,
	}, nil
}