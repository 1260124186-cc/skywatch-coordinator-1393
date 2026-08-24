package httpapi_test

import (
	"github.com/1260124186-cc/skywatch-coordinator/internal/httpapi"
	"github.com/1260124186-cc/skywatch-coordinator/internal/service"
	"github.com/1260124186-cc/skywatch-coordinator/internal/store"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealth(t *testing.T) {
	c := service.NewCoordinator(store.NewMemoryStore())
	server := httpapi.NewServer(c)
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	server.Routes().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("got %d", rec.Code)
	}
}
