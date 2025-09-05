package jscal

// Status values for Event status property (RFC 8984 Section 5.1.3)
const (
	StatusConfirmed = "confirmed"
	StatusCancelled = "cancelled"
	StatusTentative = "tentative"
)

// FreeBusyStatus values (RFC 8984 Section 4.4.2)
const (
	FreeBusyFree        = "free"
	FreeBusyBusy        = "busy"
	FreeBusyTentative   = "tentative"
	FreeBusyUnavailable = "unavailable"
)

// Privacy values (RFC 8984 Section 4.4.3)
const (
	PrivacyPublic  = "public"
	PrivacyPrivate = "private"
	PrivacySecret  = "secret"
)

// ParticipationStatus values (RFC 8984 Section 4.4.6)
const (
	ParticipationNeedsAction = "needs-action"
	ParticipationAccepted    = "accepted"
	ParticipationDeclined    = "declined"
	ParticipationTentative   = "tentative"
	ParticipationDelegated   = "delegated"
)

// Participant roles (RFC 8984 Section 4.4.6)
const (
	RoleOwner         = "owner"
	RoleAttendee      = "attendee"
	RoleOptional      = "optional"
	RoleInformational = "informational"
	RoleChair         = "chair"
	RoleContact       = "contact"
)

// Participant kinds (RFC 8984 Section 4.4.6)
const (
	KindIndividual = "individual"
	KindGroup      = "group"
	KindResource   = "resource"
	KindLocation   = "location"
	KindUnknown    = "unknown"
)

// ScheduleAgent values (RFC 8984 Section 4.4.6)
const (
	ScheduleAgentServer = "server"
	ScheduleAgentClient = "client"
	ScheduleAgentNone   = "none"
)

// Method values (RFC 8984 Section 4.1.8)
const (
	MethodPublish        = "publish"
	MethodRequest        = "request"
	MethodReply          = "reply"
	MethodAdd            = "add"
	MethodCancel         = "cancel"
	MethodRefresh        = "refresh"
	MethodCounter        = "counter"
	MethodDeclineCounter = "declineCounter"
)

// Progress values for tasks (used in Participant)
const (
	ProgressNeedsAction = "needs-action"
	ProgressInProcess   = "in-process"
	ProgressCompleted   = "completed"
	ProgressFailed      = "failed"
	ProgressCancelled   = "cancelled"
)

// Common MIME types for descriptionContentType
const (
	MIMETextPlain    = "text/plain"
	MIMETextHTML     = "text/html"
	MIMETextMarkdown = "text/markdown"
)

// Priority range constraints (RFC 8984 Section 4.4.1)
const (
	PriorityMin     = 0
	PriorityMax     = 9
	PriorityDefault = 0
)

// Alert trigger relationships
const (
	AlertTriggerStart = "start"
	AlertTriggerEnd   = "end"
)

// RelativeTo values for Location and OffsetTrigger
const (
	RelativeToStart = "start"
	RelativeToEnd   = "end"
)

// Alert actions
const (
	AlertActionDisplay = "display"
	AlertActionEmail   = "email"
)

// Link relation types
const (
	LinkRelationAlternate   = "alternate"
	LinkRelationIcon        = "icon"
	LinkRelationAttachment  = "attachment"
	LinkRelationDescribedBy = "describedby"
	LinkRelationEnclosure   = "enclosure"
)

// Recurrence frequency values
const (
	FrequencyYearly   = "yearly"
	FrequencyMonthly  = "monthly"
	FrequencyWeekly   = "weekly"
	FrequencyDaily    = "daily"
	FrequencyHourly   = "hourly"
	FrequencyMinutely = "minutely"
	FrequencySecondly = "secondly"
)

// Days of the week
const (
	DayMonday    = "mo"
	DayTuesday   = "tu"
	DayWednesday = "we"
	DayThursday  = "th"
	DayFriday    = "fr"
	DaySaturday  = "sa"
	DaySunday    = "su"
)

// Relation types for relatedTo
const (
	RelationTypeParent  = "parent"
	RelationTypeChild   = "child"
	RelationTypeSibling = "sibling"
	RelationTypeNext    = "next"
	RelationTypePrior   = "prior"
)

// Skip values for RecurrenceRule
const (
	SkipOmit     = "omit"
	SkipBackward = "backward"
	SkipForward  = "forward"
)
