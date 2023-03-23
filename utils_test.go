package web

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_joinPaths(t *testing.T) {
	type args struct {
		absolutePath string
		relativePath string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "/",
			args: args{"/", "/"},
			want: "/",
		},
		{
			name: "/api/user",
			args: args{"/api", "/user"},
			want: "/api/user",
		},
		{
			name: "/api//detail",
			args: args{"/api//", "/detail"},
			want: "/api/detail",
		},
		{
			name: "/api/v1/user/",
			args: args{"/api", "/v1/user/"},
			want: "/api/v1/user/",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, joinPaths(tt.args.absolutePath, tt.args.relativePath), "joinPaths(%v, %v)", tt.args.absolutePath, tt.args.relativePath)
		})
	}
}
