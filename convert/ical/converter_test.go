package ical

import (
	"strings"
	"testing"
	"time"

	"github.com/airtrafik/jscal"
)

func TestConverterDetect(t *testing.T) {
	converter := New()

	// Test iCalendar detection
	icalData := `BEGIN:VCALENDAR
VERSION:2.0
PRODID:-//Test//Test//EN
BEGIN:VEVENT
UID:test@example.com
SUMMARY:Test Event
DTSTART:20250301T140000Z
DTEND:20250301T150000Z
END:VEVENT
END:VCALENDAR`

	if !converter.Detect([]byte(icalData)) {
		t.Error("Failed to detect valid iCalendar data")
	}

	// Test non-iCalendar data
	jsonData := `{"@type": "Event", "uid": "test", "title": "Test"}`
	if converter.Detect([]byte(jsonData)) {
		t.Error("Incorrectly detected JSON as iCalendar")
	}

	// Test partial iCalendar patterns
	partialData := `DTSTART:20250301T140000Z
SUMMARY:Test Event
UID:test@example.com`

	if !converter.Detect([]byte(partialData)) {
		t.Error("Failed to detect iCalendar patterns")
	}
}

func TestSimpleEventConversion(t *testing.T) {
	converter := New()

	icalData := `BEGIN:VCALENDAR
VERSION:2.0
PRODID:-//Test//Test//EN
BEGIN:VEVENT
UID:simple-test@example.com
SUMMARY:Simple Test Event
DESCRIPTION:This is a test event
DTSTART:20250301T140000Z
DTEND:20250301T150000Z
CREATED:20250201T120000Z
LAST-MODIFIED:20250215T090000Z
SEQUENCE:1
STATUS:CONFIRMED
LOCATION:Test Room
CATEGORIES:Test,Meeting
END:VEVENT
END:VCALENDAR`

	events, err := converter.ParseAll([]byte(icalData))
	if err != nil {
		t.Fatalf("Failed to convert iCalendar to JSCalendar: %v", err)
	}

	if len(events) != 1 {
		t.Fatalf("Expected 1 event, got %d", len(events))
	}

	event := events[0]

	// Check basic properties
	if event.UID != "simple-test@example.com" {
		t.Errorf("Expected UID 'simple-test@example.com', got '%s'", event.UID)
	}

	if event.Title == nil || *event.Title != "Simple Test Event" {
		t.Errorf("Expected title 'Simple Test Event', got %v", event.Title)
	}

	if event.Description == nil || *event.Description != "This is a test event" {
		t.Errorf("Expected description 'This is a test event', got %v", event.Description)
	}

	// Check timestamps
	expectedStart := time.Date(2025, 3, 1, 14, 0, 0, 0, time.UTC)
	expectedStartLDT := jscal.NewLocalDateTime(expectedStart)
	if event.Start == nil || !event.Start.Equal(expectedStartLDT) {
		t.Errorf("Expected start time %v, got %v", expectedStartLDT, event.Start)
	}

	// Check duration
	if event.Duration == nil || *event.Duration != "PT1H" {
		t.Errorf("Expected duration 'PT1H', got %v", event.Duration)
	}

	// Check sequence
	if event.Sequence == nil || *event.Sequence != 1 {
		t.Errorf("Expected sequence 1, got %v", event.Sequence)
	}

	// Check status
	if event.Status == nil || *event.Status != "confirmed" {
		t.Errorf("Expected status 'confirmed', got %v", event.Status)
	}

	// Check location
	if len(event.Locations) != 1 {
		t.Errorf("Expected 1 location, got %d", len(event.Locations))
	} else {
		location := event.Locations["1"]
		if location == nil || location.Name == nil || *location.Name != "Test Room" {
			t.Errorf("Expected location name 'Test Room', got %v", location)
		}
	}

	// Check categories
	if len(event.Categories) != 2 {
		t.Errorf("Expected 2 categories, got %d", len(event.Categories))
	}
	if !event.Categories["Test"] || !event.Categories["Meeting"] {
		t.Errorf("Expected categories 'Test' and 'Meeting', got %v", event.Categories)
	}
}

func TestAllDayEventConversion(t *testing.T) {
	converter := New()

	icalData := `BEGIN:VCALENDAR
VERSION:2.0
PRODID:-//Test//Test//EN
BEGIN:VEVENT
UID:allday-test@example.com
SUMMARY:All Day Event
DTSTART;VALUE=DATE:20251225
DTEND;VALUE=DATE:20251226
TRANSP:TRANSPARENT
END:VEVENT
END:VCALENDAR`

	events, err := converter.ParseAll([]byte(icalData))
	if err != nil {
		t.Fatalf("Failed to convert all-day event: %v", err)
	}

	if len(events) != 1 {
		t.Fatalf("Expected 1 event, got %d", len(events))
	}

	event := events[0]

	// Check all-day flag
	if !event.IsAllDay() {
		t.Error("Expected event to be marked as all-day")
	}

	// Check free/busy status from TRANSP
	if event.FreeBusyStatus == nil || *event.FreeBusyStatus != "free" {
		t.Errorf("Expected freeBusyStatus 'free' from TRANSP:TRANSPARENT, got %v", event.FreeBusyStatus)
	}
}

func TestEventWithParticipants(t *testing.T) {
	converter := New()

	icalData := `BEGIN:VCALENDAR
VERSION:2.0
PRODID:-//Test//Test//EN
BEGIN:VEVENT
UID:meeting-test@example.com
SUMMARY:Team Meeting
DTSTART:20250301T140000Z
DTEND:20250301T150000Z
ORGANIZER;CN=John Doe:mailto:john.doe@example.com
ATTENDEE;CN=John Doe;ROLE=CHAIR;PARTSTAT=ACCEPTED:mailto:john.doe@example.com
ATTENDEE;CN=Jane Smith;ROLE=REQ-PARTICIPANT;PARTSTAT=TENTATIVE:mailto:jane.smith@example.com
ATTENDEE;CN=Bob Johnson;ROLE=OPT-PARTICIPANT;PARTSTAT=NEEDS-ACTION:mailto:bob.johnson@example.com
END:VEVENT
END:VCALENDAR`

	events, err := converter.ParseAll([]byte(icalData))
	if err != nil {
		t.Fatalf("Failed to convert event with participants: %v", err)
	}

	event := events[0]

	// Should have 3 participants (organizer is also an attendee)
	if len(event.Participants) != 3 {
		t.Fatalf("Expected 3 participants, got %d", len(event.Participants))
	}

	// Check organizer
	organizer := event.Participants["john.doe@example.com"]
	if organizer == nil {
		t.Fatal("Organizer not found in participants")
	}

	if organizer.Name == nil || *organizer.Name != "John Doe" {
		t.Errorf("Expected organizer name 'John Doe', got %v", organizer.Name)
	}

	if !organizer.Roles["owner"] || !organizer.Roles["chair"] {
		t.Errorf("Expected organizer to have owner and chair roles, got %v", organizer.Roles)
	}

	if organizer.ParticipationStatus == nil || *organizer.ParticipationStatus != "accepted" {
		t.Errorf("Expected organizer participation status 'accepted', got %v", organizer.ParticipationStatus)
	}

	// Check optional participant
	optional := event.Participants["bob.johnson@example.com"]
	if optional == nil {
		t.Fatal("Optional participant not found")
	}

	if !optional.Roles["optional"] {
		t.Errorf("Expected Bob to have optional role, got %v", optional.Roles)
	}

	if optional.ParticipationStatus == nil || *optional.ParticipationStatus != "needs-action" {
		t.Errorf("Expected participation status 'needs-action', got %v", optional.ParticipationStatus)
	}
}

func TestRecurringEventConversion(t *testing.T) {
	converter := New()

	icalData := `BEGIN:VCALENDAR
VERSION:2.0
PRODID:-//Test//Test//EN
BEGIN:VEVENT
UID:recurring-test@example.com
SUMMARY:Daily Standup
DTSTART:20250303T090000Z
DURATION:PT30M
RRULE:FREQ=DAILY;BYDAY=MO,TU,WE,TH,FR;UNTIL=20250331T235959Z
END:VEVENT
END:VCALENDAR`

	events, err := converter.ParseAll([]byte(icalData))
	if err != nil {
		t.Fatalf("Failed to convert recurring event: %v", err)
	}

	event := events[0]

	// Check recurrence
	if !event.IsRecurring() {
		t.Error("Expected event to be recurring")
	}

	if len(event.RecurrenceRules) != 1 {
		t.Fatalf("Expected 1 recurrence rule, got %d", len(event.RecurrenceRules))
	}

	rule := event.RecurrenceRules[0]

	if rule.Frequency != "daily" {
		t.Errorf("Expected frequency 'daily', got '%s'", rule.Frequency)
	}

	if len(rule.ByDay) != 5 {
		t.Errorf("Expected 5 days in BYDAY, got %d", len(rule.ByDay))
	}

	// Check that weekdays are included
	daySet := make(map[string]bool)
	for _, nday := range rule.ByDay {
		daySet[nday.Day] = true
	}

	expectedDays := []string{"mo", "tu", "we", "th", "fr"}
	for _, day := range expectedDays {
		if !daySet[day] {
			t.Errorf("Expected day '%s' in recurrence rule", day)
		}
	}

	// Check UNTIL
	if rule.Until == nil {
		t.Error("Expected UNTIL to be set")
	} else {
		expectedUntil := time.Date(2025, 3, 31, 23, 59, 59, 0, time.UTC)
		expectedUntilLocal := jscal.NewLocalDateTime(expectedUntil)
		if !rule.Until.Equal(expectedUntilLocal) {
			t.Errorf("Expected UNTIL %v, got %v", expectedUntil, *rule.Until)
		}
	}
}

func TestRoundTripConversion(t *testing.T) {
	converter := New()

	// Create a JSCalendar event
	originalEvent := jscal.NewEvent("roundtrip-test@example.com", "Round Trip Test")
	startTime := time.Date(2025, 3, 1, 14, 0, 0, 0, time.UTC)
	originalEvent.Start = jscal.NewLocalDateTime(startTime)
	originalEvent.Duration = jscal.String("PT1H")
	originalEvent.TimeZone = jscal.String("UTC")
	originalEvent.Description = jscal.String("Test description with special chars: ,;\\n")
	originalEvent.AddCategory("Test")
	originalEvent.AddCategory("Round Trip")

	participant := jscal.NewParticipant("Test User", "test@example.com")
	participant.ParticipationStatus = jscal.String("accepted")
	originalEvent.AddParticipant("test@example.com", participant)

	// Convert to iCalendar
	icalData, err := converter.FormatAll([]*jscal.Event{originalEvent})
	if err != nil {
		t.Fatalf("Failed to convert JSCalendar to iCalendar: %v", err)
	}

	// Verify iCalendar contains expected content
	icalStr := string(icalData)
	expectedPatterns := []string{
		"BEGIN:VCALENDAR",
		"BEGIN:VEVENT",
		"UID:roundtrip-test@example.com",
		"SUMMARY:Round Trip Test",
		"DTSTART:20250301T140000Z",
		"DURATION:PT1H",
		"CATEGORIES:Round Trip\\,Test", // Categories should be sorted (comma escaped in iCal)
		"ATTENDEE",
		"test@example.com",
		"END:VEVENT",
		"END:VCALENDAR",
	}

	for _, pattern := range expectedPatterns {
		if !strings.Contains(icalStr, pattern) {
			t.Errorf("Generated iCalendar missing expected pattern: %s\nGenerated:\n%s", pattern, icalStr)
		}
	}

	// Convert back to JSCalendar
	convertedEvents, err := converter.ParseAll(icalData)
	if err != nil {
		t.Fatalf("Failed to convert iCalendar back to JSCalendar: %v", err)
	}

	if len(convertedEvents) != 1 {
		t.Fatalf("Expected 1 event after round trip, got %d", len(convertedEvents))
	}

	convertedEvent := convertedEvents[0]

	// Verify key properties survived the round trip
	if convertedEvent.UID != originalEvent.UID {
		t.Errorf("UID changed during round trip: %s -> %s", originalEvent.UID, convertedEvent.UID)
	}

	if convertedEvent.Title == nil || *convertedEvent.Title != *originalEvent.Title {
		t.Errorf("Title changed during round trip: %v -> %v", originalEvent.Title, convertedEvent.Title)
	}

	if convertedEvent.Start == nil || !convertedEvent.Start.Equal(originalEvent.Start) {
		t.Errorf("Start time changed during round trip: %v -> %v", originalEvent.Start, convertedEvent.Start)
	}

	if convertedEvent.Duration == nil || *convertedEvent.Duration != *originalEvent.Duration {
		t.Errorf("Duration changed during round trip: %v -> %v", originalEvent.Duration, convertedEvent.Duration)
	}

	// Check categories
	if len(convertedEvent.Categories) != len(originalEvent.Categories) {
		t.Errorf("Category count changed during round trip: %d -> %d",
			len(originalEvent.Categories), len(convertedEvent.Categories))
	}

	for cat := range originalEvent.Categories {
		if !convertedEvent.Categories[cat] {
			t.Errorf("Category '%s' lost during round trip", cat)
		}
	}
}
