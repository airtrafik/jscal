package jscal

import (
	"os"
	"path/filepath"
	"testing"
)

// TestRFC8984Examples tests all RFC 8984 Section 6 examples
func TestRFC8984Examples(t *testing.T) {
	examplesDir := "testdata/rfc8984/examples"

	testCases := []struct {
		name     string
		file     string
		validate func(t *testing.T, obj CalendarObject)
	}{
		{
			name:     "6.1 Simple Event",
			file:     "6.1-simple-event.json",
			validate: validateSimpleEvent,
		},
		{
			name:     "6.2 Simple Task",
			file:     "6.2-simple-task.json",
			validate: validateSimpleTask,
		},
		{
			name:     "6.3 Simple Group",
			file:     "6.3-simple-group.json",
			validate: validateSimpleGroup,
		},
		{
			name:     "6.4 All-Day Event",
			file:     "6.4-all-day-event.json",
			validate: validateAllDayEvent,
		},
		{
			name:     "6.5 Task with Due Date",
			file:     "6.5-task-with-due-date.json",
			validate: validateTaskWithDueDate,
		},
		{
			name:     "6.6 Event with End Timezone",
			file:     "6.6-event-end-timezone.json",
			validate: validateEventWithEndTimezone,
		},
		{
			name:     "6.7 Floating-Time Event",
			file:     "6.7-floating-time-event.json",
			validate: validateFloatingTimeEvent,
		},
		{
			name:     "6.8 Event with Localization",
			file:     "6.8-event-localization.json",
			validate: validateEventWithLocalization,
		},
		{
			name:     "6.9 Recurring Event with Overrides",
			file:     "6.9-recurring-overrides.json",
			validate: validateRecurringWithOverrides,
		},
		{
			name:     "6.10 Recurring Event with Participants",
			file:     "6.10-recurring-participants.json",
			validate: validateRecurringWithParticipants,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Load the test data
			data, err := os.ReadFile(filepath.Join(examplesDir, tc.file))
			if err != nil {
				t.Fatalf("Failed to read test file %s: %v", tc.file, err)
			}

			// Parse the object
			obj, err := Parse(data)
			if err != nil {
				t.Fatalf("Failed to parse %s: %v", tc.file, err)
			}

			// Validate the object
			if err := obj.Validate(); err != nil {
				t.Errorf("Validation failed for %s: %v", tc.file, err)
			}

			// Run specific validation
			tc.validate(t, obj)
		})
	}
}

// validateSimpleEvent validates RFC example 6.1
func validateSimpleEvent(t *testing.T, obj CalendarObject) {
	event, ok := obj.(*Event)
	if !ok {
		t.Fatal("Expected Event type")
	}

	// Check required fields
	if event.UID != "a8df6573-0474-496d-8496-033ad45d7fea" {
		t.Errorf("Unexpected UID: %s", event.UID)
	}

	if event.Title == nil || *event.Title != "Some event" {
		t.Error("Title mismatch")
	}

	if event.TimeZone == nil || *event.TimeZone != "America/New_York" {
		t.Error("TimeZone mismatch")
	}

	if event.Duration == nil || *event.Duration != "PT1H" {
		t.Error("Duration mismatch")
	}

	// Verify it's not marked as all-day
	if event.ShowWithoutTime != nil && *event.ShowWithoutTime {
		t.Error("Should not be all-day event")
	}
}

// validateSimpleTask validates RFC example 6.2
func validateSimpleTask(t *testing.T, obj CalendarObject) {
	task, ok := obj.(*Task)
	if !ok {
		t.Fatal("Expected Task type")
	}

	if task.UID != "2a358cee-6489-4f14-a57f-c104db4dc2f2" {
		t.Errorf("Unexpected UID: %s", task.UID)
	}

	if task.Title == nil || *task.Title != "Do something" {
		t.Error("Title mismatch")
	}

	// Should not have due date
	if task.Due != nil {
		t.Error("Simple task should not have due date")
	}
}

// validateSimpleGroup validates RFC example 6.3
func validateSimpleGroup(t *testing.T, obj CalendarObject) {
	group, ok := obj.(*Group)
	if !ok {
		t.Fatal("Expected Group type")
	}

	if group.UID != "bf0ac22b-4989-4caf-9ebd-54301b4ee51a" {
		t.Errorf("Unexpected UID: %s", group.UID)
	}

	if group.Title == nil || *group.Title != "A simple group" {
		t.Error("Title mismatch")
	}

	// Check entries
	if len(group.Entries) != 2 {
		t.Fatalf("Expected 2 entries, got %d", len(group.Entries))
	}

	// First entry should be Event
	event, ok := group.Entries[0].(*Event)
	if !ok {
		t.Fatal("First entry should be Event")
	}
	if event.UID != "a8df6573-0474-496d-8496-033ad45d7fea" {
		t.Error("First entry UID mismatch")
	}

	// Second entry should be Task
	task, ok := group.Entries[1].(*Task)
	if !ok {
		t.Fatal("Second entry should be Task")
	}
	if task.UID != "2a358cee-6489-4f14-a57f-c104db4dc2f2" {
		t.Error("Second entry UID mismatch")
	}
}

// validateAllDayEvent validates RFC example 6.4
func validateAllDayEvent(t *testing.T, obj CalendarObject) {
	event, ok := obj.(*Event)
	if !ok {
		t.Fatal("Expected Event type")
	}

	if event.Title == nil || *event.Title != "April Fool's Day" {
		t.Error("Title mismatch")
	}

	// Should be marked as all-day
	if event.ShowWithoutTime == nil || !*event.ShowWithoutTime {
		t.Error("Should be all-day event")
	}

	// Should have duration of 1 day
	if event.Duration == nil || *event.Duration != "P1D" {
		t.Error("Duration should be P1D")
	}

	// Should have yearly recurrence
	if len(event.RecurrenceRules) != 1 {
		t.Fatal("Should have one recurrence rule")
	}
	if event.RecurrenceRules[0].Frequency != "yearly" {
		t.Error("Should be yearly recurrence")
	}
}

// validateTaskWithDueDate validates RFC example 6.5
func validateTaskWithDueDate(t *testing.T, obj CalendarObject) {
	task, ok := obj.(*Task)
	if !ok {
		t.Fatal("Expected Task type")
	}

	if task.Title == nil || *task.Title != "Buy groceries" {
		t.Error("Title mismatch")
	}

	// Should have due date
	if task.Due == nil {
		t.Fatal("Should have due date")
	}

	// Should have timezone
	if task.TimeZone == nil || *task.TimeZone != "Europe/Vienna" {
		t.Error("TimeZone mismatch")
	}

	// Should have estimated duration
	if task.EstimatedDuration == nil || *task.EstimatedDuration != "PT1H" {
		t.Error("EstimatedDuration mismatch")
	}
}

// validateEventWithEndTimezone validates RFC example 6.6
func validateEventWithEndTimezone(t *testing.T, obj CalendarObject) {
	event, ok := obj.(*Event)
	if !ok {
		t.Fatal("Expected Event type")
	}

	if event.Title == nil || *event.Title != "Flight XY51 to Tokyo" {
		t.Error("Title mismatch")
	}

	// Should have locations
	if len(event.Locations) != 2 {
		t.Fatalf("Should have 2 locations, got %d", len(event.Locations))
	}

	// Check start location
	if loc, ok := event.Locations["1"]; ok {
		if loc.Name == nil || *loc.Name != "Frankfurt Airport (FRA)" {
			t.Error("Start location name mismatch")
		}
		if loc.Rel == nil || *loc.Rel != "start" {
			t.Error("Start location should have rel='start'")
		}
	} else {
		t.Error("Missing location 1")
	}

	// Check end location
	if loc, ok := event.Locations["2"]; ok {
		if loc.Name == nil || *loc.Name != "Narita International Airport (NRT)" {
			t.Error("End location name mismatch")
		}
		if loc.Rel == nil || *loc.Rel != "end" {
			t.Error("End location should have rel='end'")
		}
		if loc.TimeZone == nil || *loc.TimeZone != "Asia/Tokyo" {
			t.Error("End location should have Tokyo timezone")
		}
	} else {
		t.Error("Missing location 2")
	}
}

// validateFloatingTimeEvent validates RFC example 6.7
func validateFloatingTimeEvent(t *testing.T, obj CalendarObject) {
	event, ok := obj.(*Event)
	if !ok {
		t.Fatal("Expected Event type")
	}

	if event.Title == nil || *event.Title != "Yoga" {
		t.Error("Title mismatch")
	}

	// Should NOT have timezone (floating time)
	if event.TimeZone != nil {
		t.Error("Floating time event should not have timezone")
	}

	// Should have daily recurrence
	if len(event.RecurrenceRules) != 1 {
		t.Fatal("Should have one recurrence rule")
	}
	if event.RecurrenceRules[0].Frequency != "daily" {
		t.Error("Should be daily recurrence")
	}

	// Duration should be 30 minutes
	if event.Duration == nil || *event.Duration != "PT30M" {
		t.Error("Duration should be PT30M")
	}
}

// validateEventWithLocalization validates RFC example 6.8
func validateEventWithLocalization(t *testing.T, obj CalendarObject) {
	event, ok := obj.(*Event)
	if !ok {
		t.Fatal("Expected Event type")
	}

	if event.Title == nil || *event.Title != "Live from Music Bowl: The Band" {
		t.Error("Title mismatch")
	}

	// Should have locale
	if event.Locale == nil || *event.Locale != "en" {
		t.Error("Locale should be 'en'")
	}

	// Should have physical location
	if len(event.Locations) != 1 {
		t.Error("Should have 1 physical location")
	}

	// Should have virtual location
	if len(event.VirtualLocations) != 1 {
		t.Error("Should have 1 virtual location")
	}

	if vLoc, ok := event.VirtualLocations["vloc1"]; ok {
		if vLoc.Name == nil || *vLoc.Name != "Free live Stream from Music Bowl" {
			t.Error("Virtual location name mismatch")
		}
		if vLoc.URI != "https://stream.example.com/the_band_2020" {
			t.Error("Virtual location URI mismatch")
		}
	} else {
		t.Error("Missing virtual location vloc1")
	}

	// Should have localizations
	if len(event.Localizations) != 1 {
		t.Error("Should have 1 localization")
	}
	if _, ok := event.Localizations["de"]; !ok {
		t.Error("Should have German localization")
	}
}

// validateRecurringWithOverrides validates RFC example 6.9
func validateRecurringWithOverrides(t *testing.T, obj CalendarObject) {
	event, ok := obj.(*Event)
	if !ok {
		t.Fatal("Expected Event type")
	}

	if event.Title == nil || *event.Title != "Calculus I" {
		t.Error("Title mismatch")
	}

	// Should have weekly recurrence
	if len(event.RecurrenceRules) != 1 {
		t.Fatal("Should have one recurrence rule")
	}
	rule := event.RecurrenceRules[0]
	if rule.Frequency != "weekly" {
		t.Error("Should be weekly recurrence")
	}
	if rule.Until == nil {
		t.Error("Should have until date")
	}

	// Should have 3 overrides
	if len(event.RecurrenceOverrides) != 3 {
		t.Errorf("Should have 3 recurrence overrides, got %d", len(event.RecurrenceOverrides))
	}

	// Check for excluded date
	if override, ok := event.RecurrenceOverrides["2020-04-01T09:00:00"]; ok {
		if excluded, ok := override["excluded"].(bool); !ok || !excluded {
			t.Error("April 1 should be excluded")
		}
	} else {
		t.Error("Missing April 1 override")
	}
}

// validateRecurringWithParticipants validates RFC example 6.10
func validateRecurringWithParticipants(t *testing.T, obj CalendarObject) {
	event, ok := obj.(*Event)
	if !ok {
		t.Fatal("Expected Event type")
	}

	if event.Title == nil || *event.Title != "FooBar team meeting" {
		t.Error("Title mismatch")
	}

	// Should have weekly recurrence
	if len(event.RecurrenceRules) != 1 {
		t.Fatal("Should have one recurrence rule")
	}
	if event.RecurrenceRules[0].Frequency != "weekly" {
		t.Error("Should be weekly recurrence")
	}

	// Should have participants
	if len(event.Participants) != 2 {
		t.Errorf("Should have 2 participants, got %d", len(event.Participants))
	}

	// Check participant roles
	if p, ok := event.Participants["em9lQGZvb2GFtcGxlLmNvbQ"]; ok {
		if p.Name == nil || *p.Name != "Zoe Zelda" {
			t.Error("Participant name mismatch")
		}
		if !p.Roles["owner"] || !p.Roles["chair"] || !p.Roles["attendee"] {
			t.Error("Zoe should have owner, chair, and attendee roles")
		}
		// Check SendTo
		if p.SendTo == nil || p.SendTo["imip"] != "mailto:zoe@foobar.example.com" {
			t.Error("Zoe should have imip sendTo")
		}
	} else {
		t.Error("Missing Zoe participant")
	}

	// Should have virtual location
	if len(event.VirtualLocations) != 1 {
		t.Error("Should have 1 virtual location")
	}

	// Should have replyTo
	if event.ReplyTo == nil || event.ReplyTo["imip"] != "mailto:f245f875-7f63-4a5e-a2c8@schedule.example.com" {
		t.Error("Should have replyTo imip")
	}

	// Should have recurrence override
	if len(event.RecurrenceOverrides) != 1 {
		t.Error("Should have 1 recurrence override")
	}
}
