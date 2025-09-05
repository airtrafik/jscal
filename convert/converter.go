package convert

import "github.com/airtrafik/jscal"

// Converter defines the interface for calendar format converters
type Converter interface {
	// Single event (common case - simple names)
	Parse(data []byte) (*jscal.Event, error)
	Format(event *jscal.Event) ([]byte, error)
	
	// Multiple events (explicit with "All")
	ParseAll(data []byte) ([]*jscal.Event, error)
	FormatAll(events []*jscal.Event) ([]byte, error)
	
	// Detection
	Detect(data []byte) bool
}