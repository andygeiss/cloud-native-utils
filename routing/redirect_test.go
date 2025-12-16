package routing

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/andygeiss/cloud-native-utils/assert"
)

func TestWithPRG_StandardRedirect(t *testing.T) {
	// Arrange
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/target", http.StatusSeeOther)
	})
	wrapped := WithPRG(handler)
	req := httptest.NewRequest(http.MethodPost, "/source", nil)
	rec := httptest.NewRecorder()

	// Act
	wrapped.ServeHTTP(rec, req)

	// Assert
	assert.That(t, "status should be SeeOther", rec.Code, http.StatusSeeOther)
	assert.That(t, "Location should be /target", rec.Header().Get("Location"), "/target")
	assert.That(t, "HX-Redirect should be empty", rec.Header().Get("HX-Redirect"), "")
}

func TestWithPRG_HTMXRedirect(t *testing.T) {
	// Arrange
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/target", http.StatusSeeOther)
	})
	wrapped := WithPRG(handler)
	req := httptest.NewRequest(http.MethodPost, "/source", nil)
	req.Header.Set("HX-Request", "true")
	rec := httptest.NewRecorder()

	// Act
	wrapped.ServeHTTP(rec, req)

	// Assert
	assert.That(t, "status should be OK", rec.Code, http.StatusOK)
	assert.That(t, "HX-Redirect should be /target", rec.Header().Get("HX-Redirect"), "/target")
	assert.That(t, "Location should be empty", rec.Header().Get("Location"), "")
}

func TestWithPRG_NonRedirectStatus(t *testing.T) {
	// Arrange
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	wrapped := WithPRG(handler)
	req := httptest.NewRequest(http.MethodGet, "/page", nil)
	req.Header.Set("HX-Request", "true")
	rec := httptest.NewRecorder()

	// Act
	wrapped.ServeHTTP(rec, req)

	// Assert
	assert.That(t, "status should be OK", rec.Code, http.StatusOK)
}

func TestWithPRG_AllRedirectStatusCodes(t *testing.T) {
	redirectCodes := []int{
		http.StatusMovedPermanently,
		http.StatusFound,
		http.StatusSeeOther,
		http.StatusTemporaryRedirect,
		http.StatusPermanentRedirect,
	}

	for _, code := range redirectCodes {
		t.Run(http.StatusText(code), func(t *testing.T) {
			// Arrange
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Location", "/target")
				w.WriteHeader(code)
			})
			wrapped := WithPRG(handler)
			req := httptest.NewRequest(http.MethodPost, "/source", nil)
			req.Header.Set("HX-Request", "true")
			rec := httptest.NewRecorder()

			// Act
			wrapped.ServeHTTP(rec, req)

			// Assert
			assert.That(t, "status should be OK", rec.Code, http.StatusOK)
			assert.That(t, "HX-Redirect should be /target", rec.Header().Get("HX-Redirect"), "/target")
		})
	}
}

func TestWithPRG_RedirectWithoutLocation(t *testing.T) {
	// Arrange
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusSeeOther)
	})
	wrapped := WithPRG(handler)
	req := httptest.NewRequest(http.MethodPost, "/source", nil)
	req.Header.Set("HX-Request", "true")
	rec := httptest.NewRecorder()

	// Act
	wrapped.ServeHTTP(rec, req)

	// Assert
	assert.That(t, "status should be SeeOther", rec.Code, http.StatusSeeOther)
}

func TestRedirect_Standard(t *testing.T) {
	// Arrange
	req := httptest.NewRequest(http.MethodPost, "/source", nil)
	rec := httptest.NewRecorder()

	// Act
	Redirect(rec, req, "/target")

	// Assert
	assert.That(t, "status should be SeeOther", rec.Code, http.StatusSeeOther)
	assert.That(t, "Location should be /target", rec.Header().Get("Location"), "/target")
}

func TestRedirect_HTMX(t *testing.T) {
	// Arrange
	req := httptest.NewRequest(http.MethodPost, "/source", nil)
	req.Header.Set("HX-Request", "true")
	rec := httptest.NewRecorder()

	// Act
	Redirect(rec, req, "/target")

	// Assert
	assert.That(t, "status should be OK", rec.Code, http.StatusOK)
	assert.That(t, "HX-Redirect should be /target", rec.Header().Get("HX-Redirect"), "/target")
}

func TestRedirectWithMessage_Standard(t *testing.T) {
	// Arrange
	req := httptest.NewRequest(http.MethodPost, "/source", nil)
	rec := httptest.NewRecorder()

	// Act
	RedirectWithMessage(rec, req, "/target", "msg", "hello")

	// Assert
	assert.That(t, "status should be SeeOther", rec.Code, http.StatusSeeOther)
	assert.That(t, "Location should include message", rec.Header().Get("Location"), "/target?msg=hello")
}

func TestRedirectWithMessage_HTMX(t *testing.T) {
	// Arrange
	req := httptest.NewRequest(http.MethodPost, "/source", nil)
	req.Header.Set("HX-Request", "true")
	rec := httptest.NewRecorder()

	// Act
	RedirectWithMessage(rec, req, "/target", "error", "not found")

	// Assert
	assert.That(t, "status should be OK", rec.Code, http.StatusOK)
	assert.That(t, "HX-Redirect should include message", rec.Header().Get("HX-Redirect"), "/target?error=not+found")
}

func TestRedirectWithMessage_EmptyValue(t *testing.T) {
	// Arrange
	req := httptest.NewRequest(http.MethodPost, "/source", nil)
	rec := httptest.NewRecorder()

	// Act
	RedirectWithMessage(rec, req, "/target", "msg", "")

	// Assert
	assert.That(t, "status should be SeeOther", rec.Code, http.StatusSeeOther)
	assert.That(t, "Location should have no query", rec.Header().Get("Location"), "/target")
}

func TestIsRedirect(t *testing.T) {
	tests := []struct {
		code     int
		expected bool
	}{
		{http.StatusOK, false},
		{http.StatusCreated, false},
		{http.StatusNoContent, false},
		{http.StatusMovedPermanently, true},
		{http.StatusFound, true},
		{http.StatusSeeOther, true},
		{http.StatusTemporaryRedirect, true},
		{http.StatusPermanentRedirect, true},
		{http.StatusBadRequest, false},
		{http.StatusInternalServerError, false},
	}
	for _, tt := range tests {
		t.Run(http.StatusText(tt.code), func(t *testing.T) {
			assert.That(t, "isRedirect should match expected", isRedirect(tt.code), tt.expected)
		})
	}
}

func TestWithMessage(t *testing.T) {
	tests := []struct {
		name     string
		target   string
		key      string
		value    string
		expected string
	}{
		{"simple", "/target", "msg", "hello", "/target?msg=hello"},
		{"with spaces", "/target", "msg", "hello world", "/target?msg=hello+world"},
		{"empty value", "/target", "msg", "", "/target"},
		{"special chars", "/target", "msg", "a&b=c", "/target?msg=a%26b%3Dc"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.That(t, "withMessage should build correct URL", withMessage(tt.target, tt.key, tt.value), tt.expected)
		})
	}
}
