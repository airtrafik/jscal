package jscal

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func TestNewTask(t *testing.T) {
	uid := "task-123"
	title := "Test Task"

	task := NewTask(uid, title)

	if task.Type != "Task" {
		t.Errorf("Expected Type to be 'Task', got '%s'", task.Type)
	}

	if task.UID != uid {
		t.Errorf("Expected UID to be '%s', got '%s'", uid, task.UID)
	}

	if task.Title == nil || *task.Title != title {
		t.Errorf("Expected Title to be '%s', got '%v'", title, task.Title)
	}

	if task.Created == nil {
		t.Error("Expected Created to be set")
	}

	if task.Updated == nil {
		t.Error("Expected Updated to be set")
	}

	if task.Sequence == nil || *task.Sequence != 0 {
		t.Errorf("Expected Sequence to be 0, got %v", task.Sequence)
	}

	if task.Progress == nil || *task.Progress != ProgressNeedsAction {
		t.Errorf("Expected Progress to be '%s', got '%v'", ProgressNeedsAction, task.Progress)
	}
}

func TestTaskJSON(t *testing.T) {
	task := NewTask("task-123", "Test Task")
	task.Due = NewLocalDateTime(time.Now().Add(24 * time.Hour))
	task.EstimatedDuration = String("PT2H")

	jsonData, err := task.JSON()
	if err != nil {
		t.Fatalf("Failed to marshal to JSON: %v", err)
	}

	// Parse it back
	parsedTask, err := ParseTask(jsonData)
	if err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	if parsedTask.UID != task.UID {
		t.Errorf("Expected UID to be '%s', got '%s'", task.UID, parsedTask.UID)
	}

	if parsedTask.Title == nil || *parsedTask.Title != *task.Title {
		t.Errorf("Expected Title to be '%s', got '%v'", *task.Title, parsedTask.Title)
	}
}

func TestTaskPrettyJSON(t *testing.T) {
	task := NewTask("task-123", "Test Task")
	task.Description = String("This is a test task")
	task.Due = NewLocalDateTime(time.Date(2025, 3, 1, 14, 0, 0, 0, time.UTC))
	task.EstimatedDuration = String("PT2H")

	pretty, err := task.PrettyJSON()
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
	var decoded Task
	if err := json.Unmarshal(pretty, &decoded); err != nil {
		t.Errorf("PrettyJSON output is not valid JSON: %v", err)
	}

	// Verify content
	if decoded.UID != task.UID {
		t.Error("PrettyJSON should preserve UID")
	}
	if decoded.Title == nil || *decoded.Title != *task.Title {
		t.Error("PrettyJSON should preserve Title")
	}
}

func TestTaskClone(t *testing.T) {
	// Create a complex task with many properties
	original := NewTask("task-123", "Original Task")
	original.Description = String("Original description")
	original.Start = NewLocalDateTime(time.Date(2025, 3, 1, 9, 0, 0, 0, time.UTC))
	original.Due = NewLocalDateTime(time.Date(2025, 3, 1, 17, 0, 0, 0, time.UTC))
	original.EstimatedDuration = String("PT8H")
	original.Status = String("in-process")
	original.Privacy = String(PrivacyPrivate)
	original.Sequence = Int(2)
	original.Priority = Int(5)
	original.PercentComplete = Int(50)
	original.Progress = String(ProgressInProcess)
	original.Categories = map[string]bool{"work": true, "important": true}
	original.Keywords = map[string]bool{"project": true, "deadline": true}

	// Add participant
	participant := NewParticipant("John Doe", "john@example.com")
	original.AddParticipant("john@example.com", participant)

	// Add location
	location := NewLocation("Office")
	original.AddLocation("loc1", location)

	// Add alert
	alert := &Alert{
		Trigger: &OffsetTrigger{
			Offset: "-PT15M",
		},
	}
	original.AddAlert("alert1", alert)

	// Clone the task
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
	if !cloned.Due.Equal(original.Due) {
		t.Error("Clone() should preserve Due")
	}
	if cloned.EstimatedDuration == nil || *cloned.EstimatedDuration != *original.EstimatedDuration {
		t.Error("Clone() should preserve EstimatedDuration")
	}
	if cloned.Status == nil || *cloned.Status != *original.Status {
		t.Error("Clone() should preserve Status")
	}
	if cloned.PercentComplete == nil || *cloned.PercentComplete != *original.PercentComplete {
		t.Error("Clone() should preserve PercentComplete")
	}
	if cloned.Progress == nil || *cloned.Progress != *original.Progress {
		t.Error("Clone() should preserve Progress")
	}

	// Modify the clone and verify original is unchanged
	cloned.Title = String("Modified Title")
	if *original.Title == "Modified Title" {
		t.Error("Modifying clone should not affect original")
	}
}

func TestTaskIsCompleted(t *testing.T) {
	task := NewTask("task-123", "Test Task")

	// Initially not completed
	if task.IsCompleted() {
		t.Error("New task should not be completed")
	}

	// Mark as completed
	task.Progress = String(ProgressCompleted)
	if !task.IsCompleted() {
		t.Error("Task with completed progress should be completed")
	}

	// Mark as in-process
	task.Progress = String(ProgressInProcess)
	if task.IsCompleted() {
		t.Error("Task with in-process progress should not be completed")
	}
}

func TestTaskIsOverdue(t *testing.T) {
	task := NewTask("task-123", "Test Task")

	// No due date - not overdue
	if task.IsOverdue() {
		t.Error("Task without due date should not be overdue")
	}

	// Due date in future - not overdue
	task.Due = NewLocalDateTime(time.Now().Add(24 * time.Hour))
	if task.IsOverdue() {
		t.Error("Task with future due date should not be overdue")
	}

	// Due date in past - overdue
	task.Due = NewLocalDateTime(time.Now().Add(-24 * time.Hour))
	if !task.IsOverdue() {
		t.Error("Task with past due date should be overdue")
	}

	// Completed task - not overdue even if past due
	task.Progress = String(ProgressCompleted)
	if task.IsOverdue() {
		t.Error("Completed task should not be overdue")
	}
}

func TestTaskSetProgress(t *testing.T) {
	task := NewTask("task-123", "Test Task")

	originalUpdated := task.Updated
	originalSequence := *task.Sequence

	// Wait a bit to ensure timestamp difference
	time.Sleep(time.Millisecond)

	// Set progress
	task.SetProgress(ProgressInProcess, 25)

	if task.Progress == nil || *task.Progress != ProgressInProcess {
		t.Errorf("Expected Progress to be '%s', got '%v'", ProgressInProcess, task.Progress)
	}

	if task.PercentComplete == nil || *task.PercentComplete != 25 {
		t.Errorf("Expected PercentComplete to be 25, got %v", task.PercentComplete)
	}

	if task.ProgressUpdated == nil {
		t.Error("Expected ProgressUpdated to be set")
	}

	if task.Updated.Equal(*originalUpdated) {
		t.Error("Expected Updated timestamp to change")
	}

	if *task.Sequence != originalSequence+1 {
		t.Errorf("Expected sequence to increment from %d to %d, got %d",
			originalSequence, originalSequence+1, *task.Sequence)
	}
}

func TestTaskGetEstimatedDuration(t *testing.T) {
	task := NewTask("task-123", "Test Task")

	// No duration specified
	_, err := task.GetEstimatedDuration()
	if err == nil {
		t.Error("Expected error when no estimated duration is set")
	}

	// Set duration to 2 hours
	task.EstimatedDuration = String("PT2H")
	duration, err := task.GetEstimatedDuration()
	if err != nil {
		t.Fatalf("Failed to get estimated duration: %v", err)
	}

	expected := 2 * time.Hour
	if duration != expected {
		t.Errorf("Expected duration to be %v, got %v", expected, duration)
	}
}

func TestTaskGetTimeToComplete(t *testing.T) {
	task := NewTask("task-123", "Test Task")
	task.EstimatedDuration = String("PT4H")

	// No progress - full time remaining
	remaining, err := task.GetTimeToComplete()
	if err != nil {
		t.Fatalf("Failed to get time to complete: %v", err)
	}
	if remaining != 4*time.Hour {
		t.Errorf("Expected 4 hours remaining, got %v", remaining)
	}

	// 25% complete - 75% remaining
	task.PercentComplete = Int(25)
	remaining, err = task.GetTimeToComplete()
	if err != nil {
		t.Fatalf("Failed to get time to complete: %v", err)
	}
	if remaining != 3*time.Hour {
		t.Errorf("Expected 3 hours remaining, got %v", remaining)
	}

	// 50% complete - 50% remaining
	task.PercentComplete = Int(50)
	remaining, err = task.GetTimeToComplete()
	if err != nil {
		t.Fatalf("Failed to get time to complete: %v", err)
	}
	if remaining != 2*time.Hour {
		t.Errorf("Expected 2 hours remaining, got %v", remaining)
	}

	// Completed - no time remaining
	task.Progress = String(ProgressCompleted)
	remaining, err = task.GetTimeToComplete()
	if err != nil {
		t.Fatalf("Failed to get time to complete: %v", err)
	}
	if remaining != 0 {
		t.Errorf("Expected 0 time remaining for completed task, got %v", remaining)
	}
}

func TestTaskValidation(t *testing.T) {
	tests := []struct {
		name    string
		task    *Task
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid task",
			task: &Task{
				Type:  "Task",
				UID:   "task-123",
				Title: String("Valid Task"),
			},
			wantErr: false,
		},
		{
			name: "missing type",
			task: &Task{
				UID:   "task-123",
				Title: String("Task"),
			},
			wantErr: true,
			errMsg:  "must be 'Task'",
		},
		{
			name: "wrong type",
			task: &Task{
				Type:  "Event",
				UID:   "task-123",
				Title: String("Task"),
			},
			wantErr: true,
			errMsg:  "must be 'Task'",
		},
		{
			name: "missing UID",
			task: &Task{
				Type:  "Task",
				Title: String("Task"),
			},
			wantErr: true,
			errMsg:  "is required",
		},
		{
			name: "invalid progress",
			task: &Task{
				Type:     "Task",
				UID:      "task-123",
				Progress: String("invalid"),
			},
			wantErr: true,
			errMsg:  "invalid progress value",
		},
		{
			name: "percentComplete too high",
			task: &Task{
				Type:            "Task",
				UID:             "task-123",
				PercentComplete: Int(101),
			},
			wantErr: true,
			errMsg:  "must be between 0 and 100",
		},
		{
			name: "percentComplete negative",
			task: &Task{
				Type:            "Task",
				UID:             "task-123",
				PercentComplete: Int(-1),
			},
			wantErr: true,
			errMsg:  "must be between 0 and 100",
		},
		{
			name: "invalid estimatedDuration",
			task: &Task{
				Type:              "Task",
				UID:               "task-123",
				EstimatedDuration: String("2 hours"),
			},
			wantErr: true,
			errMsg:  "invalid ISO 8601 duration format",
		},
		{
			name: "due before start",
			task: &Task{
				Type:  "Task",
				UID:   "task-123",
				Start: NewLocalDateTime(time.Date(2025, 3, 2, 10, 0, 0, 0, time.UTC)),
				Due:   NewLocalDateTime(time.Date(2025, 3, 1, 10, 0, 0, 0, time.UTC)),
			},
			wantErr: true,
			errMsg:  "due date cannot be before start date",
		},
		{
			name: "invalid status",
			task: &Task{
				Type:   "Task",
				UID:    "task-123",
				Status: String("done"),
			},
			wantErr: true,
			errMsg:  "invalid status",
		},
		{
			name: "invalid priority",
			task: &Task{
				Type:     "Task",
				UID:      "task-123",
				Priority: Int(10),
			},
			wantErr: true,
			errMsg:  "must be between 0 and 9",
		},
		{
			name: "valid with all progress states",
			task: &Task{
				Type:            "Task",
				UID:             "task-123",
				Progress:        String(ProgressInProcess),
				PercentComplete: Int(50),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.task.Validate()
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

func TestParseTask(t *testing.T) {
	jsonData := []byte(`{
		"@type": "Task",
		"uid": "task-001",
		"title": "Complete project proposal",
		"description": "Finish writing the Q1 project proposal",
		"due": "2025-03-15T17:00:00",
		"estimatedDuration": "PT4H",
		"progress": "in-process",
		"percentComplete": 75,
		"priority": 2
	}`)

	task, err := ParseTask(jsonData)
	if err != nil {
		t.Fatalf("Failed to parse task: %v", err)
	}

	if task.UID != "task-001" {
		t.Errorf("Expected UID 'task-001', got '%s'", task.UID)
	}

	if task.Title == nil || *task.Title != "Complete project proposal" {
		t.Errorf("Expected title 'Complete project proposal', got '%v'", task.Title)
	}

	if task.PercentComplete == nil || *task.PercentComplete != 75 {
		t.Errorf("Expected 75%% complete, got %v", task.PercentComplete)
	}

	if task.Progress == nil || *task.Progress != ProgressInProcess {
		t.Errorf("Expected progress 'in-process', got '%v'", task.Progress)
	}
}

func TestParseAllTasks(t *testing.T) {
	jsonData := []byte(`[
		{
			"@type": "Task",
			"uid": "task-001",
			"title": "Task 1",
			"progress": "needs-action"
		},
		{
			"@type": "Task",
			"uid": "task-002",
			"title": "Task 2",
			"progress": "completed",
			"percentComplete": 100
		}
	]`)

	tasks, err := ParseAllTasks(jsonData)
	if err != nil {
		t.Fatalf("Failed to parse tasks: %v", err)
	}

	if len(tasks) != 2 {
		t.Fatalf("Expected 2 tasks, got %d", len(tasks))
	}

	if tasks[0].UID != "task-001" {
		t.Errorf("Expected first task UID 'task-001', got '%s'", tasks[0].UID)
	}

	if tasks[1].Progress == nil || *tasks[1].Progress != ProgressCompleted {
		t.Errorf("Expected second task to be completed, got '%v'", tasks[1].Progress)
	}

	if tasks[1].PercentComplete == nil || *tasks[1].PercentComplete != 100 {
		t.Errorf("Expected second task to be 100%% complete, got %v", tasks[1].PercentComplete)
	}
}
