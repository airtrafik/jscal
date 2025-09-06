package jscal

import (
	"testing"
)

func TestNewLocation(t *testing.T) {
	tests := []struct {
		name     string
		locName  string
		wantName string
	}{
		{
			name:     "simple location",
			locName:  "Conference Room A",
			wantName: "Conference Room A",
		},
		{
			name:     "location with special chars",
			locName:  "Room #123 (Building B)",
			wantName: "Room #123 (Building B)",
		},
		{
			name:     "empty name",
			locName:  "",
			wantName: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loc := NewLocation(tt.locName)
			if loc == nil {
				t.Fatal("NewLocation returned nil")
			}
			if loc.Name == nil {
				t.Error("Location.Name should not be nil")
			} else if *loc.Name != tt.wantName {
				t.Errorf("Location.Name = %v, want %v", *loc.Name, tt.wantName)
			}
		})
	}
}

func TestLocationWithProperties(t *testing.T) {
	loc := NewLocation("Office Building")

	// Add description
	loc.Description = String("Main office building, 3rd floor")
	if loc.Description == nil || *loc.Description != "Main office building, 3rd floor" {
		t.Error("Description not set correctly")
	}

	// Add coordinates
	loc.Coordinates = String("geo:37.386013,-122.082932")
	if loc.Coordinates == nil || *loc.Coordinates != "geo:37.386013,-122.082932" {
		t.Error("Coordinates not set correctly")
	}

	// Add time zone
	loc.TimeZone = String("America/Los_Angeles")
	if loc.TimeZone == nil || *loc.TimeZone != "America/Los_Angeles" {
		t.Error("TimeZone not set correctly")
	}

	// Add relativeTo
	loc.RelativeTo = String(RelativeToStart)
	if loc.RelativeTo == nil || *loc.RelativeTo != RelativeToStart {
		t.Error("RelativeTo not set correctly")
	}

	// Add Links
	loc.Links = map[string]*Link{
		"map": NewLink("https://maps.example.com/location/123"),
	}
	if loc.Links == nil || loc.Links["map"] == nil {
		t.Error("Links not set correctly")
	}
}

func TestNewVirtualLocation(t *testing.T) {
	tests := []struct {
		name    string
		uri     string
		wantURI string
	}{
		{
			name:    "Zoom meeting",
			uri:     "https://zoom.us/j/123456789",
			wantURI: "https://zoom.us/j/123456789",
		},
		{
			name:    "Google Meet",
			uri:     "https://meet.google.com/abc-defg-hij",
			wantURI: "https://meet.google.com/abc-defg-hij",
		},
		{
			name:    "Teams meeting",
			uri:     "https://teams.microsoft.com/l/meetup-join/123",
			wantURI: "https://teams.microsoft.com/l/meetup-join/123",
		},
		{
			name:    "Phone number",
			uri:     "tel:+1-555-123-4567",
			wantURI: "tel:+1-555-123-4567",
		},
		{
			name:    "SIP URI",
			uri:     "sip:conference@example.com",
			wantURI: "sip:conference@example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vLoc := NewVirtualLocation("Meeting", tt.uri)
			if vLoc == nil {
				t.Fatal("NewVirtualLocation returned nil")
			}
			if vLoc.URI != tt.wantURI {
				t.Errorf("VirtualLocation.URI = %v, want %v", vLoc.URI, tt.wantURI)
			}
		})
	}
}

func TestVirtualLocationWithProperties(t *testing.T) {
	vLoc := NewVirtualLocation("Weekly Team Standup", "https://zoom.us/j/987654321")

	// Check name was set correctly
	if vLoc.Name == nil || *vLoc.Name != "Weekly Team Standup" {
		t.Error("Name not set correctly")
	}

	// Add description
	vLoc.Description = String("Weekly team synchronization meeting")
	if vLoc.Description == nil || *vLoc.Description != "Weekly team synchronization meeting" {
		t.Error("Description not set correctly")
	}

	// Add features
	vLoc.Features = map[string]bool{
		"audio":        true,
		"video":        true,
		"chat":         true,
		"screen-share": true,
	}

	if len(vLoc.Features) != 4 {
		t.Error("Features not set correctly")
	}
	if !vLoc.Features["audio"] || !vLoc.Features["video"] {
		t.Error("Audio and video features should be enabled")
	}
}

func TestNewLink(t *testing.T) {
	tests := []struct {
		name     string
		href     string
		wantHref string
	}{
		{
			name:     "HTTP URL",
			href:     "http://example.com/event",
			wantHref: "http://example.com/event",
		},
		{
			name:     "HTTPS URL",
			href:     "https://example.com/event/details",
			wantHref: "https://example.com/event/details",
		},
		{
			name:     "FTP URL",
			href:     "ftp://files.example.com/agenda.pdf",
			wantHref: "ftp://files.example.com/agenda.pdf",
		},
		{
			name:     "Data URI",
			href:     "data:text/plain;base64,SGVsbG8gV29ybGQ=",
			wantHref: "data:text/plain;base64,SGVsbG8gV29ybGQ=",
		},
		{
			name:     "Relative URL",
			href:     "/event/123",
			wantHref: "/event/123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			link := NewLink(tt.href)
			if link == nil {
				t.Fatal("NewLink returned nil")
			}
			if link.Href != tt.wantHref {
				t.Errorf("Link.Href = %v, want %v", link.Href, tt.wantHref)
			}
		})
	}
}

func TestLinkWithProperties(t *testing.T) {
	link := NewLink("https://example.com/document.pdf")

	// Add content type
	link.ContentType = String("application/pdf")
	if link.ContentType == nil || *link.ContentType != "application/pdf" {
		t.Error("ContentType not set correctly")
	}

	// Add size
	link.Size = Int(2048576) // 2MB
	if link.Size == nil || *link.Size != 2048576 {
		t.Error("Size not set correctly")
	}

	// Add rel
	link.Rel = String("enclosure")
	if link.Rel == nil || *link.Rel != "enclosure" {
		t.Error("Rel not set correctly")
	}

	// Add display
	link.Display = String("badge")
	if link.Display == nil || *link.Display != "badge" {
		t.Error("Display not set correctly")
	}

	// Add title
	link.Title = String("Meeting Agenda (PDF)")
	if link.Title == nil || *link.Title != "Meeting Agenda (PDF)" {
		t.Error("Title not set correctly")
	}
}

func TestParseNDay(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantDay string
		wantOcc *int
		wantErr bool
	}{
		{
			name:    "simple monday",
			input:   "mo",
			wantDay: "mo",
			wantOcc: nil,
			wantErr: false,
		},
		{
			name:    "simple friday",
			input:   "fr",
			wantDay: "fr",
			wantOcc: nil,
			wantErr: false,
		},
		{
			name:    "first monday",
			input:   "1mo",
			wantDay: "mo",
			wantOcc: Int(1),
			wantErr: false,
		},
		{
			name:    "second tuesday",
			input:   "2tu",
			wantDay: "tu",
			wantOcc: Int(2),
			wantErr: false,
		},
		{
			name:    "third wednesday",
			input:   "3we",
			wantDay: "we",
			wantOcc: Int(3),
			wantErr: false,
		},
		{
			name:    "last friday",
			input:   "-1fr",
			wantDay: "fr",
			wantOcc: Int(-1),
			wantErr: false,
		},
		{
			name:    "second to last thursday",
			input:   "-2th",
			wantDay: "th",
			wantOcc: Int(-2),
			wantErr: false,
		},
		{
			name:    "uppercase day",
			input:   "MO",
			wantDay: "mo",
			wantOcc: nil,
			wantErr: false,
		},
		{
			name:    "uppercase with occurrence",
			input:   "1MO",
			wantDay: "mo",
			wantOcc: Int(1),
			wantErr: false,
		},
		{
			name:    "mixed case",
			input:   "2Tu",
			wantDay: "tu",
			wantOcc: Int(2),
			wantErr: false,
		},
		{
			name:    "invalid day",
			input:   "xx",
			wantDay: "",
			wantOcc: nil,
			wantErr: true,
		},
		{
			name:    "invalid format",
			input:   "monday",
			wantDay: "",
			wantOcc: nil,
			wantErr: true,
		},
		{
			name:    "empty string",
			input:   "",
			wantDay: "",
			wantOcc: nil,
			wantErr: true,
		},
		{
			name:    "number only",
			input:   "1",
			wantDay: "",
			wantOcc: nil,
			wantErr: true,
		},
		{
			name:    "invalid occurrence format",
			input:   "1.5mo",
			wantDay: "",
			wantOcc: nil,
			wantErr: true,
		},
		{
			name:    "too large occurrence",
			input:   "100mo",
			wantDay: "mo",
			wantOcc: Int(100),
			wantErr: false, // ParseNDay doesn't validate range
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nday, err := ParseNDay(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseNDay() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}
			if nday.Day != tt.wantDay {
				t.Errorf("ParseNDay().Day = %v, want %v", nday.Day, tt.wantDay)
			}
			if tt.wantOcc == nil {
				if nday.NthOfPeriod != nil {
					t.Errorf("ParseNDay().NthOfPeriod = %v, want nil", *nday.NthOfPeriod)
				}
			} else {
				if nday.NthOfPeriod == nil {
					t.Error("ParseNDay().NthOfPeriod is nil, want value")
				} else if *nday.NthOfPeriod != *tt.wantOcc {
					t.Errorf("ParseNDay().NthOfPeriod = %v, want %v", *nday.NthOfPeriod, *tt.wantOcc)
				}
			}
		})
	}
}

// Test removed as NDay doesn't have a String() method

func TestParticipantWithAllFields(t *testing.T) {
	// Test creating participant with all possible fields
	p := NewParticipant("John Doe", "john.doe@example.com")

	// Set all optional fields
	p.ParticipationStatus = String(ParticipationAccepted)
	falseValue := false
	p.ExpectReply = &falseValue
	p.Kind = String("individual")
	p.Roles = map[string]bool{
		"owner":    true,
		"chair":    true,
		"attendee": true,
	}
	p.LocationId = String("conference-room-a")
	p.Language = String("en-US")
	p.DelegatedTo = map[string]bool{"jane.doe@example.com": true}
	p.DelegatedFrom = map[string]bool{"boss@example.com": true}
	p.MemberOf = map[string]bool{"team@example.com": true, "project@example.com": true}
	p.ScheduleAgent = String(ScheduleAgentServer)
	p.ScheduleSequence = Int(1)
	p.ScheduleStatus = []string{"2.0"}
	p.InvitedBy = String("organizer@example.com")
	p.SentBy = String("assistant@example.com")

	// Verify all fields are set correctly
	if p.Name == nil || *p.Name != "John Doe" {
		t.Error("Name not set correctly")
	}
	if p.Email == nil || *p.Email != "john.doe@example.com" {
		t.Error("Email not set correctly")
	}
	if p.ParticipationStatus == nil || *p.ParticipationStatus != ParticipationAccepted {
		t.Error("ParticipationStatus not set correctly")
	}
	if p.ExpectReply == nil || *p.ExpectReply != false {
		t.Error("ExpectReply not set correctly")
	}
	if p.Kind == nil || *p.Kind != "individual" {
		t.Error("Kind not set correctly")
	}
	if len(p.Roles) != 3 || !p.Roles["owner"] || !p.Roles["chair"] || !p.Roles["attendee"] {
		t.Error("Roles not set correctly")
	}
	if p.LocationId == nil || *p.LocationId != "conference-room-a" {
		t.Error("LocationId not set correctly")
	}
	if p.Language == nil || *p.Language != "en-US" {
		t.Error("Language not set correctly")
	}
	if len(p.DelegatedTo) != 1 || !p.DelegatedTo["jane.doe@example.com"] {
		t.Error("DelegatedTo not set correctly")
	}
	if len(p.DelegatedFrom) != 1 || !p.DelegatedFrom["boss@example.com"] {
		t.Error("DelegatedFrom not set correctly")
	}
	if len(p.MemberOf) != 2 || !p.MemberOf["team@example.com"] {
		t.Error("MemberOf not set correctly")
	}
	if p.ScheduleAgent == nil || *p.ScheduleAgent != ScheduleAgentServer {
		t.Error("ScheduleAgent not set correctly")
	}
	if p.ScheduleSequence == nil || *p.ScheduleSequence != 1 {
		t.Error("ScheduleSequence not set correctly")
	}
	if len(p.ScheduleStatus) != 1 || p.ScheduleStatus[0] != "2.0" {
		t.Error("ScheduleStatus not set correctly")
	}
	if p.InvitedBy == nil || *p.InvitedBy != "organizer@example.com" {
		t.Error("InvitedBy not set correctly")
	}
	if p.SentBy == nil || *p.SentBy != "assistant@example.com" {
		t.Error("SentBy not set correctly")
	}
}
