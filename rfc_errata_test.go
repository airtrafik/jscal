package jscal

import (
	"testing"
	"time"
)

// TestRFCErrata tests the corrections from RFC 8984 errata
func TestRFCErrata(t *testing.T) {
	t.Run("Erratum 6873: RecurrenceIdTimeZone optional with RecurrenceId", func(t *testing.T) {
		// Per erratum 6873, recurrenceIdTimeZone can be null/omitted even when recurrenceId is set
		// This allows for floating time recurrence instances
		
		now := time.Now()
		recurrenceTime := NewLocalDateTime(now.Add(24 * time.Hour))
		
		// Test Event with recurrenceId but no recurrenceIdTimeZone (floating time)
		event := &Event{
			Type:         "Event",
			UID:          "test-recurrence-instance",
			Title:        String("Recurrence Instance"),
			Start:        NewLocalDateTime(now),
			RecurrenceId: recurrenceTime,
			// RecurrenceIdTimeZone is intentionally omitted (floating time)
		}
		
		// Should validate successfully
		if err := event.Validate(); err != nil {
			t.Errorf("Event with recurrenceId but no recurrenceIdTimeZone should be valid (floating time): %v", err)
		}
		
		// Test Event with both recurrenceId and recurrenceIdTimeZone
		event2 := &Event{
			Type:                 "Event",
			UID:                  "test-recurrence-instance-tz",
			Title:                String("Recurrence Instance with TZ"),
			Start:                NewLocalDateTime(now),
			RecurrenceId:         recurrenceTime,
			RecurrenceIdTimeZone: String("America/New_York"),
		}
		
		// Should also validate successfully
		if err := event2.Validate(); err != nil {
			t.Errorf("Event with both recurrenceId and recurrenceIdTimeZone should be valid: %v", err)
		}
		
		// Test Task with recurrenceId but no recurrenceIdTimeZone
		task := &Task{
			Type:         "Task",
			UID:          "test-task-recurrence",
			Title:        String("Task Recurrence Instance"),
			RecurrenceId: recurrenceTime,
			// RecurrenceIdTimeZone is intentionally omitted
		}
		
		// Should validate successfully
		if err := task.Validate(); err != nil {
			t.Errorf("Task with recurrenceId but no recurrenceIdTimeZone should be valid: %v", err)
		}
	})
	
	t.Run("Erratum 8028: Group uses title not name", func(t *testing.T) {
		// This was already discovered and fixed
		// Verify that Group has title field and it works correctly
		
		group := &Group{
			Type:  "Group",
			UID:   "test-group",
			Title: String("Test Group Title"),
		}
		
		// Should have title
		if group.Title == nil || *group.Title != "Test Group Title" {
			t.Error("Group should have title field")
		}
		
		// Should validate successfully
		if err := group.Validate(); err != nil {
			t.Errorf("Group with title should validate: %v", err)
		}
	})
	
	t.Run("Erratum 6872: Privacy property documentation", func(t *testing.T) {
		// This erratum extends the list of shareable properties for private events
		// We don't implement privacy filtering yet, but we can test that the properties exist
		
		event := &Event{
			Type:    "Event",
			UID:     "private-event",
			Title:   String("Private Event"),
			Start:   NewLocalDateTime(time.Now()),
			Privacy: String("private"),
			
			// These properties should be shareable even for private events (per erratum)
			Excluded:                Bool(false),
			ExcludedRecurrenceRules: []RecurrenceRule{},
			RecurrenceId:            NewLocalDateTime(time.Now()),
			RecurrenceIdTimeZone:    String("UTC"),
			RecurrenceRules: []RecurrenceRule{
				{
					Type:      "RecurrenceRule",
					Frequency: "daily",
				},
			},
		}
		
		// Event should validate with privacy and recurrence properties
		if err := event.Validate(); err != nil {
			t.Errorf("Private event with recurrence properties should validate: %v", err)
		}
		
		// Verify privacy values are accepted
		for _, privacy := range []string{"public", "private", "secret"} {
			event.Privacy = String(privacy)
			if err := event.Validate(); err != nil {
				t.Errorf("Event with privacy=%s should validate: %v", privacy, err)
			}
		}
	})
}

// TestFloatingTimeRecurrence tests floating time recurrence instances
func TestFloatingTimeRecurrence(t *testing.T) {
	// Create a recurring event in floating time (no timezone)
	mainEvent := &Event{
		Type:  "Event",
		UID:   "recurring-floating",
		Title: String("Daily Meditation"),
		Start: NewLocalDateTime(time.Date(2024, 1, 1, 7, 0, 0, 0, time.UTC)),
		// No TimeZone specified - this is floating time
		Duration: String("PT30M"),
		RecurrenceRules: []RecurrenceRule{
			{
				Type:      "RecurrenceRule",
				Frequency: "daily",
			},
		},
	}
	
	if err := mainEvent.Validate(); err != nil {
		t.Fatalf("Main floating time event should validate: %v", err)
	}
	
	// Create an override for a specific instance (also floating time)
	overrideInstance := &Event{
		Type:         "Event",
		UID:          "recurring-floating",
		Title:        String("Extended Meditation"),
		Start:        NewLocalDateTime(time.Date(2024, 1, 15, 7, 0, 0, 0, time.UTC)),
		Duration:     String("PT1H"),
		RecurrenceId: NewLocalDateTime(time.Date(2024, 1, 15, 7, 0, 0, 0, time.UTC)),
		// No RecurrenceIdTimeZone - this is valid per erratum 6873
	}
	
	if err := overrideInstance.Validate(); err != nil {
		t.Errorf("Floating time recurrence override should validate: %v", err)
	}
}