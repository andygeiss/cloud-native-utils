package scheduling

import (
	"fmt"
)

// TimeOfDay represents a time without date (hours and minutes).
// Minutes are stored as offset from midnight (0-1439).
type TimeOfDay int

// NewTimeOfDay creates a TimeOfDay from hours and minutes.
// Returns error if hours or minutes are out of range.
func NewTimeOfDay(hours, minutes int) (TimeOfDay, error) {
	if hours < 0 || hours > 23 {
		return 0, fmt.Errorf("hours must be 0-23, got %d", hours)
	}
	if minutes < 0 || minutes > 59 {
		return 0, fmt.Errorf("minutes must be 0-59, got %d", minutes)
	}
	return TimeOfDay(hours*60 + minutes), nil
}

// ParseTimeOfDay parses a time-of-day string in "HH:MM" format.
func ParseTimeOfDay(value string) (TimeOfDay, error) {
	var hours, minutes int
	if _, err := fmt.Sscanf(value, "%d:%d", &hours, &minutes); err != nil {
		return 0, fmt.Errorf("invalid time format: %q", value)
	}
	return NewTimeOfDay(hours, minutes)
}

// MustTimeOfDay creates a TimeOfDay from hours and minutes, panics on error.
func MustTimeOfDay(hours, minutes int) TimeOfDay {
	t, err := NewTimeOfDay(hours, minutes)
	if err != nil {
		panic(err)
	}
	return t
}

// Hours returns the hour component (0-23).
func (a TimeOfDay) Hours() int {
	return int(a) / 60
}

// Minutes returns the minute component (0-59).
func (a TimeOfDay) Minutes() int {
	return int(a) % 60
}

// String returns the time in HH:MM format.
func (a TimeOfDay) String() string {
	return fmt.Sprintf("%02d:%02d", a.Hours(), a.Minutes())
}

// TotalMinutes returns the total minutes from midnight.
func (a TimeOfDay) TotalMinutes() int {
	return int(a)
}

// Before returns true if this time is before other.
func (a TimeOfDay) Before(other TimeOfDay) bool {
	return int(a) < int(other)
}

// After returns true if this time is after other.
func (a TimeOfDay) After(other TimeOfDay) bool {
	return int(a) > int(other)
}

// Equal returns true if this time equals other.
func (a TimeOfDay) Equal(other TimeOfDay) bool {
	return int(a) == int(other)
}

// Add adds minutes to the time and returns a new TimeOfDay.
// Does not wrap around midnight.
func (a TimeOfDay) Add(minutes int) TimeOfDay {
	result := int(a) + minutes
	if result < 0 {
		return TimeOfDay(0)
	}
	if result > 24*60-1 {
		return TimeOfDay(24*60 - 1)
	}
	return TimeOfDay(result)
}
