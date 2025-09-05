package jscal

import (
	"encoding/json"
	"fmt"
	"time"
)

// Event represents a JSCalendar Event object according to RFC 8984.
// An Event describes a calendar event such as a meeting, appointment,
// or reminder with all associated metadata.
type Event struct {
	// Core metadata (Section 4.1)
	Type     string     `json:"@type"`              // Always "Event"
	UID      string     `json:"uid"`                // Unique identifier
	Created  *time.Time `json:"created,omitempty"`  // Creation timestamp
	Updated  *time.Time `json:"updated,omitempty"`  // Last modification timestamp
	Sequence *int       `json:"sequence,omitempty"` // Revision sequence number
	Method   *string    `json:"method,omitempty"`   // publish, request, reply, etc.
	ProdId   *string    `json:"prodId,omitempty"`   // Product identifier that created this

	// What and Where Properties (Section 4.2)
	Title                  *string                           `json:"title,omitempty"`                  // Event summary/title
	Description            *string                           `json:"description,omitempty"`            // Detailed description
	DescriptionContentType *string                           `json:"descriptionContentType,omitempty"` // MIME type of description
	ShowWithoutTime        *bool                             `json:"showWithoutTime,omitempty"`        // All-day event flag
	Locale                 *string                           `json:"locale,omitempty"`                 // Language tag (RFC 5646)
	Localizations          map[string]map[string]interface{} `json:"localizations,omitempty"`          // Localization patches
	Keywords               map[string]bool                   `json:"keywords,omitempty"`               // Keywords/tags
	Categories             map[string]bool                   `json:"categories,omitempty"`             // Categories
	Color                  *string                           `json:"color,omitempty"`                  // CSS color value

	// Date and Time Properties (Section 5.1)
	Start     *LocalDateTime       `json:"start,omitempty"`     // Start date-time (LocalDateTime)
	Duration  *string              `json:"duration,omitempty"`  // ISO 8601 duration
	TimeZone  *string              `json:"timeZone,omitempty"`  // IANA timezone identifier
	TimeZones map[string]*TimeZone `json:"timeZones,omitempty"` // Custom timezone definitions

	// Recurrence Properties (Section 4.3)
	RecurrenceId            *LocalDateTime                    `json:"recurrenceId,omitempty"`         // Recurrence instance identifier
	RecurrenceIdTimeZone    *string                           `json:"recurrenceIdTimeZone,omitempty"` // TimeZone for recurrenceId (optional per RFC Erratum 6873)
	RecurrenceRules         []RecurrenceRule                  `json:"recurrenceRules,omitempty"`      // Recurrence rules
	RecurrenceOverrides     map[string]map[string]interface{} `json:"recurrenceOverrides,omitempty"`  // Overrides as patches
	ExcludedRecurrenceRules []RecurrenceRule                  `json:"excludedRecurrenceRules,omitempty"`
	Excluded                *bool                             `json:"excluded,omitempty"` // Is this instance excluded?

	// Sharing and Scheduling Properties (Section 4.4)
	Priority       *int                    `json:"priority,omitempty"`       // Priority (0-9)
	FreeBusyStatus *string                 `json:"freeBusyStatus,omitempty"` // free, busy, tentative
	Privacy        *string                 `json:"privacy,omitempty"`        // public, private, secret (see RFC Erratum 6872 for shareable fields)
	ReplyTo        map[string]string       `json:"replyTo,omitempty"`        // Reply methods
	SentBy         *string                 `json:"sentBy,omitempty"`         // Email of sender
	Participants   map[string]*Participant `json:"participants,omitempty"`   // Event participants
	RequestStatus  *string                 `json:"requestStatus,omitempty"`  // Scheduling request status

	// Alerts Properties (Section 4.5)
	UseDefaultAlerts *bool             `json:"useDefaultAlerts,omitempty"` // Use default alert settings
	Alerts           map[string]*Alert `json:"alerts,omitempty"`           // Custom alerts

	// Links and Locations
	Links            map[string]*Link            `json:"links,omitempty"`            // External links
	Locations        map[string]*Location        `json:"locations,omitempty"`        // Physical locations
	VirtualLocations map[string]*VirtualLocation `json:"virtualLocations,omitempty"` // Virtual meeting locations

	// Relationships
	RelatedTo map[string]*Relation `json:"relatedTo,omitempty"` // Related events

	// Event-specific (Section 5.1)
	Status *string `json:"status,omitempty"` // confirmed, tentative, cancelled

	// Free-form properties for extensions
	Extensions map[string]interface{} `json:"-"` // Custom properties

	// Localization
	LocalizedStrings map[string]map[string]string `json:"localizedStrings,omitempty"`
}

// NewEvent creates a new JSCalendar Event with required fields
func NewEvent(uid, title string) *Event {
	now := time.Now().UTC()
	return &Event{
		Type:     "Event",
		UID:      uid,
		Title:    &title,
		Start:    NewLocalDateTime(now),
		Created:  &now,
		Updated:  &now,
		Sequence: Int(0),
	}
}

// JSON returns the Event as JSON bytes
func (e *Event) JSON() ([]byte, error) {
	return json.Marshal(e)
}

// PrettyJSON returns the Event as indented JSON bytes
func (e *Event) PrettyJSON() ([]byte, error) {
	return json.MarshalIndent(e, "", "  ")
}

// Clone creates a deep copy of the Event
func (e *Event) Clone() *Event {
	data, _ := json.Marshal(e)
	var clone Event
	_ = json.Unmarshal(data, &clone)
	return &clone
}

// IsAllDay returns true if this is an all-day event
func (e *Event) IsAllDay() bool {
	return e.ShowWithoutTime != nil && *e.ShowWithoutTime
}

// GetDuration returns the event duration
func (e *Event) GetDuration() (time.Duration, error) {
	if e.Duration == nil {
		return 0, fmt.Errorf("no duration specified")
	}

	// Parse ISO 8601 duration
	return parseISO8601Duration(*e.Duration)
}

// GetEndTime calculates the end time based on start and duration
func (e *Event) GetEndTime() (*time.Time, error) {
	if e.Start == nil {
		return nil, fmt.Errorf("no start time specified")
	}

	// Convert LocalDateTime to time.Time
	startTime := e.Start.Time()

	duration, err := e.GetDuration()
	if err != nil {
		return nil, fmt.Errorf("failed to parse duration: %w", err)
	}

	endTime := startTime.Add(duration)
	return &endTime, nil
}

// AddParticipant adds a participant to the event
func (e *Event) AddParticipant(id string, participant *Participant) {
	if e.Participants == nil {
		e.Participants = make(map[string]*Participant)
	}
	e.Participants[id] = participant
}

// AddLocation adds a location to the event
func (e *Event) AddLocation(id string, location *Location) {
	if e.Locations == nil {
		e.Locations = make(map[string]*Location)
	}
	e.Locations[id] = location
}

// AddVirtualLocation adds a virtual location to the event
func (e *Event) AddVirtualLocation(id string, virtualLocation *VirtualLocation) {
	if e.VirtualLocations == nil {
		e.VirtualLocations = make(map[string]*VirtualLocation)
	}
	e.VirtualLocations[id] = virtualLocation
}

// AddAlert adds an alert to the event
func (e *Event) AddAlert(id string, alert *Alert) {
	if e.Alerts == nil {
		e.Alerts = make(map[string]*Alert)
	}
	e.Alerts[id] = alert
}

// AddCategory adds a category to the event
func (e *Event) AddCategory(category string) {
	if e.Categories == nil {
		e.Categories = make(map[string]bool)
	}
	e.Categories[category] = true
}

// AddKeyword adds a keyword to the event
func (e *Event) AddKeyword(keyword string) {
	if e.Keywords == nil {
		e.Keywords = make(map[string]bool)
	}
	e.Keywords[keyword] = true
}

// AddLink adds a link to the event
func (e *Event) AddLink(id string, link *Link) {
	if e.Links == nil {
		e.Links = make(map[string]*Link)
	}
	e.Links[id] = link
}

// SetRecurrence sets the recurrence rules for the event
func (e *Event) SetRecurrence(rules []RecurrenceRule) {
	e.RecurrenceRules = rules
}

// IsRecurring returns true if the event has recurrence rules
func (e *Event) IsRecurring() bool {
	return len(e.RecurrenceRules) > 0
}

// Touch updates the Updated timestamp to now
func (e *Event) Touch() {
	now := time.Now().UTC()
	e.Updated = &now
	if e.Sequence != nil {
		*e.Sequence++
	} else {
		e.Sequence = Int(1)
	}
}

// GetUID returns the event's UID (implements CalendarObject)
func (e *Event) GetUID() string {
	return e.UID
}

// GetType returns the event's type (implements CalendarObject)
func (e *Event) GetType() string {
	return e.Type
}