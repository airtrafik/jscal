package jscal

import (
	"encoding/json"
	"fmt"
	"time"
)

// Task represents a JSCalendar Task object according to RFC 8984.
// A Task represents an action item, assignment, to-do item, or work item.
// It may start and be due at certain points in time, take some estimated
// time to complete, and recur, none of which is required.
type Task struct {
	// Core metadata (Section 4.1)
	Type     string     `json:"@type"`              // Always "Task"
	UID      string     `json:"uid"`                // Unique identifier
	Created  *time.Time `json:"created,omitempty"`  // Creation timestamp
	Updated  *time.Time `json:"updated,omitempty"`  // Last modification timestamp
	Sequence *int       `json:"sequence,omitempty"` // Revision sequence number
	Method   *string    `json:"method,omitempty"`   // iTIP method
	ProdId   *string    `json:"prodId,omitempty"`   // Product identifier that created this

	// What and Where Properties (Section 4.2)
	Title                  *string                           `json:"title,omitempty"`                  // Task summary/title
	Description            *string                           `json:"description,omitempty"`            // Detailed description
	DescriptionContentType *string                           `json:"descriptionContentType,omitempty"` // MIME type of description
	ShowWithoutTime        *bool                             `json:"showWithoutTime,omitempty"`        // All-day task flag
	Locale                 *string                           `json:"locale,omitempty"`                 // Language tag (RFC 5646)
	Localizations          map[string]map[string]interface{} `json:"localizations,omitempty"`          // Localization patches
	Keywords               map[string]bool                   `json:"keywords,omitempty"`               // Keywords/tags
	Categories             map[string]bool                   `json:"categories,omitempty"`             // Categories
	Color                  *string                           `json:"color,omitempty"`                  // CSS color value

	// Task-specific timing properties (Section 5.2)
	Start             *LocalDateTime       `json:"start,omitempty"`             // When the task starts
	Due               *LocalDateTime       `json:"due,omitempty"`               // When the task is due
	EstimatedDuration *string              `json:"estimatedDuration,omitempty"` // ISO 8601 duration estimate
	TimeZone          *string              `json:"timeZone,omitempty"`          // IANA timezone identifier
	TimeZones         map[string]*TimeZone `json:"timeZones,omitempty"`         // Custom timezone definitions

	// Task-specific progress properties (Section 5.2)
	PercentComplete *int       `json:"percentComplete,omitempty"` // Completion percentage (0-100)
	Progress        *string    `json:"progress,omitempty"`        // needs-action, in-process, completed, failed, cancelled
	ProgressUpdated *time.Time `json:"progressUpdated,omitempty"` // When progress was last updated

	// Recurrence Properties (Section 4.3)
	RecurrenceId            *LocalDateTime                    `json:"recurrenceId,omitempty"`         // Recurrence instance identifier
	RecurrenceIdTimeZone    *string                           `json:"recurrenceIdTimeZone,omitempty"` // TimeZone for recurrenceId (optional per RFC Erratum 6873)
	RecurrenceRules         []RecurrenceRule                  `json:"recurrenceRules,omitempty"`      // Recurrence rules
	RecurrenceOverrides     map[string]map[string]interface{} `json:"recurrenceOverrides,omitempty"`  // Overrides as patches
	ExcludedRecurrenceRules []RecurrenceRule                  `json:"excludedRecurrenceRules,omitempty"`
	Excluded                *bool                             `json:"excluded,omitempty"` // Is this instance excluded?

	// Sharing and Scheduling Properties (Section 4.4)
	Priority       *int                    `json:"priority,omitempty"`       // Priority (0-9)
	FreeBusyStatus *string                 `json:"freeBusyStatus,omitempty"` // free, busy, tentative
	Privacy        *string                 `json:"privacy,omitempty"`        // public, private, secret
	ReplyTo        map[string]string       `json:"replyTo,omitempty"`        // Reply methods
	SentBy         *string                 `json:"sentBy,omitempty"`         // Email of sender
	Participants   map[string]*Participant `json:"participants,omitempty"`   // Task participants
	RequestStatus  *string                 `json:"requestStatus,omitempty"`  // Scheduling request status

	// Alerts Properties (Section 4.5)
	UseDefaultAlerts *bool             `json:"useDefaultAlerts,omitempty"` // Use default alert settings
	Alerts           map[string]*Alert `json:"alerts,omitempty"`           // Custom alerts

	// Links and Locations
	Links            map[string]*Link            `json:"links,omitempty"`            // External links
	Locations        map[string]*Location        `json:"locations,omitempty"`        // Physical locations
	VirtualLocations map[string]*VirtualLocation `json:"virtualLocations,omitempty"` // Virtual meeting locations

	// Relationships
	RelatedTo map[string]*Relation `json:"relatedTo,omitempty"` // Related tasks/events

	// Task-specific status (Section 5.2)
	Status *string `json:"status,omitempty"` // needs-action, in-process, completed, cancelled

	// Free-form properties for extensions
	Extensions map[string]interface{} `json:"-"` // Custom properties

	// Localization
	LocalizedStrings map[string]map[string]string `json:"localizedStrings,omitempty"`
}

// NewTask creates a new JSCalendar Task with required fields
func NewTask(uid, title string) *Task {
	now := time.Now().UTC()
	return &Task{
		Type:     "Task",
		UID:      uid,
		Title:    &title,
		Created:  &now,
		Updated:  &now,
		Sequence: Int(0),
		Progress: String(ProgressNeedsAction),
	}
}

// JSON returns the Task as JSON bytes
func (t *Task) JSON() ([]byte, error) {
	return json.Marshal(t)
}

// PrettyJSON returns the Task as indented JSON bytes
func (t *Task) PrettyJSON() ([]byte, error) {
	return json.MarshalIndent(t, "", "  ")
}

// Clone creates a deep copy of the Task
func (t *Task) Clone() *Task {
	data, _ := json.Marshal(t)
	var clone Task
	_ = json.Unmarshal(data, &clone)
	return &clone
}

// IsCompleted returns true if the task is marked as completed
func (t *Task) IsCompleted() bool {
	return t.Progress != nil && *t.Progress == ProgressCompleted
}

// IsOverdue returns true if the task is past its due date
func (t *Task) IsOverdue() bool {
	if t.Due == nil || t.IsCompleted() {
		return false
	}
	
	dueTime := t.Due.Time()
	return time.Now().After(dueTime)
}

// SetProgress updates the task's progress and completion percentage
func (t *Task) SetProgress(progress string, percentComplete int) {
	t.Progress = &progress
	t.PercentComplete = &percentComplete
	now := time.Now().UTC()
	t.ProgressUpdated = &now
	t.Touch()
}

// GetEstimatedDuration returns the task's estimated duration
func (t *Task) GetEstimatedDuration() (time.Duration, error) {
	if t.EstimatedDuration == nil {
		return 0, fmt.Errorf("no estimated duration specified")
	}

	// Parse ISO 8601 duration
	return parseISO8601Duration(*t.EstimatedDuration)
}

// GetTimeToComplete calculates remaining time based on progress
func (t *Task) GetTimeToComplete() (time.Duration, error) {
	if t.IsCompleted() {
		return 0, nil
	}

	estimated, err := t.GetEstimatedDuration()
	if err != nil {
		return 0, err
	}

	if t.PercentComplete == nil || *t.PercentComplete == 0 {
		return estimated, nil
	}

	// Calculate remaining time based on percentage
	percentRemaining := 100 - *t.PercentComplete
	remainingTime := time.Duration(float64(estimated) * float64(percentRemaining) / 100.0)
	return remainingTime, nil
}

// Touch updates the Updated timestamp and increments sequence
func (t *Task) Touch() {
	now := time.Now().UTC()
	t.Updated = &now
	if t.Sequence != nil {
		*t.Sequence++
	} else {
		t.Sequence = Int(1)
	}
}

// AddParticipant adds a participant to the task
func (t *Task) AddParticipant(id string, participant *Participant) {
	if t.Participants == nil {
		t.Participants = make(map[string]*Participant)
	}
	t.Participants[id] = participant
}

// AddLocation adds a location to the task
func (t *Task) AddLocation(id string, location *Location) {
	if t.Locations == nil {
		t.Locations = make(map[string]*Location)
	}
	t.Locations[id] = location
}

// AddVirtualLocation adds a virtual location to the task
func (t *Task) AddVirtualLocation(id string, virtualLocation *VirtualLocation) {
	if t.VirtualLocations == nil {
		t.VirtualLocations = make(map[string]*VirtualLocation)
	}
	t.VirtualLocations[id] = virtualLocation
}

// AddAlert adds an alert to the task
func (t *Task) AddAlert(id string, alert *Alert) {
	if t.Alerts == nil {
		t.Alerts = make(map[string]*Alert)
	}
	t.Alerts[id] = alert
}

// AddCategory adds a category to the task
func (t *Task) AddCategory(category string) {
	if t.Categories == nil {
		t.Categories = make(map[string]bool)
	}
	t.Categories[category] = true
}

// AddKeyword adds a keyword to the task
func (t *Task) AddKeyword(keyword string) {
	if t.Keywords == nil {
		t.Keywords = make(map[string]bool)
	}
	t.Keywords[keyword] = true
}

// AddLink adds a link to the task
func (t *Task) AddLink(id string, link *Link) {
	if t.Links == nil {
		t.Links = make(map[string]*Link)
	}
	t.Links[id] = link
}

// SetRecurrence sets the recurrence rules for the task
func (t *Task) SetRecurrence(rules []RecurrenceRule) {
	t.RecurrenceRules = rules
}

// IsRecurring returns true if the task has recurrence rules
func (t *Task) IsRecurring() bool {
	return len(t.RecurrenceRules) > 0
}

// GetDueTime returns the due time as a time.Time, or error if not set
func (t *Task) GetDueTime() (*time.Time, error) {
	if t.Due == nil {
		return nil, fmt.Errorf("no due date specified")
	}
	dueTime := t.Due.Time()
	return &dueTime, nil
}

// GetStartTime returns the start time as a time.Time, or error if not set
func (t *Task) GetStartTime() (*time.Time, error) {
	if t.Start == nil {
		return nil, fmt.Errorf("no start date specified")
	}
	startTime := t.Start.Time()
	return &startTime, nil
}

// GetUID returns the task's UID (implements CalendarObject)
func (t *Task) GetUID() string {
	return t.UID
}

// GetType returns the task's type (implements CalendarObject)
func (t *Task) GetType() string {
	return t.Type
}

// Validate validates the Task according to RFC 8984
func (t *Task) Validate() error {
	if t == nil {
		return ValidationError{
			Field:   "task",
			Message: "task is nil",
		}
	}

	var errors ValidationErrors

	// Validate @type
	if t.Type != "Task" {
		errors = append(errors, ValidationError{
			Field:   "@type",
			Value:   t.Type,
			Message: "must be 'Task'",
		})
	}

	// Validate required fields
	if t.UID == "" {
		errors = append(errors, ValidationError{
			Field:   "uid",
			Value:   t.UID,
			Message: "is required",
		})
	}

	// Validate UID length (RFC 8984: max 255 octets)
	if len(t.UID) > MaxUIDLength {
		errors = append(errors, ValidationError{
			Field:   "uid",
			Value:   t.UID,
			Message: fmt.Sprintf("exceeds maximum length of %d characters", MaxUIDLength),
		})
	}

	// Validate title length if present
	if t.Title != nil && len(*t.Title) > MaxTitleLength {
		errors = append(errors, ValidationError{
			Field:   "title",
			Value:   *t.Title,
			Message: fmt.Sprintf("exceeds maximum length of %d characters", MaxTitleLength),
		})
	}

	// Validate description length if present
	if t.Description != nil && len(*t.Description) > MaxDescriptionLength {
		errors = append(errors, ValidationError{
			Field:   "description",
			Value:   *t.Description,
			Message: fmt.Sprintf("exceeds maximum length of %d characters", MaxDescriptionLength),
		})
	}

	// Validate progress
	if t.Progress != nil {
		validProgress := map[string]bool{
			ProgressNeedsAction: true,
			ProgressInProcess:   true,
			ProgressCompleted:   true,
			ProgressFailed:      true,
			ProgressCancelled:   true,
		}
		if !validProgress[*t.Progress] {
			errors = append(errors, ValidationError{
				Field:   "progress",
				Value:   *t.Progress,
				Message: "invalid progress value",
			})
		}
	}

	// Validate percentComplete
	if t.PercentComplete != nil {
		if *t.PercentComplete < 0 || *t.PercentComplete > 100 {
			errors = append(errors, ValidationError{
				Field:   "percentComplete",
				Value:   *t.PercentComplete,
				Message: "must be between 0 and 100",
			})
		}
	}

	// Validate estimatedDuration format if present
	if t.EstimatedDuration != nil {
		if !durationPattern.MatchString(*t.EstimatedDuration) {
			errors = append(errors, ValidationError{
				Field:   "estimatedDuration",
				Value:   *t.EstimatedDuration,
				Message: "invalid ISO 8601 duration format",
			})
		}
	}

	// Validate due vs start timing
	if t.Start != nil && t.Due != nil {
		startTime := t.Start.Time()
		dueTime := t.Due.Time()
		if dueTime.Before(startTime) {
			errors = append(errors, ValidationError{
				Field:   "due",
				Value:   t.Due,
				Message: "due date cannot be before start date",
			})
		}
	}

	// Validate status
	if t.Status != nil {
		validStatus := map[string]bool{
			"needs-action": true,
			"in-process":   true,
			"completed":    true,
			"cancelled":    true,
		}
		if !validStatus[*t.Status] {
			errors = append(errors, ValidationError{
				Field:   "status",
				Value:   *t.Status,
				Message: "invalid status",
			})
		}
	}

	// Validate priority
	if t.Priority != nil {
		if *t.Priority < PriorityMin || *t.Priority > PriorityMax {
			errors = append(errors, ValidationError{
				Field:   "priority",
				Value:   *t.Priority,
				Message: fmt.Sprintf("must be between %d and %d", PriorityMin, PriorityMax),
			})
		}
	}

	// Validate freeBusyStatus
	if t.FreeBusyStatus != nil {
		validFreeBusy := map[string]bool{
			FreeBusyFree:        true,
			FreeBusyBusy:        true,
			FreeBusyTentative:   true,
			FreeBusyUnavailable: true,
		}
		if !validFreeBusy[*t.FreeBusyStatus] {
			errors = append(errors, ValidationError{
				Field:   "freeBusyStatus",
				Value:   *t.FreeBusyStatus,
				Message: "invalid freeBusyStatus",
			})
		}
	}

	// Validate privacy
	if t.Privacy != nil {
		validPrivacy := map[string]bool{
			PrivacyPublic:  true,
			PrivacyPrivate: true,
			PrivacySecret:  true,
		}
		if !validPrivacy[*t.Privacy] {
			errors = append(errors, ValidationError{
				Field:   "privacy",
				Value:   *t.Privacy,
				Message: "invalid privacy",
			})
		}
	}

	// Validate sequence
	if t.Sequence != nil && *t.Sequence < 0 {
		errors = append(errors, ValidationError{
			Field:   "sequence",
			Value:   *t.Sequence,
			Message: "cannot be negative",
		})
	}

	// Validate participants
	for id, participant := range t.Participants {
		if errs := validateParticipant(id, participant); len(errs) > 0 {
			errors = append(errors, errs...)
		}
	}

	// Validate locations
	for id, location := range t.Locations {
		if errs := validateLocation(id, location); len(errs) > 0 {
			errors = append(errors, errs...)
		}
	}

	// Validate virtual locations
	for id, virtualLocation := range t.VirtualLocations {
		if errs := validateVirtualLocation(id, virtualLocation); len(errs) > 0 {
			errors = append(errors, errs...)
		}
	}

	// Validate alerts
	for id, alert := range t.Alerts {
		if errs := validateAlert(id, alert); len(errs) > 0 {
			errors = append(errors, errs...)
		}
	}

	// Validate links
	for id, link := range t.Links {
		if errs := validateLink(id, link); len(errs) > 0 {
			errors = append(errors, errs...)
		}
	}

	// Validate recurrence rules
	for i, rule := range t.RecurrenceRules {
		if errs := validateRecurrenceRule(fmt.Sprintf("recurrenceRules[%d]", i), &rule); len(errs) > 0 {
			errors = append(errors, errs...)
		}
	}

	if len(errors) > 0 {
		return errors
	}

	return nil
}