package scheduling

import (
	"testing"
	"time"

	"github.com/andygeiss/cloud-native-utils/assert"
)

func Test_NewTimeRange_With_ValidParams_Should_ReturnTimeRange(t *testing.T) {
	// Arrange
	date := time.Date(2025, 12, 21, 0, 0, 0, 0, time.UTC)
	start := MustTimeOfDay(9, 0)
	end := MustTimeOfDay(10, 0)
	// Act
	result, err := NewTimeRange(date, start, end)
	// Assert
	assert.That(t, "error", err, nil)
	assert.That(t, "start time", result.StartTime(), start)
	assert.That(t, "end time", result.EndTime(), end)
	assert.That(t, "date", result.Date().Format("2006-01-02"), "2025-12-21")
}

func Test_NewTimeRange_With_EndBeforeStart_Should_ReturnError(t *testing.T) {
	// Arrange
	date := time.Date(2025, 12, 21, 0, 0, 0, 0, time.UTC)
	start := MustTimeOfDay(10, 0)
	end := MustTimeOfDay(9, 0)
	// Act
	result, err := NewTimeRange(date, start, end)
	// Assert
	assert.That(t, "result nil", result == nil, true)
	assert.That(t, "error", err, ErrEndBeforeStart)
}

func Test_TimeRange_With_ValidRange_Should_ReturnDuration(t *testing.T) {
	// Arrange
	date := time.Date(2025, 12, 21, 0, 0, 0, 0, time.UTC)
	tr, _ := NewTimeRange(date, MustTimeOfDay(9, 0), MustTimeOfDay(11, 0))
	// Act
	result := tr.Duration()
	// Assert
	assert.That(t, "duration", result, 120)
}

func Test_TimeRange_With_ValidRange_Should_ReturnSlotCount(t *testing.T) {
	// Arrange
	date := time.Date(2025, 12, 21, 0, 0, 0, 0, time.UTC)
	tr, _ := NewTimeRange(date, MustTimeOfDay(9, 0), MustTimeOfDay(11, 0))
	// Act & Assert
	assert.That(t, "30 min slots", tr.SlotCount(30), 4)
	assert.That(t, "60 min slots", tr.SlotCount(60), 2)
	assert.That(t, "zero duration", tr.SlotCount(0), 0)
}

func Test_TimeRange_With_VariousTimes_Should_CheckContains(t *testing.T) {
	// Arrange
	date := time.Date(2025, 12, 21, 0, 0, 0, 0, time.UTC)
	tr, _ := NewTimeRange(date, MustTimeOfDay(9, 0), MustTimeOfDay(11, 0))
	// Act & Assert
	assert.That(t, "at start", tr.Contains(MustTimeOfDay(9, 0)), true)
	assert.That(t, "within", tr.Contains(MustTimeOfDay(10, 0)), true)
	assert.That(t, "at end", tr.Contains(MustTimeOfDay(11, 0)), false)
	assert.That(t, "before", tr.Contains(MustTimeOfDay(8, 0)), false)
	assert.That(t, "after", tr.Contains(MustTimeOfDay(12, 0)), false)
}

func Test_TimeRange_With_VariousRanges_Should_CheckOverlaps(t *testing.T) {
	// Arrange
	date := time.Date(2025, 12, 21, 0, 0, 0, 0, time.UTC)
	tr1, _ := NewTimeRange(date, MustTimeOfDay(9, 0), MustTimeOfDay(11, 0))
	tr2, _ := NewTimeRange(date, MustTimeOfDay(10, 0), MustTimeOfDay(12, 0)) // overlaps
	tr3, _ := NewTimeRange(date, MustTimeOfDay(11, 0), MustTimeOfDay(13, 0)) // adjacent, no overlap
	tr4, _ := NewTimeRange(date, MustTimeOfDay(7, 0), MustTimeOfDay(8, 0))   // before, no overlap
	otherDate := time.Date(2025, 12, 22, 0, 0, 0, 0, time.UTC)
	tr5, _ := NewTimeRange(otherDate, MustTimeOfDay(9, 0), MustTimeOfDay(11, 0)) // different date
	// Act & Assert
	assert.That(t, "overlaps", tr1.Overlaps(tr2), true)
	assert.That(t, "adjacent", tr1.Overlaps(tr3), false)
	assert.That(t, "before", tr1.Overlaps(tr4), false)
	assert.That(t, "different date", tr1.Overlaps(tr5), false)
}

func Test_TimeRange_With_ValidRange_Should_ReturnString(t *testing.T) {
	// Arrange
	date := time.Date(2025, 12, 21, 0, 0, 0, 0, time.UTC)
	tr, _ := NewTimeRange(date, MustTimeOfDay(9, 0), MustTimeOfDay(11, 0))
	// Act
	result := tr.String()
	// Assert
	assert.That(t, "string", result, "2025-12-21 09:00-11:00")
}
