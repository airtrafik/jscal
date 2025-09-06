// Package jscal implements RFC 8984 JSCalendar specification for Go.
//
// JSCalendar is a modern JSON-based calendar data format that provides
// a cleaner alternative to iCalendar (RFC 5545) with better support for
// modern calendar applications.
//
// This package provides:
//   - Complete JSCalendar Event type implementation
//   - JSON marshaling/unmarshaling with validation
//   - RFC 8984 compliance validation
//
// Basic usage:
//
//	// Parse JSCalendar JSON
//	event, err := jscal.ParseEvent(jsonData)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Parse multiple events
//	events, err := jscal.ParseAll(jsonData)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Validate JSCalendar compliance
//	if err := event.Validate(); err != nil {
//		log.Printf("Validation error: %v", err)
//	}
package jscal

import (
	"encoding/json"
	"fmt"
)

// Parse parses any JSCalendar object based on @type field
func Parse(data []byte) (CalendarObject, error) {
	// First, unmarshal to a map to check the @type field
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Get the @type field
	typeField, ok := raw["@type"].(string)
	if !ok {
		return nil, fmt.Errorf("missing or invalid @type field")
	}

	// Parse based on type
	switch typeField {
	case "Event":
		return ParseEvent(data)
	case "Task":
		return ParseTask(data)
	case "Group":
		return ParseGroup(data)
	default:
		return nil, fmt.Errorf("unknown @type: %s", typeField)
	}
}

// ParseAll parses multiple JSCalendar objects of any type
func ParseAll(data []byte) ([]CalendarObject, error) {
	// First, unmarshal to array of raw JSON
	var rawArray []json.RawMessage
	if err := json.Unmarshal(data, &rawArray); err != nil {
		return nil, fmt.Errorf("failed to parse JSON array: %w", err)
	}

	// Parse each object
	objects := make([]CalendarObject, 0, len(rawArray))
	for i, raw := range rawArray {
		obj, err := Parse(raw)
		if err != nil {
			return nil, fmt.Errorf("failed to parse object at index %d: %w", i, err)
		}
		objects = append(objects, obj)
	}

	return objects, nil
}

// ParseEvent parses JSCalendar JSON data into an Event
func ParseEvent(data []byte) (*Event, error) {
	var event Event
	if err := json.Unmarshal(data, &event); err != nil {
		return nil, fmt.Errorf("failed to parse JSCalendar Event JSON: %w", err)
	}

	// Validate the parsed event
	if err := event.Validate(); err != nil {
		return nil, fmt.Errorf("parsed JSCalendar Event is invalid: %w", err)
	}

	return &event, nil
}

// ParseAllEvents parses multiple JSCalendar events from JSON array
func ParseAllEvents(data []byte) ([]*Event, error) {
	var events []*Event
	if err := json.Unmarshal(data, &events); err != nil {
		return nil, fmt.Errorf("failed to parse JSCalendar Event JSON array: %w", err)
	}

	// Validate each event
	for i, event := range events {
		if err := event.Validate(); err != nil {
			return nil, fmt.Errorf("event at index %d is invalid: %w", i, err)
		}
	}

	return events, nil
}

// ParseTask parses JSCalendar JSON data into a Task
func ParseTask(data []byte) (*Task, error) {
	var task Task
	if err := json.Unmarshal(data, &task); err != nil {
		return nil, fmt.Errorf("failed to parse JSCalendar Task JSON: %w", err)
	}

	// Validate the parsed task
	if err := task.Validate(); err != nil {
		return nil, fmt.Errorf("parsed JSCalendar Task is invalid: %w", err)
	}

	return &task, nil
}

// ParseAllTasks parses multiple JSCalendar tasks from JSON array
func ParseAllTasks(data []byte) ([]*Task, error) {
	var tasks []*Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		return nil, fmt.Errorf("failed to parse JSCalendar Task JSON array: %w", err)
	}

	// Validate each task
	for i, task := range tasks {
		if err := task.Validate(); err != nil {
			return nil, fmt.Errorf("task at index %d is invalid: %w", i, err)
		}
	}

	return tasks, nil
}

// ParseGroup parses JSCalendar JSON data into a Group
func ParseGroup(data []byte) (*Group, error) {
	var group Group
	if err := json.Unmarshal(data, &group); err != nil {
		return nil, fmt.Errorf("failed to parse JSCalendar Group JSON: %w", err)
	}

	// Validate the parsed group
	if err := group.Validate(); err != nil {
		return nil, fmt.Errorf("parsed JSCalendar Group is invalid: %w", err)
	}

	return &group, nil
}

// ParseAllGroups parses multiple JSCalendar groups from JSON array
func ParseAllGroups(data []byte) ([]*Group, error) {
	var groups []*Group
	if err := json.Unmarshal(data, &groups); err != nil {
		return nil, fmt.Errorf("failed to parse JSCalendar Group JSON array: %w", err)
	}

	// Validate each group
	for i, group := range groups {
		if err := group.Validate(); err != nil {
			return nil, fmt.Errorf("group at index %d is invalid: %w", i, err)
		}
	}

	return groups, nil
}