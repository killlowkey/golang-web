package web

import (
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"testing"
)

type User struct {
	Name string `json:"Name"`
	Age  int    `json:"Age"`
}

func TestStartWeb(t *testing.T) {
	server := Default()
	server.GET("/hello", func(ctx *Context) {
		ctx.Header("Cache-Control", "max-age=600")
		ctx.String(200, "<h1>Hello world</h1>")
	})
	server.GET("/user/:name", func(ctx *Context) {
		ctx.String(200, "hello, "+ctx.Param("name"))
	})
	server.POST("/user", func(ctx *Context) {
		user := &User{}
		if err := ctx.BindJSON(user); err != nil {
			ctx.String(http.StatusBadRequest, err.Error())
			return
		}

		ctx.JSON(200, user)
	})

	// create a group with middleware
	group := server.Group("/api", func(next HandleFunc) HandleFunc {
		return func(ctx *Context) {
			log.Println("I am /api middleware")
			next(ctx)
		}
	})
	group.GET("/user/:name", func(ctx *Context) {
		ctx.JSON(200, H{
			"code": 200,
			"msg":  "hello world",
		})
	})

	group2 := group.Group("/v1", func(next HandleFunc) HandleFunc {
		return func(ctx *Context) {
			log.Println("I am /api/v1 middleware")
			next(ctx)
		}
	})

	group2.GET("/:name", func(ctx *Context) {
		ctx.JSON(200, H{
			"code": 200,
			"msg":  "hello world",
		})
	})

	_ = server.Run(":8080")
}

func TestWebWithInvalidMethod(t *testing.T) {
	testCases := []struct {
		name string
		arg  string
	}{
		{
			name: "n",
			arg:  "n",
		},
		{
			name: "test",
			arg:  "test",
		},
		{
			name: "get",
			arg:  "get",
		},
		{
			name: "post",
			arg:  "post",
		},
	}

	server := New()
	for _, tt := range testCases {
		assert.PanicsWithValue(t, "web：无效 HTTP Method", func() {
			server.addRouter(tt.arg, "/", nil, nil)
		})
	}
}
