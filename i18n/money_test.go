package i18n

import (
	"testing"

	"github.com/andygeiss/cloud-native-utils/assert"
)

func Test_MoneyFormatterDE_With_VariousAmounts_Should_Format(t *testing.T) {
	// Arrange
	formatter := NewMoneyFormatterDE()

	tests := []struct {
		name        string
		amountCents int64
		currency    string
		expected    string
	}{
		{"zero", 0, "EUR", "0,00 EUR"},
		{"small amount", 99, "EUR", "0,99 EUR"},
		{"one euro", 100, "EUR", "1,00 EUR"},
		{"decimal", 2550, "EUR", "25,50 EUR"},
		{"hundreds", 12345, "EUR", "123,45 EUR"},
		{"thousands", 123456, "EUR", "1.234,56 EUR"},
		{"ten thousands", 1234567, "EUR", "12.345,67 EUR"},
		{"hundred thousands", 12345678, "EUR", "123.456,78 EUR"},
		{"millions", 123456789, "EUR", "1.234.567,89 EUR"},
		{"negative treated as positive", -2500, "EUR", "25,00 EUR"},
		{"no currency", 2500, "", "25,00"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := formatter.Format(tt.amountCents, tt.currency)
			// Assert
			assert.That(t, "formatted amount", result, tt.expected)
		})
	}
}

func Test_MoneyFormatterEN_With_VariousAmounts_Should_Format(t *testing.T) {
	// Arrange
	formatter := NewMoneyFormatterEN()

	tests := []struct {
		name        string
		amountCents int64
		currency    string
		expected    string
	}{
		{"zero", 0, "EUR", "0.00 EUR"},
		{"small amount", 99, "EUR", "0.99 EUR"},
		{"one euro", 100, "EUR", "1.00 EUR"},
		{"decimal", 2550, "EUR", "25.50 EUR"},
		{"hundreds", 12345, "EUR", "123.45 EUR"},
		{"thousands", 123456, "EUR", "1,234.56 EUR"},
		{"ten thousands", 1234567, "EUR", "12,345.67 EUR"},
		{"hundred thousands", 12345678, "EUR", "123,456.78 EUR"},
		{"millions", 123456789, "EUR", "1,234,567.89 EUR"},
		{"negative treated as positive", -2500, "EUR", "25.00 EUR"},
		{"no currency", 2500, "", "25.00"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := formatter.Format(tt.amountCents, tt.currency)
			// Assert
			assert.That(t, "formatted amount", result, tt.expected)
		})
	}
}

func Test_FormatMoney_With_Locales_Should_FormatByLocale(t *testing.T) {
	// Arrange
	tests := []struct {
		name        string
		amountCents int64
		currency    string
		locale      string
		expected    string
	}{
		{"german locale de", 2500, "EUR", "de", "25,00 EUR"},
		{"german locale de-DE", 2500, "EUR", "de-DE", "25,00 EUR"},
		{"german locale de_DE", 2500, "EUR", "de_DE", "25,00 EUR"},
		{"english locale en", 2500, "EUR", "en", "25.00 EUR"},
		{"unknown locale defaults to en", 2500, "EUR", "fr", "25.00 EUR"},
		{"empty locale defaults to en", 2500, "EUR", "", "25.00 EUR"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := FormatMoney(tt.amountCents, tt.currency, tt.locale)
			// Assert
			assert.That(t, "formatted amount", result, tt.expected)
		})
	}
}

func Test_MoneyFormatter_With_Amount_Should_FormatWithoutCurrency(t *testing.T) {
	// Arrange
	formatter := NewMoneyFormatterDE()
	// Act
	result := formatter.FormatWithoutCurrency(2550)
	// Assert
	assert.That(t, "formatted amount without currency", result, "25,50")
}
