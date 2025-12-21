package scheduling

import (
	"testing"
	"time"

	"github.com/andygeiss/cloud-native-utils/assert"
)

func Test_NewDayHours_With_ValidTimes_Should_ReturnDayHours(t *testing.T) {
	// Arrange
	open := MustTimeOfDay(9, 0)
	close := MustTimeOfDay(17, 0)
	// Act
	result, err := NewDayHours(Monday, open, close)
	// Assert
	assert.That(t, "error", err, nil)
	assert.That(t, "weekday", result.Weekday, Monday)
	assert.That(t, "open time", result.OpenTime, open)
	assert.That(t, "close time", result.CloseTime, close)
	assert.That(t, "is closed", result.IsClosed, false)
}

func Test_NewDayHours_With_CloseBeforeOpen_Should_ReturnError(t *testing.T) {
	// Arrange
	open := MustTimeOfDay(17, 0)
	close := MustTimeOfDay(9, 0)
	// Act
	result, err := NewDayHours(Monday, open, close)
	// Assert
	assert.That(t, "result nil", result == nil, true)
	assert.That(t, "error", err, ErrCloseBeforeOpen)
}

func Test_NewClosedDay_With_Weekday_Should_ReturnClosedDay(t *testing.T) {
	// Act
	result := NewClosedDay(Sunday)
	// Assert
	assert.That(t, "weekday", result.Weekday, Sunday)
	assert.That(t, "is closed", result.IsClosed, true)
}

func Test_DayHours_With_OpenDay_Should_ReturnDuration(t *testing.T) {
	// Arrange
	open := MustTimeOfDay(9, 0)
	close := MustTimeOfDay(17, 0)
	hours, _ := NewDayHours(Monday, open, close)
	// Act
	result := hours.Duration()
	// Assert
	assert.That(t, "duration", result, 480) // 8 hours = 480 minutes
}

func Test_DayHours_With_ClosedDay_Should_ReturnZeroDuration(t *testing.T) {
	// Arrange
	hours := NewClosedDay(Monday)
	// Act
	result := hours.Duration()
	// Assert
	assert.That(t, "duration", result, 0)
}

func Test_DayHours_With_SlotDuration_Should_ReturnSlotCount(t *testing.T) {
	// Arrange
	open := MustTimeOfDay(9, 0)
	close := MustTimeOfDay(17, 0)
	hours, _ := NewDayHours(Monday, open, close)
	// Act
	result := hours.SlotCount(30)
	// Assert
	assert.That(t, "slot count", result, 16) // 480 / 30 = 16
}

func Test_DayHours_With_TimeOfDay_Should_ReturnIsOpen(t *testing.T) {
	// Arrange
	open := MustTimeOfDay(9, 0)
	close := MustTimeOfDay(17, 0)
	hours, _ := NewDayHours(Monday, open, close)
	// Act & Assert
	assert.That(t, "at open time", hours.IsOpen(MustTimeOfDay(9, 0)), true)
	assert.That(t, "during hours", hours.IsOpen(MustTimeOfDay(12, 0)), true)
	assert.That(t, "before open", hours.IsOpen(MustTimeOfDay(8, 0)), false)
	assert.That(t, "at close time", hours.IsOpen(MustTimeOfDay(17, 0)), false)
	assert.That(t, "after close", hours.IsOpen(MustTimeOfDay(18, 0)), false)
}

func Test_DayHours_With_Validation_Should_ReturnCorrectResult(t *testing.T) {
	// Arrange
	validHours, _ := NewDayHours(Monday, MustTimeOfDay(9, 0), MustTimeOfDay(17, 0))
	closedHours := NewClosedDay(Sunday)
	invalidHours := &DayHours{OpenTime: MustTimeOfDay(17, 0), CloseTime: MustTimeOfDay(9, 0)}
	// Act & Assert
	assert.That(t, "valid hours", validHours.Validate(), nil)
	assert.That(t, "closed hours", closedHours.Validate(), nil)
	assert.That(t, "invalid hours", invalidHours.Validate(), ErrCloseBeforeOpen)
}

func Test_NewOpeningHours_With_Defaults_Should_ReturnAllDaysClosed(t *testing.T) {
	// Act
	result := NewOpeningHours()
	// Assert
	assert.That(t, "days count", len(result.Days), 7)
	for i := Monday; i <= Sunday; i++ {
		assert.That(t, i.String()+" is closed", result.Days[i].IsClosed, true)
	}
}

func Test_OpeningHours_With_SetOpen_Should_UpdateDay(t *testing.T) {
	// Arrange
	hours := NewOpeningHours()
	open := MustTimeOfDay(9, 0)
	close := MustTimeOfDay(17, 0)
	// Act
	err := hours.SetOpen(Monday, open, close)
	// Assert
	assert.That(t, "error", err, nil)
	assert.That(t, "monday open", hours.Days[Monday].IsClosed, false)
	assert.That(t, "monday open time", hours.Days[Monday].OpenTime, open)
}

func Test_OpeningHours_With_SetClosed_Should_CloseDay(t *testing.T) {
	// Arrange
	hours := NewOpeningHours()
	hours.SetOpen(Monday, MustTimeOfDay(9, 0), MustTimeOfDay(17, 0))
	// Act
	hours.SetClosed(Monday)
	// Assert
	assert.That(t, "monday closed", hours.Days[Monday].IsClosed, true)
}

func Test_OpeningHours_With_Date_Should_ReturnDayHoursForDate(t *testing.T) {
	// Arrange
	hours := NewOpeningHours()
	hours.SetOpen(Monday, MustTimeOfDay(9, 0), MustTimeOfDay(17, 0))
	monday := time.Date(2025, 12, 22, 0, 0, 0, 0, time.UTC) // A Monday
	// Act
	result := hours.GetDayHoursForDate(monday)
	// Assert
	assert.That(t, "is closed", result.IsClosed, false)
	assert.That(t, "open time", result.OpenTime.String(), "09:00")
}

func Test_OpeningHours_With_Weekday_Should_ReturnIsOpenOn(t *testing.T) {
	// Arrange
	hours := NewOpeningHours()
	hours.SetOpen(Monday, MustTimeOfDay(9, 0), MustTimeOfDay(17, 0))
	// Act & Assert
	assert.That(t, "monday open", hours.IsOpenOn(Monday), true)
	assert.That(t, "tuesday closed", hours.IsOpenOn(Tuesday), false)
}
