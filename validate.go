package jscal

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

// Validation constants
const (
	MaxTitleLength       = 1024
	MaxDescriptionLength = 32768
	MaxUIDLength         = 255
)

// Regular expressions for validation
var (
	// ISO 8601 duration pattern (simplified)
	durationPattern = regexp.MustCompile(`^P(?:\d+Y)?(?:\d+M)?(?:\d+D)?(?:T(?:\d+H)?(?:\d+M)?(?:\d+(?:\.\d+)?S)?)?$`)
	
	// Color pattern (CSS color values)
	colorPattern = regexp.MustCompile(`^(?:#[0-9a-fA-F]{3,8}|rgb\(|rgba\(|hsl\(|hsla\(|[a-zA-Z]+)`)
	
	// IANA timezone pattern
	timezonePattern = regexp.MustCompile(`^[A-Za-z0-9/_+-]+$`)
)

// ValidationError represents a validation error
type ValidationError struct {
	Field   string
	Value   interface{}
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("validation error in field %s: %s", e.Field, e.Message)
}

// ValidationErrors represents multiple validation errors
type ValidationErrors []ValidationError

func (e ValidationErrors) Error() string {
	if len(e) == 0 {
		return "no validation errors"
	}
	if len(e) == 1 {
		return e[0].Error()
	}
	
	var messages []string
	for _, err := range e {
		messages = append(messages, err.Error())
	}
	return fmt.Sprintf("multiple validation errors: %s", strings.Join(messages, "; "))
}

// Validate validates the Event according to RFC 8984
func (e *Event) Validate() error {
	var errors ValidationErrors
	
	// Required fields
	if e.Type != "Event" {
		errors = append(errors, ValidationError{
			Field:   "@type",
			Value:   e.Type,
			Message: "must be 'Event'",
		})
	}
	
	if e.UID == "" {
		errors = append(errors, ValidationError{
			Field:   "uid",
			Value:   e.UID,
			Message: "is required",
		})
	} else if len(e.UID) > MaxUIDLength {
		errors = append(errors, ValidationError{
			Field:   "uid",
			Value:   e.UID,
			Message: fmt.Sprintf("exceeds maximum length of %d characters", MaxUIDLength),
		})
	}
	
	// Validate title length
	if e.Title != nil && len(*e.Title) > MaxTitleLength {
		errors = append(errors, ValidationError{
			Field:   "title",
			Value:   *e.Title,
			Message: fmt.Sprintf("exceeds maximum length of %d characters", MaxTitleLength),
		})
	}
	
	// Validate description length
	if e.Description != nil && len(*e.Description) > MaxDescriptionLength {
		errors = append(errors, ValidationError{
			Field:   "description",
			Value:   *e.Description,
			Message: fmt.Sprintf("exceeds maximum length of %d characters", MaxDescriptionLength),
		})
	}
	
	// Validate duration format
	if e.Duration != nil {
		if !durationPattern.MatchString(*e.Duration) {
			errors = append(errors, ValidationError{
				Field:   "duration",
				Value:   *e.Duration,
				Message: "invalid ISO 8601 duration format",
			})
		}
	}
	
	// Validate timezone
	if e.TimeZone != nil {
		if !timezonePattern.MatchString(*e.TimeZone) {
			errors = append(errors, ValidationError{
				Field:   "timeZone",
				Value:   *e.TimeZone,
				Message: "invalid IANA timezone identifier",
			})
		}
	}
	
	// Validate color
	if e.Color != nil {
		if !colorPattern.MatchString(*e.Color) {
			errors = append(errors, ValidationError{
				Field:   "color",
				Value:   *e.Color,
				Message: "invalid CSS color value",
			})
		}
	}
	
	// Validate status
	if e.Status != nil {
		validStatuses := map[string]bool{
			"confirmed": true,
			"tentative": true,
			"cancelled": true,
		}
		if !validStatuses[*e.Status] {
			errors = append(errors, ValidationError{
				Field:   "status",
				Value:   *e.Status,
				Message: "must be one of: confirmed, tentative, cancelled",
			})
		}
	}
	
	// Validate freeBusyStatus
	if e.FreeBusyStatus != nil {
		validStatuses := map[string]bool{
			"free":      true,
			"busy":      true,
			"tentative": true,
		}
		if !validStatuses[*e.FreeBusyStatus] {
			errors = append(errors, ValidationError{
				Field:   "freeBusyStatus",
				Value:   *e.FreeBusyStatus,
				Message: "must be one of: free, busy, tentative",
			})
		}
	}
	
	// Validate privacy
	if e.Privacy != nil {
		validPrivacyLevels := map[string]bool{
			"public":  true,
			"private": true,
			"secret":  true,
		}
		if !validPrivacyLevels[*e.Privacy] {
			errors = append(errors, ValidationError{
				Field:   "privacy",
				Value:   *e.Privacy,
				Message: "must be one of: public, private, secret",
			})
		}
	}
	
	// Validate sequence
	if e.Sequence != nil && *e.Sequence < 0 {
		errors = append(errors, ValidationError{
			Field:   "sequence",
			Value:   *e.Sequence,
			Message: "must be non-negative",
		})
	}
	
	// Validate participants
	for id, participant := range e.Participants {
		if errs := validateParticipant(id, participant); len(errs) > 0 {
			errors = append(errors, errs...)
		}
	}
	
	// Validate locations
	for id, location := range e.Locations {
		if errs := validateLocation(id, location); len(errs) > 0 {
			errors = append(errors, errs...)
		}
	}
	
	// Validate virtual locations
	for id, virtualLocation := range e.VirtualLocations {
		if errs := validateVirtualLocation(id, virtualLocation); len(errs) > 0 {
			errors = append(errors, errs...)
		}
	}
	
	// Validate alerts
	for id, alert := range e.Alerts {
		if errs := validateAlert(id, alert); len(errs) > 0 {
			errors = append(errors, errs...)
		}
	}
	
	// Validate links
	for id, link := range e.Links {
		if errs := validateLink(id, link); len(errs) > 0 {
			errors = append(errors, errs...)
		}
	}
	
	// Validate recurrence rules
	for i, rule := range e.RecurrenceRules {
		if errs := validateRecurrenceRule(fmt.Sprintf("recurrenceRules[%d]", i), &rule); len(errs) > 0 {
			errors = append(errors, errs...)
		}
	}
	
	if len(errors) > 0 {
		return errors
	}
	
	return nil
}

func validateParticipant(id string, p *Participant) ValidationErrors {
	var errors ValidationErrors
	
	if p == nil {
		errors = append(errors, ValidationError{
			Field:   fmt.Sprintf("participants[%s]", id),
			Value:   nil,
			Message: "participant cannot be nil",
		})
		return errors
	}
	
	// Validate email format if present
	if p.Email != nil && *p.Email != "" {
		if !strings.Contains(*p.Email, "@") {
			errors = append(errors, ValidationError{
				Field:   fmt.Sprintf("participants[%s].email", id),
				Value:   *p.Email,
				Message: "invalid email format",
			})
		}
	}
	
	// Validate participation status
	if p.ParticipationStatus != nil {
		validStatuses := map[string]bool{
			"needs-action": true,
			"accepted":     true,
			"declined":     true,
			"tentative":    true,
			"delegated":    true,
		}
		if !validStatuses[*p.ParticipationStatus] {
			errors = append(errors, ValidationError{
				Field:   fmt.Sprintf("participants[%s].participationStatus", id),
				Value:   *p.ParticipationStatus,
				Message: "must be one of: needs-action, accepted, declined, tentative, delegated",
			})
		}
	}
	
	// Validate progress
	if p.Progress != nil {
		validProgress := map[string]bool{
			"needs-action": true,
			"in-process":   true,
			"completed":    true,
			"failed":       true,
			"cancelled":    true,
		}
		if !validProgress[*p.Progress] {
			errors = append(errors, ValidationError{
				Field:   fmt.Sprintf("participants[%s].progress", id),
				Value:   *p.Progress,
				Message: "must be one of: needs-action, in-process, completed, failed, cancelled",
			})
		}
	}
	
	// Validate percentage complete
	if p.PercentComplete != nil {
		if *p.PercentComplete < 0 || *p.PercentComplete > 100 {
			errors = append(errors, ValidationError{
				Field:   fmt.Sprintf("participants[%s].percentComplete", id),
				Value:   *p.PercentComplete,
				Message: "must be between 0 and 100",
			})
		}
	}
	
	return errors
}

func validateLocation(id string, l *Location) ValidationErrors {
	var errors ValidationErrors
	
	if l == nil {
		errors = append(errors, ValidationError{
			Field:   fmt.Sprintf("locations[%s]", id),
			Value:   nil,
			Message: "location cannot be nil",
		})
		return errors
	}
	
	// Validate coordinates format (geo: URI)
	if l.Coordinates != nil {
		if !strings.HasPrefix(*l.Coordinates, "geo:") {
			errors = append(errors, ValidationError{
				Field:   fmt.Sprintf("locations[%s].coordinates", id),
				Value:   *l.Coordinates,
				Message: "must be a geo: URI",
			})
		}
	}
	
	return errors
}

func validateVirtualLocation(id string, vl *VirtualLocation) ValidationErrors {
	var errors ValidationErrors
	
	if vl == nil {
		errors = append(errors, ValidationError{
			Field:   fmt.Sprintf("virtualLocations[%s]", id),
			Value:   nil,
			Message: "virtual location cannot be nil",
		})
		return errors
	}
	
	if vl.Type != "VirtualLocation" {
		errors = append(errors, ValidationError{
			Field:   fmt.Sprintf("virtualLocations[%s].@type", id),
			Value:   vl.Type,
			Message: "must be 'VirtualLocation'",
		})
	}
	
	if vl.URI == "" {
		errors = append(errors, ValidationError{
			Field:   fmt.Sprintf("virtualLocations[%s].uri", id),
			Value:   vl.URI,
			Message: "uri is required",
		})
	} else {
		if _, err := url.Parse(vl.URI); err != nil {
			errors = append(errors, ValidationError{
				Field:   fmt.Sprintf("virtualLocations[%s].uri", id),
				Value:   vl.URI,
				Message: "invalid URI format",
			})
		}
	}
	
	// Validate features
	validFeatures := map[string]bool{
		"audio":  true,
		"video":  true,
		"chat":   true,
		"screen": true,
		"phone":  true,
	}
	for _, feature := range vl.Features {
		if !validFeatures[feature] {
			errors = append(errors, ValidationError{
				Field:   fmt.Sprintf("virtualLocations[%s].features", id),
				Value:   feature,
				Message: "invalid feature type",
			})
		}
	}
	
	return errors
}

func validateAlert(id string, a *Alert) ValidationErrors {
	var errors ValidationErrors
	
	if a == nil {
		errors = append(errors, ValidationError{
			Field:   fmt.Sprintf("alerts[%s]", id),
			Value:   nil,
			Message: "alert cannot be nil",
		})
		return errors
	}
	
	if a.Type != "Alert" {
		errors = append(errors, ValidationError{
			Field:   fmt.Sprintf("alerts[%s].@type", id),
			Value:   a.Type,
			Message: "must be 'Alert'",
		})
	}
	
	// Validate action
	if a.Action != nil {
		validActions := map[string]bool{
			"display": true,
			"email":   true,
		}
		if !validActions[*a.Action] {
			errors = append(errors, ValidationError{
				Field:   fmt.Sprintf("alerts[%s].action", id),
				Value:   *a.Action,
				Message: "must be one of: display, email",
			})
		}
	}
	
	return errors
}

func validateLink(id string, l *Link) ValidationErrors {
	var errors ValidationErrors
	
	if l == nil {
		errors = append(errors, ValidationError{
			Field:   fmt.Sprintf("links[%s]", id),
			Value:   nil,
			Message: "link cannot be nil",
		})
		return errors
	}
	
	if l.Href == "" {
		errors = append(errors, ValidationError{
			Field:   fmt.Sprintf("links[%s].href", id),
			Value:   l.Href,
			Message: "href is required",
		})
	} else {
		if _, err := url.Parse(l.Href); err != nil {
			errors = append(errors, ValidationError{
				Field:   fmt.Sprintf("links[%s].href", id),
				Value:   l.Href,
				Message: "invalid URL format",
			})
		}
	}
	
	return errors
}

func validateRecurrenceRule(fieldPrefix string, rr *RecurrenceRule) ValidationErrors {
	var errors ValidationErrors
	
	if rr == nil {
		errors = append(errors, ValidationError{
			Field:   fieldPrefix,
			Value:   nil,
			Message: "recurrence rule cannot be nil",
		})
		return errors
	}
	
	if rr.Type != "RecurrenceRule" {
		errors = append(errors, ValidationError{
			Field:   fmt.Sprintf("%s.@type", fieldPrefix),
			Value:   rr.Type,
			Message: "must be 'RecurrenceRule'",
		})
	}
	
	// Validate frequency
	validFrequencies := map[string]bool{
		"yearly":   true,
		"monthly":  true,
		"weekly":   true,
		"daily":    true,
		"hourly":   true,
		"minutely": true,
		"secondly": true,
	}
	if !validFrequencies[rr.Frequency] {
		errors = append(errors, ValidationError{
			Field:   fmt.Sprintf("%s.frequency", fieldPrefix),
			Value:   rr.Frequency,
			Message: "must be one of: yearly, monthly, weekly, daily, hourly, minutely, secondly",
		})
	}
	
	// Validate interval
	if rr.Interval != nil && *rr.Interval < 1 {
		errors = append(errors, ValidationError{
			Field:   fmt.Sprintf("%s.interval", fieldPrefix),
			Value:   *rr.Interval,
			Message: "must be positive",
		})
	}
	
	// Validate count
	if rr.Count != nil && *rr.Count < 1 {
		errors = append(errors, ValidationError{
			Field:   fmt.Sprintf("%s.count", fieldPrefix),
			Value:   *rr.Count,
			Message: "must be positive",
		})
	}
	
	// Validate first day of week
	if rr.FirstDayOfWeek != nil {
		if *rr.FirstDayOfWeek < 0 || *rr.FirstDayOfWeek > 6 {
			errors = append(errors, ValidationError{
				Field:   fmt.Sprintf("%s.firstDayOfWeek", fieldPrefix),
				Value:   *rr.FirstDayOfWeek,
				Message: "must be between 0 (Monday) and 6 (Sunday)",
			})
		}
	}
	
	// Validate byDay
	for i, nday := range rr.ByDay {
		validDays := map[string]bool{
			"mo": true, "tu": true, "we": true, "th": true,
			"fr": true, "sa": true, "su": true,
		}
		if !validDays[nday.Day] {
			errors = append(errors, ValidationError{
				Field:   fmt.Sprintf("%s.byDay[%d].day", fieldPrefix, i),
				Value:   nday.Day,
				Message: "must be one of: mo, tu, we, th, fr, sa, su",
			})
		}
	}
	
	return errors
}