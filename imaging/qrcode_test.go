package imaging

import (
	"strings"
	"testing"

	"github.com/andygeiss/cloud-native-utils/assert"
)

func Test_QRCodeGenerator_With_URLContent_Should_ReturnValidDataURL(t *testing.T) {
	// Arrange
	gen := NewQRCodeGenerator()

	// Act
	dataURL, err := gen.GenerateDataURL("https://example.com")

	// Assert
	assert.That(t, "error must be nil", err, nil)
	assert.That(t, "must start with data:image/png;base64,", strings.HasPrefix(dataURL, "data:image/png;base64,"), true)
	assert.That(t, "must have content after prefix", len(dataURL) > len("data:image/png;base64,"), true)
}

func Test_QRCodeGenerator_With_TextContent_Should_ReturnValidPNG(t *testing.T) {
	// Arrange
	gen := NewQRCodeGenerator()

	// Act
	png, err := gen.GeneratePNG("test content")

	// Assert
	assert.That(t, "error must be nil", err, nil)
	assert.That(t, "PNG data must not be empty", len(png) > 0, true)
	// PNG magic bytes: 89 50 4E 47
	assert.That(t, "must have PNG magic byte 1", png[0], byte(0x89))
	assert.That(t, "must have PNG magic byte 2", png[1], byte(0x50))
	assert.That(t, "must have PNG magic byte 3", png[2], byte(0x4E))
	assert.That(t, "must have PNG magic byte 4", png[3], byte(0x47))
}

func Test_QRCodeGenerator_With_DifferentSizes_Should_CreateLargerImage(t *testing.T) {
	// Arrange
	smallGen := NewQRCodeGenerator().WithSize(100)
	largeGen := NewQRCodeGenerator().WithSize(400)

	// Act
	smallPNG, _ := smallGen.GeneratePNG("test")
	largePNG, _ := largeGen.GeneratePNG("test")

	// Assert
	assert.That(t, "larger size should produce larger PNG", len(largePNG) > len(smallPNG), true)
}

func Test_QRCodeGenerator_With_RecoveryLevels_Should_AcceptAllLevels(t *testing.T) {
	// Arrange
	levels := []RecoveryLevel{RecoveryLow, RecoveryMedium, RecoveryHigh, RecoveryHighest}

	for _, level := range levels {
		// Act
		gen := NewQRCodeGenerator().WithRecoveryLevel(level)
		dataURL, err := gen.GenerateDataURL("test")
		// Assert
		assert.That(t, "error must be nil", err, nil)
		assert.That(t, "must return valid data URL", strings.HasPrefix(dataURL, "data:image/png;base64,"), true)
	}
}

func Test_GenerateQRCodeDataURL_With_URLContent_Should_ReturnValidDataURL(t *testing.T) {
	// Arrange & Act
	dataURL, err := GenerateQRCodeDataURL("https://example.com/booking/123")

	// Assert
	assert.That(t, "error must be nil", err, nil)
	assert.That(t, "must start with data:image/png;base64,", strings.HasPrefix(dataURL, "data:image/png;base64,"), true)
}

func Test_QRCodeGenerator_With_MethodChaining_Should_Work(t *testing.T) {
	// Arrange & Act
	gen := NewQRCodeGenerator().
		WithSize(300).
		WithRecoveryLevel(RecoveryHigh)

	dataURL, err := gen.GenerateDataURL("chained test")

	// Assert
	assert.That(t, "error must be nil", err, nil)
	assert.That(t, "must return valid data URL", strings.HasPrefix(dataURL, "data:image/png;base64,"), true)
}
