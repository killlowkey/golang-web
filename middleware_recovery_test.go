package web

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPanicMiddleware(t *testing.T) {
	server := New()
	server.Use(recovery())
	server.GET("/", func(ctx *Context) { panic("发生 panic") })
	server.GET("/user/:name", func(ctx *Context) { panic("发生 panic") })
	server.GET("/api/*", func(ctx *Context) { panic("发生 panic") })

	testCases := []struct {
		name    string
		path    string
		wantRes map[string]any
	}{
		{
			name: "/",
			path: "/",
			wantRes: map[string]any{
				"code":  float64(500),
				"path":  "/",
				"route": "",
				"msg":   "发生 panic",
			},
		},
		{
			name: "/user/:name",
			path: "/user/ray",
			wantRes: map[string]any{
				"code":  float64(500),
				"path":  "/user/ray",
				"route": "/user/:name",
				"msg":   "发生 panic",
			},
		},
		{
			name: "/api/*",
			path: "/api/user",
			wantRes: map[string]any{
				"code":  float64(500),
				"path":  "/api/user",
				"route": "/api/*",
				"msg":   "发生 panic",
			},
		},
	}

	for _, tt := range testCases {
		request, err := http.NewRequest(http.MethodGet, tt.path, nil)
		assert.NoError(t, err)

		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)
		data, err := io.ReadAll(response.Body)
		assert.NoError(t, err)

		res := map[string]any{}
		err = json.Unmarshal(data, &res)
		assert.NoError(t, err)

		assert.Equal(t, tt.wantRes, res)
	}
}
