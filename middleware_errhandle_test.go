package web

import (
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestErrorHandle(t *testing.T) {
	server := New()
	server.Use(errorHandle())

	request, err := http.NewRequest(http.MethodGet, "/", nil)
	assert.NoError(t, err)

	response := httptest.NewRecorder()
	server.ServeHTTP(response, request)
	assert.Equal(t, http.StatusNotFound, response.Code)
	data, err := io.ReadAll(response.Body)
	assert.NoError(t, err)
	assert.Equal(t, notFound, string(data))
}
