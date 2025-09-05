// Package main demonstrates basic usage of the jscal library.
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/airtrafik/jscal"
)

func main() {
	fmt.Println("JSCalendar Library - Basic Usage Examples")
	fmt.Println("========================================")

	// Example 1: Creating a new event
	fmt.Println("\n1. Creating a new event:")
	event := createBasicEvent()
	printEvent(event)

	// Example 2: Converting to JSON
	fmt.Println("\n2. Converting to JSON:")
	jsonData, err := event.PrettyJSON()
	if err != nil {
		log.Fatalf("Error converting to JSON: %v", err)
	}
	fmt.Println(string(jsonData))

	// Example 3: Parsing from JSON
	fmt.Println("\n3. Parsing from JSON:")
	parsedEvent, err := jscal.Parse(jsonData)
	if err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}
	fmt.Printf("Parsed event: %s (UID: %s)\n", *parsedEvent.Title, parsedEvent.UID)

	// Example 4: Creating a recurring event
	fmt.Println("\n4. Creating a recurring event:")
	recurringEvent := createRecurringEvent()
	printEvent(recurringEvent)

	// Example 5: Creating an all-day event
	fmt.Println("\n5. Creating an all-day event:")
	allDayEvent := createAllDayEvent()
	printEvent(allDayEvent)

	// Example 6: Validation
	fmt.Println("\n6. Validating events:")
	events := []*jscal.Event{event, recurringEvent, allDayEvent}
	for i, e := range events {
		if err := e.Validate(); err != nil {
			fmt.Printf("Event %d validation failed: %v\n", i+1, err)
		} else {
			fmt.Printf("Event %d validation passed ✓\n", i+1)
		}
	}
}

func createBasicEvent() *jscal.Event {
	// Create a new event
	event := jscal.NewEvent("meeting-2025-03-01@example.com", "Team Weekly Meeting")

	// Set basic properties
	description := "Weekly team sync to discuss progress and plan next steps"
	event.Description = &description

	// Set timing (1 hour meeting starting at 2 PM UTC)
	startTime := time.Date(2025, 3, 1, 14, 0, 0, 0, time.UTC)
	event.Start = &startTime
	event.Duration = jscal.String("PT1H")
	event.TimeZone = jscal.String("UTC")

	// Set status and privacy
	event.Status = jscal.String("confirmed")
	event.FreeBusyStatus = jscal.String("busy")
	event.Privacy = jscal.String("public")

	// Add categories
	event.AddCategory("Work")
	event.AddCategory("Meeting")

	// Add location
	location := jscal.NewLocation("Conference Room A")
	location.Description = jscal.String("Main conference room on the 5th floor")
	event.AddLocation("main-room", location)

	// Add virtual location (video meeting)
	virtualLocation := jscal.NewVirtualLocation("Zoom Meeting", "https://zoom.us/j/123456789")
	virtualLocation.Features = []string{"audio", "video", "screen"}
	event.AddVirtualLocation("zoom", virtualLocation)

	// Add participants
	organizer := jscal.NewParticipant("Alice Johnson", "alice@example.com")
	organizer.Roles = map[string]bool{"owner": true, "chair": true, "attendee": true}
	organizer.ParticipationStatus = jscal.String("accepted")
	event.AddParticipant("alice@example.com", organizer)

	attendee1 := jscal.NewParticipant("Bob Smith", "bob@example.com")
	attendee1.ParticipationStatus = jscal.String("accepted")
	event.AddParticipant("bob@example.com", attendee1)

	attendee2 := jscal.NewParticipant("Carol Davis", "carol@example.com")
	attendee2.Roles = map[string]bool{"optional": true}
	attendee2.ParticipationStatus = jscal.String("tentative")
	event.AddParticipant("carol@example.com", attendee2)

	return event
}

func createRecurringEvent() *jscal.Event {
	event := jscal.NewEvent("standup-daily@example.com", "Daily Standup")

	description := "Daily team standup meeting"
	event.Description = &description

	// Set timing (30 minutes starting at 9 AM UTC)
	startTime := time.Date(2025, 3, 3, 9, 0, 0, 0, time.UTC)
	event.Start = &startTime
	event.Duration = jscal.String("PT30M")
	event.TimeZone = jscal.String("UTC")

	// Add recurrence rule: Daily, weekdays only, for 4 weeks
	recurrenceRule := jscal.RecurrenceRule{
		Type:      "RecurrenceRule",
		Frequency: "daily",
		ByDay: []jscal.NDay{
			{Day: "mo"}, {Day: "tu"}, {Day: "we"}, {Day: "th"}, {Day: "fr"},
		},
		Count: jscal.Int(20), // 4 weeks × 5 weekdays
	}
	event.SetRecurrence([]jscal.RecurrenceRule{recurrenceRule})

	event.AddCategory("Work")
	event.AddCategory("Standup")

	return event
}

func createAllDayEvent() *jscal.Event {
	event := jscal.NewEvent("holiday-christmas@example.com", "Christmas Day")

	description := "Public holiday - office closed"
	event.Description = &description

	// Set as all-day event
	startTime := time.Date(2025, 12, 25, 0, 0, 0, 0, time.UTC)
	event.Start = &startTime
	event.Duration = jscal.String("P1D") // 1 day duration
	event.ShowWithoutTime = jscal.Bool(true)

	event.Status = jscal.String("confirmed")
	event.FreeBusyStatus = jscal.String("free") // Holiday = free time

	event.AddCategory("Holiday")
	event.AddCategory("Personal")

	return event
}

func printEvent(event *jscal.Event) {
	fmt.Printf("Event: %s\n", *event.Title)
	fmt.Printf("  UID: %s\n", event.UID)
	if event.Description != nil {
		fmt.Printf("  Description: %s\n", *event.Description)
	}
	if event.Start != nil {
		fmt.Printf("  Start: %s\n", event.Start.Format(time.RFC3339))
		if event.Duration != nil {
			fmt.Printf("  Duration: %s\n", *event.Duration)
		}
	}
	if event.IsAllDay() {
		fmt.Printf("  All-day event: yes\n")
	}
	if event.IsRecurring() {
		fmt.Printf("  Recurring: yes (%d rules)\n", len(event.RecurrenceRules))
	}
	if len(event.Categories) > 0 {
		fmt.Printf("  Categories: ")
		first := true
		for cat := range event.Categories {
			if !first {
				fmt.Print(", ")
			}
			fmt.Print(cat)
			first = false
		}
		fmt.Println()
	}
	if len(event.Participants) > 0 {
		fmt.Printf("  Participants: %d\n", len(event.Participants))
	}
}