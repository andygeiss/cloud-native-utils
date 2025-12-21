package scheduling

import (
	"testing"
	"time"

	"github.com/andygeiss/cloud-native-utils/assert"
)

func Test_NewSlot_With_ValidParams_Should_CreateSlot(t *testing.T) {
	// Arrange
	date := time.Date(2025, 12, 21, 0, 0, 0, 0, time.UTC)
	start := MustTimeOfDay(9, 0)
	end := MustTimeOfDay(9, 30)
	// Act
	result := NewSlot(date, start, end, SlotStateAvailable)
	// Assert
	assert.That(t, "id", result.ID, "2025-12-21-09:00")
	assert.That(t, "start time", result.StartTime, start)
	assert.That(t, "end time", result.EndTime, end)
	assert.That(t, "state", result.State, SlotStateAvailable)
}

func Test_Slot_With_DifferentStates_Should_ReturnIsAvailable(t *testing.T) {
	// Arrange
	date := time.Date(2025, 12, 21, 0, 0, 0, 0, time.UTC)
	available := NewSlot(date, MustTimeOfDay(9, 0), MustTimeOfDay(9, 30), SlotStateAvailable)
	booked := NewSlot(date, MustTimeOfDay(9, 0), MustTimeOfDay(9, 30), SlotStateBooked)
	held := NewSlot(date, MustTimeOfDay(9, 0), MustTimeOfDay(9, 30), SlotStateHeld)
	// Act & Assert
	assert.That(t, "available", available.IsAvailable(), true)
	assert.That(t, "booked", booked.IsAvailable(), false)
	assert.That(t, "held", held.IsAvailable(), false)
}

func Test_Slot_With_TimeRange_Should_ReturnDuration(t *testing.T) {
	// Arrange
	date := time.Date(2025, 12, 21, 0, 0, 0, 0, time.UTC)
	slot := NewSlot(date, MustTimeOfDay(9, 0), MustTimeOfDay(9, 30), SlotStateAvailable)
	// Act
	result := slot.Duration()
	// Assert
	assert.That(t, "duration", result, 30)
}

func Test_NewDaySlots_With_OpenDay_Should_CreateSlots(t *testing.T) {
	// Arrange
	date := time.Date(2025, 12, 21, 0, 0, 0, 0, time.UTC)
	dayHours, _ := NewDayHours(Sunday, MustTimeOfDay(9, 0), MustTimeOfDay(12, 0))
	// Act
	result := NewDaySlots(date, dayHours, 30)
	// Assert
	assert.That(t, "slot count", len(result.Slots), 6) // 3 hours = 6 x 30 min slots
	assert.That(t, "first slot ID", result.Slots[0].ID, "2025-12-21-09:00")
	assert.That(t, "last slot ID", result.Slots[5].ID, "2025-12-21-11:30")
}

func Test_NewDaySlots_With_ClosedDay_Should_ReturnEmptySlots(t *testing.T) {
	// Arrange
	date := time.Date(2025, 12, 21, 0, 0, 0, 0, time.UTC)
	dayHours := NewClosedDay(Sunday)
	// Act
	result := NewDaySlots(date, dayHours, 30)
	// Assert
	assert.That(t, "slot count", len(result.Slots), 0)
}

func Test_DaySlots_With_BookedSlot_Should_ReturnAvailableCount(t *testing.T) {
	// Arrange
	date := time.Date(2025, 12, 21, 0, 0, 0, 0, time.UTC)
	dayHours, _ := NewDayHours(Sunday, MustTimeOfDay(9, 0), MustTimeOfDay(11, 0))
	daySlots := NewDaySlots(date, dayHours, 30)
	daySlots.Slots[0].State = SlotStateBooked
	// Act
	result := daySlots.AvailableCount()
	// Assert
	assert.That(t, "available count", result, 3) // 4 slots total - 1 booked
}

func Test_DaySlots_With_SlotID_Should_ReturnGetSlot(t *testing.T) {
	// Arrange
	date := time.Date(2025, 12, 21, 0, 0, 0, 0, time.UTC)
	dayHours, _ := NewDayHours(Sunday, MustTimeOfDay(9, 0), MustTimeOfDay(11, 0))
	daySlots := NewDaySlots(date, dayHours, 30)
	// Act
	found := daySlots.GetSlot("2025-12-21-09:30")
	notFound := daySlots.GetSlot("invalid-id")
	// Assert
	assert.That(t, "found slot", found != nil, true)
	assert.That(t, "found slot start", found.StartTime.String(), "09:30")
	assert.That(t, "not found slot", notFound == nil, true)
}

func Test_DaySlots_With_ValidSlotID_Should_MarkBooked(t *testing.T) {
	// Arrange
	date := time.Date(2025, 12, 21, 0, 0, 0, 0, time.UTC)
	dayHours, _ := NewDayHours(Sunday, MustTimeOfDay(9, 0), MustTimeOfDay(11, 0))
	daySlots := NewDaySlots(date, dayHours, 30)
	// Act
	success := daySlots.MarkBooked("2025-12-21-09:00")
	failure := daySlots.MarkBooked("invalid-id")
	// Assert
	assert.That(t, "success", success, true)
	assert.That(t, "slot state", daySlots.Slots[0].State, SlotStateBooked)
	assert.That(t, "failure", failure, false)
}

func Test_DaySlots_With_ValidSlotID_Should_MarkHeld(t *testing.T) {
	// Arrange
	date := time.Date(2025, 12, 21, 0, 0, 0, 0, time.UTC)
	dayHours, _ := NewDayHours(Sunday, MustTimeOfDay(9, 0), MustTimeOfDay(11, 0))
	daySlots := NewDaySlots(date, dayHours, 30)
	// Act
	success := daySlots.MarkHeld("2025-12-21-09:00")
	// Assert
	assert.That(t, "success", success, true)
	assert.That(t, "slot state", daySlots.Slots[0].State, SlotStateHeld)
}
