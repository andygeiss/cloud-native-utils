package scheduling

import (
	"testing"

	"github.com/andygeiss/cloud-native-utils/assert"
)

func Test_NewTimeOfDay_With_VariousInputs_Should_ReturnCorrectResult(t *testing.T) {
	// Arrange
	testCases := []struct {
		name        string
		hours       int
		minutes     int
		expectError bool
		expected    TimeOfDay
	}{
		{"valid midnight", 0, 0, false, TimeOfDay(0)},
		{"valid noon", 12, 0, false, TimeOfDay(720)},
		{"valid 23:59", 23, 59, false, TimeOfDay(1439)},
		{"invalid hours negative", -1, 0, true, 0},
		{"invalid hours > 23", 24, 0, true, 0},
		{"invalid minutes negative", 0, -1, true, 0},
		{"invalid minutes > 59", 0, 60, true, 0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			result, err := NewTimeOfDay(tc.hours, tc.minutes)
			// Assert
			if tc.expectError {
				assert.That(t, "error", err != nil, true)
			} else {
				assert.That(t, "error", err, nil)
				assert.That(t, "result", result, tc.expected)
			}
		})
	}
}

func Test_ParseTimeOfDay_With_VariousInputs_Should_ReturnCorrectResult(t *testing.T) {
	// Arrange
	testCases := []struct {
		name        string
		value       string
		expectError bool
		expected    TimeOfDay
	}{
		{"valid 09:00", "09:00", false, TimeOfDay(540)},
		{"valid 14:30", "14:30", false, TimeOfDay(870)},
		{"invalid no colon", "900", true, 0},
		{"invalid empty", "", true, 0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			result, err := ParseTimeOfDay(tc.value)
			// Assert
			if tc.expectError {
				assert.That(t, "error", err != nil, true)
			} else {
				assert.That(t, "error", err, nil)
				assert.That(t, "result", result, tc.expected)
			}
		})
	}
}

func Test_MustTimeOfDay_With_InvalidInput_Should_Panic(t *testing.T) {
	// Arrange
	defer func() {
		r := recover()
		// Assert
		assert.That(t, "recovered", r != nil, true)
	}()
	// Act
	MustTimeOfDay(25, 0) // Should panic
}

func Test_TimeOfDay_With_ValidTime_Should_ReturnHours(t *testing.T) {
	// Arrange
	tod := MustTimeOfDay(14, 30)
	// Act
	result := tod.Hours()
	// Assert
	assert.That(t, "hours", result, 14)
}

func Test_TimeOfDay_With_ValidTime_Should_ReturnMinutes(t *testing.T) {
	// Arrange
	tod := MustTimeOfDay(14, 30)
	// Act
	result := tod.Minutes()
	// Assert
	assert.That(t, "minutes", result, 30)
}

func Test_TimeOfDay_With_ValidTime_Should_ReturnString(t *testing.T) {
	// Arrange
	tod := MustTimeOfDay(9, 5)
	// Act
	result := tod.String()
	// Assert
	assert.That(t, "string", result, "09:05")
}

func Test_TimeOfDay_With_ValidTime_Should_ReturnTotalMinutes(t *testing.T) {
	// Arrange
	tod := MustTimeOfDay(2, 30)
	// Act
	result := tod.TotalMinutes()
	// Assert
	assert.That(t, "total minutes", result, 150)
}

func Test_TimeOfDay_With_TwoTimes_Should_CompareBefore(t *testing.T) {
	// Arrange
	early := MustTimeOfDay(9, 0)
	late := MustTimeOfDay(10, 0)
	// Act & Assert
	assert.That(t, "early before late", early.Before(late), true)
	assert.That(t, "late before early", late.Before(early), false)
	assert.That(t, "same time", early.Before(early), false)
}

func Test_TimeOfDay_With_TwoTimes_Should_CompareAfter(t *testing.T) {
	// Arrange
	early := MustTimeOfDay(9, 0)
	late := MustTimeOfDay(10, 0)
	// Act & Assert
	assert.That(t, "late after early", late.After(early), true)
	assert.That(t, "early after late", early.After(late), false)
	assert.That(t, "same time", early.After(early), false)
}

func Test_TimeOfDay_With_TwoTimes_Should_CompareEqual(t *testing.T) {
	// Arrange
	t1 := MustTimeOfDay(9, 30)
	t2 := MustTimeOfDay(9, 30)
	t3 := MustTimeOfDay(10, 0)
	// Act & Assert
	assert.That(t, "equal times", t1.Equal(t2), true)
	assert.That(t, "different times", t1.Equal(t3), false)
}

func Test_TimeOfDay_With_MinutesToAdd_Should_ReturnAddedTime(t *testing.T) {
	// Arrange
	tod := MustTimeOfDay(9, 30)
	// Act & Assert
	assert.That(t, "add 30 minutes", tod.Add(30).String(), "10:00")
	assert.That(t, "add negative minutes", tod.Add(-60).String(), "08:30")
	assert.That(t, "clamp to 0", MustTimeOfDay(0, 30).Add(-60).String(), "00:00")
	assert.That(t, "clamp to max", MustTimeOfDay(23, 30).Add(60).String(), "23:59")
}
