package scheduling

import (
	"errors"
	"time"
)

// Errors for opening hours validation.
var (
	// ErrCloseBeforeOpen is returned when close time is not after open time.
	ErrCloseBeforeOpen = errors.New("close time must be after open time")
)

// DayHours represents the opening hours for a single day.
type DayHours struct {
	// CloseTime is when the facility closes.
	CloseTime TimeOfDay `json:"close_time"`
	// IsClosed indicates if the facility is closed this day.
	IsClosed bool `json:"is_closed"`
	// OpenTime is when the facility opens.
	OpenTime TimeOfDay `json:"open_time"`
	// Weekday is the day of the week.
	Weekday Weekday `json:"weekday"`
}

// NewDayHours creates opening hours for a specific day.
// Returns error if close time is not after open time.
func NewDayHours(weekday Weekday, openTime, closeTime TimeOfDay) (*DayHours, error) {
	if !closeTime.After(openTime) {
		return nil, ErrCloseBeforeOpen
	}
	return &DayHours{
		CloseTime: closeTime,
		IsClosed:  false,
		OpenTime:  openTime,
		Weekday:   weekday,
	}, nil
}

// NewClosedDay creates a closed day for a specific weekday.
func NewClosedDay(weekday Weekday) *DayHours {
	return &DayHours{
		IsClosed: true,
		Weekday:  weekday,
	}
}

// Duration returns the number of minutes the facility is open.
func (a *DayHours) Duration() int {
	if a.IsClosed {
		return 0
	}
	return a.CloseTime.TotalMinutes() - a.OpenTime.TotalMinutes()
}

// SlotCount returns the number of slots available given slot duration.
func (a *DayHours) SlotCount(slotDurationMinutes int) int {
	if a.IsClosed || slotDurationMinutes <= 0 {
		return 0
	}
	return a.Duration() / slotDurationMinutes
}

// IsOpen returns true if the facility is open at the given time.
func (a *DayHours) IsOpen(t TimeOfDay) bool {
	if a.IsClosed {
		return false
	}
	return !t.Before(a.OpenTime) && t.Before(a.CloseTime)
}

// Validate checks if the day hours are valid.
func (a *DayHours) Validate() error {
	if a.IsClosed {
		return nil
	}
	if !a.CloseTime.After(a.OpenTime) {
		return ErrCloseBeforeOpen
	}
	return nil
}

// OpeningHours represents the weekly opening hours for a facility.
type OpeningHours struct {
	// Days maps weekday to opening hours.
	Days map[Weekday]*DayHours `json:"days"`
	// FacilityID is the facility these hours belong to.
	// Note: In generic context, FacilityID might not be relevant, but keeping it for compatibility or removing it?
	// I'll remove FacilityID as it's domain specific.
	// But wait, the original code had it. If I remove it, I need to update the usage.
	// Let's keep it generic. OpeningHours is just a schedule.
	// If the user needs to associate it with a facility, they can wrap it.
}

// NewOpeningHours creates a new opening hours schedule.
// All days default to closed.
func NewOpeningHours() *OpeningHours {
	days := make(map[Weekday]*DayHours)
	for i := Monday; i <= Sunday; i++ {
		days[i] = NewClosedDay(i)
	}
	return &OpeningHours{
		Days: days,
	}
}

// GetDayHours returns the opening hours for a specific weekday.
func (a *OpeningHours) GetDayHours(weekday Weekday) *DayHours {
	if hours, ok := a.Days[weekday]; ok {
		return hours
	}
	return NewClosedDay(weekday)
}

// GetDayHoursForDate returns the opening hours for a specific date.
func (a *OpeningHours) GetDayHoursForDate(date time.Time) *DayHours {
	weekday := WeekdayFromGoWeekday(date.Weekday())
	return a.GetDayHours(weekday)
}

// IsOpenOn returns true if the facility is open on the given weekday.
func (a *OpeningHours) IsOpenOn(weekday Weekday) bool {
	hours := a.GetDayHours(weekday)
	return !hours.IsClosed
}

// SetDayHours sets the opening hours for a specific day.
func (a *OpeningHours) SetDayHours(hours *DayHours) error {
	if err := hours.Validate(); err != nil {
		return err
	}
	a.Days[hours.Weekday] = hours
	return nil
}

// SetOpen sets the facility as open for a specific day.
func (a *OpeningHours) SetOpen(weekday Weekday, openTime, closeTime TimeOfDay) error {
	hours, err := NewDayHours(weekday, openTime, closeTime)
	if err != nil {
		return err
	}
	a.Days[weekday] = hours
	return nil
}

// SetClosed sets the facility as closed for a specific day.
func (a *OpeningHours) SetClosed(weekday Weekday) {
	a.Days[weekday] = NewClosedDay(weekday)
}

// Validate checks if all opening hours are valid.
func (a *OpeningHours) Validate() error {
	for _, hours := range a.Days {
		if err := hours.Validate(); err != nil {
			return err
		}
	}
	return nil
}
