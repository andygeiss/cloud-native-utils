package scheduling

import (
	"testing"
	"time"

	"github.com/andygeiss/cloud-native-utils/assert"
)

func Test_Weekday_With_ValidValue_Should_ReturnString(t *testing.T) {
	// Arrange
	testCases := []struct {
		weekday  Weekday
		expected string
	}{
		{Monday, "Monday"},
		{Tuesday, "Tuesday"},
		{Wednesday, "Wednesday"},
		{Thursday, "Thursday"},
		{Friday, "Friday"},
		{Saturday, "Saturday"},
		{Sunday, "Sunday"},
		{Weekday(99), "Unknown"},
	}

	for _, tc := range testCases {
		t.Run(tc.expected, func(t *testing.T) {
			// Act
			result := tc.weekday.String()
			// Assert
			assert.That(t, "string", result, tc.expected)
		})
	}
}

func Test_Weekday_With_ValidValue_Should_ReturnGoWeekday(t *testing.T) {
	// Arrange
	testCases := []struct {
		weekday  Weekday
		expected time.Weekday
	}{
		{Monday, time.Monday},
		{Tuesday, time.Tuesday},
		{Wednesday, time.Wednesday},
		{Thursday, time.Thursday},
		{Friday, time.Friday},
		{Saturday, time.Saturday},
		{Sunday, time.Sunday},
	}

	for _, tc := range testCases {
		t.Run(tc.weekday.String(), func(t *testing.T) {
			// Act
			result := tc.weekday.GoWeekday()
			// Assert
			assert.That(t, "go weekday", result, tc.expected)
		})
	}
}

func Test_WeekdayFromGoWeekday_With_ValidValue_Should_ReturnWeekday(t *testing.T) {
	// Arrange
	testCases := []struct {
		goWeekday time.Weekday
		expected  Weekday
	}{
		{time.Monday, Monday},
		{time.Tuesday, Tuesday},
		{time.Wednesday, Wednesday},
		{time.Thursday, Thursday},
		{time.Friday, Friday},
		{time.Saturday, Saturday},
		{time.Sunday, Sunday},
	}

	for _, tc := range testCases {
		t.Run(tc.expected.String(), func(t *testing.T) {
			// Act
			result := WeekdayFromGoWeekday(tc.goWeekday)
			// Assert
			assert.That(t, "weekday", result, tc.expected)
		})
	}
}
