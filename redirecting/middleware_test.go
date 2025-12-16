package redirecting_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/redirecting"
)

func TestWithPRG_StandardRedirect(t *testing.T) {
	// Arrange
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/target", http.StatusSeeOther)
	})
	wrapped := redirecting.WithPRG(handler)
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
	wrapped := redirecting.WithPRG(handler)
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
	wrapped := redirecting.WithPRG(handler)
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
			wrapped := redirecting.WithPRG(handler)
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
	wrapped := redirecting.WithPRG(handler)
	req := httptest.NewRequest(http.MethodPost, "/source", nil)
	req.Header.Set("HX-Request", "true")
	rec := httptest.NewRecorder()

	// Act
	wrapped.ServeHTTP(rec, req)

	// Assert
	assert.That(t, "status should be SeeOther", rec.Code, http.StatusSeeOther)
}
