package jscal

import (
	"testing"
	"time"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		wantUID  string
		wantType string
		wantErr  bool
	}{
		{
			name: "parse Event",
			data: []byte(`{
				"@type": "Event",
				"uid": "event-1",
				"title": "Test Event",
				"start": "2024-01-01T10:00:00"
			}`),
			wantUID:  "event-1",
			wantType: "Event",
			wantErr:  false,
		},
		{
			name: "parse Task",
			data: []byte(`{
				"@type": "Task",
				"uid": "task-1",
				"title": "Test Task",
				"due": "2024-01-01T10:00:00"
			}`),
			wantUID:  "task-1",
			wantType: "Task",
			wantErr:  false,
		},
		{
			name: "parse Group",
			data: []byte(`{
				"@type": "Group",
				"uid": "group-1",
				"name": "Test Group",
				"entries": []
			}`),
			wantUID:  "group-1",
			wantType: "Group",
			wantErr:  false,
		},
		{
			name: "missing @type field",
			data: []byte(`{
				"uid": "test-1",
				"title": "Test"
			}`),
			wantErr: true,
		},
		{
			name: "unknown @type",
			data: []byte(`{
				"@type": "UnknownType",
				"uid": "test-1"
			}`),
			wantErr: true,
		},
		{
			name:    "invalid JSON",
			data:    []byte(`{"@type": "Event", invalid`),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj, err := Parse(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				if obj.GetUID() != tt.wantUID {
					t.Errorf("Parse() UID = %v, want %v", obj.GetUID(), tt.wantUID)
				}
				if obj.GetType() != tt.wantType {
					t.Errorf("Parse() Type = %v, want %v", obj.GetType(), tt.wantType)
				}
			}
		})
	}
}

func TestParseAll(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		want    int
		wantErr bool
	}{
		{
			name:    "empty array",
			data:    []byte(`[]`),
			want:    0,
			wantErr: false,
		},
		{
			name: "mixed types",
			data: []byte(`[
				{
					"@type": "Event",
					"uid": "event-1",
					"title": "Test Event",
					"start": "2024-01-01T10:00:00"
				},
				{
					"@type": "Task",
					"uid": "task-1",
					"title": "Test Task",
					"due": "2024-01-01T10:00:00"
				},
				{
					"@type": "Group",
					"uid": "group-1",
					"name": "Test Group",
					"entries": []
				}
			]`),
			want:    3,
			wantErr: false,
		},
		{
			name: "all Events",
			data: []byte(`[
				{
					"@type": "Event",
					"uid": "event-1",
					"title": "Event 1",
					"start": "2024-01-01T10:00:00"
				},
				{
					"@type": "Event",
					"uid": "event-2",
					"title": "Event 2",
					"start": "2024-01-02T10:00:00"
				}
			]`),
			want:    2,
			wantErr: false,
		},
		{
			name: "all Tasks",
			data: []byte(`[
				{
					"@type": "Task",
					"uid": "task-1",
					"title": "Task 1",
					"due": "2024-01-01T10:00:00"
				},
				{
					"@type": "Task",
					"uid": "task-2",
					"title": "Task 2",
					"due": "2024-01-02T10:00:00"
				}
			]`),
			want:    2,
			wantErr: false,
		},
		{
			name:    "invalid JSON array",
			data:    []byte(`[{"@type": "Event", invalid`),
			wantErr: true,
		},
		{
			name: "object with missing @type",
			data: []byte(`[
				{
					"@type": "Event",
					"uid": "event-1",
					"title": "Event 1",
					"start": "2024-01-01T10:00:00"
				},
				{
					"uid": "unknown-1",
					"title": "Unknown"
				}
			]`),
			wantErr: true,
		},
		{
			name:    "not an array",
			data:    []byte(`{"@type": "Event", "uid": "test-1"}`),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			objects, err := ParseAll(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && len(objects) != tt.want {
				t.Errorf("ParseAll() returned %d objects, want %d", len(objects), tt.want)
			}
		})
	}
}

func TestParseAllMixedTypes(t *testing.T) {
	// Create a mixed collection
	event := NewEvent("event-1", "Team Meeting")
	event.Start = NewLocalDateTime(time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC))

	task := NewTask("task-1", "Complete Project")
	task.Due = NewLocalDateTime(time.Date(2024, 1, 15, 17, 0, 0, 0, time.UTC))
	task.PercentComplete = Int(50)

	group := NewGroup("group-1", "January Events")

	// Marshal them to JSON
	eventJSON, _ := event.JSON()
	taskJSON, _ := task.JSON()
	groupJSON, _ := group.JSON()

	// Create a JSON array string
	arrayJSON := []byte(`[` + string(eventJSON) + `,` + string(taskJSON) + `,` + string(groupJSON) + `]`)

	// Parse the mixed array
	objects, err := ParseAll(arrayJSON)
	if err != nil {
		t.Fatalf("ParseAll() failed: %v", err)
	}

	if len(objects) != 3 {
		t.Fatalf("ParseAll() returned %d objects, want 3", len(objects))
	}

	// Verify types
	expectedTypes := []string{"Event", "Task", "Group"}
	for i, obj := range objects {
		if obj.GetType() != expectedTypes[i] {
			t.Errorf("Object %d has type %s, want %s", i, obj.GetType(), expectedTypes[i])
		}
	}

	// Verify we can type-assert them back
	if _, ok := objects[0].(*Event); !ok {
		t.Error("First object should be an Event")
	}
	if _, ok := objects[1].(*Task); !ok {
		t.Error("Second object should be a Task")
	}
	if _, ok := objects[2].(*Group); !ok {
		t.Error("Third object should be a Group")
	}
}

func TestParseValidation(t *testing.T) {
	// Test that Parse validates the parsed object
	invalidEvent := []byte(`{
		"@type": "Event",
		"uid": "",
		"title": "Invalid Event",
		"start": "2024-01-01T10:00:00"
	}`)

	_, err := Parse(invalidEvent)
	if err == nil {
		t.Error("Parse() should fail for invalid Event with empty UID")
	}

	invalidTask := []byte(`{
		"@type": "Task",
		"uid": "",
		"title": "Invalid Task",
		"due": "2024-01-01T10:00:00"
	}`)

	_, err = Parse(invalidTask)
	if err == nil {
		t.Error("Parse() should fail for invalid Task with empty UID")
	}

	invalidGroup := []byte(`{
		"@type": "Group",
		"uid": "",
		"name": "Invalid Group"
	}`)

	_, err = Parse(invalidGroup)
	if err == nil {
		t.Error("Parse() should fail for invalid Group with empty UID")
	}
}
