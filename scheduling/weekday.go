package scheduling

import (
	"time"
)

// Weekday represents a day of the week for opening hours.
type Weekday int

// Weekday constants for opening hours.
const (
	Monday Weekday = iota
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
	Sunday
)

// String returns the string representation of a weekday.
func (a Weekday) String() string {
	switch a {
	case Monday:
		return "Monday"
	case Tuesday:
		return "Tuesday"
	case Wednesday:
		return "Wednesday"
	case Thursday:
		return "Thursday"
	case Friday:
		return "Friday"
	case Saturday:
		return "Saturday"
	case Sunday:
		return "Sunday"
	default:
		return "Unknown"
	}
}

// GoWeekday converts facility Weekday to time.Weekday.
func (a Weekday) GoWeekday() time.Weekday {
	switch a {
	case Monday:
		return time.Monday
	case Tuesday:
		return time.Tuesday
	case Wednesday:
		return time.Wednesday
	case Thursday:
		return time.Thursday
	case Friday:
		return time.Friday
	case Saturday:
		return time.Saturday
	case Sunday:
		return time.Sunday
	default:
		return time.Monday
	}
}

// WeekdayFromGoWeekday converts time.Weekday to facility Weekday.
func WeekdayFromGoWeekday(w time.Weekday) Weekday {
	switch w {
	case time.Monday:
		return Monday
	case time.Tuesday:
		return Tuesday
	case time.Wednesday:
		return Wednesday
	case time.Thursday:
		return Thursday
	case time.Friday:
		return Friday
	case time.Saturday:
		return Saturday
	case time.Sunday:
		return Sunday
	default:
		return Monday
	}
}
