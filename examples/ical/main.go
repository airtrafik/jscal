// Package main demonstrates converting between iCalendar and JSCalendar formats.
package main

import (
	"fmt"
	"log"

	"github.com/airtrafik/jscal"
	"github.com/airtrafik/jscal/convert/ical"
)

func main() {
	fmt.Println("JSCalendar Library - iCalendar Conversion Examples")
	fmt.Println("=================================================")

	// Example 1: Convert iCalendar to JSCalendar
	fmt.Println("\n1. Converting iCalendar to JSCalendar:")
	icalData := sampleICalendarData()
	fmt.Println("Input iCalendar data:")
	fmt.Println(icalData)

	converter := ical.New()
	events, err := converter.ParseAll([]byte(icalData))
	if err != nil {
		log.Fatalf("Error converting from iCalendar: %v", err)
	}

	fmt.Printf("\nConverted to %d JSCalendar event(s):\n", len(events))
	for i, event := range events {
		fmt.Printf("\nEvent %d:\n", i+1)
		printEventSummary(event)
	}

	// Example 2: Convert JSCalendar back to iCalendar
	fmt.Println("\n2. Converting JSCalendar back to iCalendar:")
	icalOutput, err := converter.FormatAll(events)
	if err != nil {
		log.Fatalf("Error converting to iCalendar: %v", err)
	}

	fmt.Println("Output iCalendar data:")
	fmt.Println(string(icalOutput))

	// Example 3: Show JSON representation
	fmt.Println("\n3. JSCalendar JSON representation:")
	if len(events) > 0 {
		jsonData, err := events[0].PrettyJSON()
		if err != nil {
			log.Fatalf("Error converting to JSON: %v", err)
		}
		fmt.Println(string(jsonData))
	}
}

func sampleICalendarData() string {
	return `BEGIN:VCALENDAR
VERSION:2.0
PRODID:-//Example Corp//CalDAV Client//EN
METHOD:PUBLISH
BEGIN:VEVENT
UID:quarterly-review-2025@example.com
DTSTART:20250315T140000Z
DTEND:20250315T160000Z
DTSTAMP:20250301T120000Z
CREATED:20250301T100000Z
LAST-MODIFIED:20250310T090000Z
SEQUENCE:2
SUMMARY:Quarterly Business Review
DESCRIPTION:Review of Q1 performance and planning for Q2.\n\nAgenda:\n- Q1 results review\n- Q2 goals and objectives\n- Budget allocation\n- Team updates
LOCATION:Executive Conference Room
URL:https://company.com/qbr-2025-q1
CATEGORIES:Business,Review,Executive
STATUS:CONFIRMED
CLASS:CONFIDENTIAL
TRANSP:OPAQUE
ORGANIZER;CN=CEO:mailto:ceo@example.com
ATTENDEE;CN=CEO;ROLE=CHAIR;PARTSTAT=ACCEPTED:mailto:ceo@example.com
ATTENDEE;CN=VP Sales;ROLE=REQ-PARTICIPANT;PARTSTAT=ACCEPTED:mailto:vp.sales@example.com
ATTENDEE;CN=VP Marketing;ROLE=REQ-PARTICIPANT;PARTSTAT=TENTATIVE:mailto:vp.marketing@example.com
ATTENDEE;CN=CFO;ROLE=REQ-PARTICIPANT;PARTSTAT=ACCEPTED:mailto:cfo@example.com
ATTENDEE;CN=CTO;ROLE=REQ-PARTICIPANT;PARTSTAT=NEEDS-ACTION:mailto:cto@example.com
ATTENDEE;CN=Board Observer;ROLE=NON-PARTICIPANT;PARTSTAT=ACCEPTED:mailto:observer@board.com
BEGIN:VALARM
ACTION:EMAIL
DESCRIPTION:QBR Reminder
SUMMARY:Quarterly Business Review in 1 hour
ATTENDEE:mailto:ceo@example.com
TRIGGER:-PT1H
END:VALARM
BEGIN:VALARM
ACTION:DISPLAY
DESCRIPTION:QBR starts in 15 minutes
TRIGGER:-PT15M
END:VALARM
END:VEVENT
END:VCALENDAR`
}

func printEventSummary(event *jscal.Event) {
	fmt.Printf("  Title: %s\n", getStringValue(event.Title))
	fmt.Printf("  UID: %s\n", event.UID)
	if event.Start != nil {
		fmt.Printf("  Start: %s\n", event.Start.Format("2006-01-02 15:04:05 UTC"))
		if endTime, err := event.GetEndTime(); err == nil {
			fmt.Printf("  End: %s\n", endTime.Format("2006-01-02 15:04:05 UTC"))
		}
	}
	if event.Duration != nil {
		fmt.Printf("  Duration: %s\n", *event.Duration)
	}
	fmt.Printf("  Status: %s\n", getStringValue(event.Status))
	fmt.Printf("  Privacy: %s\n", getStringValue(event.Privacy))

	if event.Description != nil {
		desc := *event.Description
		if len(desc) > 100 {
			desc = desc[:97] + "..."
		}
		fmt.Printf("  Description: %s\n", desc)
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

	if len(event.Locations) > 0 {
		fmt.Printf("  Locations: ")
		first := true
		for _, loc := range event.Locations {
			if loc.Name != nil {
				if !first {
					fmt.Print(", ")
				}
				fmt.Print(*loc.Name)
				first = false
			}
		}
		fmt.Println()
	}

	if len(event.Participants) > 0 {
		fmt.Printf("  Participants (%d):\n", len(event.Participants))
		for email, participant := range event.Participants {
			name := email
			if participant.Name != nil {
				name = *participant.Name
			}

			roles := "attendee"
			if participant.Roles != nil {
				var roleList []string
				for role, present := range participant.Roles {
					if present {
						roleList = append(roleList, role)
					}
				}
				if len(roleList) > 0 {
					roles = ""
					for i, role := range roleList {
						if i > 0 {
							roles += ", "
						}
						roles += role
					}
				}
			}

			status := "unknown"
			if participant.ParticipationStatus != nil {
				status = *participant.ParticipationStatus
			}

			fmt.Printf("    - %s (%s) [%s] - %s\n", name, email, roles, status)
		}
	}

	fmt.Printf("  Sequence: %s\n", getIntValue(event.Sequence))
	if event.Created != nil {
		fmt.Printf("  Created: %s\n", event.Created.Format("2006-01-02 15:04:05 UTC"))
	}
	if event.Updated != nil {
		fmt.Printf("  Updated: %s\n", event.Updated.Format("2006-01-02 15:04:05 UTC"))
	}
}

func getStringValue(s *string) string {
	if s == nil {
		return "not set"
	}
	return *s
}

func getIntValue(i *int) string {
	if i == nil {
		return "not set"
	}
	return fmt.Sprintf("%d", *i)
}
