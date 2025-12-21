package i18n

import (
	"time"
)

// DateFormat defines common date format patterns.
type DateFormat string

const (
	// DateFormatISO is the ISO 8601 format: 2006-01-02
	DateFormatISO DateFormat = "2006-01-02"
	// DateFormatEU is the European format: 02.01.2006
	DateFormatEU DateFormat = "02.01.2006"
	// DateFormatUS is the US format: 01/02/2006
	DateFormatUS DateFormat = "January 2, 2006"
	// DateFormatLong is the long format: January 2, 2006
	DateFormatLong DateFormat = "January 2, 2006"
)

// DateFormatter formats dates according to locale conventions.
type DateFormatter struct {
	// format is the date format pattern.
	format DateFormat
	// location is the timezone for formatting.
	location *time.Location
}

// NewDateFormatterDE creates a German date formatter.
// Format: 02.01.2006 (e.g., "25.12.2024")
func NewDateFormatterDE() *DateFormatter {
	loc, _ := time.LoadLocation("Europe/Berlin")
	return &DateFormatter{
		format:   DateFormatEU,
		location: loc,
	}
}

// NewDateFormatterEN creates an English (US) date formatter.
// Format: January 2, 2006 (e.g., "December 25, 2024")
func NewDateFormatterEN() *DateFormatter {
	loc, _ := time.LoadLocation("UTC")
	return &DateFormatter{
		format:   DateFormatLong,
		location: loc,
	}
}

// NewDateFormatterISO creates an ISO date formatter.
// Format: 2006-01-02 (e.g., "2024-12-25")
func NewDateFormatterISO() *DateFormatter {
	return &DateFormatter{
		format:   DateFormatISO,
		location: time.UTC,
	}
}

// WithLocation sets a custom timezone for the formatter.
func (f *DateFormatter) WithLocation(loc *time.Location) *DateFormatter {
	f.location = loc
	return f
}

// WithFormat sets a custom format pattern.
func (f *DateFormatter) WithFormat(format DateFormat) *DateFormatter {
	f.format = format
	return f
}

// Format formats a time value according to the formatter's locale settings.
func (f *DateFormatter) Format(t time.Time) string {
	if f.location != nil {
		t = t.In(f.location)
	}
	return t.Format(string(f.format))
}

// FormatDate is a convenience function to format a date with locale awareness.
// Supported languages: "de" (German), "en" (English), default is ISO format.
func FormatDate(t time.Time, lang string) string {
	switch lang {
	case "de":
		return NewDateFormatterDE().Format(t)
	case "en":
		return NewDateFormatterEN().Format(t)
	default:
		return NewDateFormatterISO().Format(t)
	}
}

// FormatDateISO formats a date in ISO 8601 format (2006-01-02).
func FormatDateISO(t time.Time) string {
	return NewDateFormatterISO().Format(t)
}

// ParseDateISO parses a date string in ISO 8601 format.
// Returns the current time if parsing fails.
func ParseDateISO(dateStr string) time.Time {
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return time.Now().UTC()
	}
	return t
}
