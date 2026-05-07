package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kiing-dom/api-rate-limiter/internal/rate_limiter"
)

type mockStore struct{ allowed bool }
type mockLimiter struct{ allowed bool }

func (m *mockLimiter) Allow(_ string) bool { return m.allowed }
func (m *mockStore) GetRateLimiter(_ string) rate_limiter.RateLimiter {
	return &mockLimiter{allowed: m.allowed}
}
func TestHTTPHandler_Allows(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-API-KEY", "test-key")
	w := httptest.NewRecorder()

	RateLimitHandler(&mockStore{allowed: true}).ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, but got %d", w.Code)
	}
}

func TestHTTPHandler_Rejects(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-API-KEY", "test-key")
	w := httptest.NewRecorder()

	RateLimitHandler(&mockStore{allowed: false}).ServeHTTP(w, req)

	if w.Code != http.StatusTooManyRequests {
		t.Fatalf("expected 429 Too Many Requests, but got %d", w.Code)
	}
}
