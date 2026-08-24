package domain

import "time"

type TimeWindow struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

func NewTimeWindow(start, end time.Time) TimeWindow {
	return TimeWindow{Start: start.UTC(), End: end.UTC()}
}
func (w TimeWindow) Valid() bool                { return w.End.After(w.Start) }
func (w TimeWindow) Duration() time.Duration    { return w.End.Sub(w.Start) }
func (w TimeWindow) Contains(at time.Time) bool { return !at.Before(w.Start) && !at.After(w.End) }
func (w TimeWindow) ContainsWindow(other TimeWindow) bool {
	return w.Contains(other.Start) && w.Contains(other.End)
}
func (w TimeWindow) Overlaps(other TimeWindow) bool {
	return w.Start.Before(other.End) && other.Start.Before(w.End)
}
func (w TimeWindow) Touches(other TimeWindow) bool {
	return w.End.Equal(other.Start) || other.End.Equal(w.Start)
}
func (w TimeWindow) Intersect(other TimeWindow) (TimeWindow, bool) {
	if !w.Overlaps(other) {
		return TimeWindow{}, false
	}
	start, end := w.Start, w.End
	if other.Start.After(start) {
		start = other.Start
	}
	if other.End.Before(end) {
		end = other.End
	}
	return NewTimeWindow(start, end), true
}
func (w TimeWindow) Clamp(at time.Time) time.Time {
	if at.Before(w.Start) {
		return w.Start
	}
	if at.After(w.End) {
		return w.End
	}
	return at
}
func (w TimeWindow) Midpoint() time.Time { return w.Start.Add(w.Duration() / 2) }
func (w TimeWindow) Hours() float64      { return w.Duration().Hours() }

type WindowSegment struct {
	Window TimeWindow `json:"window"`
	Label  string     `json:"label"`
}

func SplitWindow(window TimeWindow, parts int) []WindowSegment {
	if !window.Valid() || parts <= 0 {
		return nil
	}
	segments := make([]WindowSegment, 0, parts)
	delta := window.Duration() / time.Duration(parts)
	cursor := window.Start
	for index := 0; index < parts; index++ {
		next := cursor.Add(delta)
		if index == parts-1 {
			next = window.End
		}
		segments = append(segments, WindowSegment{Window: NewTimeWindow(cursor, next), Label: string(rune('A' + index))})
		cursor = next
	}
	return segments
}

func MergeWindows(windows []TimeWindow) []TimeWindow {
	if len(windows) == 0 {
		return nil
	}
	result := make([]TimeWindow, 0, len(windows))
	for _, candidate := range windows {
		if !candidate.Valid() {
			continue
		}
		if len(result) == 0 {
			result = append(result, candidate)
			continue
		}
		last := &result[len(result)-1]
		if last.Overlaps(candidate) || last.Touches(candidate) {
			if candidate.End.After(last.End) {
				last.End = candidate.End
			}
			continue
		}
		result = append(result, candidate)
	}
	return result
}
