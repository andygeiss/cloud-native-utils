package i18n

import (
	"embed"
	"strings"
	"sync"

	"gopkg.in/yaml.v3"
)

// Translations holds all loaded translations keyed by language code.
type Translations struct {
	data map[string]map[string]any // lang -> nested key structure
	mu   sync.RWMutex
}

// NewTranslations creates a new Translations instance.
func NewTranslations() *Translations {
	return &Translations{
		data: make(map[string]map[string]any),
	}
}

// Load parses a YAML file from the embedded FS and registers it under the given language code.
func (a *Translations) Load(efs embed.FS, lang, path string) error {
	content, err := efs.ReadFile(path)
	if err != nil {
		return err
	}
	var data map[string]any
	if err := yaml.Unmarshal(content, &data); err != nil {
		return err
	}
	a.mu.Lock()
	a.data[lang] = data
	a.mu.Unlock()
	return nil
}

// T returns the translation for a dot-separated key in the given language.
// Falls back to the key itself if not found.
// Example: T("de", "nav.dashboard") -> "Dashboard"
func (a *Translations) T(lang, key string) string {
	a.mu.RLock()
	defer a.mu.RUnlock()
	langData, ok := a.data[lang]
	if !ok {
		return key
	}
	// Navigate the nested structure using dot-separated keys
	parts := strings.Split(key, ".")
	var current any = langData
	for _, part := range parts {
		switch v := current.(type) {
		case map[string]any:
			current, ok = v[part]
			if !ok {
				return key
			}
		default:
			return key
		}
	}
	// Return the final value as string
	if str, ok := current.(string); ok {
		return str
	}
	return key
}

// TMap returns a map of translations for a given language and list of keys.
// This is useful for passing a pre-resolved translation map to templates.
// Example: TMap("de", "nav.dashboard", "nav.availability", "action.logout")
// Returns: map[string]string{"nav.dashboard": "Dashboard", "nav.availability": "Verfügbarkeit", ...}
func (a *Translations) TMap(lang string, keys ...string) map[string]string {
	result := make(map[string]string, len(keys))
	for _, key := range keys {
		result[key] = a.T(lang, key)
	}
	return result
}
