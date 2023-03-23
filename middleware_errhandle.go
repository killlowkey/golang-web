package web

import "net/http"

const (
	notFound   = "404 not found"
	notAllowed = "405 not allowed"
)

type ErrHandleMiddleWareBuild struct {
	handlers map[int]HandleFunc
}

func NewErrHandleMiddleWareBuild() *ErrHandleMiddleWareBuild {
	return &ErrHandleMiddleWareBuild{
		handlers: make(map[int]HandleFunc, 64),
	}
}

func (e *ErrHandleMiddleWareBuild) Register(status int, handle HandleFunc) {
	e.handlers[status] = handle
}

func (e *ErrHandleMiddleWareBuild) Build() Middleware {
	return func(next HandleFunc) HandleFunc {
		return func(ctx *Context) {
			defer func() {
				if handler, ok := e.handlers[ctx.RespStatus]; ok {
					handler(ctx)
				}
			}()
			next(ctx)
		}
	}
}

func errorHandle() Middleware {
	build := NewErrHandleMiddleWareBuild()
	build.Register(http.StatusNotFound, func(ctx *Context) {
		ctx.String(http.StatusNotFound, notFound)
	})
	build.Register(http.StatusMethodNotAllowed, func(ctx *Context) {
		ctx.String(http.StatusMethodNotAllowed, notAllowed)
	})
	return build.Build()
}
