package i18n

import (
"embed"
"testing"

"github.com/andygeiss/cloud-native-utils/assert"
)

//go:embed testdata
var testEfs embed.FS

func Test_T_With_ExistingKey_Should_ReturnTranslation(t *testing.T) {
	// Arrange
	trans := NewTranslations()
	trans.Load(testEfs, "en", "testdata/en.yaml")
	// Act
	result := trans.T("en", "nav.dashboard")
	// Assert
	assert.That(t, "should return translated value", result, "Dashboard")
}

func Test_T_With_NonexistentKey_Should_ReturnKey(t *testing.T) {
	// Arrange
	trans := NewTranslations()
	trans.Load(testEfs, "en", "testdata/en.yaml")
	// Act
	result := trans.T("en", "nonexistent.key")
	// Assert
	assert.That(t, "should return key itself", result, "nonexistent.key")
}

func Test_T_With_UnknownLanguage_Should_ReturnKey(t *testing.T) {
	// Arrange
	trans := NewTranslations()
	trans.Load(testEfs, "en", "testdata/en.yaml")
	// Act
	result := trans.T("xx", "nav.dashboard")
	// Assert
	assert.That(t, "should return key for unknown language", result, "nav.dashboard")
}

func Test_T_With_GermanLanguage_Should_ReturnGermanTranslation(t *testing.T) {
	// Arrange
	trans := NewTranslations()
	trans.Load(testEfs, "de", "testdata/de.yaml")
	// Act
	result := trans.T("de", "nav.availability")
	// Assert
	assert.That(t, "should return German translation", result, "Verfügbarkeit")
}

func Test_TMap_With_MultipleKeys_Should_ReturnAllTranslations(t *testing.T) {
	// Arrange
	trans := NewTranslations()
	trans.Load(testEfs, "de", "testdata/de.yaml")
	// Act
	result := trans.TMap("de", "nav.dashboard", "nav.availability", "action.logout")
	// Assert
	assert.That(t, "nav.dashboard should be Dashboard", result["nav.dashboard"], "Dashboard")
	assert.That(t, "nav.availability should be Verfügbarkeit", result["nav.availability"], "Verfügbarkeit")
	assert.That(t, "action.logout should be Abmelden", result["action.logout"], "Abmelden")
}
