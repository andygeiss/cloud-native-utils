package scheduling

import (
	"testing"
	"time"

	"github.com/andygeiss/cloud-native-utils/assert"
)

func Test_OrphanGapDetector_With_NoOrphanGaps_Should_ReturnEmpty(t *testing.T) {
	// Arrange
	// Window: 09:00-17:00, Selection: 09:00-10:00 (starts at window start)
	dayHours, _ := NewDayHours(Monday, MustTimeOfDay(9, 0), MustTimeOfDay(17, 0))
	date := time.Date(2025, 12, 21, 0, 0, 0, 0, time.UTC)
	selection, _ := NewTimeRange(date, MustTimeOfDay(9, 0), MustTimeOfDay(10, 0))
	detector := NewOrphanGapDetector(dayHours, nil, 30)
	// Act
	result := detector.Detect(selection)
	// Assert
	assert.That(t, "no gaps", len(result), 0)
}

func Test_OrphanGapDetector_With_StartOfWindowGap_Should_Detect(t *testing.T) {
	// Arrange
	// Window: 09:00-17:00, Selection: 09:30-10:30 (leaves 30 min orphan at start)
	dayHours, _ := NewDayHours(Monday, MustTimeOfDay(9, 0), MustTimeOfDay(17, 0))
	date := time.Date(2025, 12, 21, 0, 0, 0, 0, time.UTC)
	selection, _ := NewTimeRange(date, MustTimeOfDay(9, 30), MustTimeOfDay(10, 30))
	detector := NewOrphanGapDetector(dayHours, nil, 30)
	// Act
	result := detector.Detect(selection)
	// Assert
	assert.That(t, "has gaps", result != nil, true)
	assert.That(t, "gap count", len(result), 1)
	assert.That(t, "gap reason", result[0].Reason, OrphanGapReasonStartOfWindow)
	assert.That(t, "gap start", result[0].StartTime.String(), "09:00")
	assert.That(t, "gap end", result[0].EndTime.String(), "09:30")
}

func Test_OrphanGapDetector_With_EndOfWindowGap_Should_Detect(t *testing.T) {
	// Arrange
	// Window: 09:00-17:00, Selection: 16:00-16:30 (leaves 30 min orphan at end)
	dayHours, _ := NewDayHours(Monday, MustTimeOfDay(9, 0), MustTimeOfDay(17, 0))
	date := time.Date(2025, 12, 21, 0, 0, 0, 0, time.UTC)
	selection, _ := NewTimeRange(date, MustTimeOfDay(16, 0), MustTimeOfDay(16, 30))
	detector := NewOrphanGapDetector(dayHours, nil, 30)
	// Act
	result := detector.Detect(selection)
	// Assert
	assert.That(t, "has gaps", result != nil, true)
	assert.That(t, "gap count", len(result), 1)
	assert.That(t, "gap reason", result[0].Reason, OrphanGapReasonEndOfWindow)
	assert.That(t, "gap start", result[0].StartTime.String(), "16:30")
	assert.That(t, "gap end", result[0].EndTime.String(), "17:00")
}

func Test_OrphanGapDetector_With_GapBetweenBookings_Should_Detect(t *testing.T) {
	// Arrange
	// Window: 09:00-17:00, Existing booking: 09:00-10:00
	// Selection: 10:30-11:00 (leaves 30 min orphan between booking end and selection start)
	dayHours, _ := NewDayHours(Monday, MustTimeOfDay(9, 0), MustTimeOfDay(17, 0))
	date := time.Date(2025, 12, 21, 0, 0, 0, 0, time.UTC)
	existingBooking, _ := NewTimeRange(date, MustTimeOfDay(9, 0), MustTimeOfDay(10, 0))
	selection, _ := NewTimeRange(date, MustTimeOfDay(10, 30), MustTimeOfDay(11, 0))
	detector := NewOrphanGapDetector(dayHours, []*TimeRange{existingBooking}, 30)
	// Act
	result := detector.Detect(selection)
	// Assert
	assert.That(t, "has gaps", result != nil, true)
	assert.That(t, "gap reason", result[0].Reason, OrphanGapReasonBetweenBookings)
	assert.That(t, "gap start", result[0].StartTime.String(), "10:00")
	assert.That(t, "gap end", result[0].EndTime.String(), "10:30")
}

func Test_OrphanGapDetector_With_NilSelection_Should_ReturnNil(t *testing.T) {
	// Arrange
	dayHours, _ := NewDayHours(Monday, MustTimeOfDay(9, 0), MustTimeOfDay(17, 0))
	detector := NewOrphanGapDetector(dayHours, nil, 30)
	// Act
	result := detector.Detect(nil)
	// Assert
	assert.That(t, "no gaps", result == nil, true)
}

func Test_OrphanGapDetector_With_ClosedDay_Should_ReturnEmpty(t *testing.T) {
	// Arrange
	dayHours := NewClosedDay(Monday)
	date := time.Date(2025, 12, 21, 0, 0, 0, 0, time.UTC)
	selection, _ := NewTimeRange(date, MustTimeOfDay(9, 0), MustTimeOfDay(10, 0))
	detector := NewOrphanGapDetector(dayHours, nil, 30)
	// Act
	result := detector.Detect(selection)
	// Assert
	assert.That(t, "no gaps", len(result), 0)
}

func Test_OrphanGapDetector_With_Selection_Should_ReturnHasOrphanGaps(t *testing.T) {
	// Arrange
	dayHours, _ := NewDayHours(Monday, MustTimeOfDay(9, 0), MustTimeOfDay(17, 0))
	date := time.Date(2025, 12, 21, 0, 0, 0, 0, time.UTC)
	selectionWithGap, _ := NewTimeRange(date, MustTimeOfDay(9, 30), MustTimeOfDay(10, 30))
	selectionNoGap, _ := NewTimeRange(date, MustTimeOfDay(9, 0), MustTimeOfDay(10, 0))
	detector := NewOrphanGapDetector(dayHours, nil, 30)
	// Act & Assert
	assert.That(t, "has gaps", detector.HasOrphanGaps(selectionWithGap), true)
	assert.That(t, "no gaps", detector.HasOrphanGaps(selectionNoGap), false)
}

func Test_OrphanGapDetector_With_OrphanGap_Should_ReturnFirstOrphanGap(t *testing.T) {
	// Arrange
	dayHours, _ := NewDayHours(Monday, MustTimeOfDay(9, 0), MustTimeOfDay(17, 0))
	date := time.Date(2025, 12, 21, 0, 0, 0, 0, time.UTC)
	selection, _ := NewTimeRange(date, MustTimeOfDay(9, 30), MustTimeOfDay(10, 30))
	detector := NewOrphanGapDetector(dayHours, nil, 30)
	// Act
	result := detector.GetFirstOrphanGap(selection)
	// Assert
	assert.That(t, "has gap", result != nil, true)
	assert.That(t, "gap reason", result.Reason, OrphanGapReasonStartOfWindow)
}

func Test_OrphanGapDetector_With_NoOrphanGap_Should_ReturnNilFirstOrphanGap(t *testing.T) {
	// Arrange
	dayHours, _ := NewDayHours(Monday, MustTimeOfDay(9, 0), MustTimeOfDay(17, 0))
	date := time.Date(2025, 12, 21, 0, 0, 0, 0, time.UTC)
	selection, _ := NewTimeRange(date, MustTimeOfDay(9, 0), MustTimeOfDay(10, 0))
	detector := NewOrphanGapDetector(dayHours, nil, 30)
	// Act
	result := detector.GetFirstOrphanGap(selection)
	// Assert
	assert.That(t, "no gap", result == nil, true)
}
