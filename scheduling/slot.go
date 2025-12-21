package scheduling

import (
	"time"
)

// SlotState represents the availability state of a time slot.
type SlotState string

// Slot states.
const (
	SlotStateAvailable   SlotState = "available"
	SlotStateBooked      SlotState = "booked"
	SlotStateHeld        SlotState = "held"
	SlotStateMaintenance SlotState = "maintenance"
)

// Slot represents a single bookable time slot.
type Slot struct {
	Date      time.Time `json:"date"`
	EndTime   TimeOfDay `json:"end_time"`
	ID        string    `json:"id"`
	StartTime TimeOfDay `json:"start_time"`
	State     SlotState `json:"state"`
}

// NewSlot creates a new slot.
func NewSlot(date time.Time, startTime, endTime TimeOfDay, state SlotState) *Slot {
	dateStr := date.Format("2006-01-02")
	return &Slot{
		Date:      date,
		EndTime:   endTime,
		ID:        dateStr + "-" + startTime.String(),
		StartTime: startTime,
		State:     state,
	}
}

// IsAvailable returns true if the slot can be booked.
func (a *Slot) IsAvailable() bool {
	return a.State == SlotStateAvailable
}

// Duration returns the slot duration in minutes.
func (a *Slot) Duration() int {
	return a.EndTime.TotalMinutes() - a.StartTime.TotalMinutes()
}

// DaySlots represents all slots for a single day.
type DaySlots struct {
	Date  time.Time `json:"date"`
	Slots []*Slot   `json:"slots"`
}

// NewDaySlots creates slots for a day based on opening hours and slot duration.
func NewDaySlots(date time.Time, dayHours *DayHours, slotDurationMinutes int) *DaySlots {
	ds := &DaySlots{
		Date:  date,
		Slots: make([]*Slot, 0),
	}
	if dayHours.IsClosed {
		return ds
	}
	current := dayHours.OpenTime
	for current.Before(dayHours.CloseTime) {
		end := current.Add(slotDurationMinutes)
		if end.After(dayHours.CloseTime) {
			break
		}
		slot := NewSlot(date, current, end, SlotStateAvailable)
		ds.Slots = append(ds.Slots, slot)
		current = end
	}
	return ds
}

// AvailableCount returns the number of available slots.
func (a *DaySlots) AvailableCount() int {
	count := 0
	for _, slot := range a.Slots {
		if slot.IsAvailable() {
			count++
		}
	}
	return count
}

// GetSlot returns the slot by ID.
func (a *DaySlots) GetSlot(id string) *Slot {
	for _, slot := range a.Slots {
		if slot.ID == id {
			return slot
		}
	}
	return nil
}

// MarkBooked marks a slot as booked.
func (a *DaySlots) MarkBooked(slotID string) bool {
	slot := a.GetSlot(slotID)
	if slot == nil {
		return false
	}
	slot.State = SlotStateBooked
	return true
}

// MarkHeld marks a slot as held.
func (a *DaySlots) MarkHeld(slotID string) bool {
	slot := a.GetSlot(slotID)
	if slot == nil {
		return false
	}
	slot.State = SlotStateHeld
	return true
}
