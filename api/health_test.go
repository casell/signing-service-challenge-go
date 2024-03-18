package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestError(t *testing.T) {
	s := &Server{}
	rec := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, "/api/v0/health", nil)
	assert.Nil(t, err)
	s.Health(rec, req)
	assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
}

func TestOK(t *testing.T) {
	s := &Server{}
	rec := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/api/v0/health", nil)
	assert.Nil(t, err)
	s.Health(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
}
