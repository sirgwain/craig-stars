package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sirgwain/craig-stars/version"
	"github.com/stretchr/testify/assert"
)

func TestVersionHeaderOnEarlyError(t *testing.T) {
	handler := withVersionHeader(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
	}))
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest(http.MethodPost, "/api/grpc/test", nil))
	assert.Equal(t, version.Semver, response.Header().Get("X-App-Version"))
}
