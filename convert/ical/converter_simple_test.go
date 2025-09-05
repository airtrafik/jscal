package ical

import (
	"testing"
	"time"

	"github.com/airtrafik/jscal"
)

func TestBasicConversion(t *testing.T) {
	converter := New()
	
	icalData := `BEGIN:VCALENDAR
VERSION:2.0
PRODID:-//Test//Test//EN
BEGIN:VEVENT
UID:test-event@example.com
SUMMARY:Test Event
DTSTART:20250301T140000Z
DTEND:20250301T150000Z
DESCRIPTION:Test description
END:VEVENT
END:VCALENDAR`
	
	// Test ParseAll
	events, err := converter.ParseAll([]byte(icalData))
	if err != nil {
		t.Fatalf("Failed to parse iCalendar: %v", err)
	}
	
	if len(events) != 1 {
		t.Fatalf("Expected 1 event, got %d", len(events))
	}
	
	event := events[0]
	
	// Check basic properties
	if event.UID != "test-event@example.com" {
		t.Errorf("Expected UID 'test-event@example.com', got '%s'", event.UID)
	}
	
	if event.Title == nil || *event.Title != "Test Event" {
		t.Errorf("Expected title 'Test Event', got %v", event.Title)
	}
	
	if event.Description == nil || *event.Description != "Test description" {
		t.Errorf("Expected description 'Test description', got %v", event.Description)
	}
}

func TestRoundTrip(t *testing.T) {
	converter := New()
	
	// Create a JSCalendar event
	event := jscal.NewEvent("roundtrip@example.com", "Round Trip Test")
	event.Description = jscal.String("Testing round trip conversion")
	startTime := time.Date(2025, 3, 15, 10, 0, 0, 0, time.UTC)
	event.Start = &startTime
	event.Duration = jscal.String("PT1H")
	
	// Convert to iCalendar
	icalData, err := converter.Format(event)
	if err != nil {
		t.Fatalf("Failed to format event: %v", err)
	}
	
	// Convert back to JSCalendar
	parsedEvent, err := converter.Parse(icalData)
	if err != nil {
		t.Fatalf("Failed to parse formatted iCalendar: %v", err)
	}
	
	// Check key properties survived
	if parsedEvent.UID != event.UID {
		t.Errorf("UID changed: %s -> %s", event.UID, parsedEvent.UID)
	}
	
	if parsedEvent.Title == nil || *parsedEvent.Title != *event.Title {
		t.Errorf("Title changed")
	}
	
	if parsedEvent.Description == nil || *parsedEvent.Description != *event.Description {
		t.Errorf("Description changed")
	}
}