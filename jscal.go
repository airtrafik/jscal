// Package jscal implements RFC 8984 JSCalendar specification for Go.
// 
// JSCalendar is a modern JSON-based calendar data format that provides
// a cleaner alternative to iCalendar (RFC 5545) with better support for
// modern calendar applications.
//
// This package provides:
//   - Complete JSCalendar Event type implementation
//   - JSON marshaling/unmarshaling with validation
//   - RFC 8984 compliance validation
//
// Basic usage:
//
//	// Parse JSCalendar JSON
//	event, err := jscal.Parse(jsonData)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Parse multiple events
//	events, err := jscal.ParseAll(jsonData)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Validate JSCalendar compliance
//	if err := event.Validate(); err != nil {
//		log.Printf("Validation error: %v", err)
//	}
package jscal

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Event represents a JSCalendar Event object according to RFC 8984.
// An Event describes a calendar event such as a meeting, appointment,
// or reminder with all associated metadata.
type Event struct {
	// Core metadata
	Type        string     `json:"@type"`                  // Always "Event"
	UID         string     `json:"uid"`                     // Unique identifier
	Created     *time.Time `json:"created,omitempty"`       // Creation timestamp
	Updated     *time.Time `json:"updated,omitempty"`       // Last modification timestamp
	Sequence    *int       `json:"sequence,omitempty"`      // Revision sequence number
	
	// Human-readable content
	Title       *string    `json:"title,omitempty"`         // Event summary/title
	Description *string    `json:"description,omitempty"`   // Detailed description
	
	// Scheduling
	Start       *time.Time `json:"start,omitempty"`         // Start date-time
	Duration    *string    `json:"duration,omitempty"`      // ISO 8601 duration (e.g. "PT1H30M")
	TimeZone    *string    `json:"timeZone,omitempty"`      // IANA timezone identifier
	ShowWithoutTime *bool  `json:"showWithoutTime,omitempty"` // All-day event flag
	
	// Recurrence
	RecurrenceRules     []RecurrenceRule      `json:"recurrenceRules,omitempty"`
	RecurrenceOverrides map[string]*Event     `json:"recurrenceOverrides,omitempty"`
	ExcludedRecurrenceRules []RecurrenceRule `json:"excludedRecurrenceRules,omitempty"`
	
	// People and resources
	Participants map[string]*Participant `json:"participants,omitempty"`
	
	// Locations
	Locations        map[string]*Location        `json:"locations,omitempty"`
	VirtualLocations map[string]*VirtualLocation `json:"virtualLocations,omitempty"`
	
	// Alerts and notifications
	Alerts map[string]*Alert `json:"alerts,omitempty"`
	
	// Categorization and organization
	Categories   map[string]bool   `json:"categories,omitempty"`
	Keywords     map[string]bool   `json:"keywords,omitempty"`
	Color        *string           `json:"color,omitempty"`     // CSS color value
	
	// Status and workflow
	Status       *string    `json:"status,omitempty"`        // confirmed, tentative, cancelled
	FreeBusyStatus *string  `json:"freeBusyStatus,omitempty"` // free, busy, tentative
	Privacy      *string    `json:"privacy,omitempty"`       // public, private, secret
	
	// Collaboration
	Method       *string    `json:"method,omitempty"`        // publish, request, reply, etc.
	ProdId       *string    `json:"prodId,omitempty"`        // Product identifier that created this
	
	// Links and attachments
	Links        map[string]*Link  `json:"links,omitempty"`
	
	// Relationships to other events
	RelatedTo    map[string]*Relation `json:"relatedTo,omitempty"`
	
	// Free-form properties for extensions
	Extensions   map[string]interface{} `json:"-"` // Custom properties
	
	// Localization
	LocalizedStrings map[string]map[string]string `json:"localizedStrings,omitempty"`
}

// NewEvent creates a new JSCalendar Event with required fields
func NewEvent(uid, title string) *Event {
	now := time.Now().UTC()
	return &Event{
		Type:    "Event",
		UID:     uid,
		Title:   &title,
		Created: &now,
		Updated: &now,
		Sequence: Int(0),
	}
}

// Parse parses JSCalendar JSON data into an Event
func Parse(data []byte) (*Event, error) {
	var event Event
	if err := json.Unmarshal(data, &event); err != nil {
		return nil, fmt.Errorf("failed to parse JSCalendar JSON: %w", err)
	}
	
	// Validate the parsed event
	if err := event.Validate(); err != nil {
		return nil, fmt.Errorf("parsed JSCalendar is invalid: %w", err)
	}
	
	return &event, nil
}

// ParseAll parses multiple JSCalendar events from JSON array
func ParseAll(data []byte) ([]*Event, error) {
	var events []*Event
	if err := json.Unmarshal(data, &events); err != nil {
		return nil, fmt.Errorf("failed to parse JSCalendar JSON array: %w", err)
	}
	
	// Validate each event
	for i, event := range events {
		if err := event.Validate(); err != nil {
			return nil, fmt.Errorf("event at index %d is invalid: %w", i, err)
		}
	}
	
	return events, nil
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

// parseISO8601Duration parses an ISO 8601 duration string to time.Duration
func parseISO8601Duration(duration string) (time.Duration, error) {
	// Simple parser for ISO 8601 durations like PT1H, P1D, PT30M
	if !strings.HasPrefix(duration, "P") {
		return 0, fmt.Errorf("invalid ISO 8601 duration: must start with P")
	}
	
	duration = duration[1:] // Remove P
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
	
	// Parse date part
	if datePart != "" {
		// Simple parsing for days
		if strings.HasSuffix(datePart, "D") {
			days := 0
			_, _ = fmt.Sscanf(datePart, "%dD", &days)
			result += time.Duration(days) * 24 * time.Hour
		}
	}
	
	// Parse time part
	if timePart != "" {
		// Parse hours
		if idx := strings.Index(timePart, "H"); idx > 0 {
			hours := 0
			_, _ = fmt.Sscanf(timePart[:idx], "%d", &hours)
			result += time.Duration(hours) * time.Hour
			timePart = timePart[idx+1:]
		}
		
		// Parse minutes
		if idx := strings.Index(timePart, "M"); idx > 0 {
			minutes := 0
			_, _ = fmt.Sscanf(timePart[:idx], "%d", &minutes)
			result += time.Duration(minutes) * time.Minute
			timePart = timePart[idx+1:]
		}
		
		// Parse seconds
		if idx := strings.Index(timePart, "S"); idx > 0 {
			seconds := 0.0
			_, _ = fmt.Sscanf(timePart[:idx], "%f", &seconds)
			result += time.Duration(seconds * float64(time.Second))
		}
	}
	
	return result, nil
}

// GetEndTime calculates the end time based on start and duration
func (e *Event) GetEndTime() (*time.Time, error) {
	if e.Start == nil {
		return nil, fmt.Errorf("no start time specified")
	}
	
	duration, err := e.GetDuration()
	if err != nil {
		return nil, fmt.Errorf("failed to parse duration: %w", err)
	}
	
	endTime := e.Start.Add(duration)
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