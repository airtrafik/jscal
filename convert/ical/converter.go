// Package ical provides conversion between iCalendar (RFC 5545) and JSCalendar formats.
package ical

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/airtrafik/jscal"
	"github.com/airtrafik/jscal/convert"
	ics "github.com/arran4/golang-ical"
)

// Converter handles iCalendar <-> JSCalendar conversions using golang-ical library
type Converter struct{}

// Ensure Converter implements the convert.Converter interface
var _ convert.Converter = (*Converter)(nil)

// New creates a new iCalendar converter
func New() *Converter {
	return &Converter{}
}

// Parse converts iCalendar data to a single JSCalendar event
func (c *Converter) Parse(data []byte) (*jscal.Event, error) {
	events, err := c.ParseAll(data)
	if err != nil {
		return nil, err
	}

	if len(events) == 0 {
		return nil, fmt.Errorf("no events found in iCalendar data")
	}

	if len(events) > 1 {
		return nil, fmt.Errorf("multiple events found, use ParseAll instead")
	}

	return events[0], nil
}

// Format converts a single JSCalendar event to iCalendar format
func (c *Converter) Format(event *jscal.Event) ([]byte, error) {
	return c.FormatAll([]*jscal.Event{event})
}

// ParseAll converts iCalendar data to JSCalendar events
func (c *Converter) ParseAll(data []byte) ([]*jscal.Event, error) {
	cal, err := ics.ParseCalendar(strings.NewReader(string(data)))
	if err != nil {
		return nil, fmt.Errorf("failed to parse iCalendar: %w", err)
	}

	var events []*jscal.Event

	for _, vevent := range cal.Events() {
		event, err := convertICalEventToJSCal(vevent)
		if err != nil {
			return nil, fmt.Errorf("failed to convert event: %w", err)
		}
		events = append(events, event)
	}

	return events, nil
}

// FormatAll converts JSCalendar events to iCalendar format
func (c *Converter) FormatAll(events []*jscal.Event) ([]byte, error) {
	if len(events) == 0 {
		return nil, fmt.Errorf("no events to convert")
	}

	cal := ics.NewCalendar()
	cal.SetProductId("-//AirTrafik//JSCal Go Library//EN")
	cal.SetVersion("2.0")

	for _, event := range events {
		vevent, err := convertJSCalEventToICal(event)
		if err != nil {
			return nil, fmt.Errorf("failed to convert event %s: %w", event.UID, err)
		}
		cal.AddVEvent(vevent)
	}

	return []byte(cal.Serialize()), nil
}

// Detect returns true if the data appears to be iCalendar format
func (c *Converter) Detect(data []byte) bool {
	dataStr := strings.TrimSpace(string(data))

	// Check for iCalendar header
	if strings.HasPrefix(dataStr, "BEGIN:VCALENDAR") {
		return true
	}

	// Check for common iCalendar patterns
	patterns := []string{
		"BEGIN:VEVENT",
		"DTSTART:",
		"DTEND:",
		"SUMMARY:",
		"UID:",
	}

	found := 0
	for _, pattern := range patterns {
		if strings.Contains(dataStr, pattern) {
			found++
		}
	}

	// If we found at least 3 patterns, it's likely iCalendar
	return found >= 3
}

// convertICalEventToJSCal converts an iCalendar event to JSCalendar
func convertICalEventToJSCal(vevent *ics.VEvent) (*jscal.Event, error) {
	event := &jscal.Event{
		Type: "Event",
	}

	// UID (required)
	if uid := vevent.Id(); uid != "" {
		event.UID = uid
	} else {
		return nil, fmt.Errorf("event missing UID")
	}

	// Title/Summary
	if prop := vevent.GetProperty(ics.ComponentPropertySummary); prop != nil {
		title := prop.Value
		event.Title = &title
	}

	// Description
	if prop := vevent.GetProperty(ics.ComponentPropertyDescription); prop != nil {
		desc := prop.Value
		// Handle escaped characters
		desc = strings.ReplaceAll(desc, "\\n", "\n")
		desc = strings.ReplaceAll(desc, "\\,", ",")
		desc = strings.ReplaceAll(desc, "\\;", ";")
		desc = strings.ReplaceAll(desc, "\\\\", "\\")
		event.Description = &desc
	}

	// Start time
	if dtstart := vevent.GetProperty(ics.ComponentPropertyDtStart); dtstart != nil {
		startTime, isAllDay, timezone := parseICalDateTime(dtstart)
		if !startTime.IsZero() {
			event.Start = jscal.NewLocalDateTime(startTime)
			if isAllDay {
				event.ShowWithoutTime = jscal.Bool(true)
			}
			if timezone != "" && timezone != "UTC" {
				event.TimeZone = &timezone
			}
		}
	}

	// End time or Duration
	if dtend := vevent.GetProperty(ics.ComponentPropertyDtEnd); dtend != nil {
		endTime, _, _ := parseICalDateTime(dtend)
		if !endTime.IsZero() && event.Start != nil {
			// Convert LocalDateTime to time.Time for duration calculation
			startTime := event.Start.Time()
			duration := endTime.Sub(startTime)
			durationStr := formatISO8601Duration(duration)
			event.Duration = &durationStr
		}
	} else if dur := vevent.GetProperty(ics.ComponentPropertyDuration); dur != nil {
		// Parse and convert duration
		duration := parseICalDuration(dur.Value)
		durationStr := formatISO8601Duration(duration)
		event.Duration = &durationStr
	}

	// Created
	if created := vevent.GetProperty(ics.ComponentPropertyCreated); created != nil {
		createdTime, _, _ := parseICalDateTime(created)
		if !createdTime.IsZero() {
			event.Created = &createdTime
		}
	}

	// Last Modified
	if modified := vevent.GetProperty(ics.ComponentPropertyLastModified); modified != nil {
		modifiedTime, _, _ := parseICalDateTime(modified)
		if !modifiedTime.IsZero() {
			event.Updated = &modifiedTime
		}
	}

	// Sequence
	if seq := vevent.GetProperty(ics.ComponentPropertySequence); seq != nil {
		if seqNum := parseInt(seq.Value); seqNum >= 0 {
			event.Sequence = &seqNum
		}
	}

	// Status
	if status := vevent.GetProperty(ics.ComponentPropertyStatus); status != nil {
		statusLower := strings.ToLower(status.Value)
		event.Status = &statusLower
	}

	// Categories
	if categories := vevent.GetProperty(ics.ComponentPropertyCategories); categories != nil {
		cats := strings.Split(categories.Value, ",")
		event.Categories = make(map[string]bool)
		for _, cat := range cats {
			event.Categories[strings.TrimSpace(cat)] = true
		}
	}

	// Location
	if location := vevent.GetProperty(ics.ComponentPropertyLocation); location != nil {
		if event.Locations == nil {
			event.Locations = make(map[string]*jscal.Location)
		}
		loc := jscal.NewLocation(location.Value)
		event.Locations["1"] = loc
	}

	// Transparency -> FreeBusyStatus
	if transp := vevent.GetProperty(ics.ComponentPropertyTransp); transp != nil {
		var freeBusy string
		if strings.ToUpper(transp.Value) == "TRANSPARENT" {
			freeBusy = "free"
		} else {
			freeBusy = "busy"
		}
		event.FreeBusyStatus = &freeBusy
	}

	// Class -> Privacy
	if class := vevent.GetProperty(ics.ComponentPropertyClass); class != nil {
		privacy := strings.ToLower(class.Value)
		if privacy == "confidential" {
			privacy = "private"
		}
		event.Privacy = &privacy
	}

	// URL -> Links
	if url := vevent.GetProperty(ics.ComponentPropertyUrl); url != nil {
		if event.Links == nil {
			event.Links = make(map[string]*jscal.Link)
		}
		link := jscal.NewLink(url.Value)
		event.Links["1"] = link
	}

	// Process Attendees and Organizer
	processParticipants(vevent, event)

	// Process Recurrence Rules
	processRecurrenceRules(vevent, event)

	return event, nil
}

// convertJSCalEventToICal converts a JSCalendar event to iCalendar
func convertJSCalEventToICal(event *jscal.Event) (*ics.VEvent, error) {
	vevent := ics.NewEvent(event.UID)

	// Set timestamp
	vevent.SetProperty(ics.ComponentPropertyDtstamp, time.Now().Format("20060102T150405Z"))

	// Title/Summary
	if event.Title != nil {
		vevent.SetSummary(*event.Title)
	}

	// Description
	if event.Description != nil {
		desc := *event.Description
		// Escape special characters for iCalendar
		desc = strings.ReplaceAll(desc, "\\", "\\\\")
		desc = strings.ReplaceAll(desc, ";", "\\;")
		desc = strings.ReplaceAll(desc, ",", "\\,")
		desc = strings.ReplaceAll(desc, "\n", "\\n")
		vevent.SetDescription(desc)
	}

	// Start time
	if event.Start != nil {
		// Convert LocalDateTime to time.Time
		startTime := event.Start.Time()

		if event.IsAllDay() {
			// All-day event - use date format
			vevent.SetAllDayStartAt(startTime)
		} else {
			// Timed event
			vevent.SetStartAt(startTime)
		}
	}

	// Duration (prefer DURATION property over DTEND for better round-trip)
	if event.Duration != nil && *event.Duration != "" {
		vevent.AddProperty("DURATION", *event.Duration)
	}

	// Created
	if event.Created != nil {
		vevent.SetProperty(ics.ComponentPropertyCreated, event.Created.Format("20060102T150405Z"))
	}

	// Last Modified
	if event.Updated != nil {
		vevent.SetProperty(ics.ComponentPropertyLastModified, event.Updated.Format("20060102T150405Z"))
	}

	// Sequence
	if event.Sequence != nil {
		vevent.SetSequence(*event.Sequence)
	}

	// Status
	if event.Status != nil {
		vevent.SetStatus(ics.ObjectStatus(strings.ToUpper(*event.Status)))
	}

	// Categories (sorted for consistency)
	if len(event.Categories) > 0 {
		var cats []string
		for cat := range event.Categories {
			cats = append(cats, cat)
		}
		sort.Strings(cats)
		prop := ics.IANAProperty{
			BaseProperty: ics.BaseProperty{
				IANAToken: string(ics.ComponentPropertyCategories),
				Value:     strings.Join(cats, ","),
			},
		}
		vevent.Properties = append(vevent.Properties, prop)
	}

	// Location
	for _, location := range event.Locations {
		if location.Name != nil {
			vevent.SetLocation(*location.Name)
			break // Only use first location
		}
	}

	// FreeBusyStatus -> Transparency
	if event.FreeBusyStatus != nil {
		if *event.FreeBusyStatus == "free" {
			vevent.AddProperty(ics.ComponentPropertyTransp, "TRANSPARENT")
		} else {
			vevent.AddProperty(ics.ComponentPropertyTransp, "OPAQUE")
		}
	}

	// Privacy -> Class
	if event.Privacy != nil {
		class := strings.ToUpper(*event.Privacy)
		if class == "PRIVATE" {
			class = "CONFIDENTIAL"
		}
		vevent.AddProperty(ics.ComponentPropertyClass, class)
	}

	// URL
	for _, link := range event.Links {
		vevent.SetURL(link.Href)
		break // Only use first link
	}

	// Convert participants
	convertParticipants(event, vevent)

	// Convert recurrence rules
	convertRecurrenceRules(event, vevent)

	return vevent, nil
}

// Helper functions

func parseICalDateTime(prop *ics.IANAProperty) (time.Time, bool, string) {
	value := prop.Value
	params := prop.ICalParameters

	var timezone string
	isAllDay := false

	// Check for VALUE=DATE (all-day event)
	if params != nil {
		if val, ok := params["VALUE"]; ok && len(val) > 0 && val[0] == "DATE" {
			isAllDay = true
		}
		if tzid, ok := params["TZID"]; ok && len(tzid) > 0 {
			timezone = tzid[0]
		}
	}

	// Remove Z suffix for UTC times
	if strings.HasSuffix(value, "Z") {
		value = value[:len(value)-1]
		timezone = "UTC"
	}

	// Try different date-time formats
	formats := []string{
		"20060102T150405",     // Basic format
		"20060102T150405Z",    // Basic UTC format
		"2006-01-02T15:04:05", // Extended format
		"20060102",            // Date only
		"2006-01-02",          // Date only with hyphens
	}

	for _, format := range formats {
		if t, err := time.Parse(format, value); err == nil {
			// Apply timezone if specified
			if timezone != "" && timezone != "UTC" {
				if loc, err := time.LoadLocation(timezone); err == nil {
					t = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), loc)
				}
			}
			return t, isAllDay, timezone
		}
	}

	return time.Time{}, isAllDay, timezone
}

func formatISO8601Duration(d time.Duration) string {
	if d == 0 {
		return "PT0S"
	}

	result := "P"

	days := int(d.Hours() / 24)
	if days > 0 {
		result += fmt.Sprintf("%dD", days)
		d -= time.Duration(days) * 24 * time.Hour
	}

	if d > 0 {
		result += "T"

		hours := int(d.Hours())
		if hours > 0 {
			result += fmt.Sprintf("%dH", hours)
			d -= time.Duration(hours) * time.Hour
		}

		minutes := int(d.Minutes())
		if minutes > 0 {
			result += fmt.Sprintf("%dM", minutes)
			d -= time.Duration(minutes) * time.Minute
		}

		seconds := d.Seconds()
		if seconds > 0 {
			result += fmt.Sprintf("%.0fS", seconds)
		}
	}

	return result
}

func parseISO8601Duration(duration string) time.Duration {
	// Simple parser for ISO 8601 durations
	var result time.Duration

	if !strings.HasPrefix(duration, "P") {
		return 0
	}

	duration = duration[1:] // Remove P

	// Split on T
	parts := strings.Split(duration, "T")

	// Parse date part if present
	if len(parts) > 0 && parts[0] != "" {
		datePart := parts[0]
		// Parse days
		if idx := strings.Index(datePart, "D"); idx > 0 {
			days := parseInt(datePart[:idx])
			result += time.Duration(days) * 24 * time.Hour
		}
	}

	// Parse time part if present
	if len(parts) > 1 && parts[1] != "" {
		timePart := parts[1]

		// Parse hours
		if idx := strings.Index(timePart, "H"); idx > 0 {
			hours := parseInt(timePart[:idx])
			result += time.Duration(hours) * time.Hour
			timePart = timePart[idx+1:]
		}

		// Parse minutes
		if idx := strings.Index(timePart, "M"); idx > 0 {
			minutes := parseInt(timePart[:idx])
			result += time.Duration(minutes) * time.Minute
			timePart = timePart[idx+1:]
		}

		// Parse seconds
		if idx := strings.Index(timePart, "S"); idx > 0 {
			seconds := parseInt(timePart[:idx])
			result += time.Duration(seconds) * time.Second
		}
	}

	return result
}

func parseICalDuration(value string) time.Duration {
	// Handle iCalendar duration format (e.g., "P1DT2H30M" or "PT1H30M")
	return parseISO8601Duration(value)
}

func parseInt(s string) int {
	var i int
	fmt.Sscanf(s, "%d", &i)
	return i
}

func processParticipants(vevent *ics.VEvent, event *jscal.Event) {
	// Process organizer
	if organizer := vevent.GetProperty(ics.ComponentPropertyOrganizer); organizer != nil {
		if event.Participants == nil {
			event.Participants = make(map[string]*jscal.Participant)
		}

		email := strings.TrimPrefix(organizer.Value, "mailto:")
		participant := jscal.NewParticipant("", email)
		participant.Roles = map[string]bool{"owner": true, "attendee": true}

		if cn := organizer.ICalParameters["CN"]; len(cn) > 0 {
			participant.Name = &cn[0]
		}

		event.Participants[email] = participant
	}

	// Process attendees
	attendeeProps := vevent.Properties
	for i := range attendeeProps {
		attendee := &attendeeProps[i]
		if attendee.IANAToken != string(ics.ComponentPropertyAttendee) {
			continue
		}
		if event.Participants == nil {
			event.Participants = make(map[string]*jscal.Participant)
		}

		email := strings.TrimPrefix(attendee.Value, "mailto:")

		// Check if participant already exists (might be the organizer)
		participant, exists := event.Participants[email]
		if !exists {
			participant = jscal.NewParticipant("", email)
		}

		// Common Name
		if cn := attendee.ICalParameters["CN"]; len(cn) > 0 {
			participant.Name = &cn[0]
		}

		// Participation Status
		if partstat := attendee.ICalParameters["PARTSTAT"]; len(partstat) > 0 {
			status := strings.ToLower(partstat[0])
			status = strings.ReplaceAll(status, "-", "-")
			participant.ParticipationStatus = &status
		}

		// Role - merge with existing roles if participant exists
		if role := attendee.ICalParameters["ROLE"]; len(role) > 0 {
			if participant.Roles == nil {
				participant.Roles = make(map[string]bool)
			}
			switch strings.ToUpper(role[0]) {
			case "CHAIR":
				participant.Roles["chair"] = true
				participant.Roles["attendee"] = true
			case "REQ-PARTICIPANT":
				participant.Roles["attendee"] = true
			case "OPT-PARTICIPANT":
				participant.Roles["optional"] = true
			case "NON-PARTICIPANT":
				participant.Roles["informational"] = true
			default:
				participant.Roles["attendee"] = true
			}
		}

		event.Participants[email] = participant
	}
}

func convertParticipants(event *jscal.Event, vevent *ics.VEvent) {
	for email, participant := range event.Participants {
		mailto := email
		if !strings.HasPrefix(mailto, "mailto:") {
			mailto = "mailto:" + mailto
		}

		// Check if this is the organizer
		if participant.Roles != nil && participant.Roles["owner"] {
			organizerParams := make(map[string][]string)
			if participant.Name != nil {
				organizerParams["CN"] = []string{*participant.Name}
			}

			prop := ics.IANAProperty{
				BaseProperty: ics.BaseProperty{
					IANAToken:      string(ics.ComponentPropertyOrganizer),
					Value:          mailto,
					ICalParameters: organizerParams,
				},
			}
			vevent.Properties = append(vevent.Properties, prop)
		}

		// Add as attendee
		params := make(map[string][]string)

		if participant.Name != nil {
			params["CN"] = []string{*participant.Name}
		}

		if participant.ParticipationStatus != nil {
			params["PARTSTAT"] = []string{strings.ToUpper(*participant.ParticipationStatus)}
		}

		if participant.Roles != nil {
			role := "REQ-PARTICIPANT" // default
			if participant.Roles["chair"] {
				role = "CHAIR"
			} else if participant.Roles["optional"] {
				role = "OPT-PARTICIPANT"
			} else if participant.Roles["informational"] {
				role = "NON-PARTICIPANT"
			}
			params["ROLE"] = []string{role}
		}

		// Add attendee property
		prop := ics.IANAProperty{
			BaseProperty: ics.BaseProperty{
				IANAToken:      string(ics.ComponentPropertyAttendee),
				Value:          mailto,
				ICalParameters: params,
			},
		}
		vevent.Properties = append(vevent.Properties, prop)
	}
}

func processRecurrenceRules(vevent *ics.VEvent, event *jscal.Event) {
	// Process RRULE
	if rrule := vevent.GetProperty(ics.ComponentPropertyRrule); rrule != nil {
		rule := parseRRule(rrule.Value)
		if rule != nil {
			event.RecurrenceRules = append(event.RecurrenceRules, *rule)
		}
	}
}

func convertRecurrenceRules(event *jscal.Event, vevent *ics.VEvent) {
	for _, rule := range event.RecurrenceRules {
		rrule := formatRRule(&rule)
		if rrule != "" {
			vevent.AddProperty(ics.ComponentPropertyRrule, rrule)
		}
	}
}

func parseRRule(rruleValue string) *jscal.RecurrenceRule {
	rule := &jscal.RecurrenceRule{
		Type: "RecurrenceRule",
	}

	// Split the RRULE into parts
	parts := strings.Split(rruleValue, ";")

	for _, part := range parts {
		kv := strings.SplitN(part, "=", 2)
		if len(kv) != 2 {
			continue
		}

		key := strings.ToUpper(strings.TrimSpace(kv[0]))
		value := strings.TrimSpace(kv[1])

		switch key {
		case "FREQ":
			rule.Frequency = strings.ToLower(value)
		case "INTERVAL":
			if interval := parseInt(value); interval > 0 {
				rule.Interval = &interval
			}
		case "COUNT":
			if count := parseInt(value); count > 0 {
				rule.Count = &count
			}
		case "UNTIL":
			// Parse UNTIL date
			if until, _, _ := parseICalDateTime(&ics.IANAProperty{BaseProperty: ics.BaseProperty{Value: value}}); !until.IsZero() {
				rule.Until = jscal.NewLocalDateTime(until)
			}
		case "BYDAY":
			days := strings.Split(value, ",")
			for _, day := range days {
				if nday, err := jscal.ParseNDay(strings.TrimSpace(day)); err == nil {
					rule.ByDay = append(rule.ByDay, *nday)
				}
			}
		}
	}

	return rule
}

func formatRRule(rule *jscal.RecurrenceRule) string {
	var parts []string

	// FREQ is required
	if rule.Frequency != "" {
		parts = append(parts, fmt.Sprintf("FREQ=%s", strings.ToUpper(rule.Frequency)))
	}

	// INTERVAL
	if rule.Interval != nil && *rule.Interval > 1 {
		parts = append(parts, fmt.Sprintf("INTERVAL=%d", *rule.Interval))
	}

	// COUNT or UNTIL (mutually exclusive)
	if rule.Count != nil {
		parts = append(parts, fmt.Sprintf("COUNT=%d", *rule.Count))
	} else if rule.Until != nil {
		parts = append(parts, fmt.Sprintf("UNTIL=%s", rule.Until.Time().Format("20060102T150405Z")))
	}

	// BYDAY
	if len(rule.ByDay) > 0 {
		var dayStrings []string
		for _, nday := range rule.ByDay {
			dayStr := strings.ToUpper(nday.Day)
			if nday.NthOfPeriod != nil {
				dayStr = fmt.Sprintf("%d%s", *nday.NthOfPeriod, dayStr)
			}
			dayStrings = append(dayStrings, dayStr)
		}
		parts = append(parts, fmt.Sprintf("BYDAY=%s", strings.Join(dayStrings, ",")))
	}

	return strings.Join(parts, ";")
}
