package jscal

import (
	"encoding/json"
	"testing"
	"time"
)

func TestLocalDateTimeTypeAlias(t *testing.T) {
	// Test that LocalDateTime is a type alias of time.Time
	now := time.Now()
	ldt := NewLocalDateTime(now)

	// Should be able to convert back to time.Time
	converted := ldt.Time()
	if converted.IsZero() {
		t.Error("Converted time should not be zero")
	}

	// Direct casting should work
	casted := time.Time(*ldt)
	if casted.IsZero() {
		t.Error("Casted time should not be zero")
	}
}

func TestLocalDateTimeString(t *testing.T) {
	tests := []struct {
		name     string
		input    time.Time
		expected string
	}{
		{
			name:     "UTC time",
			input:    time.Date(2025, 3, 1, 14, 30, 45, 0, time.UTC),
			expected: "2025-03-01T14:30:45",
		},
		{
			name:     "Time with nanoseconds",
			input:    time.Date(2025, 3, 1, 14, 30, 45, 123456789, time.UTC),
			expected: "2025-03-01T14:30:45.123456789",
		},
		{
			name:     "Different timezone (should still strip TZ)",
			input:    time.Date(2025, 3, 1, 14, 30, 45, 0, time.FixedZone("EST", -5*3600)),
			expected: "2025-03-01T14:30:45",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ldt := NewLocalDateTime(tt.input)
			result := ldt.String()
			if result != tt.expected {
				t.Errorf("String() = %s, want %s", result, tt.expected)
			}
		})
	}
}

func TestLocalDateTimeMarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    time.Time
		expected string
	}{
		{
			name:     "Basic datetime",
			input:    time.Date(2025, 3, 1, 14, 30, 45, 0, time.UTC),
			expected: `"2025-03-01T14:30:45"`,
		},
		{
			name:     "With nanoseconds",
			input:    time.Date(2025, 3, 1, 14, 30, 45, 123000000, time.UTC),
			expected: `"2025-03-01T14:30:45.123"`,
		},
		{
			name:     "Different timezone",
			input:    time.Date(2025, 3, 1, 14, 30, 45, 0, time.FixedZone("PST", -8*3600)),
			expected: `"2025-03-01T14:30:45"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ldt := NewLocalDateTime(tt.input)
			data, err := json.Marshal(ldt)
			if err != nil {
				t.Fatalf("MarshalJSON() error = %v", err)
			}
			if string(data) != tt.expected {
				t.Errorf("MarshalJSON() = %s, want %s", string(data), tt.expected)
			}
		})
	}
}

func TestLocalDateTimeUnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "Basic datetime",
			input:   `"2025-03-01T14:30:45"`,
			wantErr: false,
		},
		{
			name:    "With fractional seconds",
			input:   `"2025-03-01T14:30:45.123"`,
			wantErr: false,
		},
		{
			name:    "With Z timezone (should be handled)",
			input:   `"2025-03-01T14:30:45Z"`,
			wantErr: false,
		},
		{
			name:    "With offset (should be handled)",
			input:   `"2025-03-01T14:30:45+05:00"`,
			wantErr: false,
		},
		{
			name:    "Invalid format",
			input:   `"not a datetime"`,
			wantErr: true,
		},
		{
			name:    "Invalid JSON",
			input:   `not json`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ldt LocalDateTime
			err := json.Unmarshal([]byte(tt.input), &ldt)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLocalDateTimeRoundTrip(t *testing.T) {
	// Test that marshaling and unmarshaling preserves the value
	original := time.Date(2025, 3, 1, 14, 30, 45, 123456789, time.UTC)
	ldt1 := NewLocalDateTime(original)

	// Marshal to JSON
	data, err := json.Marshal(ldt1)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	// Unmarshal back
	var ldt2 LocalDateTime
	err = json.Unmarshal(data, &ldt2)
	if err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	// Compare the values
	if !ldt1.Equal(&ldt2) {
		t.Errorf("Round trip failed: %v != %v", ldt1, ldt2)
	}
}

func TestLocalDateTimeEqual(t *testing.T) {
	t1 := time.Date(2025, 3, 1, 14, 30, 45, 123456789, time.UTC)
	t2 := time.Date(2025, 3, 1, 14, 30, 45, 123456789, time.FixedZone("EST", -5*3600))
	t3 := time.Date(2025, 3, 1, 14, 30, 46, 123456789, time.UTC)

	ldt1 := NewLocalDateTime(t1)
	ldt2 := NewLocalDateTime(t2)
	ldt3 := NewLocalDateTime(t3)

	// Same time, different timezone - should be equal
	if !ldt1.Equal(ldt2) {
		t.Error("Equal times in different timezones should be equal")
	}

	// Different times - should not be equal
	if ldt1.Equal(ldt3) {
		t.Error("Different times should not be equal")
	}

	// Nil handling
	if ldt1.Equal(nil) {
		t.Error("Non-nil should not equal nil")
	}

	var nilLDT *LocalDateTime
	if !nilLDT.Equal(nil) {
		t.Error("Nil should equal nil")
	}
}

func TestLocalDateTimeIsZero(t *testing.T) {
	// Zero value
	var ldt1 LocalDateTime
	if !ldt1.IsZero() {
		t.Error("Zero value should be zero")
	}

	// Non-zero value
	ldt2 := NewLocalDateTime(time.Now())
	if ldt2.IsZero() {
		t.Error("Non-zero value should not be zero")
	}

	// Nil pointer
	var ldt3 *LocalDateTime
	if !ldt3.IsZero() {
		t.Error("Nil pointer should be zero")
	}
}

func TestLocalDateTimeBeforeAfter(t *testing.T) {
	t1 := time.Date(2025, 3, 1, 14, 30, 45, 0, time.UTC)
	t2 := time.Date(2025, 3, 1, 14, 30, 46, 0, time.UTC)
	t3 := time.Date(2025, 3, 1, 14, 30, 45, 0, time.FixedZone("EST", -5*3600))

	ldt1 := NewLocalDateTime(t1)
	ldt2 := NewLocalDateTime(t2)
	ldt3 := NewLocalDateTime(t3)

	// Before tests
	if !ldt1.Before(ldt2) {
		t.Error("Earlier time should be before later time")
	}
	if ldt2.Before(ldt1) {
		t.Error("Later time should not be before earlier time")
	}
	if ldt1.Before(ldt3) || ldt3.Before(ldt1) {
		t.Error("Same local time should not be before itself")
	}

	// After tests
	if !ldt2.After(ldt1) {
		t.Error("Later time should be after earlier time")
	}
	if ldt1.After(ldt2) {
		t.Error("Earlier time should not be after later time")
	}
	if ldt1.After(ldt3) || ldt3.After(ldt1) {
		t.Error("Same local time should not be after itself")
	}

	// Nil handling
	if ldt1.Before(nil) || ldt1.After(nil) {
		t.Error("Non-nil should not be before or after nil")
	}

	var nilLDT *LocalDateTime
	if nilLDT.Before(ldt1) || nilLDT.After(ldt1) {
		t.Error("Nil should not be before or after non-nil")
	}
}

func TestLocalDateTimeAdd(t *testing.T) {
	t1 := time.Date(2025, 3, 1, 14, 30, 45, 0, time.UTC)
	ldt1 := NewLocalDateTime(t1)

	// Add 1 hour
	ldt2 := ldt1.Add(time.Hour)
	expected := time.Date(2025, 3, 1, 15, 30, 45, 0, time.UTC)
	if ldt2.Hour() != 15 {
		t.Errorf("Add(1 hour) failed: got hour %d, want 15", ldt2.Hour())
	}
	if !ldt2.Equal(NewLocalDateTime(expected)) {
		t.Error("Add(1 hour) result mismatch")
	}

	// Add negative duration
	ldt3 := ldt1.Add(-30 * time.Minute)
	if ldt3.Hour() != 14 || ldt3.Minute() != 0 {
		t.Errorf("Add(-30 min) failed: got %02d:%02d", ldt3.Hour(), ldt3.Minute())
	}
}

func TestLocalDateTimeSub(t *testing.T) {
	t1 := time.Date(2025, 3, 1, 14, 30, 45, 0, time.UTC)
	t2 := time.Date(2025, 3, 1, 15, 30, 45, 0, time.UTC)

	ldt1 := NewLocalDateTime(t1)
	ldt2 := NewLocalDateTime(t2)

	duration := ldt2.Sub(*ldt1)
	if duration != time.Hour {
		t.Errorf("Sub() = %v, want 1 hour", duration)
	}

	duration = ldt1.Sub(*ldt2)
	if duration != -time.Hour {
		t.Errorf("Sub() = %v, want -1 hour", duration)
	}
}

func TestLocalDateTimeFormat(t *testing.T) {
	t1 := time.Date(2025, 3, 1, 14, 30, 45, 0, time.UTC)
	ldt := NewLocalDateTime(t1)

	tests := []struct {
		layout   string
		expected string
	}{
		{time.RFC3339, "2025-03-01T14:30:45Z"},
		{"2006-01-02", "2025-03-01"},
		{"15:04:05", "14:30:45"},
		{time.Kitchen, "2:30PM"},
	}

	for _, tt := range tests {
		result := ldt.Format(tt.layout)
		if result != tt.expected {
			t.Errorf("Format(%s) = %s, want %s", tt.layout, result, tt.expected)
		}
	}
}

func TestLocalDateTimeAccessors(t *testing.T) {
	t1 := time.Date(2025, 3, 15, 14, 30, 45, 123456789, time.UTC)
	ldt := NewLocalDateTime(t1)

	if ldt.Year() != 2025 {
		t.Errorf("Year() = %d, want 2025", ldt.Year())
	}
	if ldt.Month() != time.March {
		t.Errorf("Month() = %v, want March", ldt.Month())
	}
	if ldt.Day() != 15 {
		t.Errorf("Day() = %d, want 15", ldt.Day())
	}
	if ldt.Hour() != 14 {
		t.Errorf("Hour() = %d, want 14", ldt.Hour())
	}
	if ldt.Minute() != 30 {
		t.Errorf("Minute() = %d, want 30", ldt.Minute())
	}
	if ldt.Second() != 45 {
		t.Errorf("Second() = %d, want 45", ldt.Second())
	}
	if ldt.Nanosecond() != 123456789 {
		t.Errorf("Nanosecond() = %d, want 123456789", ldt.Nanosecond())
	}
}

func TestParseLocalDateTime(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"Without timezone", "2025-03-01T14:30:45", false},
		{"With fractional seconds", "2025-03-01T14:30:45.123", false},
		{"With Z timezone", "2025-03-01T14:30:45Z", false},
		{"With offset", "2025-03-01T14:30:45+05:00", false},
		{"RFC3339", "2025-03-01T14:30:45Z", false},
		{"RFC3339Nano", "2025-03-01T14:30:45.123456789Z", false},
		{"Invalid format", "not-a-datetime", true},
		{"Date only", "2025-03-01", true},
		{"Time only", "14:30:45", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ldt, err := ParseLocalDateTime(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseLocalDateTime() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && ldt == nil {
				t.Error("ParseLocalDateTime() returned nil without error")
			}
		})
	}
}

func TestLocalDateTimeNilHandling(t *testing.T) {
	var nilLDT *LocalDateTime

	// Time() on nil
	if !nilLDT.Time().IsZero() {
		t.Error("Time() on nil should return zero time")
	}

	// IsZero() on nil
	if !nilLDT.IsZero() {
		t.Error("IsZero() on nil should return true")
	}

	// Equal() with nil
	ldt := NewLocalDateTime(time.Now())
	if ldt.Equal(nil) {
		t.Error("Non-nil should not equal nil")
	}
	if !nilLDT.Equal(nil) {
		t.Error("Nil should equal nil")
	}

	// Before/After with nil
	if nilLDT.Before(ldt) || nilLDT.After(ldt) {
		t.Error("Nil should not be before or after non-nil")
	}
	if ldt.Before(nil) || ldt.After(nil) {
		t.Error("Non-nil should not be before or after nil")
	}
}

func TestLocalDateTimeStructInJSON(t *testing.T) {
	type TestStruct struct {
		Start *LocalDateTime `json:"start,omitempty"`
		End   *LocalDateTime `json:"end,omitempty"`
		Name  string         `json:"name"`
	}

	// Test marshaling
	ts := TestStruct{
		Start: NewLocalDateTime(time.Date(2025, 3, 1, 14, 30, 0, 0, time.UTC)),
		End:   NewLocalDateTime(time.Date(2025, 3, 1, 15, 30, 0, 0, time.UTC)),
		Name:  "Test Event",
	}

	data, err := json.Marshal(ts)
	if err != nil {
		t.Fatalf("Failed to marshal struct: %v", err)
	}

	expected := `{"start":"2025-03-01T14:30:00","end":"2025-03-01T15:30:00","name":"Test Event"}`
	if string(data) != expected {
		t.Errorf("Marshal struct = %s, want %s", string(data), expected)
	}

	// Test unmarshaling
	var ts2 TestStruct
	err = json.Unmarshal(data, &ts2)
	if err != nil {
		t.Fatalf("Failed to unmarshal struct: %v", err)
	}

	if !ts.Start.Equal(ts2.Start) || !ts.End.Equal(ts2.End) || ts.Name != ts2.Name {
		t.Error("Unmarshaled struct doesn't match original")
	}
}
