package jscal

import (
	"encoding/json"
	"fmt"
	"time"
)

// Group represents a JSCalendar Group object according to RFC 8984.
// A Group is a collection of Event and/or Task objects.
// Typically, objects are grouped by topic (e.g., by keywords) or calendar membership.
type Group struct {
	// Core metadata (Section 4.1)
	Type     string     `json:"@type"`              // Always "Group"
	UID      string     `json:"uid"`                // Unique identifier
	Created  *time.Time `json:"created,omitempty"`  // Creation timestamp
	Updated  *time.Time `json:"updated,omitempty"`  // Last modification timestamp
	Sequence *int       `json:"sequence,omitempty"` // Revision sequence number
	Method   *string    `json:"method,omitempty"`   // iTIP method
	ProdId   *string    `json:"prodId,omitempty"`   // Product identifier that created this

	// Group properties
	Title       *string                 `json:"title,omitempty"`       // Group title
	Description *string                 `json:"description,omitempty"` // Group description
	Locale      *string                 `json:"locale,omitempty"`      // Language tag (RFC 5646)
	Keywords    map[string]bool         `json:"keywords,omitempty"`    // Keywords/tags
	Categories  map[string]bool         `json:"categories,omitempty"`  // Categories
	Color       *string                 `json:"color,omitempty"`       // CSS color value
	Links       map[string]*Link        `json:"links,omitempty"`       // External links

	// Group-specific (Section 5.3)
	Entries []CalendarObject `json:"entries"`          // Array of Event/Task objects
	Source  *string          `json:"source,omitempty"` // Source of the group

	// Free-form properties for extensions
	Extensions map[string]interface{} `json:"-"` // Custom properties
}

// CalendarObject interface for objects that can be in a Group
type CalendarObject interface {
	GetUID() string
	GetType() string
	Validate() error
}

// NewGroup creates a new JSCalendar Group with required fields
func NewGroup(uid, title string) *Group {
	now := time.Now().UTC()
	return &Group{
		Type:     "Group",
		UID:      uid,
		Title:    &title,
		Created:  &now,
		Updated:  &now,
		Sequence: Int(0),
		Entries:  []CalendarObject{},
	}
}

// JSON returns the Group as JSON bytes
func (g *Group) JSON() ([]byte, error) {
	return json.Marshal(g)
}

// PrettyJSON returns the Group as indented JSON bytes
func (g *Group) PrettyJSON() ([]byte, error) {
	return json.MarshalIndent(g, "", "  ")
}

// Clone creates a deep copy of the Group
func (g *Group) Clone() *Group {
	data, _ := json.Marshal(g)
	var clone Group
	_ = json.Unmarshal(data, &clone)
	return &clone
}

// AddEntry adds an Event or Task to the group
func (g *Group) AddEntry(entry CalendarObject) error {
	if entry == nil {
		return fmt.Errorf("cannot add nil entry to group")
	}
	
	// Validate the entry type
	entryType := entry.GetType()
	if entryType != "Event" && entryType != "Task" {
		return fmt.Errorf("invalid entry type '%s': must be Event or Task", entryType)
	}
	
	// Check for duplicate UIDs
	uid := entry.GetUID()
	for _, existing := range g.Entries {
		if existing.GetUID() == uid {
			return fmt.Errorf("entry with UID '%s' already exists in group", uid)
		}
	}
	
	g.Entries = append(g.Entries, entry)
	g.Touch()
	return nil
}

// RemoveEntry removes an entry from the group by UID
func (g *Group) RemoveEntry(uid string) error {
	for i, entry := range g.Entries {
		if entry.GetUID() == uid {
			// Remove the entry
			g.Entries = append(g.Entries[:i], g.Entries[i+1:]...)
			g.Touch()
			return nil
		}
	}
	return fmt.Errorf("entry with UID '%s' not found in group", uid)
}

// GetEntry retrieves an entry by UID
func (g *Group) GetEntry(uid string) CalendarObject {
	for _, entry := range g.Entries {
		if entry.GetUID() == uid {
			return entry
		}
	}
	return nil
}

// GetEvents returns all Event entries in the group
func (g *Group) GetEvents() []*Event {
	var events []*Event
	for _, entry := range g.Entries {
		if event, ok := entry.(*Event); ok {
			events = append(events, event)
		}
	}
	return events
}

// GetTasks returns all Task entries in the group
func (g *Group) GetTasks() []*Task {
	var tasks []*Task
	for _, entry := range g.Entries {
		if task, ok := entry.(*Task); ok {
			tasks = append(tasks, task)
		}
	}
	return tasks
}

// CountEntries returns the total number of entries
func (g *Group) CountEntries() int {
	return len(g.Entries)
}

// CountEvents returns the number of Event entries
func (g *Group) CountEvents() int {
	count := 0
	for _, entry := range g.Entries {
		if entry.GetType() == "Event" {
			count++
		}
	}
	return count
}

// CountTasks returns the number of Task entries
func (g *Group) CountTasks() int {
	count := 0
	for _, entry := range g.Entries {
		if entry.GetType() == "Task" {
			count++
		}
	}
	return count
}

// Touch updates the Updated timestamp and increments sequence
func (g *Group) Touch() {
	now := time.Now().UTC()
	g.Updated = &now
	if g.Sequence != nil {
		*g.Sequence++
	} else {
		g.Sequence = Int(1)
	}
}

// AddKeyword adds a keyword to the group
func (g *Group) AddKeyword(keyword string) {
	if g.Keywords == nil {
		g.Keywords = make(map[string]bool)
	}
	g.Keywords[keyword] = true
}

// AddCategory adds a category to the group
func (g *Group) AddCategory(category string) {
	if g.Categories == nil {
		g.Categories = make(map[string]bool)
	}
	g.Categories[category] = true
}

// AddLink adds a link to the group
func (g *Group) AddLink(id string, link *Link) {
	if g.Links == nil {
		g.Links = make(map[string]*Link)
	}
	g.Links[id] = link
}

// GetUID returns the group's UID (implements CalendarObject)
func (g *Group) GetUID() string {
	return g.UID
}

// GetType returns the group's type (implements CalendarObject)
func (g *Group) GetType() string {
	return g.Type
}

// UnmarshalJSON implements custom JSON unmarshaling for Group to handle polymorphic entries
func (g *Group) UnmarshalJSON(data []byte) error {
	// Create an alias to avoid infinite recursion
	type Alias Group
	
	// Create auxiliary struct with json.RawMessage for entries
	aux := &struct {
		Entries []json.RawMessage `json:"entries,omitempty"`
		*Alias
	}{
		Alias: (*Alias)(g),
	}
	
	// First unmarshal everything except entries
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	
	// If no entries, we're done
	if len(aux.Entries) == 0 {
		g.Entries = []CalendarObject{}
		return nil
	}
	
	// Process each entry based on its @type
	g.Entries = make([]CalendarObject, 0, len(aux.Entries))
	for i, rawEntry := range aux.Entries {
		// Check the @type field to determine the concrete type
		var typeCheck struct {
			Type string `json:"@type"`
		}
		if err := json.Unmarshal(rawEntry, &typeCheck); err != nil {
			return fmt.Errorf("failed to determine type of entry %d: %w", i, err)
		}
		
		// Unmarshal into the appropriate type
		var entry CalendarObject
		switch typeCheck.Type {
		case "Event":
			var e Event
			if err := json.Unmarshal(rawEntry, &e); err != nil {
				return fmt.Errorf("failed to unmarshal Event at index %d: %w", i, err)
			}
			entry = &e
		case "Task":
			var t Task
			if err := json.Unmarshal(rawEntry, &t); err != nil {
				return fmt.Errorf("failed to unmarshal Task at index %d: %w", i, err)
			}
			entry = &t
		case "Group":
			var subGroup Group
			if err := json.Unmarshal(rawEntry, &subGroup); err != nil {
				return fmt.Errorf("failed to unmarshal nested Group at index %d: %w", i, err)
			}
			entry = &subGroup
		default:
			return fmt.Errorf("unknown entry type at index %d: %s", i, typeCheck.Type)
		}
		
		g.Entries = append(g.Entries, entry)
	}
	
	return nil
}

// Validate validates the Group according to RFC 8984
func (g *Group) Validate() error {
	if g == nil {
		return ValidationError{
			Field:   "group",
			Message: "group is nil",
		}
	}

	var errors ValidationErrors

	// Validate @type
	if g.Type != "Group" {
		errors = append(errors, ValidationError{
			Field:   "@type",
			Value:   g.Type,
			Message: "must be 'Group'",
		})
	}

	// Validate required fields
	if g.UID == "" {
		errors = append(errors, ValidationError{
			Field:   "uid",
			Value:   g.UID,
			Message: "is required",
		})
	}

	// Validate UID length (RFC 8984: max 255 octets)
	if len(g.UID) > MaxUIDLength {
		errors = append(errors, ValidationError{
			Field:   "uid",
			Value:   g.UID,
			Message: fmt.Sprintf("exceeds maximum length of %d characters", MaxUIDLength),
		})
	}

	// Validate title length if present
	if g.Title != nil && len(*g.Title) > MaxTitleLength {
		errors = append(errors, ValidationError{
			Field:   "title",
			Value:   *g.Title,
			Message: fmt.Sprintf("exceeds maximum length of %d characters", MaxTitleLength),
		})
	}

	// Validate description length if present
	if g.Description != nil && len(*g.Description) > MaxDescriptionLength {
		errors = append(errors, ValidationError{
			Field:   "description",
			Value:   *g.Description,
			Message: fmt.Sprintf("exceeds maximum length of %d characters", MaxDescriptionLength),
		})
	}

	// Validate sequence
	if g.Sequence != nil && *g.Sequence < 0 {
		errors = append(errors, ValidationError{
			Field:   "sequence",
			Value:   *g.Sequence,
			Message: "cannot be negative",
		})
	}

	// Validate color format if present
	if g.Color != nil {
		if !colorPattern.MatchString(*g.Color) {
			errors = append(errors, ValidationError{
				Field:   "color",
				Value:   *g.Color,
				Message: "invalid color format",
			})
		}
	}

	// Validate links
	for id, link := range g.Links {
		if errs := validateLink(id, link); len(errs) > 0 {
			errors = append(errors, errs...)
		}
	}

	// Validate entries
	if g.Entries == nil {
		g.Entries = []CalendarObject{}
	}

	// Check for duplicate UIDs in entries
	uidMap := make(map[string]bool)
	for i, entry := range g.Entries {
		if entry == nil {
			errors = append(errors, ValidationError{
				Field:   fmt.Sprintf("entries[%d]", i),
				Value:   nil,
				Message: "entry cannot be nil",
			})
			continue
		}

		uid := entry.GetUID()
		if uidMap[uid] {
			errors = append(errors, ValidationError{
				Field:   fmt.Sprintf("entries[%d]", i),
				Value:   uid,
				Message: fmt.Sprintf("duplicate UID '%s' in group entries", uid),
			})
		}
		uidMap[uid] = true

		// Validate entry type
		entryType := entry.GetType()
		if entryType != "Event" && entryType != "Task" {
			errors = append(errors, ValidationError{
				Field:   fmt.Sprintf("entries[%d].@type", i),
				Value:   entryType,
				Message: "must be 'Event' or 'Task'",
			})
		}

		// Validate the entry itself
		if err := entry.Validate(); err != nil {
			if valErrors, ok := err.(ValidationErrors); ok {
				for _, valErr := range valErrors {
					// Prefix the field with entries[i]
					valErr.Field = fmt.Sprintf("entries[%d].%s", i, valErr.Field)
					errors = append(errors, valErr)
				}
			} else {
				errors = append(errors, ValidationError{
					Field:   fmt.Sprintf("entries[%d]", i),
					Value:   entry,
					Message: err.Error(),
				})
			}
		}
	}

	// Check for circular references (a group cannot contain itself)
	for i, entry := range g.Entries {
		if entry.GetUID() == g.UID {
			errors = append(errors, ValidationError{
				Field:   fmt.Sprintf("entries[%d]", i),
				Value:   entry.GetUID(),
				Message: "group cannot contain itself",
			})
		}
	}

	if len(errors) > 0 {
		return errors
	}

	return nil
}