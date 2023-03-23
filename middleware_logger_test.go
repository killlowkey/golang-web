package web

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type LoggerTest struct {
	body *bytes.Buffer
}

func (c LoggerTest) Write(l *LogEntry) error {
	msg := fmt.Sprintf("%s-%s-%d", l.Method, l.Path, l.RespStatus)
	_, _ = c.body.Write([]byte(msg))
	return nil
}

func TestAccesslog(t *testing.T) {
	server := New()
	// 写入到内存，进行数据校验
	buf := &bytes.Buffer{}
	build := NewLogMiddlewareBuild(LoggerTest{buf})
	server.Use(build.Build())

	server.GET("/user", func(ctx *Context) {
		ctx.Status(http.StatusCreated)
	})
	server.GET("/api/user", func(ctx *Context) {
		ctx.String(200, "hello")
	})

	testCases := []struct {
		name    string
		path    string
		wantRes string
	}{
		{
			name:    "/user",
			path:    "/user",
			wantRes: "GET-/user-201",
		},
		{
			name:    "/api/user",
			path:    "/api/user",
			wantRes: "GET-/api/user-200",
		},
	}

	for _, tt := range testCases {
		request, _ := http.NewRequest(http.MethodGet, tt.path, nil)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		data, err := io.ReadAll(buf)
		assert.NoError(t, err)
		assert.Equal(t, tt.wantRes, string(data))
	}
}
