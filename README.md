# jscal - JSCalendar (RFC 8984) for Go

[![Go Reference](https://pkg.go.dev/badge/github.com/airtrafik/jscal.svg)](https://pkg.go.dev/github.com/airtrafik/jscal)
[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/go-%3E%3D1.23-blue.svg)](go.mod)

A Go library implementing [RFC 8984 (JSCalendar)](https://datatracker.ietf.org/doc/html/rfc8984) with converters for common calendar formats.

JSCalendar is a modern JSON-based calendar data format that provides a cleaner, more structured alternative to iCalendar (RFC 5545) with better support for modern calendar applications.

## Features

- ✅ **Full RFC 8984 Compliance** - Complete implementation including all verified errata
- ✅ Bidirectional iCalendar (RFC 5545) conversion
- ✅ JSON marshaling/unmarshaling with validation
- ✅ Recurrence rules (structured, not RRULE strings)
- ✅ Participants, locations, virtual locations
- ✅ Time zones as separate fields
- ✅ Localized strings with multi-language support
- ✅ CLI tool for conversions
- ✅ Zero external dependencies (core package)

## RFC 8984 Compliance

This library is **fully compliant with RFC 8984** including all verified errata. All 10 examples from RFC 8984 Section 6 are included in our test suite and pass validation.

See [RFC_COMPLIANCE.md](RFC_COMPLIANCE.md) for detailed compliance information including:
- Complete feature implementation status
- RFC errata handling (Errata 6872, 6873, 8028)
- Test coverage of all RFC examples
- Known limitations and design decisions

## Installation

### Library

```bash
go get github.com/airtrafik/jscal
```

### CLI Tool

```bash
go install github.com/airtrafik/jscal/cmd/jscal@latest
```

## Quick Start

### Basic Usage

```go
package main

import (
    "fmt"
    "log"
    "time"

    "github.com/airtrafik/jscal"
)

func main() {
    // Create a new event
    event := jscal.NewEvent("meeting-123", "Team Meeting")
    event.Description = jscal.String("Weekly team sync")

    // Set timing
    startTime := time.Date(2025, 3, 1, 14, 0, 0, 0, time.UTC)
    event.Start = &startTime
    event.Duration = jscal.String("PT1H") // 1 hour

    // Add participants
    participant := jscal.NewParticipant("Alice", "alice@example.com")
    participant.Roles = map[string]bool{"owner": true, "chair": true}
    event.AddParticipant("alice@example.com", participant)

    // Convert to JSON
    jsonData, err := event.PrettyJSON()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(string(jsonData))

    // Validate
    if err := event.Validate(); err != nil {
        log.Printf("Validation error: %v", err)
    }
}
```

### Converting from iCalendar

The iCalendar converter is a separate module to keep dependencies isolated:

```bash
go get github.com/airtrafik/jscal/convert/ical
```

```go
package main

import (
    "fmt"
    "log"

    "github.com/airtrafik/jscal/convert/ical"
)

func main() {
    icalData := `BEGIN:VCALENDAR
VERSION:2.0
BEGIN:VEVENT
UID:event-123@example.com
SUMMARY:Team Meeting
DTSTART:20250301T140000Z
DTEND:20250301T150000Z
DESCRIPTION:Weekly team sync meeting
LOCATION:Conference Room A
ORGANIZER;CN=Alice:mailto:alice@example.com
ATTENDEE;CN=Bob;PARTSTAT=ACCEPTED:mailto:bob@example.com
END:VEVENT
END:VCALENDAR`

    // Convert iCalendar to JSCalendar
    converter := ical.New()
    events, err := converter.ParseAll([]byte(icalData))
    if err != nil {
        log.Fatal(err)
    }

    // Work with JSCalendar events
    for _, event := range events {
        fmt.Printf("Event: %s\n", *event.Title)
        fmt.Printf("Start: %s\n", event.Start.Format(time.RFC3339))
    }

    // Convert back to iCalendar
    icalOutput, err := converter.FormatAll(events)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(string(icalOutput))
}
```

## CLI Usage

The `jscal` command-line tool provides easy conversion between formats:

```bash
# Convert iCalendar to JSCalendar
jscal convert calendar.ics calendar.json

# Convert JSCalendar to iCalendar
jscal convert event.json event.ics

# Validate JSCalendar files
jscal validate events.json

# Pretty-print JSCalendar
jscal format event.json
```

## API Documentation

### Core Types

- `Event` - Main JSCalendar event type
- `Participant` - Event participants (attendees, organizers)
- `Location` - Physical locations
- `VirtualLocation` - Virtual meeting locations
- `RecurrenceRule` - Structured recurrence rules
- `Alert` - Reminders and notifications

### Key Functions

```go
// Create a new event
event := jscal.NewEvent(uid, title)

// Parse JSCalendar JSON
event, err := jscal.Parse(jsonData)

// Parse multiple events
events, err := jscal.ParseAll(jsonData)

// Validate RFC 8984 compliance
err := event.Validate()

// Convert to/from iCalendar
converter := ical.New()

// Parse single event
event, err := converter.Parse(icalData)

// Parse multiple events
events, err := converter.ParseAll(icalData)

// Format back to iCalendar
icalData, err := converter.Format(event)
icalData, err := converter.FormatAll(events)
```

## Format Support

| Format | Import | Export | Status |
|--------|--------|--------|--------|
| JSCalendar (RFC 8984) | ✅ | ✅ | Complete |
| iCalendar (RFC 5545) | ✅ | ✅ | Complete |
| Google Calendar API | 🚧 | 🚧 | Planned |
| Microsoft Graph API | 🚧 | 🚧 | Planned |

## Examples

See the [`examples/`](examples/) directory for more comprehensive examples:

- [`examples/basic/`](examples/basic/) - Basic event creation and manipulation
- [`examples/ical/`](examples/ical/) - iCalendar conversion examples

## Architecture

### Modular Converters

Each converter is a separate Go module with its own `go.mod` file to isolate dependencies:

```
jscal/                           # Core library (no external deps)
├── convert/
│   ├── ical/                   # iCalendar converter module
│   │   ├── go.mod              # Uses github.com/arran4/golang-ical
│   │   └── converter.go
│   ├── google/                 # Google Calendar converter (future)
│   │   └── go.mod              # Will have Google API deps
│   └── outlook/                # Outlook converter (future)
│       └── go.mod              # Will have Microsoft Graph deps
```

This design ensures the core `jscal` package remains dependency-free while converters can use specialized libraries.

## Design Principles

- **No external dependencies** in the core package
- **Modular converters** with isolated dependencies
- **Clean, idiomatic Go** API design
- **RFC 8984 compliance** with comprehensive validation
- **Efficient parsing** and conversion
- **Extensible architecture** for adding new converters

## Performance

Target performance metrics:
- Parse JSCalendar: < 100μs per event
- Parse iCalendar: < 500μs per event
- Convert between formats: < 1ms per event
- Memory usage: < 10KB per event

## Testing

```bash
# Run all tests (including converters)
make test

# Run tests and generate coverage reports
make cover

# Run tests with verbose output
make test-verbose

# View coverage report in browser
make cover-view

# Run linter
make lint

# Run linter with auto-fix
make lint-fix

# Run everything (clean, lint, test, cover, build)
make all

# Run benchmarks
go test -bench=. ./...
```

The Makefile automatically handles testing across all modules including the main package and all converters (like `convert/ical`). Coverage reports are generated separately for each module:
- Main module: `./coverage/main.html`
- iCal converter: `./coverage/ical.html`

## Contributing

Contributions are welcome! Please feel free to submit issues and pull requests.

### Commit Message Format

This project uses [Conventional Commits](https://www.conventionalcommits.org/) format for all commit messages. This provides a clear and consistent commit history that can be used for automated versioning and changelog generation.

#### Format
```
<type>(<scope>): <subject>

[optional body]

[optional footer(s)]
```

#### Types
- `feat`: A new feature
- `fix`: A bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, missing semicolons, etc.)
- `refactor`: Code changes that neither fix a bug nor add a feature
- `perf`: Performance improvements
- `test`: Adding or updating tests
- `chore`: Changes to build process or auxiliary tools

#### Scopes
- Main package: no scope or `core`
- Converters: `ical`, `google`, `outlook`, etc.
- CLI: `cli` or `cmd`
- Examples: `examples`

#### Examples

**Main Package:**
```bash
# Feature
git commit -m "feat: add recurring event support with RRULE parsing"

# Bug fix
git commit -m "fix: correctly handle all-day events across timezones"

# Performance
git commit -m "perf: optimize JSON marshaling for large event collections"

# Tests
git commit -m "test: add validation tests for RFC 8984 compliance"
```

**iCalendar Converter (scoped):**
```bash
# Feature
git commit -m "feat(ical): add support for VALARM components"

# Bug fix  
git commit -m "fix(ical): properly escape special characters in CATEGORIES"

# Refactor
git commit -m "refactor(ical): simplify participant role mapping logic"

# Tests
git commit -m "test(ical): add round-trip tests for recurring events"
```

**CLI Tool:**
```bash
# Feature
git commit -m "feat(cli): add batch conversion support for multiple files"

# Documentation
git commit -m "docs(cli): update usage examples for convert command"
```

**Breaking Changes:**
```bash
# Use BREAKING CHANGE footer for breaking API changes
git commit -m "feat(ical)!: change Parse to return single event

BREAKING CHANGE: Parse() now returns a single event instead of a slice.
Use ParseAll() for multiple events."
```

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Acknowledgments

This library implements [RFC 8984](https://datatracker.ietf.org/doc/html/rfc8984) as defined by the IETF CALEXT Working Group.

## Roadmap

- [x] Core JSCalendar types
- [x] JSON marshaling/unmarshaling
- [x] RFC 8984 validation
- [x] iCalendar converter
- [x] CLI tool
- [ ] Google Calendar API converter
- [ ] Microsoft Graph API converter
- [ ] CalDAV support
