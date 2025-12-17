package redirecting_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/redirecting"
)

func Test_Redirect_With_HTMXRequest_Should_SetHXRedirectHeader(t *testing.T) {
	// Arrange
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/source", nil)
	req.Header.Set("HX-Request", "true")

	// Act
	redirecting.Redirect(rec, req, "/target")

	// Assert
	assert.That(t, "status should be OK", rec.Code, http.StatusOK)
	assert.That(t, "HX-Redirect should be /target", rec.Header().Get("HX-Redirect"), "/target")
}

func Test_Redirect_With_StandardRequest_Should_SetLocationHeader(t *testing.T) {
	// Arrange
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/source", nil)

	// Act
	redirecting.Redirect(rec, req, "/target")

	// Assert
	assert.That(t, "status should be SeeOther", rec.Code, http.StatusSeeOther)
	assert.That(t, "Location should be /target", rec.Header().Get("Location"), "/target")
}

func Test_RedirectWithMessage_With_EmptyValue_Should_NotIncludeQuery(t *testing.T) {
	// Arrange
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/source", nil)

	// Act
	redirecting.RedirectWithMessage(rec, req, "/target", "error", "")

	// Assert
	assert.That(t, "Location should not include query", rec.Header().Get("Location"), "/target")
}

func Test_RedirectWithMessage_With_HTMXRequest_Should_SetHXRedirectWithMessage(t *testing.T) {
	// Arrange
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/source", nil)
	req.Header.Set("HX-Request", "true")

	// Act
	redirecting.RedirectWithMessage(rec, req, "/target", "success", "done")

	// Assert
	assert.That(t, "status should be OK", rec.Code, http.StatusOK)
	assert.That(t, "HX-Redirect should include message", rec.Header().Get("HX-Redirect"), "/target?success=done")
}

func Test_RedirectWithMessage_With_StandardRequest_Should_SetLocationWithMessage(t *testing.T) {
	// Arrange
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/source", nil)

	// Act
	redirecting.RedirectWithMessage(rec, req, "/target", "error", "something went wrong")

	// Assert
	assert.That(t, "status should be SeeOther", rec.Code, http.StatusSeeOther)
	assert.That(t, "Location should include message", rec.Header().Get("Location"), "/target?error=something+went+wrong")
}

func Test_WithPRG_With_AllRedirectStatusCodes_Should_ConvertToHTMXRedirect(t *testing.T) {
	// Arrange
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

func Test_WithPRG_With_HTMXRequest_Should_ConvertToHTMXRedirect(t *testing.T) {
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

func Test_WithPRG_With_NonRedirectStatus_Should_PassThrough(t *testing.T) {
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

func Test_WithPRG_With_RedirectWithoutLocation_Should_KeepOriginalStatus(t *testing.T) {
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

func Test_WithPRG_With_StandardRequest_Should_KeepOriginalRedirect(t *testing.T) {
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
