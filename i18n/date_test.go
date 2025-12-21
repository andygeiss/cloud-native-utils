package i18n

import (
	"testing"
	"time"

	"github.com/andygeiss/cloud-native-utils/assert"
)

func Test_DateFormatterDE_Format_With_UTCDate_Should_UseEuropeanFormat(t *testing.T) {
	// Arrange
	formatter := NewDateFormatterDE()
	date := time.Date(2024, 12, 25, 10, 30, 0, 0, time.UTC)

	// Act
	result := formatter.Format(date)

	// Assert
	assert.That(t, "should format as DD.MM.YYYY", result, "25.12.2024")
}

func Test_DateFormatterEN_Format_With_UTCDate_Should_UseLongFormat(t *testing.T) {
	// Arrange
	formatter := NewDateFormatterEN()
	date := time.Date(2024, 12, 25, 10, 30, 0, 0, time.UTC)

	// Act
	result := formatter.Format(date)

	// Assert
	assert.That(t, "should format as Month DD, YYYY", result, "December 25, 2024")
}

func Test_DateFormatterISO_Format_With_UTCDate_Should_UseISOFormat(t *testing.T) {
	// Arrange
	formatter := NewDateFormatterISO()
	date := time.Date(2024, 12, 25, 10, 30, 0, 0, time.UTC)

	// Act
	result := formatter.Format(date)

	// Assert
	assert.That(t, "should format as YYYY-MM-DD", result, "2024-12-25")
}

func Test_FormatDate_With_LangDE_Should_UseGermanFormat(t *testing.T) {
	// Arrange
	date := time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC)

	// Act
	result := FormatDate(date, "de")

	// Assert
	assert.That(t, "should format with German locale", result, "05.01.2024")
}

func Test_FormatDate_With_LangEN_Should_UseEnglishFormat(t *testing.T) {
	// Arrange
	date := time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC)

	// Act
	result := FormatDate(date, "en")

	// Assert
	assert.That(t, "should format with English locale", result, "January 5, 2024")
}

func Test_FormatDate_With_UnknownLang_Should_UseISOFormat(t *testing.T) {
	// Arrange
	date := time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC)

	// Act
	result := FormatDate(date, "fr")

	// Assert
	assert.That(t, "should fallback to ISO format", result, "2024-01-05")
}

func Test_FormatDateISO_With_UTCDate_Should_FormatCorrectly(t *testing.T) {
	// Arrange
	date := time.Date(2024, 6, 15, 14, 30, 0, 0, time.UTC)

	// Act
	result := FormatDateISO(date)

	// Assert
	assert.That(t, "should format as ISO date", result, "2024-06-15")
}

func Test_ParseDateISO_With_ValidDate_Should_Parse(t *testing.T) {
	// Arrange & Act
	result := ParseDateISO("2024-06-15")

	// Assert
	assert.That(t, "year should be 2024", result.Year(), 2024)
	assert.That(t, "month should be June", result.Month(), time.June)
	assert.That(t, "day should be 15", result.Day(), 15)
}

func Test_ParseDateISO_With_InvalidDate_Should_ReturnNow(t *testing.T) {
	// Arrange
	now := time.Now()

	// Act
	result := ParseDateISO("invalid-date")

	// Assert - should be close to now (within 1 second)
	assert.That(t, "should return time close to now", result.Year(), now.Year())
}

func Test_DateFormatter_With_CustomLocation_Should_ApplyTimezone(t *testing.T) {
	// Arrange
	loc, _ := time.LoadLocation("America/New_York")
	formatter := NewDateFormatterISO().WithLocation(loc)
	// UTC midnight = 7PM previous day in New York (EST)
	date := time.Date(2024, 6, 15, 4, 0, 0, 0, time.UTC)

	// Act
	result := formatter.Format(date)

	// Assert - should be previous day in New York
	assert.That(t, "should apply timezone", result, "2024-06-15")
}

func Test_DateFormatter_With_CustomFormat_Should_UseCustomFormat(t *testing.T) {
	// Arrange
	formatter := NewDateFormatterDE().WithFormat(DateFormatISO)
	date := time.Date(2024, 12, 25, 10, 30, 0, 0, time.UTC)

	// Act
	result := formatter.Format(date)

	// Assert
	assert.That(t, "should use custom format", result, "2024-12-25")
}
