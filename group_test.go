package jscal

import (
	"strings"
	"testing"
	"time"
)

func TestNewGroup(t *testing.T) {
	uid := "group-123"
	title := "Test Group"

	group := NewGroup(uid, title)

	if group.Type != "Group" {
		t.Errorf("Expected Type to be 'Group', got '%s'", group.Type)
	}

	if group.UID != uid {
		t.Errorf("Expected UID to be '%s', got '%s'", uid, group.UID)
	}

	if group.Title == nil || *group.Title != title {
		t.Errorf("Expected Title to be '%s', got '%v'", title, group.Title)
	}

	if group.Created == nil {
		t.Error("Expected Created to be set")
	}

	if group.Updated == nil {
		t.Error("Expected Updated to be set")
	}

	if group.Sequence == nil || *group.Sequence != 0 {
		t.Errorf("Expected Sequence to be 0, got %v", group.Sequence)
	}

	if group.Entries == nil {
		t.Error("Expected Entries to be initialized")
	}

	if len(group.Entries) != 0 {
		t.Errorf("Expected empty Entries, got %d entries", len(group.Entries))
	}
}

func TestGroupJSON(t *testing.T) {
	group := NewGroup("group-123", "Test Group")
	group.Description = String("A test group")

	// Test basic group without entries first
	jsonData, err := group.JSON()
	if err != nil {
		t.Fatalf("Failed to marshal to JSON: %v", err)
	}

	// Parse it back
	parsedGroup, err := ParseGroup(jsonData)
	if err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	if parsedGroup.UID != group.UID {
		t.Errorf("Expected UID to be '%s', got '%s'", group.UID, parsedGroup.UID)
	}

	if parsedGroup.Title == nil || *parsedGroup.Title != *group.Title {
		t.Errorf("Expected Title to be '%s', got '%v'", *group.Title, parsedGroup.Title)
	}

	// Note: Groups with entries would need custom JSON unmarshaling
	// to properly deserialize entries as Event/Task objects.
	// This is a known limitation that would be addressed in production.
}

func TestGroupAddEntry(t *testing.T) {
	group := NewGroup("group-123", "Test Group")

	// Add an event
	event := NewEvent("event-1", "Event 1")
	err := group.AddEntry(event)
	if err != nil {
		t.Fatalf("Failed to add event: %v", err)
	}

	if len(group.Entries) != 1 {
		t.Errorf("Expected 1 entry, got %d", len(group.Entries))
	}

	// Add a task
	task := NewTask("task-1", "Task 1")
	err = group.AddEntry(task)
	if err != nil {
		t.Fatalf("Failed to add task: %v", err)
	}

	if len(group.Entries) != 2 {
		t.Errorf("Expected 2 entries, got %d", len(group.Entries))
	}

	// Try to add duplicate UID
	duplicateEvent := NewEvent("event-1", "Duplicate Event")
	err = group.AddEntry(duplicateEvent)
	if err == nil {
		t.Error("Expected error when adding duplicate UID")
	}
	if !strings.Contains(err.Error(), "already exists") {
		t.Errorf("Expected error about duplicate UID, got: %v", err)
	}

	// Try to add nil entry
	err = group.AddEntry(nil)
	if err == nil {
		t.Error("Expected error when adding nil entry")
	}
}

func TestGroupRemoveEntry(t *testing.T) {
	group := NewGroup("group-123", "Test Group")

	// Add entries
	event := NewEvent("event-1", "Event 1")
	task := NewTask("task-1", "Task 1")
	_ = group.AddEntry(event)
	_ = group.AddEntry(task)

	// Remove event
	err := group.RemoveEntry("event-1")
	if err != nil {
		t.Fatalf("Failed to remove entry: %v", err)
	}

	if len(group.Entries) != 1 {
		t.Errorf("Expected 1 entry after removal, got %d", len(group.Entries))
	}

	// Try to remove non-existent entry
	err = group.RemoveEntry("non-existent")
	if err == nil {
		t.Error("Expected error when removing non-existent entry")
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("Expected error about entry not found, got: %v", err)
	}
}

func TestGroupGetEntry(t *testing.T) {
	group := NewGroup("group-123", "Test Group")

	// Add entries
	event := NewEvent("event-1", "Event 1")
	task := NewTask("task-1", "Task 1")
	_ = group.AddEntry(event)
	_ = group.AddEntry(task)

	// Get event
	entry := group.GetEntry("event-1")
	if entry == nil {
		t.Fatal("Failed to get entry")
	}
	if entry.GetUID() != "event-1" {
		t.Errorf("Expected UID 'event-1', got '%s'", entry.GetUID())
	}
	if entry.GetType() != "Event" {
		t.Errorf("Expected type 'Event', got '%s'", entry.GetType())
	}

	// Get non-existent entry
	entry = group.GetEntry("non-existent")
	if entry != nil {
		t.Error("Expected nil for non-existent entry")
	}
}

func TestGroupGetEventsAndTasks(t *testing.T) {
	group := NewGroup("group-123", "Test Group")

	// Add mixed entries
	event1 := NewEvent("event-1", "Event 1")
	event2 := NewEvent("event-2", "Event 2")
	task1 := NewTask("task-1", "Task 1")
	task2 := NewTask("task-2", "Task 2")

	_ = group.AddEntry(event1)
	_ = group.AddEntry(task1)
	_ = group.AddEntry(event2)
	_ = group.AddEntry(task2)

	// Get events
	events := group.GetEvents()
	if len(events) != 2 {
		t.Errorf("Expected 2 events, got %d", len(events))
	}
	if events[0].UID != "event-1" || events[1].UID != "event-2" {
		t.Error("Events not retrieved correctly")
	}

	// Get tasks
	tasks := group.GetTasks()
	if len(tasks) != 2 {
		t.Errorf("Expected 2 tasks, got %d", len(tasks))
	}
	if tasks[0].UID != "task-1" || tasks[1].UID != "task-2" {
		t.Error("Tasks not retrieved correctly")
	}
}

func TestGroupCounts(t *testing.T) {
	group := NewGroup("group-123", "Test Group")

	// Initially empty
	if group.CountEntries() != 0 {
		t.Errorf("Expected 0 entries, got %d", group.CountEntries())
	}
	if group.CountEvents() != 0 {
		t.Errorf("Expected 0 events, got %d", group.CountEvents())
	}
	if group.CountTasks() != 0 {
		t.Errorf("Expected 0 tasks, got %d", group.CountTasks())
	}

	// Add entries
	_ = group.AddEntry(NewEvent("event-1", "Event 1"))
	_ = group.AddEntry(NewTask("task-1", "Task 1"))
	_ = group.AddEntry(NewEvent("event-2", "Event 2"))

	if group.CountEntries() != 3 {
		t.Errorf("Expected 3 total entries, got %d", group.CountEntries())
	}
	if group.CountEvents() != 2 {
		t.Errorf("Expected 2 events, got %d", group.CountEvents())
	}
	if group.CountTasks() != 1 {
		t.Errorf("Expected 1 task, got %d", group.CountTasks())
	}
}

func TestGroupClone(t *testing.T) {
	// Create a complex group
	original := NewGroup("group-123", "Original Group")
	original.Description = String("Original description")
	original.Categories = map[string]bool{"calendar": true, "work": true}
	original.Keywords = map[string]bool{"project": true, "team": true}
	original.Color = String("#FF5733")

	// Add entries
	event := NewEvent("event-1", "Event 1")
	task := NewTask("task-1", "Task 1")
	_ = original.AddEntry(event)
	_ = original.AddEntry(task)

	// Add link
	link := NewLink("https://example.com/group")
	original.AddLink("group-link", link)

	// Clone the group
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
	if len(cloned.Categories) != len(original.Categories) {
		t.Error("Clone() should preserve Categories")
	}
	if len(cloned.Keywords) != len(original.Keywords) {
		t.Error("Clone() should preserve Keywords")
	}
	if len(cloned.Entries) != len(original.Entries) {
		t.Error("Clone() should preserve Entries count")
	}

	// Modify the clone and verify original is unchanged
	cloned.Title = String("Modified Title")
	if *original.Title == "Modified Title" {
		t.Error("Modifying clone should not affect original")
	}
}

func TestGroupValidation(t *testing.T) {
	tests := []struct {
		name    string
		group   *Group
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid group",
			group: &Group{
				Type:    "Group",
				UID:     "group-123",
				Title:   String("Valid Group"),
				Entries: []CalendarObject{},
			},
			wantErr: false,
		},
		{
			name: "missing type",
			group: &Group{
				UID:     "group-123",
				Title:   String("Group"),
				Entries: []CalendarObject{},
			},
			wantErr: true,
			errMsg:  "must be 'Group'",
		},
		{
			name: "wrong type",
			group: &Group{
				Type:    "Event",
				UID:     "group-123",
				Title:   String("Group"),
				Entries: []CalendarObject{},
			},
			wantErr: true,
			errMsg:  "must be 'Group'",
		},
		{
			name: "missing UID",
			group: &Group{
				Type:    "Group",
				Title:   String("Group"),
				Entries: []CalendarObject{},
			},
			wantErr: true,
			errMsg:  "is required",
		},
		{
			name: "group with valid entries",
			group: &Group{
				Type:  "Group",
				UID:   "group-123",
				Title: String("Group with entries"),
				Entries: []CalendarObject{
					NewEvent("event-1", "Event 1"),
					NewTask("task-1", "Task 1"),
				},
			},
			wantErr: false,
		},
		{
			name: "group with duplicate UIDs",
			group: &Group{
				Type:  "Group",
				UID:   "group-123",
				Title: String("Group with duplicates"),
				Entries: []CalendarObject{
					NewEvent("same-id", "Event 1"),
					NewTask("same-id", "Task 1"),
				},
			},
			wantErr: true,
			errMsg:  "duplicate UID",
		},
		{
			name: "group containing itself",
			group: &Group{
				Type:  "Group",
				UID:   "group-123",
				Title: String("Circular reference"),
				Entries: []CalendarObject{
					&Group{Type: "Group", UID: "group-123"}, // Same UID as parent
				},
			},
			wantErr: true,
			errMsg:  "cannot contain itself",
		},
		{
			name: "negative sequence",
			group: &Group{
				Type:     "Group",
				UID:      "group-123",
				Sequence: Int(-1),
				Entries:  []CalendarObject{},
			},
			wantErr: true,
			errMsg:  "cannot be negative",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.group.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && tt.errMsg != "" {
				if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("Validate() error message = %v, want to contain %v", err.Error(), tt.errMsg)
				}
			}
		})
	}
}

func TestParseGroup(t *testing.T) {
	jsonData := []byte(`{
		"@type": "Group",
		"uid": "group-001",
		"title": "Project Tasks",
		"description": "All tasks for Q1 project",
		"categories": {"work": true, "project": true},
		"entries": []
	}`)

	group, err := ParseGroup(jsonData)
	if err != nil {
		t.Fatalf("Failed to parse group: %v", err)
	}

	if group.UID != "group-001" {
		t.Errorf("Expected UID 'group-001', got '%s'", group.UID)
	}

	if group.Title == nil || *group.Title != "Project Tasks" {
		t.Errorf("Expected title 'Project Tasks', got '%v'", group.Title)
	}

	if len(group.Categories) != 2 {
		t.Errorf("Expected 2 categories, got %d", len(group.Categories))
	}
}

func TestParseAllGroups(t *testing.T) {
	jsonData := []byte(`[
		{
			"@type": "Group",
			"uid": "group-001",
			"title": "Group 1",
			"entries": []
		},
		{
			"@type": "Group",
			"uid": "group-002",
			"title": "Group 2",
			"entries": []
		}
	]`)

	groups, err := ParseAllGroups(jsonData)
	if err != nil {
		t.Fatalf("Failed to parse groups: %v", err)
	}

	if len(groups) != 2 {
		t.Fatalf("Expected 2 groups, got %d", len(groups))
	}

	if groups[0].UID != "group-001" {
		t.Errorf("Expected first group UID 'group-001', got '%s'", groups[0].UID)
	}

	if groups[1].UID != "group-002" {
		t.Errorf("Expected second group UID 'group-002', got '%s'", groups[1].UID)
	}
}

func TestGroupTouch(t *testing.T) {
	group := NewGroup("group-123", "Test Group")

	originalUpdated := group.Updated
	originalSequence := *group.Sequence

	// Wait a bit to ensure timestamp difference
	time.Sleep(time.Millisecond)

	// Add an entry which should trigger Touch
	event := NewEvent("event-1", "Event 1")
	_ = group.AddEntry(event)

	if group.Updated.Equal(*originalUpdated) {
		t.Error("Expected Updated timestamp to change")
	}

	if *group.Sequence != originalSequence+1 {
		t.Errorf("Expected sequence to increment from %d to %d, got %d",
			originalSequence, originalSequence+1, *group.Sequence)
	}
}

func TestGroupHelperMethods(t *testing.T) {
	group := NewGroup("group-123", "Test Group")

	// Test AddKeyword
	group.AddKeyword("important")
	group.AddKeyword("urgent")
	if len(group.Keywords) != 2 {
		t.Errorf("Expected 2 keywords, got %d", len(group.Keywords))
	}
	if !group.Keywords["important"] || !group.Keywords["urgent"] {
		t.Error("Keywords not added correctly")
	}

	// Test AddCategory
	group.AddCategory("work")
	group.AddCategory("project")
	if len(group.Categories) != 2 {
		t.Errorf("Expected 2 categories, got %d", len(group.Categories))
	}
	if !group.Categories["work"] || !group.Categories["project"] {
		t.Error("Categories not added correctly")
	}

	// Test AddLink
	link := NewLink("https://example.com/group")
	group.AddLink("main", link)
	if len(group.Links) != 1 {
		t.Errorf("Expected 1 link, got %d", len(group.Links))
	}
	if group.Links["main"] == nil || group.Links["main"].Href != "https://example.com/group" {
		t.Error("Link not added correctly")
	}
}

func TestGroupCalendarObjectInterface(t *testing.T) {
	group := NewGroup("group-123", "Test Group")

	// Group should implement CalendarObject
	var obj CalendarObject = group

	if obj.GetUID() != "group-123" {
		t.Errorf("GetUID() = %s, want 'group-123'", obj.GetUID())
	}

	if obj.GetType() != "Group" {
		t.Errorf("GetType() = %s, want 'Group'", obj.GetType())
	}

	err := obj.Validate()
	if err != nil {
		t.Errorf("Validate() error = %v", err)
	}
}
