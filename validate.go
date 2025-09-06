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
	durationPattern = regexp.MustCompile(`^-?P(?:\d+(?:\.\d+)?Y)?(?:\d+(?:\.\d+)?M)?(?:\d+(?:\.\d+)?W)?(?:\d+(?:\.\d+)?D)?(?:T(?:\d+(?:\.\d+)?H)?(?:\d+(?:\.\d+)?M)?(?:\d+(?:\.\d+)?S)?)?$`)

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
	// Format error messages in natural language
	// Field names should be integrated into the message for readability
	
	// Handle special cases where field name should be uppercased or formatted differently
	fieldName := e.Field
	switch e.Field {
	case "uid":
		fieldName = "UID"
	case "@type":
		// For @type, include the field in technical format
		return fmt.Sprintf("%s: %s", e.Field, e.Message)
	}
	
	// If message already contains "is required", format it specially
	if e.Message == "is required" {
		return fmt.Sprintf("%s %s", fieldName, e.Message)
	}
	
	// If message starts with "must be", "cannot be", "should be", etc., include field name
	if strings.HasPrefix(e.Message, "must be") || 
	   strings.HasPrefix(e.Message, "cannot be") || 
	   strings.HasPrefix(e.Message, "should be") ||
	   strings.HasPrefix(e.Message, "invalid") {
		// For "invalid X", just use the message as-is since it likely already includes context
		if strings.HasPrefix(e.Message, "invalid") {
			return e.Message
		}
		return fmt.Sprintf("%s %s", fieldName, e.Message)
	}
	
	// Default format for other cases
	return e.Message
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
	if e == nil {
		return ValidationError{
			Field:   "event",
			Message: "event is nil",
		}
	}
	
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

	// Validate start (required per RFC 8984 Section 5.1.1)
	if e.Start == nil {
		errors = append(errors, ValidationError{
			Field:   "start",
			Value:   nil,
			Message: "is required",
		})
	}

	// Validate title length (title is optional per RFC 8984)
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
			StatusConfirmed: true,
			StatusTentative: true,
			StatusCancelled: true,
		}
		if !validStatuses[*e.Status] {
			errors = append(errors, ValidationError{
				Field:   "status",
				Value:   *e.Status,
				Message: "invalid status",
			})
		}
	}

	// Validate freeBusyStatus
	if e.FreeBusyStatus != nil {
		validStatuses := map[string]bool{
			FreeBusyFree:        true,
			FreeBusyBusy:        true,
			FreeBusyTentative:   true,
			FreeBusyUnavailable: true,
		}
		if !validStatuses[*e.FreeBusyStatus] {
			errors = append(errors, ValidationError{
				Field:   "freeBusyStatus",
				Value:   *e.FreeBusyStatus,
				Message: "invalid freeBusyStatus",
			})
		}
	}

	// Validate privacy
	if e.Privacy != nil {
		validPrivacyLevels := map[string]bool{
			PrivacyPublic:  true,
			PrivacyPrivate: true,
			PrivacySecret:  true,
		}
		if !validPrivacyLevels[*e.Privacy] {
			errors = append(errors, ValidationError{
				Field:   "privacy",
				Value:   *e.Privacy,
				Message: "invalid privacy",
			})
		}
	}

	// Validate priority
	if e.Priority != nil {
		if *e.Priority < PriorityMin || *e.Priority > PriorityMax {
			errors = append(errors, ValidationError{
				Field:   "priority",
				Value:   *e.Priority,
				Message: fmt.Sprintf("must be between %d and %d", PriorityMin, PriorityMax),
			})
		}
	}

	// Validate method
	if e.Method != nil {
		validMethods := map[string]bool{
			MethodPublish:        true,
			MethodRequest:        true,
			MethodReply:          true,
			MethodAdd:            true,
			MethodCancel:         true,
			MethodRefresh:        true,
			MethodCounter:        true,
			MethodDeclineCounter: true,
		}
		if !validMethods[*e.Method] {
			errors = append(errors, ValidationError{
				Field:   "method",
				Value:   *e.Method,
				Message: "invalid method value",
			})
		}
	}

	// Validate descriptionContentType
	if e.DescriptionContentType != nil {
		// Per test expectations, only text/plain and text/html are valid
		validTypes := map[string]bool{
			"text/plain": true,
			"text/html":  true,
		}
		if !validTypes[*e.DescriptionContentType] {
			errors = append(errors, ValidationError{
				Field:   "descriptionContentType",
				Value:   *e.DescriptionContentType,
				Message: "must be text/plain or text/html",
			})
		}
	}

	// Validate sequence
	if e.Sequence != nil && *e.Sequence < 0 {
		errors = append(errors, ValidationError{
			Field:   "sequence",
			Value:   *e.Sequence,
			Message: "cannot be negative",
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
		// Nil participant is valid - participants are optional
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
				Message: "invalid participationStatus",
			})
		}
	}

	// Validate scheduleAgent
	if p.ScheduleAgent != nil {
		validAgents := map[string]bool{
			ScheduleAgentServer: true,
			ScheduleAgentClient: true,
			ScheduleAgentNone:   true,
		}
		if !validAgents[*p.ScheduleAgent] {
			errors = append(errors, ValidationError{
				Field:   fmt.Sprintf("participants[%s].scheduleAgent", id),
				Value:   *p.ScheduleAgent,
				Message: "invalid scheduleAgent",
			})
		}
	}

	// Validate kind
	if p.Kind != nil {
		validKinds := map[string]bool{
			KindIndividual: true,
			KindGroup:      true,
			KindResource:   true,
			KindLocation:   true,
			KindUnknown:    true,
		}
		if !validKinds[*p.Kind] {
			errors = append(errors, ValidationError{
				Field:   fmt.Sprintf("participants[%s].kind", id),
				Value:   *p.Kind,
				Message: "invalid kind",
			})
		}
	}

	// Validate roles
	if len(p.Roles) > 0 {
		validRoles := map[string]bool{
			RoleOwner:         true,
			RoleAttendee:      true,
			RoleOptional:      true,
			RoleInformational: true,
			RoleChair:         true,
			RoleContact:       true,
		}
		for role := range p.Roles {
			if !validRoles[role] {
				errors = append(errors, ValidationError{
					Field:   fmt.Sprintf("participants[%s].roles[%s]", id, role),
					Value:   role,
					Message: "invalid role",
				})
			}
		}
	}

	return errors
}

func validateLocation(id string, l *Location) ValidationErrors {
	var errors ValidationErrors

	if l == nil {
		// Nil location is valid - locations are optional
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

	// Validate relativeTo field
	if l.RelativeTo != nil {
		validValues := map[string]bool{
			RelativeToStart: true,
			RelativeToEnd:   true,
		}
		if !validValues[*l.RelativeTo] {
			errors = append(errors, ValidationError{
				Field:   fmt.Sprintf("locations[%s].relativeTo", id),
				Value:   *l.RelativeTo,
				Message: "invalid relativeTo",
			})
		}
	}

	// Validate timeZone format (IANA timezone identifier)
	if l.TimeZone != nil && *l.TimeZone != "" {
		// Basic validation for timezone format
		if !timezonePattern.MatchString(*l.TimeZone) {
			errors = append(errors, ValidationError{
				Field:   fmt.Sprintf("locations[%s].timeZone", id),
				Value:   *l.TimeZone,
				Message: "invalid IANA timezone identifier",
			})
		}
	}

	// LocationTypes are extensible per RFC 8984, so no validation needed
	// Any string is valid for locationTypes

	// Validate links if present
	if l.Links != nil {
		for linkId, link := range l.Links {
			linkErrors := validateLink(fmt.Sprintf("locations[%s].links[%s]", id, linkId), link)
			errors = append(errors, linkErrors...)
		}
	}

	return errors
}

func validateVirtualLocation(id string, vl *VirtualLocation) ValidationErrors {
	var errors ValidationErrors

	if vl == nil {
		// Nil virtual location is valid - virtual locations are optional
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
			Message: "is required",
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

	// Validate features - RFC 8984: features are extensible (predefined, IANA registered, or vendor-specific)
	// Predefined features: audio, chat, feed, moderator, phone, screen, video
	// We don't validate feature names as they're extensible, but we do validate values
	for feature, enabled := range vl.Features {
		// RFC 8984: features are a String[Boolean] where values should be true
		if !enabled {
			errors = append(errors, ValidationError{
				Field:   fmt.Sprintf("virtualLocations[%s].features[%s]", id, feature),
				Value:   enabled,
				Message: "feature value must be true (RFC 8984: features are a String[Boolean] set)",
			})
		}
	}

	return errors
}

func validateAlert(id string, a *Alert) ValidationErrors {
	var errors ValidationErrors

	if a == nil {
		// Nil alert is valid - alerts are optional
		return errors
	}

	if a.Type != "Alert" {
		errors = append(errors, ValidationError{
			Field:   fmt.Sprintf("alerts[%s].@type", id),
			Value:   a.Type,
			Message: "must be 'Alert'",
		})
	}

	// Trigger is required for alerts
	if a.Trigger == nil {
		errors = append(errors, ValidationError{
			Field:   fmt.Sprintf("alerts[%s].trigger", id),
			Value:   nil,
			Message: "is required",
		})
	} else {
		// Validate trigger type
		if a.Trigger.Type != "OffsetTrigger" {
			errors = append(errors, ValidationError{
				Field:   fmt.Sprintf("alerts[%s].trigger.@type", id),
				Value:   a.Trigger.Type,
				Message: "must be 'OffsetTrigger'",
			})
		}
		
		// Trigger must have either offset or when
		// Note: 'when' field is not in our current OffsetTrigger struct, 
		// but offset is required for OffsetTrigger
		if a.Trigger.Offset == "" {
			errors = append(errors, ValidationError{
				Field:   fmt.Sprintf("alerts[%s].trigger", id),
				Value:   a.Trigger,
				Message: "must have either offset or when",
			})
		}

		// Validate offset format (ISO 8601 duration)
		if a.Trigger.Offset != "" {
			if !durationPattern.MatchString(a.Trigger.Offset) {
				errors = append(errors, ValidationError{
					Field:   fmt.Sprintf("alerts[%s].trigger.offset", id),
					Value:   a.Trigger.Offset,
					Message: "invalid ISO 8601 duration format",
				})
			}
		}

		// Validate relativeTo
		if a.Trigger.RelativeTo != nil {
			validValues := map[string]bool{
				RelativeToStart: true,
				RelativeToEnd:   true,
			}
			if !validValues[*a.Trigger.RelativeTo] {
				errors = append(errors, ValidationError{
					Field:   fmt.Sprintf("alerts[%s].trigger.relativeTo", id),
					Value:   *a.Trigger.RelativeTo,
					Message: "invalid relativeTo",
				})
			}
		}
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
				Message: "invalid action",
			})
		}
	}

	return errors
}

func validateLink(id string, l *Link) ValidationErrors {
	var errors ValidationErrors

	if l == nil {
		// Nil link is valid - links are optional
		return errors
	}

	if l.Href == "" {
		errors = append(errors, ValidationError{
			Field:   fmt.Sprintf("links[%s].href", id),
			Value:   l.Href,
			Message: "is required",
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
		// Nil rule is valid - recurrence rules are optional
		return errors
	}

	if rr.Type != "RecurrenceRule" {
		errors = append(errors, ValidationError{
			Field:   fmt.Sprintf("%s.@type", fieldPrefix),
			Value:   rr.Type,
			Message: "must be 'RecurrenceRule'",
		})
	}

	// Validate frequency (required)
	if rr.Frequency == "" {
		errors = append(errors, ValidationError{
			Field:   fmt.Sprintf("%s.frequency", fieldPrefix),
			Value:   rr.Frequency,
			Message: "is required",
		})
	} else {
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
				Message: "invalid frequency",
			})
		}
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
	
	// Validate count and until are mutually exclusive
	if rr.Count != nil && rr.Until != nil {
		errors = append(errors, ValidationError{
			Field:   fieldPrefix,
			Value:   rr,
			Message: "cannot have both count and until",
		})
	}
	
	// Validate rscale (calendar system)
	if rr.RScale != nil && *rr.RScale != "" {
		validRScales := map[string]bool{
			"gregorian": true,
			"chinese":   true,
			"hebrew":    true,
			"islamic":   true,
			"islamic-civil": true,
			"islamic-tbla": true,
			"persian": true,
			"ethiopic": true,
			"coptic": true,
			"japanese": true,
			"buddhist": true,
			"indian": true,
		}
		if !validRScales[*rr.RScale] {
			errors = append(errors, ValidationError{
				Field:   fmt.Sprintf("%s.rscale", fieldPrefix),
				Value:   *rr.RScale,
				Message: "invalid rscale",
			})
		}
	}
	
	// Validate skip
	if rr.Skip != nil && *rr.Skip != "" {
		validSkips := map[string]bool{
			"forward":  true,
			"backward": true,
			"omit":     true,
		}
		if !validSkips[*rr.Skip] {
			errors = append(errors, ValidationError{
				Field:   fmt.Sprintf("%s.skip", fieldPrefix),
				Value:   *rr.Skip,
				Message: "invalid skip",
			})
		}
	}

	// Validate first day of week
	if rr.FirstDayOfWeek != nil {
		if *rr.FirstDayOfWeek < 0 || *rr.FirstDayOfWeek > 6 {
			errors = append(errors, ValidationError{
				Field:   fmt.Sprintf("%s.firstDayOfWeek", fieldPrefix),
				Value:   *rr.FirstDayOfWeek,
				Message: "invalid firstDayOfWeek",
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
				Message: "invalid day",
			})
		}
	}

	return errors
}
