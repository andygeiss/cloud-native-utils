package scheduling

import (
	"errors"
	"fmt"
	"time"
)

var ErrEndBeforeStart = errors.New("end time must be after start time")

// TimeRange represents a contiguous time range for a selection.
// This is a value object - immutable once created.
type TimeRange struct {
	date      time.Time
	endTime   TimeOfDay
	startTime TimeOfDay
}

// NewTimeRange creates a new time range.
// Returns an error if end time is not after start time.
func NewTimeRange(date time.Time, startTime, endTime TimeOfDay) (*TimeRange, error) {
	if !endTime.After(startTime) {
		return nil, ErrEndBeforeStart
	}
	return &TimeRange{
		date:      date,
		endTime:   endTime,
		startTime: startTime,
	}, nil
}

// Date returns the date of the time range.
func (a *TimeRange) Date() time.Time {
	return a.date
}

// EndTime returns the end time of the range.
func (a *TimeRange) EndTime() TimeOfDay {
	return a.endTime
}

// StartTime returns the start time of the range.
func (a *TimeRange) StartTime() TimeOfDay {
	return a.startTime
}

// Duration returns the duration in minutes.
func (a *TimeRange) Duration() int {
	return a.endTime.TotalMinutes() - a.startTime.TotalMinutes()
}

// SlotCount returns the number of slots in this range.
func (a *TimeRange) SlotCount(slotDurationMinutes int) int {
	if slotDurationMinutes <= 0 {
		return 0
	}
	return a.Duration() / slotDurationMinutes
}

// Contains returns true if the given time is within the range.
func (a *TimeRange) Contains(t TimeOfDay) bool {
	return !t.Before(a.startTime) && t.Before(a.endTime)
}

// Overlaps returns true if this range overlaps with another.
func (a *TimeRange) Overlaps(other *TimeRange) bool {
	if a.date.Format("2006-01-02") != other.date.Format("2006-01-02") {
		return false
	}
	return a.startTime.Before(other.endTime) && other.startTime.Before(a.endTime)
}

// String returns a human-readable representation of the time range.
func (a *TimeRange) String() string {
	return fmt.Sprintf("%s %s-%s", a.date.Format("2006-01-02"), a.startTime.String(), a.endTime.String())
}
