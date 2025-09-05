package jscal

import (
	"testing"
	"time"
)

func TestNewEvent(t *testing.T) {
	uid := "test-event-123"
	title := "Test Event"
	
	event := NewEvent(uid, title)
	
	if event.Type != "Event" {
		t.Errorf("Expected Type to be 'Event', got '%s'", event.Type)
	}
	
	if event.UID != uid {
		t.Errorf("Expected UID to be '%s', got '%s'", uid, event.UID)
	}
	
	if event.Title == nil || *event.Title != title {
		t.Errorf("Expected Title to be '%s', got '%v'", title, event.Title)
	}
	
	if event.Created == nil {
		t.Error("Expected Created to be set")
	}
	
	if event.Updated == nil {
		t.Error("Expected Updated to be set")
	}
	
	if event.Sequence == nil || *event.Sequence != 0 {
		t.Errorf("Expected Sequence to be 0, got %v", event.Sequence)
	}
}

func TestEventJSON(t *testing.T) {
	event := NewEvent("test-123", "Test Event")
	
	jsonData, err := event.JSON()
	if err != nil {
		t.Fatalf("Failed to marshal to JSON: %v", err)
	}
	
	// Parse it back
	parsedEvent, err := Parse(jsonData)
	if err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}
	
	if parsedEvent.UID != event.UID {
		t.Errorf("Expected UID to be '%s', got '%s'", event.UID, parsedEvent.UID)
	}
	
	if parsedEvent.Title == nil || *parsedEvent.Title != *event.Title {
		t.Errorf("Expected Title to be '%s', got '%v'", *event.Title, parsedEvent.Title)
	}
}

func TestEventIsAllDay(t *testing.T) {
	event := NewEvent("test-123", "All Day Event")
	
	// Should not be all-day by default
	if event.IsAllDay() {
		t.Error("Expected event to not be all-day by default")
	}
	
	// Set as all-day
	event.ShowWithoutTime = Bool(true)
	if !event.IsAllDay() {
		t.Error("Expected event to be all-day")
	}
}

func TestEventDuration(t *testing.T) {
	event := NewEvent("test-123", "Test Event")
	
	// Set duration to 1 hour
	event.Duration = String("PT1H")
	
	duration, err := event.GetDuration()
	if err != nil {
		t.Fatalf("Failed to get duration: %v", err)
	}
	
	expected := time.Hour
	if duration != expected {
		t.Errorf("Expected duration to be %v, got %v", expected, duration)
	}
}

func TestEventEndTime(t *testing.T) {
	event := NewEvent("test-123", "Test Event")
	
	startTime := time.Date(2025, 3, 1, 14, 0, 0, 0, time.UTC)
	event.Start = &startTime
	event.Duration = String("PT1H")
	
	endTime, err := event.GetEndTime()
	if err != nil {
		t.Fatalf("Failed to get end time: %v", err)
	}
	
	expectedEnd := startTime.Add(time.Hour)
	if !endTime.Equal(expectedEnd) {
		t.Errorf("Expected end time to be %v, got %v", expectedEnd, *endTime)
	}
}

func TestEventParticipants(t *testing.T) {
	event := NewEvent("test-123", "Test Event")
	
	participant := NewParticipant("John Doe", "john.doe@example.com")
	event.AddParticipant("john.doe@example.com", participant)
	
	if len(event.Participants) != 1 {
		t.Errorf("Expected 1 participant, got %d", len(event.Participants))
	}
	
	retrieved := event.Participants["john.doe@example.com"]
	if retrieved == nil {
		t.Fatal("Participant not found")
	}
	
	if retrieved.Name == nil || *retrieved.Name != "John Doe" {
		t.Errorf("Expected participant name to be 'John Doe', got %v", retrieved.Name)
	}
	
	if retrieved.Email == nil || *retrieved.Email != "john.doe@example.com" {
		t.Errorf("Expected participant email to be 'john.doe@example.com', got %v", retrieved.Email)
	}
}

func TestEventCategories(t *testing.T) {
	event := NewEvent("test-123", "Test Event")
	
	event.AddCategory("Work")
	event.AddCategory("Meeting")
	
	if len(event.Categories) != 2 {
		t.Errorf("Expected 2 categories, got %d", len(event.Categories))
	}
	
	if !event.Categories["Work"] {
		t.Error("Expected 'Work' category to be true")
	}
	
	if !event.Categories["Meeting"] {
		t.Error("Expected 'Meeting' category to be true")
	}
}

func TestEventTouch(t *testing.T) {
	event := NewEvent("test-123", "Test Event")
	
	originalUpdated := event.Updated
	originalSequence := *event.Sequence
	
	// Wait a bit to ensure timestamp difference
	time.Sleep(time.Millisecond)
	
	event.Touch()
	
	if event.Updated.Equal(*originalUpdated) {
		t.Error("Expected Updated timestamp to change")
	}
	
	if *event.Sequence != originalSequence+1 {
		t.Errorf("Expected sequence to increment from %d to %d, got %d", 
			originalSequence, originalSequence+1, *event.Sequence)
	}
}

func TestHelperFunctions(t *testing.T) {
	// Test String helper
	s := String("test")
	if s == nil || *s != "test" {
		t.Error("String helper failed")
	}
	
	// Test Int helper
	i := Int(42)
	if i == nil || *i != 42 {
		t.Error("Int helper failed")
	}
	
	// Test Bool helper
	b := Bool(true)
	if b == nil || *b != true {
		t.Error("Bool helper failed")
	}
}

func TestFormatDayOfWeek(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"MONDAY", "mo"},
		{"MO", "mo"},
		{"Tuesday", "tu"},
		{"WE", "we"},
		{"thursday", "th"},
		{"FR", "fr"},
		{"saturday", "sa"},
		{"SU", "su"},
	}
	
	for _, test := range tests {
		result := FormatDayOfWeek(test.input)
		if result != test.expected {
			t.Errorf("FormatDayOfWeek(%s) = %s, expected %s", 
				test.input, result, test.expected)
		}
	}
}