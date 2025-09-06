package jscal

import (
	"strings"
	"testing"
	"time"
)

func TestValidationError(t *testing.T) {
	err := ValidationError{
		Field:   "start",
		Message: "is required",
	}

	expected := "start is required"
	if err.Error() != expected {
		t.Errorf("ValidationError.Error() = %v, want %v", err.Error(), expected)
	}
}

func TestValidationErrors(t *testing.T) {
	errs := ValidationErrors{
		{Field: "uid", Message: "UID is required"},
		{Field: "title", Message: "is required"},
		{Field: "start", Message: "is required"},
	}

	errStr := errs.Error()
	if !strings.Contains(errStr, "UID is required") {
		t.Error("ValidationErrors should contain uid error")
	}
	if !strings.Contains(errStr, "title is required") {
		t.Error("ValidationErrors should contain title error")
	}
	if !strings.Contains(errStr, "start is required") {
		t.Error("ValidationErrors should contain start error")
	}
}

func TestValidationErrorsEmpty(t *testing.T) {
	var errs ValidationErrors
	expected := "no validation errors"
	if errs.Error() != expected {
		t.Errorf("Empty ValidationErrors.Error() = %v, want %v", errs.Error(), expected)
	}
}

func TestValidateEvent(t *testing.T) {
	// Helper to create test LocalDateTime
	testStart := NewLocalDateTime(time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC))

	tests := []struct {
		name    string
		event   *Event
		wantErr bool
		errMsg  string
	}{
		{
			name:    "nil event",
			event:   nil,
			wantErr: true,
			errMsg:  "event is nil",
		},
		{
			name:    "empty event",
			event:   &Event{},
			wantErr: true,
			errMsg:  "UID is required",
		},
		{
			name: "missing title",
			event: &Event{
				Type:  "Event",
				UID:   "test-123",
				Start: testStart,
			},
			wantErr: false, // Title is optional per RFC 8984
		},
		{
			name: "missing start",
			event: &Event{
				Type:  "Event",
				UID:   "test-123",
				Title: String("Test Event"),
			},
			wantErr: true,
			errMsg:  "start is required",
		},
		{
			name: "valid minimal event",
			event: &Event{
				Type:  "Event",
				UID:   "test-123",
				Title: String("Test Event"),
				Start: testStart,
			},
			wantErr: false,
		},
		{
			name: "invalid status",
			event: &Event{
				UID:    "test-123",
				Title:  String("Test Event"),
				Status: String("invalid-status"),
			},
			wantErr: true,
			errMsg:  "invalid status",
		},
		{
			name: "valid status confirmed",
			event: &Event{
				Type:   "Event",
				UID:    "test-123",
				Title:  String("Test Event"),
				Start:  testStart,
				Status: String(StatusConfirmed),
			},
			wantErr: false,
		},
		{
			name: "valid status cancelled",
			event: &Event{
				Type:   "Event",
				UID:    "test-123",
				Title:  String("Test Event"),
				Start:  testStart,
				Status: String(StatusCancelled),
			},
			wantErr: false,
		},
		{
			name: "valid status tentative",
			event: &Event{
				Type:   "Event",
				UID:    "test-123",
				Title:  String("Test Event"),
				Start:  testStart,
				Status: String(StatusTentative),
			},
			wantErr: false,
		},
		{
			name: "invalid privacy",
			event: &Event{
				UID:     "test-123",
				Title:   String("Test Event"),
				Privacy: String("top-secret"),
			},
			wantErr: true,
			errMsg:  "invalid privacy",
		},
		{
			name: "valid privacy public",
			event: &Event{
				Type:    "Event",
				UID:     "test-123",
				Title:   String("Test Event"),
				Start:   testStart,
				Privacy: String(PrivacyPublic),
			},
			wantErr: false,
		},
		{
			name: "valid privacy private",
			event: &Event{
				Type:    "Event",
				UID:     "test-123",
				Title:   String("Test Event"),
				Start:   testStart,
				Privacy: String(PrivacyPrivate),
			},
			wantErr: false,
		},
		{
			name: "valid privacy secret",
			event: &Event{
				Type:    "Event",
				UID:     "test-123",
				Title:   String("Test Event"),
				Start:   testStart,
				Privacy: String(PrivacySecret),
			},
			wantErr: false,
		},
		{
			name: "invalid freeBusyStatus",
			event: &Event{
				UID:            "test-123",
				Title:          String("Test Event"),
				FreeBusyStatus: String("maybe-busy"),
			},
			wantErr: true,
			errMsg:  "invalid freeBusyStatus",
		},
		{
			name: "valid freeBusyStatus free",
			event: &Event{
				Type:           "Event",
				UID:            "test-123",
				Title:          String("Test Event"),
				Start:          testStart,
				FreeBusyStatus: String(FreeBusyFree),
			},
			wantErr: false,
		},
		{
			name: "valid freeBusyStatus busy",
			event: &Event{
				Type:           "Event",
				UID:            "test-123",
				Title:          String("Test Event"),
				Start:          testStart,
				FreeBusyStatus: String(FreeBusyBusy),
			},
			wantErr: false,
		},
		{
			name: "invalid priority too low",
			event: &Event{
				UID:      "test-123",
				Title:    String("Test Event"),
				Priority: Int(-1),
			},
			wantErr: true,
			errMsg:  "priority must be between 0 and 9",
		},
		{
			name: "invalid priority too high",
			event: &Event{
				UID:      "test-123",
				Title:    String("Test Event"),
				Priority: Int(10),
			},
			wantErr: true,
			errMsg:  "priority must be between 0 and 9",
		},
		{
			name: "valid priority 0",
			event: &Event{
				Type:     "Event",
				UID:      "test-123",
				Title:    String("Test Event"),
				Start:    testStart,
				Priority: Int(0),
			},
			wantErr: false,
		},
		{
			name: "valid priority 5",
			event: &Event{
				Type:     "Event",
				UID:      "test-123",
				Title:    String("Test Event"),
				Start:    testStart,
				Priority: Int(5),
			},
			wantErr: false,
		},
		{
			name: "valid priority 9",
			event: &Event{
				Type:     "Event",
				UID:      "test-123",
				Title:    String("Test Event"),
				Start:    testStart,
				Priority: Int(9),
			},
			wantErr: false,
		},
		{
			name: "invalid descriptionContentType",
			event: &Event{
				UID:                    "test-123",
				Title:                  String("Test Event"),
				DescriptionContentType: String("text/rtf"),
			},
			wantErr: true,
			errMsg:  "descriptionContentType must be text/plain or text/html",
		},
		{
			name: "valid descriptionContentType text/plain",
			event: &Event{
				Type:                   "Event",
				UID:                    "test-123",
				Title:                  String("Test Event"),
				Start:                  testStart,
				DescriptionContentType: String("text/plain"),
			},
			wantErr: false,
		},
		{
			name: "valid descriptionContentType text/html",
			event: &Event{
				Type:                   "Event",
				UID:                    "test-123",
				Title:                  String("Test Event"),
				Start:                  testStart,
				DescriptionContentType: String("text/html"),
			},
			wantErr: false,
		},
		{
			name: "invalid sequence negative",
			event: &Event{
				UID:      "test-123",
				Title:    String("Test Event"),
				Sequence: Int(-1),
			},
			wantErr: true,
			errMsg:  "sequence cannot be negative",
		},
		{
			name: "valid sequence 0",
			event: &Event{
				Type:     "Event",
				UID:      "test-123",
				Title:    String("Test Event"),
				Start:    testStart,
				Sequence: Int(0),
			},
			wantErr: false,
		},
		{
			name: "valid sequence positive",
			event: &Event{
				Type:     "Event",
				UID:      "test-123",
				Title:    String("Test Event"),
				Start:    testStart,
				Sequence: Int(5),
			},
			wantErr: false,
		},
		// Field length validation tests
		{
			name: "uid too long",
			event: &Event{
				Type:  "Event",
				UID:   strings.Repeat("a", 256), // > MaxUIDLength (255)
				Title: String("Test Event"),
				Start: testStart,
			},
			wantErr: true,
			errMsg:  "exceeds maximum length",
		},
		{
			name: "title too long",
			event: &Event{
				Type:  "Event",
				UID:   "test-123",
				Title: String(strings.Repeat("a", 1025)), // > MaxTitleLength (1024)
				Start: testStart,
			},
			wantErr: true,
			errMsg:  "exceeds maximum length",
		},
		{
			name: "description too long",
			event: &Event{
				Type:        "Event",
				UID:         "test-123",
				Title:       String("Test Event"),
				Start:       testStart,
				Description: String(strings.Repeat("a", 32769)), // > MaxDescriptionLength (32768)
			},
			wantErr: true,
			errMsg:  "exceeds maximum length",
		},
		// Format validation tests
		{
			name: "invalid duration format",
			event: &Event{
				Type:     "Event",
				UID:      "test-123",
				Title:    String("Test Event"),
				Start:    testStart,
				Duration: String("P1Z"), // Invalid ISO 8601 duration
			},
			wantErr: true,
			errMsg:  "invalid ISO 8601 duration format",
		},
		// Method validation tests
		{
			name: "invalid method",
			event: &Event{
				Type:   "Event",
				UID:    "test-123",
				Title:  String("Test Event"),
				Start:  testStart,
				Method: String("invalid-method"),
			},
			wantErr: true,
			errMsg:  "invalid method",
		},
		{
			name: "valid method publish",
			event: &Event{
				Type:   "Event",
				UID:    "test-123",
				Title:  String("Test Event"),
				Start:  testStart,
				Method: String("publish"),
			},
			wantErr: false,
		},
		{
			name: "valid method request",
			event: &Event{
				Type:   "Event",
				UID:    "test-123",
				Title:  String("Test Event"),
				Start:  testStart,
				Method: String("request"),
			},
			wantErr: false,
		},
		// RecurrenceOverrides validation (currently not validated, just testing it doesn't error)
		{
			name: "event with recurrenceOverrides",
			event: &Event{
				Type:  "Event",
				UID:   "test-123",
				Title: String("Test Event"),
				Start: testStart,
				RecurrenceOverrides: map[string]map[string]interface{}{
					"2024-01-01T10:00:00": {
						"title": "Override Title",
					},
				},
			},
			wantErr: false, // Currently no validation for overrides
		},
		// Alert validation tests
		{
			name: "alert with invalid action",
			event: &Event{
				Type:  "Event",
				UID:   "test-123",
				Title: String("Test Event"),
				Start: testStart,
				Alerts: map[string]*Alert{
					"alert1": {
						Type: "Alert",
						Trigger: &OffsetTrigger{
							Type:   "OffsetTrigger",
							Offset: "-PT15M",
						},
						Action: String("invalidaction"),
					},
				},
			},
			wantErr: true,
			errMsg:  "invalid action",
		},
		{
			name: "alert with invalid relativeTo",
			event: &Event{
				Type:  "Event",
				UID:   "test-123",
				Title: String("Test Event"),
				Start: testStart,
				Alerts: map[string]*Alert{
					"alert1": {
						Type: "Alert",
						Trigger: &OffsetTrigger{
							Type:       "OffsetTrigger",
							Offset:     "-PT15M",
							RelativeTo: String("middle"),
						},
					},
				},
			},
			wantErr: true,
			errMsg:  "invalid relativeTo",
		},
		// RecurrenceRule validation tests
		{
			name: "recurrence rule with invalid rscale",
			event: &Event{
				Type:  "Event",
				UID:   "test-123",
				Title: String("Test Event"),
				Start: testStart,
				RecurrenceRules: []RecurrenceRule{
					{
						Type:      "RecurrenceRule",
						Frequency: "daily",
						RScale:    String("invalid-calendar"),
					},
				},
			},
			wantErr: true,
			errMsg:  "invalid rscale",
		},
		{
			name: "recurrence rule with invalid skip",
			event: &Event{
				Type:  "Event",
				UID:   "test-123",
				Title: String("Test Event"),
				Start: testStart,
				RecurrenceRules: []RecurrenceRule{
					{
						Type:      "RecurrenceRule",
						Frequency: "daily",
						Skip:      String("invalid-skip"),
					},
				},
			},
			wantErr: true,
			errMsg:  "invalid skip",
		},
		{
			name: "recurrence rule with both count and until",
			event: &Event{
				Type:  "Event",
				UID:   "test-123",
				Title: String("Test Event"),
				Start: testStart,
				RecurrenceRules: []RecurrenceRule{
					{
						Type:      "RecurrenceRule",
						Frequency: "daily",
						Count:     Int(10),
						Until:     NewLocalDateTime(time.Now()),
					},
				},
			},
			wantErr: true,
			errMsg:  "cannot have both count and until",
		},
		{
			name: "recurrence rule with negative interval",
			event: &Event{
				Type:  "Event",
				UID:   "test-123",
				Title: String("Test Event"),
				Start: testStart,
				RecurrenceRules: []RecurrenceRule{
					{
						Type:      "RecurrenceRule",
						Frequency: "daily",
						Interval:  Int(-1),
					},
				},
			},
			wantErr: true,
			errMsg:  "must be positive",
		},
		{
			name: "recurrence rule with negative count",
			event: &Event{
				Type:  "Event",
				UID:   "test-123",
				Title: String("Test Event"),
				Start: testStart,
				RecurrenceRules: []RecurrenceRule{
					{
						Type:      "RecurrenceRule",
						Frequency: "daily",
						Count:     Int(-5),
					},
				},
			},
			wantErr: true,
			errMsg:  "must be positive",
		},
		{
			name: "recurrence rule with invalid firstDayOfWeek",
			event: &Event{
				Type:  "Event",
				UID:   "test-123",
				Title: String("Test Event"),
				Start: testStart,
				RecurrenceRules: []RecurrenceRule{
					{
						Type:           "RecurrenceRule",
						Frequency:      "weekly",
						FirstDayOfWeek: Int(7), // Must be 0-6
					},
				},
			},
			wantErr: true,
			errMsg:  "invalid firstDayOfWeek",
		},
		{
			name: "recurrence rule with invalid byDay",
			event: &Event{
				Type:  "Event",
				UID:   "test-123",
				Title: String("Test Event"),
				Start: testStart,
				RecurrenceRules: []RecurrenceRule{
					{
						Type:      "RecurrenceRule",
						Frequency: "weekly",
						ByDay: []NDay{
							{Day: "xx"}, // Invalid day code
						},
					},
				},
			},
			wantErr: true,
			errMsg:  "invalid day",
		},
		{
			name: "recurrence rule with wrong type",
			event: &Event{
				Type:  "Event",
				UID:   "test-123",
				Title: String("Test Event"),
				Start: testStart,
				RecurrenceRules: []RecurrenceRule{
					{
						Type:      "WrongType",
						Frequency: "daily",
					},
				},
			},
			wantErr: true,
			errMsg:  "must be 'RecurrenceRule'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.event.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Event.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
				t.Errorf("Event.Validate() error = %v, want error containing %v", err, tt.errMsg)
			}
		})
	}
}

func TestValidateParticipant(t *testing.T) {
	tests := []struct {
		name        string
		participant *Participant
		wantErr     bool
		errMsg      string
	}{
		{
			name:        "nil participant",
			participant: nil,
			wantErr:     false, // nil is valid (optional)
		},
		{
			name:        "empty participant",
			participant: &Participant{},
			wantErr:     false, // all fields are optional
		},
		{
			name: "invalid participation status",
			participant: &Participant{
				Name:                String("John Doe"),
				ParticipationStatus: String("maybe"),
			},
			wantErr: true,
			errMsg:  "invalid participationStatus",
		},
		{
			name: "valid participation status accepted",
			participant: &Participant{
				Name:                String("John Doe"),
				ParticipationStatus: String(ParticipationAccepted),
			},
			wantErr: false,
		},
		{
			name: "valid participation status declined",
			participant: &Participant{
				Name:                String("John Doe"),
				ParticipationStatus: String(ParticipationDeclined),
			},
			wantErr: false,
		},
		{
			name: "valid participation status tentative",
			participant: &Participant{
				Name:                String("John Doe"),
				ParticipationStatus: String(ParticipationTentative),
			},
			wantErr: false,
		},
		{
			name: "valid participation status needs-action",
			participant: &Participant{
				Name:                String("John Doe"),
				ParticipationStatus: String(ParticipationNeedsAction),
			},
			wantErr: false,
		},
		{
			name: "valid participation status delegated",
			participant: &Participant{
				Name:                String("John Doe"),
				ParticipationStatus: String(ParticipationDelegated),
			},
			wantErr: false,
		},
		// ExpectReply is a *bool, so there's no invalid case to test
		// Removing invalid test case that doesn't make sense
		{
			name: "valid expectReply true",
			participant: &Participant{
				Name:        String("John Doe"),
				ExpectReply: Bool(true),
			},
			wantErr: false,
		},
		{
			name: "valid expectReply false",
			participant: &Participant{
				Name:        String("John Doe"),
				ExpectReply: Bool(false),
			},
			wantErr: false,
		},
		{
			name: "invalid scheduleAgent",
			participant: &Participant{
				Name:          String("John Doe"),
				ScheduleAgent: String("robot"),
			},
			wantErr: true,
			errMsg:  "invalid scheduleAgent",
		},
		{
			name: "valid scheduleAgent server",
			participant: &Participant{
				Name:          String("John Doe"),
				ScheduleAgent: String(ScheduleAgentServer),
			},
			wantErr: false,
		},
		{
			name: "valid scheduleAgent client",
			participant: &Participant{
				Name:          String("John Doe"),
				ScheduleAgent: String(ScheduleAgentClient),
			},
			wantErr: false,
		},
		{
			name: "valid scheduleAgent none",
			participant: &Participant{
				Name:          String("John Doe"),
				ScheduleAgent: String(ScheduleAgentNone),
			},
			wantErr: false,
		},
		// Email validation tests
		{
			name: "invalid email format",
			participant: &Participant{
				Name:  String("John Doe"),
				Email: String("notanemail"), // Missing @
			},
			wantErr: true,
			errMsg:  "invalid email format",
		},
		{
			name: "valid email format",
			participant: &Participant{
				Name:  String("John Doe"),
				Email: String("john@example.com"),
			},
			wantErr: false,
		},
		// Roles validation tests
		{
			name: "invalid role",
			participant: &Participant{
				Name: String("John Doe"),
				Roles: map[string]bool{
					"invalid-role": true,
				},
			},
			wantErr: true,
			errMsg:  "invalid role",
		},
		{
			name: "valid roles",
			participant: &Participant{
				Name: String("John Doe"),
				Roles: map[string]bool{
					"owner":    true,
					"attendee": true,
					"chair":    true,
				},
			},
			wantErr: false,
		},
		// Kind validation tests
		{
			name: "invalid kind",
			participant: &Participant{
				Name: String("John Doe"),
				Kind: String("invalid-kind"),
			},
			wantErr: true,
			errMsg:  "invalid kind",
		},
		{
			name: "valid kind individual",
			participant: &Participant{
				Name: String("John Doe"),
				Kind: String("individual"),
			},
			wantErr: false,
		},
		{
			name: "valid kind group",
			participant: &Participant{
				Name: String("Group"),
				Kind: String("group"),
			},
			wantErr: false,
		},
		// ScheduleAgent with email validation
		{
			name: "scheduleAgent with email invalid",
			participant: &Participant{
				Name:          String("John Doe"),
				Email:         String("john@example.com"),
				ScheduleAgent: String("invalid"),
			},
			wantErr: true,
			errMsg:  "invalid scheduleAgent",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateParticipant("test-participant", tt.participant)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateParticipant() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
				t.Errorf("validateParticipant() error = %v, want error containing %v", err, tt.errMsg)
			}
		})
	}
}

func TestValidateEventWithParticipants(t *testing.T) {
	event := NewEvent("test-123", "Test Event")

	// Add valid participant
	validParticipant := NewParticipant("John Doe", "john@example.com")
	validParticipant.ParticipationStatus = String(ParticipationAccepted)
	event.AddParticipant("john@example.com", validParticipant)

	// Should validate successfully
	if err := event.Validate(); err != nil {
		t.Errorf("Event with valid participant should validate: %v", err)
	}

	// Add invalid participant
	invalidParticipant := NewParticipant("Jane Doe", "jane@example.com")
	invalidParticipant.ParticipationStatus = String("invalid-status")
	event.AddParticipant("jane@example.com", invalidParticipant)

	// Should fail validation
	if err := event.Validate(); err == nil {
		t.Error("Event with invalid participant should not validate")
	}
}

func TestValidateLocation(t *testing.T) {
	tests := []struct {
		name     string
		location *Location
		wantErr  bool
		errMsg   string
	}{
		{
			name:     "nil location",
			location: nil,
			wantErr:  false,
		},
		{
			name:     "empty location",
			location: &Location{},
			wantErr:  false, // all fields are optional
		},
		{
			name: "invalid relativeTo",
			location: &Location{
				Name:       String("Conference Room"),
				RelativeTo: String("invalid-value"),
			},
			wantErr: true,
			errMsg:  "invalid relativeTo",
		},
		{
			name: "valid relativeTo start",
			location: &Location{
				Name:       String("Conference Room"),
				RelativeTo: String(RelativeToStart),
			},
			wantErr: false,
		},
		{
			name: "valid relativeTo end",
			location: &Location{
				Name:       String("Conference Room"),
				RelativeTo: String(RelativeToEnd),
			},
			wantErr: false,
		},
		{
			name: "location with coordinates",
			location: &Location{
				Name:        String("Office"),
				Coordinates: String("geo:37.386013,-122.082932"),
			},
			wantErr: false,
		},
		{
			name: "location with timeZone",
			location: &Location{
				Name:     String("Office"),
				TimeZone: String("America/Los_Angeles"),
			},
			wantErr: false,
		},
		// Coordinates validation tests
		{
			name: "invalid coordinates format",
			location: &Location{
				Name:        String("Office"),
				Coordinates: String("invalid-coordinates"), // Not a geo: URI
			},
			wantErr: true,
			errMsg:  "coordinates must be a geo: URI",
		},
		{
			name: "valid geo URI",
			location: &Location{
				Name:        String("Office"),
				Coordinates: String("geo:37.386013,-122.082932"),
			},
			wantErr: false,
		},
		// Links validation tests
		{
			name: "location with links",
			location: &Location{
				Name: String("Office"),
				Links: map[string]*Link{
					"website": {
						Href: "https://example.com",
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateLocation("test-location", tt.location)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateLocation() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
				t.Errorf("validateLocation() error = %v, want error containing %v", err, tt.errMsg)
			}
		})
	}
}

func TestValidateEventWithLocations(t *testing.T) {
	event := NewEvent("test-123", "Test Event")

	// Add valid location
	validLocation := NewLocation("Conference Room A")
	validLocation.RelativeTo = String(RelativeToStart)
	event.AddLocation("loc1", validLocation)

	// Should validate successfully
	if err := event.Validate(); err != nil {
		t.Errorf("Event with valid location should validate: %v", err)
	}

	// Add invalid location
	invalidLocation := NewLocation("Conference Room B")
	invalidLocation.RelativeTo = String("invalid-relative")
	event.AddLocation("loc2", invalidLocation)

	// Should fail validation
	if err := event.Validate(); err == nil {
		t.Error("Event with invalid location should not validate")
	}
}

func TestValidateVirtualLocation(t *testing.T) {
	tests := []struct {
		name     string
		location *VirtualLocation
		wantErr  bool
		errMsg   string
	}{
		{
			name:     "nil virtual location",
			location: nil,
			wantErr:  false,
		},
		{
			name:     "empty virtual location",
			location: &VirtualLocation{},
			wantErr:  true,
			errMsg:   "uri is required",
		},
		{
			name: "valid virtual location with uri",
			location: &VirtualLocation{
				Type: "VirtualLocation",
				URI:  "https://zoom.us/j/123456789",
			},
			wantErr: false,
		},
		{
			name: "virtual location with name",
			location: &VirtualLocation{
				Type: "VirtualLocation",
				URI:  "https://meet.google.com/abc-defg-hij",
				Name: String("Team Meeting"),
			},
			wantErr: false,
		},
		{
			name: "virtual location with description",
			location: &VirtualLocation{
				Type:        "VirtualLocation",
				URI:         "https://teams.microsoft.com/l/meetup-join/123",
				Name:        String("Project Sync"),
				Description: String("Weekly project synchronization meeting"),
			},
			wantErr: false,
		},
		{
			name: "virtual location with features",
			location: &VirtualLocation{
				Type: "VirtualLocation",
				URI:  "https://zoom.us/j/987654321",
				Features: map[string]bool{
					"audio":        true,
					"video":        true,
					"chat":         true,
					"screen-share": true,
				},
			},
			wantErr: false,
		},
		// Invalid features test (must be true per RFC)
		{
			name: "virtual location with invalid features",
			location: &VirtualLocation{
				Type: "VirtualLocation",
				URI:  "https://zoom.us/j/987654321",
				Features: map[string]bool{
					"audio": false, // RFC requires all values to be true
				},
			},
			wantErr: true,
			errMsg:  "feature value must be true",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateVirtualLocation("test-virtual", tt.location)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateVirtualLocation() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
				t.Errorf("validateVirtualLocation() error = %v, want error containing %v", err, tt.errMsg)
			}
		})
	}
}

func TestValidateAlert(t *testing.T) {
	tests := []struct {
		name    string
		alert   *Alert
		wantErr bool
		errMsg  string
	}{
		{
			name:    "nil alert",
			alert:   nil,
			wantErr: false,
		},
		{
			name: "empty alert",
			alert: &Alert{
				Type: "Alert",
			},
			wantErr: true,
			errMsg:  "trigger is required",
		},
		{
			name: "alert with offset trigger",
			alert: &Alert{
				Type: "Alert",
				Trigger: &OffsetTrigger{
					Type:   "OffsetTrigger",
					Offset: "-PT15M",
				},
			},
			wantErr: false,
		},
		{
			name: "alert with neither offset nor when",
			alert: &Alert{
				Type: "Alert",
				Trigger: &OffsetTrigger{
					Type: "OffsetTrigger",
				},
			},
			wantErr: true,
			errMsg:  "must have either offset or when",
		},
		{
			name: "invalid relativeTo",
			alert: &Alert{
				Type: "Alert",
				Trigger: &OffsetTrigger{
					Type:       "OffsetTrigger",
					Offset:     "-PT15M",
					RelativeTo: String("invalid"),
				},
			},
			wantErr: true,
			errMsg:  "invalid relativeTo",
		},
		{
			name: "valid relativeTo start",
			alert: &Alert{
				Type: "Alert",
				Trigger: &OffsetTrigger{
					Type:       "OffsetTrigger",
					Offset:     "-PT15M",
					RelativeTo: String(RelativeToStart),
				},
			},
			wantErr: false,
		},
		{
			name: "valid relativeTo end",
			alert: &Alert{
				Type: "Alert",
				Trigger: &OffsetTrigger{
					Type:       "OffsetTrigger",
					Offset:     "-PT15M",
					RelativeTo: String(RelativeToEnd),
				},
			},
			wantErr: false,
		},
		{
			name: "alert with action display",
			alert: &Alert{
				Type: "Alert",
				Trigger: &OffsetTrigger{
					Type:   "OffsetTrigger",
					Offset: "-PT15M",
				},
				Action: String("display"),
			},
			wantErr: false,
		},
		{
			name: "alert with action email",
			alert: &Alert{
				Type: "Alert",
				Trigger: &OffsetTrigger{
					Type:   "OffsetTrigger",
					Offset: "-PT15M",
				},
				Action: String("email"),
			},
			wantErr: false,
		},
		{
			name: "alert with acknowledged time",
			alert: &Alert{
				Type: "Alert",
				Trigger: &OffsetTrigger{
					Type:   "OffsetTrigger",
					Offset: "-PT15M",
				},
				Acknowledged: TimePtr(time.Date(2025, 3, 1, 13, 45, 0, 0, time.UTC)),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateAlert("test-alert", tt.alert)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateAlert() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
				t.Errorf("validateAlert() error = %v, want error containing %v", err, tt.errMsg)
			}
		})
	}
}

func TestValidateEventWithAlerts(t *testing.T) {
	event := NewEvent("test-123", "Test Event")
	event.Start = NewLocalDateTime(time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC))

	// Add valid alert
	validAlert := &Alert{
		Type: "Alert",
		Trigger: &OffsetTrigger{
			Type:   "OffsetTrigger",
			Offset: "-PT15M",
		},
	}
	event.AddAlert("alert1", validAlert)

	// Should validate successfully
	if err := event.Validate(); err != nil {
		t.Errorf("Event with valid alert should validate: %v", err)
	}

	// Add invalid alert
	invalidAlert := &Alert{
		Type: "Alert",
		Trigger: &OffsetTrigger{
			Type:       "OffsetTrigger",
			Offset:     "invalid-format",
			RelativeTo: String("invalid-relative"),
		},
	}
	event.AddAlert("alert2", invalidAlert)

	// Should fail validation
	if err := event.Validate(); err == nil {
		t.Error("Event with invalid alert should not validate")
	}
}

func TestValidateLink(t *testing.T) {
	tests := []struct {
		name    string
		link    *Link
		wantErr bool
		errMsg  string
	}{
		{
			name:    "nil link",
			link:    nil,
			wantErr: false,
		},
		{
			name:    "empty link",
			link:    &Link{},
			wantErr: true,
			errMsg:  "href is required",
		},
		{
			name: "link with href only",
			link: &Link{
				Href: "https://example.com/event",
			},
			wantErr: false,
		},
		{
			name: "link with invalid URL",
			link: &Link{
				Href: "://invalid-url",
			},
			wantErr: true,
			errMsg:  "invalid URL format",
		},
		{
			name: "link with content type",
			link: &Link{
				Href:        "https://example.com/document.pdf",
				ContentType: String("application/pdf"),
			},
			wantErr: false,
		},
		{
			name: "link with size",
			link: &Link{
				Href: "https://example.com/image.jpg",
				Size: Int(1048576), // 1MB
			},
			wantErr: false,
		},
		{
			name: "link with rel",
			link: &Link{
				Href: "https://example.com/icon.png",
				Rel:  String("icon"),
			},
			wantErr: false,
		},
		{
			name: "link with display",
			link: &Link{
				Href:    "https://example.com/info",
				Display: String("badge"),
			},
			wantErr: false,
		},
		{
			name: "link with title",
			link: &Link{
				Href:  "https://example.com/details",
				Title: String("Event Details"),
			},
			wantErr: false,
		},
		{
			name: "complete link",
			link: &Link{
				Href:        "https://example.com/resource",
				ContentType: String("text/html"),
				Size:        Int(2048),
				Rel:         String("alternate"),
				Display:     String("fullsize"),
				Title:       String("Alternative View"),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateLink("test-link", tt.link)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateLink() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
				t.Errorf("validateLink() error = %v, want error containing %v", err, tt.errMsg)
			}
		})
	}
}

func TestValidateRecurrenceRule(t *testing.T) {
	tests := []struct {
		name    string
		rule    *RecurrenceRule
		wantErr bool
		errMsg  string
	}{
		{
			name:    "nil rule",
			rule:    nil,
			wantErr: false,
		},
		{
			name:    "empty rule",
			rule:    &RecurrenceRule{Type: "RecurrenceRule"},
			wantErr: true,
			errMsg:  "frequency is required",
		},
		{
			name: "invalid frequency",
			rule: &RecurrenceRule{
				Type:      "RecurrenceRule",
				Frequency: "sometimes",
			},
			wantErr: true,
			errMsg:  "invalid frequency",
		},
		{
			name: "valid frequency daily",
			rule: &RecurrenceRule{
				Type:      "RecurrenceRule",
				Frequency: FrequencyDaily,
			},
			wantErr: false,
		},
		{
			name: "valid frequency weekly",
			rule: &RecurrenceRule{
				Type:      "RecurrenceRule",
				Frequency: FrequencyWeekly,
			},
			wantErr: false,
		},
		{
			name: "valid frequency monthly",
			rule: &RecurrenceRule{
				Type:      "RecurrenceRule",
				Frequency: FrequencyMonthly,
			},
			wantErr: false,
		},
		{
			name: "valid frequency yearly",
			rule: &RecurrenceRule{
				Type:      "RecurrenceRule",
				Frequency: FrequencyYearly,
			},
			wantErr: false,
		},
		{
			name: "rule with count",
			rule: &RecurrenceRule{
				Type:      "RecurrenceRule",
				Frequency: FrequencyDaily,
				Count:     Int(10),
			},
			wantErr: false,
		},
		{
			name: "rule with until",
			rule: &RecurrenceRule{
				Type:      "RecurrenceRule",
				Frequency: FrequencyWeekly,
				Until:     NewLocalDateTime(time.Time{}),
			},
			wantErr: false,
		},
		{
			name: "rule with both count and until",
			rule: &RecurrenceRule{
				Type:      "RecurrenceRule",
				Frequency: FrequencyDaily,
				Count:     Int(10),
				Until:     NewLocalDateTime(time.Time{}),
			},
			wantErr: true,
			errMsg:  "cannot have both count and until",
		},
		{
			name: "rule with interval",
			rule: &RecurrenceRule{
				Type:      "RecurrenceRule",
				Frequency: FrequencyWeekly,
				Interval:  Int(2),
			},
			wantErr: false,
		},
		{
			name: "invalid rscale",
			rule: &RecurrenceRule{
				Type:      "RecurrenceRule",
				Frequency: FrequencyDaily,
				RScale:    String("invalid-calendar"),
			},
			wantErr: true,
			errMsg:  "invalid rscale",
		},
		{
			name: "valid rscale gregorian",
			rule: &RecurrenceRule{
				Type:      "RecurrenceRule",
				Frequency: FrequencyMonthly,
				RScale:    String("gregorian"),
			},
			wantErr: false,
		},
		{
			name: "valid rscale chinese",
			rule: &RecurrenceRule{
				Type:      "RecurrenceRule",
				Frequency: FrequencyYearly,
				RScale:    String("chinese"),
			},
			wantErr: false,
		},
		{
			name: "valid rscale hebrew",
			rule: &RecurrenceRule{
				Type:      "RecurrenceRule",
				Frequency: FrequencyMonthly,
				RScale:    String("hebrew"),
			},
			wantErr: false,
		},
		{
			name: "valid rscale islamic",
			rule: &RecurrenceRule{
				Type:      "RecurrenceRule",
				Frequency: FrequencyMonthly,
				RScale:    String("islamic"),
			},
			wantErr: false,
		},
		{
			name: "invalid skip",
			rule: &RecurrenceRule{
				Type:      "RecurrenceRule",
				Frequency: FrequencyMonthly,
				Skip:      String("ignore"),
			},
			wantErr: true,
			errMsg:  "invalid skip",
		},
		{
			name: "valid skip forward",
			rule: &RecurrenceRule{
				Type:      "RecurrenceRule",
				Frequency: FrequencyMonthly,
				Skip:      String(SkipForward),
			},
			wantErr: false,
		},
		{
			name: "valid skip backward",
			rule: &RecurrenceRule{
				Type:      "RecurrenceRule",
				Frequency: FrequencyMonthly,
				Skip:      String(SkipBackward),
			},
			wantErr: false,
		},
		{
			name: "valid skip omit",
			rule: &RecurrenceRule{
				Type:      "RecurrenceRule",
				Frequency: FrequencyMonthly,
				Skip:      String(SkipOmit),
			},
			wantErr: false,
		},
		{
			name: "invalid firstDayOfWeek",
			rule: &RecurrenceRule{
				Type:           "RecurrenceRule",
				Frequency:      FrequencyWeekly,
				FirstDayOfWeek: Int(7), // 0-6 are valid
			},
			wantErr: true,
			errMsg:  "invalid firstDayOfWeek",
		},
		{
			name: "valid firstDayOfWeek monday",
			rule: &RecurrenceRule{
				Type:           "RecurrenceRule",
				Frequency:      FrequencyWeekly,
				FirstDayOfWeek: Int(0), // 0 = Monday
			},
			wantErr: false,
		},
		{
			name: "rule with byDay",
			rule: &RecurrenceRule{
				Type:      "RecurrenceRule",
				Frequency: FrequencyWeekly,
				ByDay: []NDay{
					{Day: "mo"},
					{Day: "we"},
					{Day: "fr"},
				},
			},
			wantErr: false,
		},
		{
			name: "rule with byDay with occurrence",
			rule: &RecurrenceRule{
				Type:      "RecurrenceRule",
				Frequency: FrequencyMonthly,
				ByDay: []NDay{
					{Day: "mo", NthOfPeriod: Int(1)},  // First Monday
					{Day: "fr", NthOfPeriod: Int(-1)}, // Last Friday
				},
			},
			wantErr: false,
		},
		{
			name: "rule with byMonth",
			rule: &RecurrenceRule{
				Type:      "RecurrenceRule",
				Frequency: FrequencyYearly,
				ByMonth:   []string{"1", "7", "12"},
			},
			wantErr: false,
		},
		{
			name: "rule with byMonthDay",
			rule: &RecurrenceRule{
				Type:       "RecurrenceRule",
				Frequency:  FrequencyMonthly,
				ByMonthDay: []int{1, 15, -1}, // 1st, 15th, last day
			},
			wantErr: false,
		},
		{
			name: "complex rule",
			rule: &RecurrenceRule{
				Type:      "RecurrenceRule",
				Frequency: FrequencyWeekly,
				Interval:  Int(2),
				ByDay:     []NDay{{Day: "mo"}, {Day: "we"}, {Day: "fr"}},
				Count:     Int(20),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateRecurrenceRule("test-rule", tt.rule)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateRecurrenceRule() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
				t.Errorf("validateRecurrenceRule() error = %v, want error containing %v", err, tt.errMsg)
			}
		})
	}
}

func TestValidateEventWithRecurrenceRules(t *testing.T) {
	event := NewEvent("test-123", "Test Event")

	// Add valid recurrence rule
	validRule := &RecurrenceRule{
		Type:      "RecurrenceRule",
		Frequency: FrequencyWeekly,
		ByDay:     []NDay{{Day: "mo"}, {Day: "we"}, {Day: "fr"}},
	}
	event.RecurrenceRules = []RecurrenceRule{*validRule}

	// Should validate successfully
	if err := event.Validate(); err != nil {
		t.Errorf("Event with valid recurrence rule should validate: %v", err)
	}

	// Add invalid recurrence rule
	invalidRule := &RecurrenceRule{
		Type:      "RecurrenceRule",
		Frequency: "invalid-frequency",
	}
	event.RecurrenceRules = append(event.RecurrenceRules, *invalidRule)

	// Should fail validation
	if err := event.Validate(); err == nil {
		t.Error("Event with invalid recurrence rule should not validate")
	}
}
