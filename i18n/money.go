package i18n

// MoneyFormatter formats monetary amounts for different locales.
// Amounts are expected in the smallest currency unit (e.g., cents for EUR/USD).
type MoneyFormatter struct {
	// DecimalSeparator separates the whole and fractional parts (e.g., "," for DE, "." for EN).
	DecimalSeparator string
	// ThousandSeparator groups digits (e.g., "." for DE, "," for EN).
	ThousandSeparator string
	// CurrencyPosition determines where the currency symbol appears.
	// "suffix" places it after the amount (e.g., "25,00 EUR"), "prefix" before (e.g., "EUR 25.00").
	CurrencyPosition string
}

// NewMoneyFormatterDE creates a formatter for German locale.
// Format: 1.234,56 EUR
func NewMoneyFormatterDE() *MoneyFormatter {
	return &MoneyFormatter{
		DecimalSeparator:  ",",
		ThousandSeparator: ".",
		CurrencyPosition:  "suffix",
	}
}

// NewMoneyFormatterEN creates a formatter for English locale.
// Format: 1,234.56 EUR
func NewMoneyFormatterEN() *MoneyFormatter {
	return &MoneyFormatter{
		DecimalSeparator:  ".",
		ThousandSeparator: ",",
		CurrencyPosition:  "suffix",
	}
}

// Format formats an amount in the smallest currency unit (e.g., cents).
// If currency is empty, no currency suffix/prefix is added.
func (a *MoneyFormatter) Format(amountCents int64, currency string) string {
	if amountCents < 0 {
		amountCents = -amountCents
	}

	whole := amountCents / 100
	frac := amountCents % 100

	// Format whole part with thousand separators
	wholeStr := formatWholeWithThousands(whole, a.ThousandSeparator)

	// Format fractional part (always 2 digits)
	fracStr := padLeft(int(frac), 2)

	result := wholeStr + a.DecimalSeparator + fracStr

	// Add currency if provided
	if currency != "" {
		if a.CurrencyPosition == "prefix" {
			result = currency + " " + result
		} else {
			result = result + " " + currency
		}
	}

	return result
}

// FormatWithoutCurrency formats an amount without the currency symbol.
func (a *MoneyFormatter) FormatWithoutCurrency(amountCents int64) string {
	return a.Format(amountCents, "")
}

// formatWholeWithThousands formats the whole part with thousand separators.
func formatWholeWithThousands(n int64, sep string) string {
	if n == 0 {
		return "0"
	}

	// Build digits in reverse
	digits := make([]byte, 0, 20)
	count := 0
	for n > 0 {
		if count > 0 && count%3 == 0 && sep != "" {
			digits = append(digits, sep...)
		}
		digits = append(digits, byte('0'+n%10))
		n /= 10
		count++
	}

	// Reverse
	for i, j := 0, len(digits)-1; i < j; i, j = i+1, j-1 {
		digits[i], digits[j] = digits[j], digits[i]
	}

	return string(digits)
}

// padLeft pads a number with leading zeros to the specified width.
func padLeft(n, width int) string {
	s := intToString(n)
	for len(s) < width {
		s = "0" + s
	}
	return s
}

// intToString converts an integer to a string without using strconv.
func intToString(n int) string {
	if n == 0 {
		return "0"
	}
	result := ""
	for n > 0 {
		result = string(rune('0'+n%10)) + result
		n /= 10
	}
	return result
}

// FormatMoney is a convenience function that formats an amount using the specified locale.
// Supported locales: "de", "en". Defaults to "en" for unknown locales.
func FormatMoney(amountCents int64, currency, locale string) string {
	var formatter *MoneyFormatter
	switch locale {
	case "de", "de-DE", "de_DE":
		formatter = NewMoneyFormatterDE()
	default:
		formatter = NewMoneyFormatterEN()
	}
	return formatter.Format(amountCents, currency)
}
