package web

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type testRequests []struct {
	path       string
	nilHandler bool
	route      string
	ps         Params
}

func TestTreeAddAndGet(t *testing.T) {
	n := &node{}
	handler := func(ctx *Context) {}

	n.addRoute("/a", nil, handler)
	n.addRoute("/a/b", nil, handler)
	n.addRoute("/a/b/c", nil, handler)
	n.addRoute("/user/:name", nil, handler)
	n.addRoute("/user/:name/:id", nil, handler)
	n.addRoute("/api/*", nil, handler)
	n.addRoute("/1", nil, handler)
	n.addRoute("/v1/*/user", nil, handler)
	n.addRoute("/v2/:name/detail", nil, handler)

	checkRequests(t, n, testRequests{
		{"/", true, "", nil},
		{"/a", false, "/a", nil},
		{"/a/b", false, "/a/b", nil},
		{"/a/b/c", false, "/a/b/c", nil},
		{"/user/ray", false, "/user/:name", Params{
			Param{"name", "ray"},
		}},
		{"/user/ray/10", false, "/user/:name/:id", Params{
			Param{"name", "ray"},
			Param{"id", "10"},
		}},
		{"/api/user", false, "/api/*", nil},
		{"/1", false, "/1", nil},
		{"/v1/ray", true, "", nil},
		{"/v1/ray/user", false, "/v1/*/user", nil},
		{"/v2/ray/detail", false, "/v2/:name/detail", Params{
			Param{"name", "ray"},
		}},
	})
}

func TestTreePanic(t *testing.T) {
	n := &node{}
	handler := func(ctx *Context) {}

	// 路径为空
	assert.PanicsWithValue(t, "web：path 路径不允许为空", func() {
		n.addRoute("", nil, handler)
	})

	// 路径不以 / 开头
	assert.PanicsWithValue(t, "web：path 必须以 / 开头", func() {
		n.addRoute("a", nil, handler)
	})

	// 路径以 / 结尾
	assert.PanicsWithValue(t, "web：path 路径不允许为 / 结尾", func() {
		n.addRoute("/a/", nil, handler)
	})

	// 重复注册路由
	assert.PanicsWithValue(t, "web：不允许重复注册路由[/]", func() {
		n.addRoute("/", nil, handler)
		n.addRoute("/", nil, handler)
	})
	assert.PanicsWithValue(t, "web：不允许重复注册路由[/abc]", func() {
		n.addRoute("/abc", nil, handler)
		n.addRoute("/abc", nil, handler)
	})

	// 通配符路由与参数路由冲突
	assert.PanicsWithValue(t, "web：通配符与参数匹配冲突，只允许存在一个", func() {
		n.addRoute("/api/*/user", nil, handler)
		n.addRoute("/api/:name/user", nil, handler)
	})

	// 无效路径
	assert.PanicsWithValue(t, "web：拒绝 /a//b/c 形式的路由", func() {
		n.addRoute("/a//b/c", nil, handler)
	})
}

func checkRequests(t *testing.T, tree *node, requests testRequests) {
	for _, r := range requests {
		info, ok := tree.findRoute(r.path)
		assert.Equal(t, r.nilHandler, !ok)
		if !ok {
			return
		}
		assert.Equal(t, r.nilHandler, info.handler == nil)
		assert.Equal(t, r.ps, info.params)
		assert.Equal(t, r.route, info.fullPath)
	}
}

func TestParamGet(t *testing.T) {
	params := Params{
		Param{"name", "ray"},
		Param{"age", "10"},
		Param{"sex", "1"},
	}

	testCases := []struct {
		name    string
		key     string
		find    bool
		wantVal string
	}{
		{
			name: "nil",
			key:  "",
			find: false,
		},
		{
			name:    "name",
			key:     "name",
			find:    true,
			wantVal: "ray",
		},
		{
			name:    "age",
			key:     "age",
			find:    true,
			wantVal: "10",
		},
	}

	for _, tt := range testCases {
		got, ok := params.Get(tt.key)
		assert.Equal(t, tt.find, ok)
		if !ok {
			return
		}
		assert.Equal(t, tt.wantVal, got)
	}
}
