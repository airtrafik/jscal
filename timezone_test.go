package jscal

import (
	"encoding/json"
	"testing"
	"time"
)

func TestNewTimeZone(t *testing.T) {
	tests := []struct {
		name     string
		tzid     string
		wantName string
	}{
		{
			name:     "UTC timezone",
			tzid:     "UTC",
			wantName: "Coordinated Universal Time",
		},
		{
			name:     "America/New_York timezone",
			tzid:     "America/New_York",
			wantName: "Eastern Time",
		},
		{
			name:     "Custom timezone",
			tzid:     "Custom/Zone",
			wantName: "Custom/Zone",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tz := NewTimeZone(tt.tzid)
			if tz.TzId != tt.tzid {
				t.Errorf("NewTimeZone() TzId = %v, want %v", tz.TzId, tt.tzid)
			}
			if tz.Type == nil {
				t.Error("NewTimeZone() Type is nil")
			} else if *tz.Type != "TimeZone" {
				t.Errorf("NewTimeZone() Type = %v, want %v", *tz.Type, "TimeZone")
			}
		})
	}
}

func TestTimeZoneWithRules(t *testing.T) {
	tz := NewTimeZone("America/New_York")

	// Add standard time rule
	standardRule := TimeZoneRule{
		Names:      map[string]string{"standard": "EST"},
		OffsetFrom: "-04:00",
		OffsetTo:   "-05:00",
		RecurrenceRules: []RecurrenceRule{{
			Frequency: "yearly",
			ByMonth:   []string{"11"},
			ByDay:     []NDay{{Day: "su", NthOfPeriod: Int(1)}},
		}},
		Start: NewLocalDateTime(time.Date(2024, 11, 3, 2, 0, 0, 0, time.UTC)),
	}
	tz.Standard = append(tz.Standard, standardRule)

	// Add daylight saving time rule
	dstRule := TimeZoneRule{
		Names:      map[string]string{"daylight": "EDT"},
		OffsetFrom: "-05:00",
		OffsetTo:   "-04:00",
		RecurrenceRules: []RecurrenceRule{{
			Frequency: "yearly",
			ByMonth:   []string{"3"},
			ByDay:     []NDay{{Day: "su", NthOfPeriod: Int(2)}},
		}},
		Start: NewLocalDateTime(time.Date(2024, 3, 10, 2, 0, 0, 0, time.UTC)),
	}
	tz.Daylight = append(tz.Daylight, dstRule)

	// Verify rules were added
	if len(tz.Standard) != 1 {
		t.Errorf("Expected 1 standard rule, got %d", len(tz.Standard))
	}
	if len(tz.Daylight) != 1 {
		t.Errorf("Expected 1 daylight rule, got %d", len(tz.Daylight))
	}

	// Verify standard time rule
	if tz.Standard[0].Names == nil || tz.Standard[0].Names["standard"] != "EST" {
		t.Error("Standard time rule name mismatch")
	}
	if tz.Standard[0].OffsetTo != "-05:00" {
		t.Error("Standard time offset mismatch")
	}

	// Verify DST rule
	if tz.Daylight[0].Names == nil || tz.Daylight[0].Names["daylight"] != "EDT" {
		t.Error("DST rule name mismatch")
	}
	if tz.Daylight[0].OffsetTo != "-04:00" {
		t.Error("DST offset mismatch")
	}

	// Verify recurrence rules
	if len(tz.Standard[0].RecurrenceRules) == 0 || tz.Standard[0].RecurrenceRules[0].Frequency != "yearly" {
		t.Error("Standard time recurrence rule mismatch")
	}
	if len(tz.Daylight[0].RecurrenceRules) == 0 || tz.Daylight[0].RecurrenceRules[0].ByMonth[0] != "3" {
		t.Error("DST recurrence rule month mismatch")
	}
}

func TestTimeZoneJSON(t *testing.T) {
	tz := NewTimeZone("America/Chicago")
	tz.Aliases = []string{"US/Central", "America/Chicago"}
	tz.URL = String("https://example.com/tz/chicago")

	// Test that it can be marshaled and unmarshaled
	event := NewEvent("test-tz", "Test Event")
	event.Start = NewLocalDateTime(time.Date(2025, 6, 15, 14, 0, 0, 0, time.UTC))
	event.TimeZone = String(tz.TzId)

	// Store custom timezone definition
	if event.TimeZones == nil {
		event.TimeZones = make(map[string]*TimeZone)
	}
	event.TimeZones[tz.TzId] = tz

	// Marshal and verify
	data, err := event.JSON()
	if err != nil {
		t.Fatalf("Failed to marshal event with timezone: %v", err)
	}

	// Unmarshal and verify
	var decoded Event
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal event with timezone: %v", err)
	}

	if decoded.TimeZones == nil || decoded.TimeZones["America/Chicago"] == nil {
		t.Error("TimeZone not preserved in JSON round-trip")
	}

	if len(decoded.TimeZones["America/Chicago"].Aliases) != 2 {
		t.Error("TimeZone aliases not preserved")
	}

	if decoded.TimeZones["America/Chicago"].URL == nil ||
		*decoded.TimeZones["America/Chicago"].URL != "https://example.com/tz/chicago" {
		t.Error("TimeZone URL not preserved")
	}
}
