package jscal

import (
	"encoding/json"
	"strings"
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
	parsedEvent, err := ParseEvent(jsonData)
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

func TestEventPrettyJSON(t *testing.T) {
	event := NewEvent("test-123", "Test Event")
	event.Description = String("This is a test event")
	event.Start = NewLocalDateTime(time.Date(2025, 3, 1, 14, 0, 0, 0, time.UTC))
	event.Duration = String("PT1H")
	
	pretty, err := event.PrettyJSON()
	if err != nil {
		t.Fatalf("PrettyJSON() error = %v", err)
	}
	
	// Check that it's properly formatted with indentation
	lines := strings.Split(string(pretty), "\n")
	if len(lines) < 5 {
		t.Error("PrettyJSON should produce multiple lines")
	}
	
	// Check for indentation
	hasIndentation := false
	for _, line := range lines {
		if strings.HasPrefix(line, "  ") || strings.HasPrefix(line, "\t") {
			hasIndentation = true
			break
		}
	}
	if !hasIndentation {
		t.Error("PrettyJSON should include indentation")
	}
	
	// Verify it's valid JSON
	var decoded Event
	if err := json.Unmarshal(pretty, &decoded); err != nil {
		t.Errorf("PrettyJSON output is not valid JSON: %v", err)
	}
	
	// Verify content
	if decoded.UID != event.UID {
		t.Error("PrettyJSON should preserve UID")
	}
	if decoded.Title == nil || *decoded.Title != *event.Title {
		t.Error("PrettyJSON should preserve Title")
	}
}

func TestEventClone(t *testing.T) {
	// Create a complex event with many properties
	original := NewEvent("test-123", "Original Event")
	original.Description = String("Original description")
	original.Start = NewLocalDateTime(time.Date(2025, 3, 1, 14, 0, 0, 0, time.UTC))
	original.Duration = String("PT1H")
	original.Status = String(StatusConfirmed)
	original.Privacy = String(PrivacyPrivate)
	original.Sequence = Int(2)
	original.Priority = Int(5)
	original.Categories = map[string]bool{"meeting": true, "important": true}
	original.Keywords = map[string]bool{"project": true, "deadline": true}
	
	// Add participant
	participant := NewParticipant("John Doe", "john@example.com")
	participant.ParticipationStatus = String(ParticipationAccepted)
	original.AddParticipant("john@example.com", participant)
	
	// Add location
	location := NewLocation("Conference Room")
	original.AddLocation("loc1", location)
	
	// Add alert
	alert := &Alert{
		Trigger: &OffsetTrigger{
			Offset: "-PT15M",
		},
	}
	original.AddAlert("alert1", alert)
	
	// Add link
	link := NewLink("https://example.com/event")
	original.AddLink("link1", link)
	
	// Add recurrence
	rule := &RecurrenceRule{
		Frequency: FrequencyWeekly,
		ByDay:     []NDay{{Day: "mo"}, {Day: "we"}},
	}
	original.RecurrenceRules = []RecurrenceRule{*rule}
	
	// Clone the event
	cloned := original.Clone()
	
	// Verify it's a different instance
	if cloned == original {
		t.Error("Clone() should return a new instance")
	}
	
	// Verify all fields are copied
	if cloned.UID != original.UID {
		t.Error("Clone() should preserve UID")
	}
	if cloned.Title == nil || *cloned.Title != *original.Title {
		t.Error("Clone() should preserve Title")
	}
	if cloned.Description == nil || *cloned.Description != *original.Description {
		t.Error("Clone() should preserve Description")
	}
	if !cloned.Start.Equal(original.Start) {
		t.Error("Clone() should preserve Start")
	}
	if cloned.Duration == nil || *cloned.Duration != *original.Duration {
		t.Error("Clone() should preserve Duration")
	}
	if cloned.Status == nil || *cloned.Status != *original.Status {
		t.Error("Clone() should preserve Status")
	}
	if cloned.Privacy == nil || *cloned.Privacy != *original.Privacy {
		t.Error("Clone() should preserve Privacy")
	}
	if cloned.Sequence == nil || *cloned.Sequence != *original.Sequence {
		t.Error("Clone() should preserve Sequence")
	}
	if cloned.Priority == nil || *cloned.Priority != *original.Priority {
		t.Error("Clone() should preserve Priority")
	}
	
	// Verify collections are deep copied
	if len(cloned.Categories) != len(original.Categories) {
		t.Error("Clone() should preserve Categories")
	}
	// Keywords is now a map
	if len(cloned.Keywords) != len(original.Keywords) {
		t.Error("Clone() should preserve Keywords")
	}
	if len(cloned.Participants) != len(original.Participants) {
		t.Error("Clone() should preserve Participants")
	}
	if len(cloned.Locations) != len(original.Locations) {
		t.Error("Clone() should preserve Locations")
	}
	if len(cloned.Alerts) != len(original.Alerts) {
		t.Error("Clone() should preserve Alerts")
	}
	if len(cloned.Links) != len(original.Links) {
		t.Error("Clone() should preserve Links")
	}
	if len(cloned.RecurrenceRules) != len(original.RecurrenceRules) {
		t.Error("Clone() should preserve RecurrenceRules")
	}
	
	// Modify the clone and verify original is unchanged
	cloned.Title = String("Modified Title")
	if *original.Title == "Modified Title" {
		t.Error("Modifying clone should not affect original")
	}
	
	// Modify collections in clone
	cloned.Categories["new-category"] = true
	if original.Categories["new-category"] {
		t.Error("Modifying clone's Categories should not affect original")
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
	event.Start = NewLocalDateTime(startTime)
	event.Duration = String("PT1H")
	event.TimeZone = String("UTC")

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

func TestAddLocation(t *testing.T) {
	event := NewEvent("test-123", "Test Event")
	
	// Add first location
	loc1 := NewLocation("Conference Room A")
	loc1.Description = String("Main conference room")
	event.AddLocation("main", loc1)
	
	if len(event.Locations) != 1 {
		t.Error("AddLocation should add location")
	}
	if event.Locations["main"] == nil {
		t.Error("Location should be accessible by ID")
	}
	if event.Locations["main"].Name == nil || *event.Locations["main"].Name != "Conference Room A" {
		t.Error("Location details should be preserved")
	}
	
	// Add second location
	loc2 := NewLocation("Conference Room B")
	event.AddLocation("backup", loc2)
	
	if len(event.Locations) != 2 {
		t.Error("Should have 2 locations")
	}
	
	// Override existing location
	loc3 := NewLocation("Conference Room C")
	event.AddLocation("main", loc3)
	
	if len(event.Locations) != 2 {
		t.Error("Overriding should not increase count")
	}
	if event.Locations["main"].Name == nil || *event.Locations["main"].Name != "Conference Room C" {
		t.Error("Location should be overridden")
	}
}

func TestAddVirtualLocation(t *testing.T) {
	event := NewEvent("test-123", "Test Event")
	
	// Add first virtual location
	vLoc1 := NewVirtualLocation("Team Standup", "https://zoom.us/j/123456789")
	event.AddVirtualLocation("zoom", vLoc1)
	
	if len(event.VirtualLocations) != 1 {
		t.Error("AddVirtualLocation should add virtual location")
	}
	if event.VirtualLocations["zoom"] == nil {
		t.Error("Virtual location should be accessible by ID")
	}
	if event.VirtualLocations["zoom"].URI != "https://zoom.us/j/123456789" {
		t.Error("Virtual location URI should be preserved")
	}
	
	// Add second virtual location
	vLoc2 := NewVirtualLocation("Meeting", "https://meet.google.com/abc-defg-hij")
	event.AddVirtualLocation("meet", vLoc2)
	
	if len(event.VirtualLocations) != 2 {
		t.Error("Should have 2 virtual locations")
	}
	
	// Add virtual location with features
	vLoc3 := NewVirtualLocation("Team Meeting", "https://teams.microsoft.com/meeting")
	vLoc3.Features = map[string]bool{"audio": true, "video": true, "screen-share": true}
	event.AddVirtualLocation("teams", vLoc3)
	
	if len(event.VirtualLocations) != 3 {
		t.Error("Should have 3 virtual locations")
	}
	if !event.VirtualLocations["teams"].Features["screen-share"] {
		t.Error("Virtual location features should be preserved")
	}
}

func TestAddAlert(t *testing.T) {
	event := NewEvent("test-123", "Test Event")
	
	// Add first alert
	alert1 := &Alert{
		Trigger: &OffsetTrigger{
			Offset: "-PT15M",
		},
		Action: String("display"),
	}
	event.AddAlert("15min", alert1)
	
	if len(event.Alerts) != 1 {
		t.Error("AddAlert should add alert")
	}
	if event.Alerts["15min"] == nil {
		t.Error("Alert should be accessible by ID")
	}
	if event.Alerts["15min"].Trigger.Offset != "-PT15M" {
		t.Error("Alert details should be preserved")
	}
	
	// Add second alert
	alert2 := &Alert{
		Trigger: &OffsetTrigger{
			Offset: "-PT1H",
		},
		Action: String("email"),
	}
	event.AddAlert("1hour", alert2)
	
	if len(event.Alerts) != 2 {
		t.Error("Should have 2 alerts")
	}
	
	// Add alert with absolute time
	alert3 := &Alert{
		Trigger: &OffsetTrigger{
			Offset: "-PT30M",
		},
	}
	event.AddAlert("absolute", alert3)
	
	if len(event.Alerts) != 3 {
		t.Error("Should have 3 alerts")
	}
}

func TestAddKeyword(t *testing.T) {
	event := NewEvent("test-123", "Test Event")
	
	// Add keywords
	event.AddKeyword("urgent")
	if len(event.Keywords) != 1 || !event.Keywords["urgent"] {
		t.Error("AddKeyword should add keyword")
	}
	
	// Add more keywords
	event.AddKeyword("project")
	event.AddKeyword("deadline")
	
	if len(event.Keywords) != 3 {
		t.Error("Should have 3 keywords")
	}
	
	// Verify keywords exist
	if !event.Keywords["urgent"] || !event.Keywords["project"] || !event.Keywords["deadline"] {
		t.Error("Keywords not added correctly")
	}
	
	// Add duplicate keyword (should not duplicate in map)
	event.AddKeyword("urgent")
	if len(event.Keywords) != 3 {
		t.Error("Map-based keywords should not have duplicates")
	}
}

func TestAddLink(t *testing.T) {
	event := NewEvent("test-123", "Test Event")
	
	// Add first link
	link1 := NewLink("https://example.com/event")
	link1.Title = String("Event Details")
	event.AddLink("details", link1)
	
	if len(event.Links) != 1 {
		t.Error("AddLink should add link")
	}
	if event.Links["details"] == nil {
		t.Error("Link should be accessible by ID")
	}
	if event.Links["details"].Href != "https://example.com/event" {
		t.Error("Link href should be preserved")
	}
	
	// Add second link
	link2 := NewLink("https://example.com/agenda.pdf")
	link2.ContentType = String("application/pdf")
	link2.Size = Int(204800)
	event.AddLink("agenda", link2)
	
	if len(event.Links) != 2 {
		t.Error("Should have 2 links")
	}
	
	// Add link with all properties
	link3 := NewLink("https://example.com/icon.png")
	link3.ContentType = String("image/png")
	link3.Size = Int(4096)
	link3.Rel = String("icon")
	link3.Display = String("badge")
	link3.Title = String("Event Icon")
	event.AddLink("icon", link3)
	
	if len(event.Links) != 3 {
		t.Error("Should have 3 links")
	}
	if event.Links["icon"].Rel == nil || *event.Links["icon"].Rel != "icon" {
		t.Error("Link rel should be preserved")
	}
}

func TestSetRecurrence(t *testing.T) {
	event := NewEvent("test-123", "Test Event")
	
	// Set simple daily recurrence
	rule1 := &RecurrenceRule{
		Frequency: FrequencyDaily,
		Count:     Int(10),
	}
	event.RecurrenceRules = []RecurrenceRule{*rule1}
	
	if !event.IsRecurring() {
		t.Error("Event should be recurring")
	}
	if len(event.RecurrenceRules) != 1 {
		t.Error("Should have 1 recurrence rule")
	}
	
	// Set weekly recurrence with specific days
	rule2 := &RecurrenceRule{
		Frequency: FrequencyWeekly,
		ByDay:     []NDay{{Day: "mo"}, {Day: "we"}, {Day: "fr"}},
		Until:     NewLocalDateTime(time.Time{}),
	}
	event.RecurrenceRules = []RecurrenceRule{*rule2}
	
	if len(event.RecurrenceRules) != 1 {
		t.Error("SetRecurrence should replace existing rules")
	}
	if event.RecurrenceRules[0].Frequency != FrequencyWeekly {
		t.Error("New rule should be set")
	}
	
	// Set recurrence with overrides
	rule3 := &RecurrenceRule{
		Frequency: FrequencyMonthly,
		ByMonthDay: []int{15},
	}
	event.RecurrenceRules = []RecurrenceRule{*rule3}
	// Manually add overrides
	event.RecurrenceOverrides = map[string]map[string]interface{}{
		"20250415T140000": {
			"uid":      "test-123-20250415",
			"title":    "Special Instance",
			"sequence": 1,
		},
	}
	
	if len(event.RecurrenceOverrides) != 1 {
		t.Error("Should have 1 recurrence override")
	}
	if title, ok := event.RecurrenceOverrides["20250415T140000"]["title"].(string); !ok || title != "Special Instance" {
		t.Error("Override should be preserved")
	}
	
	// Clear recurrence
	event.RecurrenceRules = nil
	event.RecurrenceOverrides = nil
	if event.IsRecurring() {
		t.Error("Event should not be recurring after clearing")
	}
	if len(event.RecurrenceRules) != 0 {
		t.Error("RecurrenceRules should be empty")
	}
	if len(event.RecurrenceOverrides) != 0 {
		t.Error("RecurrenceOverrides should be empty")
	}
}

func TestIsRecurring(t *testing.T) {
	event := NewEvent("test-123", "Test Event")
	
	// Initially not recurring
	if event.IsRecurring() {
		t.Error("New event should not be recurring")
	}
	
	// Add recurrence rule
	rule := &RecurrenceRule{
		Frequency: FrequencyWeekly,
	}
	event.RecurrenceRules = []RecurrenceRule{*rule}
	
	if !event.IsRecurring() {
		t.Error("Event with recurrence rule should be recurring")
	}
	
	// Empty rules slice
	event.RecurrenceRules = []RecurrenceRule{}
	if event.IsRecurring() {
		t.Error("Event with empty rules should not be recurring")
	}
	
	// Nil rules
	event.RecurrenceRules = nil
	if event.IsRecurring() {
		t.Error("Event with nil rules should not be recurring")
	}
}