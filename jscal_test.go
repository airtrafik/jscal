package jscal

import (
	"testing"
)

func TestParseAllEvents(t *testing.T) {
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
			name: "single event",
			data: []byte(`[{
				"@type": "Event",
				"uid": "test-1",
				"title": "Event 1",
				"start": "2024-01-01T10:00:00"
			}]`),
			want:    1,
			wantErr: false,
		},
		{
			name: "multiple events",
			data: []byte(`[
				{
					"@type": "Event",
					"uid": "test-1",
					"title": "Event 1",
					"start": "2024-01-01T10:00:00"
				},
				{
					"@type": "Event",
					"uid": "test-2",
					"title": "Event 2",
					"start": "2024-01-02T14:30:00"
				},
				{
					"@type": "Event",
					"uid": "test-3",
					"title": "Event 3",
					"start": "2024-01-03T09:15:00"
				}
			]`),
			want:    3,
			wantErr: false,
		},
		{
			name:    "invalid JSON",
			data:    []byte(`[{"invalid json`),
			want:    0,
			wantErr: true,
		},
		{
			name:    "not an array",
			data:    []byte(`{"uid": "test-1", "title": "Event 1"}`),
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			events, err := ParseAllEvents(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseAllEvents() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && len(events) != tt.want {
				t.Errorf("ParseAll() returned %d events, want %d", len(events), tt.want)
			}
		})
	}
}

func TestHelperFunctions(t *testing.T) {
	// Test String helper
	s := String("test")
	if s == nil || *s != "test" {
		t.Error("String helper failed")
	}

	// Test Int helper
	i := Int(42)
	if i == nil || *i != 42 {
		t.Error("Int helper failed")
	}

	// Test Bool helper
	b := Bool(true)
	if b == nil || *b != true {
		t.Error("Bool helper failed")
	}
}

func TestFormatDayOfWeek(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"MONDAY", "mo"},
		{"MO", "mo"},
		{"Tuesday", "tu"},
		{"WE", "we"},
		{"thursday", "th"},
		{"FR", "fr"},
		{"saturday", "sa"},
		{"SU", "su"},
	}

	for _, test := range tests {
		result := FormatDayOfWeek(test.input)
		if result != test.expected {
			t.Errorf("FormatDayOfWeek(%s) = %s, expected %s",
				test.input, result, test.expected)
		}
	}
}