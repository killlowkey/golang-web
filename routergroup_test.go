package web

import (
	"net/http"
	"testing"
)

func TestRouterGroup_add(t *testing.T) {
	server := New()
	handle := func(ctx *Context) {}

	server.GET("/", handle)
	group := server.Group("/api")
	group.GET("/user", handle)
	group.GET("/age", handle)
	group.GET("/a/b/c", handle)
	group.GET("/:name/detail", handle)

	requests := testRequests{
		{"/", false, "", nil},
		{"/api", true, "", nil},
		{"/api/user", false, "/api/user", nil},
		{"/api/age", false, "/api/age", nil},
		{"/api/a/b/c", false, "/api/a/b/c", nil},
		{"/a/b/c", true, "", nil},
		{"/ray/detail", false, "/:name/detail", Params{
			Param{"name", "ray"},
		}},
	}

	checkRequests(t, server.trees.trees[http.MethodGet].root, requests)
}
