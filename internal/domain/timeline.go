package domain

import "time"

type TimelineEventKind string

const (
	EventCampaignOpened TimelineEventKind = "campaign-opened"
	EventStationLinked  TimelineEventKind = "station-linked"
	EventShiftOpened    TimelineEventKind = "shift-opened"
	EventFrameSubmitted TimelineEventKind = "frame-submitted"
	EventFrameReviewed  TimelineEventKind = "frame-reviewed"
	EventShiftClosed    TimelineEventKind = "shift-closed"
	EventReleased       TimelineEventKind = "released"
)

type TimelineEvent struct {
	Kind      TimelineEventKind `json:"kind"`
	At        time.Time         `json:"at"`
	SubjectID string            `json:"subjectId"`
	Message   string            `json:"message"`
}

type Timeline struct {
	Events []TimelineEvent `json:"events"`
}

func (t *Timeline) Add(kind TimelineEventKind, at time.Time, subjectID, message string) {
	t.Events = append(t.Events, TimelineEvent{Kind: kind, At: at.UTC(), SubjectID: subjectID, Message: message})
}
func (t Timeline) First() (TimelineEvent, bool) {
	if len(t.Events) == 0 {
		return TimelineEvent{}, false
	}
	return t.Events[0], true
}
func (t Timeline) Last() (TimelineEvent, bool) {
	if len(t.Events) == 0 {
		return TimelineEvent{}, false
	}
	return t.Events[len(t.Events)-1], true
}
func (t Timeline) Between(start, end time.Time) []TimelineEvent {
	out := make([]TimelineEvent, 0)
	for _, event := range t.Events {
		if !event.At.Before(start) && !event.At.After(end) {
			out = append(out, event)
		}
	}
	return out
}
func (t Timeline) Count(kind TimelineEventKind) int {
	if kind == "" {
		return len(t.Events)
	}
	count := 0
	for _, event := range t.Events {
		if event.Kind == kind {
			count++
		}
	}
	return count
}
func (t Timeline) HasSubject(subjectID string) bool {
	for _, event := range t.Events {
		if event.SubjectID == subjectID {
			return true
		}
	}
	return false
}
func (t Timeline) IsOrdered() bool {
	for index := 1; index < len(t.Events); index++ {
		if t.Events[index].At.Before(t.Events[index-1].At) {
			return false
		}
	}
	return true
}
func (t Timeline) Duration() time.Duration {
	first, ok := t.First()
	if !ok {
		return 0
	}
	last, _ := t.Last()
	return last.At.Sub(first.At)
}
func (t Timeline) Kinds() []TimelineEventKind {
	kinds := make([]TimelineEventKind, 0, len(t.Events))
	for _, event := range t.Events {
		kinds = append(kinds, event.Kind)
	}
	return kinds
}
