package imaging

import (
	"encoding/base64"

	"github.com/skip2/go-qrcode"
)

// RecoveryLevel defines the error correction level for QR codes.
type RecoveryLevel int

const (
	// RecoveryLow recovers from ~7% data loss.
	RecoveryLow RecoveryLevel = iota
	// RecoveryMedium recovers from ~15% data loss.
	RecoveryMedium
	// RecoveryHigh recovers from ~25% data loss.
	RecoveryHigh
	// RecoveryHighest recovers from ~30% data loss.
	RecoveryHighest
)

// QRCodeGenerator generates QR codes as data URLs or raw bytes.
type QRCodeGenerator struct {
	// size is the width/height of the QR code in pixels.
	size int
	// level is the error recovery level.
	level RecoveryLevel
}

// NewQRCodeGenerator creates a new QR code generator with default settings.
// Default size is 200x200 pixels with medium recovery level.
func NewQRCodeGenerator() *QRCodeGenerator {
	return &QRCodeGenerator{
		size:  200,
		level: RecoveryMedium,
	}
}

// WithSize sets the QR code size in pixels.
func (g *QRCodeGenerator) WithSize(size int) *QRCodeGenerator {
	g.size = size
	return g
}

// WithRecoveryLevel sets the error recovery level.
func (g *QRCodeGenerator) WithRecoveryLevel(level RecoveryLevel) *QRCodeGenerator {
	g.level = level
	return g
}

// toQRCodeLevel converts our RecoveryLevel to the library's type.
func (g *QRCodeGenerator) toQRCodeLevel() qrcode.RecoveryLevel {
	switch g.level {
	case RecoveryLow:
		return qrcode.Low
	case RecoveryMedium:
		return qrcode.Medium
	case RecoveryHigh:
		return qrcode.High
	case RecoveryHighest:
		return qrcode.Highest
	default:
		return qrcode.Medium
	}
}

// GenerateDataURL generates a QR code as a base64-encoded PNG data URL.
// The returned string can be used directly in HTML img src attributes.
// Example output: "data:image/png;base64,iVBORw0KGgoAAAANS..."
func (g *QRCodeGenerator) GenerateDataURL(content string) (string, error) {
	png, err := g.GeneratePNG(content)
	if err != nil {
		return "", err
	}
	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(png), nil
}

// GeneratePNG generates a QR code as raw PNG bytes.
func (g *QRCodeGenerator) GeneratePNG(content string) ([]byte, error) {
	return qrcode.Encode(content, g.toQRCodeLevel(), g.size)
}

// GenerateQRCodeDataURL is a convenience function to generate a QR code data URL
// with default settings (200x200, medium recovery).
func GenerateQRCodeDataURL(content string) (string, error) {
	return NewQRCodeGenerator().GenerateDataURL(content)
}
