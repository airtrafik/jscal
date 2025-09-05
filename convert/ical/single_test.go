package ical

import (
	"strings"
	"testing"

	"github.com/airtrafik/jscal"
)

func TestParseSingleEvent(t *testing.T) {
	converter := New()

	// Test parsing a single event
	icalData := `BEGIN:VCALENDAR
VERSION:2.0
PRODID:-//Test//Test//EN
BEGIN:VEVENT
UID:single-event@example.com
SUMMARY:Single Event Test
DTSTART:20250301T140000Z
DTEND:20250301T150000Z
END:VEVENT
END:VCALENDAR`

	event, err := converter.Parse([]byte(icalData))
	if err != nil {
		t.Fatalf("Failed to parse single event: %v", err)
	}

	if event.UID != "single-event@example.com" {
		t.Errorf("Expected UID 'single-event@example.com', got '%s'", event.UID)
	}

	if event.Title == nil || *event.Title != "Single Event Test" {
		t.Errorf("Expected title 'Single Event Test', got %v", event.Title)
	}
}

func TestParseMultipleEventsError(t *testing.T) {
	converter := New()

	// Test that Parse fails when there are multiple events
	icalData := `BEGIN:VCALENDAR
VERSION:2.0
PRODID:-//Test//Test//EN
BEGIN:VEVENT
UID:event1@example.com
SUMMARY:Event 1
DTSTART:20250301T140000Z
DTEND:20250301T150000Z
END:VEVENT
BEGIN:VEVENT
UID:event2@example.com
SUMMARY:Event 2
DTSTART:20250302T140000Z
DTEND:20250302T150000Z
END:VEVENT
END:VCALENDAR`

	_, err := converter.Parse([]byte(icalData))
	if err == nil {
		t.Fatal("Expected error when parsing multiple events with Parse()")
	}

	if !strings.Contains(err.Error(), "multiple events found") {
		t.Errorf("Expected error about multiple events, got: %v", err)
	}

	// Verify ParseAll works for multiple events
	events, err := converter.ParseAll([]byte(icalData))
	if err != nil {
		t.Fatalf("ParseAll should handle multiple events: %v", err)
	}

	if len(events) != 2 {
		t.Errorf("Expected 2 events from ParseAll, got %d", len(events))
	}
}

func TestFormatSingleEvent(t *testing.T) {
	converter := New()

	// Create a single event
	event := jscal.NewEvent("format-test@example.com", "Format Test Event")

	// Format single event
	icalData, err := converter.Format(event)
	if err != nil {
		t.Fatalf("Failed to format single event: %v", err)
	}

	icalStr := string(icalData)

	// Check that output contains expected iCalendar structure
	expectedPatterns := []string{
		"BEGIN:VCALENDAR",
		"BEGIN:VEVENT",
		"UID:format-test@example.com",
		"SUMMARY:Format Test Event",
		"END:VEVENT",
		"END:VCALENDAR",
	}

	for _, pattern := range expectedPatterns {
		if !strings.Contains(icalStr, pattern) {
			t.Errorf("Formatted iCalendar missing expected pattern: %s", pattern)
		}
	}
}

func TestRoundTripSingleEvent(t *testing.T) {
	converter := New()

	originalIcal := `BEGIN:VCALENDAR
VERSION:2.0
PRODID:-//Test//Test//EN
BEGIN:VEVENT
UID:roundtrip-single@example.com
SUMMARY:Round Trip Single Event
DESCRIPTION:Testing single event round trip
DTSTART:20250315T100000Z
DTEND:20250315T110000Z
LOCATION:Test Location
CATEGORIES:Test,RoundTrip
STATUS:CONFIRMED
END:VEVENT
END:VCALENDAR`

	// Parse single event
	event, err := converter.Parse([]byte(originalIcal))
	if err != nil {
		t.Fatalf("Failed to parse single event: %v", err)
	}

	// Format back to iCalendar
	formattedIcal, err := converter.Format(event)
	if err != nil {
		t.Fatalf("Failed to format single event: %v", err)
	}

	// Parse the formatted result
	reparsedEvent, err := converter.Parse(formattedIcal)
	if err != nil {
		t.Fatalf("Failed to reparse formatted event: %v", err)
	}

	// Verify key properties survived the round trip
	if reparsedEvent.UID != event.UID {
		t.Errorf("UID changed in round trip: %s -> %s", event.UID, reparsedEvent.UID)
	}

	if reparsedEvent.Title == nil || event.Title == nil || *reparsedEvent.Title != *event.Title {
		t.Errorf("Title changed in round trip")
	}

	if reparsedEvent.Description == nil || event.Description == nil || *reparsedEvent.Description != *event.Description {
		t.Errorf("Description changed in round trip")
	}
}
