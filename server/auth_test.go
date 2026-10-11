package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-pkgz/auth/v2/token"
	"github.com/stretchr/testify/assert"
)

func Test_server_userSessionCtx(t *testing.T) {
	s := &server{}
	var got userSession
	handler := s.userSessionCtx(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got = s.contextUserSession(r)
	}))

	// a token that never got a database id is rejected, not treated as user 0
	noID := newTokenUser(&token.User{Name: "newbie"})
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, token.SetUserInfo(httptest.NewRequest(http.MethodGet, "/", nil), *noID.User))
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.Equal(t, userSession{}, got)

	withID := newTokenUser(&token.User{Name: "player"})
	withID.setDatabaseID(12)
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, token.SetUserInfo(httptest.NewRequest(http.MethodGet, "/", nil), *withID.User))
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, int64(12), got.ID)
}
