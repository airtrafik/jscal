# RFC 8984 Test Data

This directory contains the official examples from RFC 8984 Section 6, which demonstrates various JSCalendar features.

## Examples

### Basic Objects
- **6.1-simple-event.json** - Basic one-time event with timezone and duration
- **6.2-simple-task.json** - Simple to-do item with minimal properties  
- **6.3-simple-group.json** - Group containing both Event and Task objects

### Advanced Features
- **6.4-all-day-event.json** - Recurring yearly all-day event (April Fool's Day)
- **6.5-task-with-due-date.json** - Task with due date and estimated duration
- **6.6-event-end-timezone.json** - Event with different start/end timezones (flight)
- **6.7-floating-time-event.json** - Daily recurring event in floating time (yoga)

### Complex Features (Not Yet Fully Supported)
- **6.8-event-localization.json** - Event with multiple locations and localizations
- **6.9-recurring-overrides.json** - Recurring event with exceptions and modifications  
- **6.10-recurring-participants.json** - Recurring meeting with participants and overrides

## Modifications from RFC

The following changes were made to the RFC examples to create valid JSON:
1. Removed `"...": ""` placeholder properties
2. Added required `uid` field where missing (using descriptive IDs)
3. Added required `updated` field where missing
4. Ensured all JSON is properly formatted
5. **Fixed 6.3-simple-group.json**: Changed `"name"` to `"title"` - the RFC example incorrectly uses "name" but per RFC 8984 Section 5.3.1, Group objects use "title" not "name"

## Current Support Status

✅ **Fully Supported**:
- Simple Event (6.1)
- Simple Task (6.2)
- Simple Group (6.3)
- Task with Due Date (6.5)

⚠️ **Partially Supported** (missing some features):
- All-Day Event (6.4) - needs recurrenceRules in Event struct
- Event with End Timezone (6.6) - needs Location.rel property
- Floating-Time Event (6.7) - needs recurrenceRules in Event struct

❌ **Not Yet Supported** (requires implementation):
- Event with Localization (6.8) - needs virtualLocations, localizations
- Recurring with Overrides (6.9) - needs recurrenceOverrides support
- Recurring with Participants (6.10) - needs replyTo, sendTo, recurrenceOverrides

## Testing

These examples are used in validation tests to ensure RFC 8984 compliance.
See `testdata/validation/rfc_compliance_test.go` for the test implementation.