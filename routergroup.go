package web

import (
	"net/http"
)

var (
	anyMethods = []string{
		http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch,
		http.MethodHead, http.MethodOptions, http.MethodDelete, http.MethodConnect,
		http.MethodTrace,
	}
)

type IRoutes interface {
	// Use register middleware
	Use(middlewares ...Middleware)

	// Handle register router
	Handle(method, path string, handler HandleFunc) IRoutes
}

type IRouter interface {
	IRoutes

	// Group create a new RouterGroup with prefix and handleFuncs
	Group(prefix string, middlewares ...Middleware) *RouterGroup
}

var _ IRouter = (*RouterGroup)(nil)

type RouterGroup struct {
	middlewares []Middleware
	handler     HandleFunc // 业务处理器
	basePath    string     // group prefix
	server      *HttpServer
}

func (r *RouterGroup) Use(middlewares ...Middleware) {
	r.middlewares = append(r.middlewares, middlewares...)
}

func (r *RouterGroup) Handle(method, path string, handler HandleFunc) IRoutes {
	if handler == nil {
		panic("web：Handle 处理路由为 nil")
	}

	return r.register(method, path, handler)
}

func (r *RouterGroup) register(method, path string, handler HandleFunc) IRoutes {
	// 计算路径
	absolutePath := r.calculateAbsolutePath(path)
	// 注册路由
	r.server.addRouter(method, absolutePath, r.middlewares, handler)
	return r
}

func (r *RouterGroup) calculateAbsolutePath(relativePath string) string {
	return joinPaths(r.basePath, relativePath)
}

func (r *RouterGroup) Any(path string, handler HandleFunc) IRoutes {
	for _, method := range anyMethods {
		r.Handle(method, path, handler)
	}
	return r
}

func (r *RouterGroup) GET(path string, handler HandleFunc) IRoutes {
	return r.Handle(http.MethodGet, path, handler)
}

func (r *RouterGroup) POST(path string, handler HandleFunc) IRoutes {
	return r.Handle(http.MethodPost, path, handler)
}

func (r *RouterGroup) DELETE(path string, handler HandleFunc) IRoutes {
	return r.Handle(http.MethodDelete, path, handler)
}

func (r *RouterGroup) PUT(path string, handler HandleFunc) IRoutes {
	return r.Handle(http.MethodPut, path, handler)
}

func (r *RouterGroup) OPTIONS(path string, handler HandleFunc) IRoutes {
	return r.Handle(http.MethodOptions, path, handler)
}

func (r *RouterGroup) PATCH(path string, handler HandleFunc) IRoutes {
	return r.Handle(http.MethodPatch, path, handler)
}

func (r *RouterGroup) HEAD(path string, handler HandleFunc) IRoutes {
	return r.Handle(http.MethodHead, path, handler)
}

func (r *RouterGroup) TRACE(path string, handler HandleFunc) IRoutes {
	return r.Handle(http.MethodTrace, path, handler)
}

func (r *RouterGroup) CONNECT(path string, handler HandleFunc) IRoutes {
	return r.Handle(http.MethodConnect, path, handler)
}

func (r *RouterGroup) Group(prefix string, middlewares ...Middleware) *RouterGroup {
	return &RouterGroup{
		middlewares: append(r.middlewares, middlewares...),
		basePath:    joinPaths(r.basePath, prefix),
		server:      r.server,
	}
}
