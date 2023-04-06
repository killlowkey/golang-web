package web

import (
	"html/template"
	"log"
	"net/http"
	"sync"
)

// HandleFunc 请求处理器
type HandleFunc func(ctx *Context)

// 校验 HttpServer 是否符合 IRouter 接口
var _ IRouter = (*HttpServer)(nil)

type HttpServer struct {
	*RouterGroup

	// 路由树
	trees *Trees
	// context pool
	pool sync.Pool

	middlewares []Middleware

	templ *template.Template
}

func New() *HttpServer {
	s := &HttpServer{
		RouterGroup: &RouterGroup{
			middlewares: nil,
			handler:     nil,
			basePath:    "/",
		},
		trees: &Trees{
			trees: make(map[string]*Tree, 9),
		},
	}
	s.RouterGroup.server = s
	s.pool.New = func() any {
		return &Context{}
	}
	return s
}

// Default 使用默认的 Middleware 创建 HttpServer
func Default() *HttpServer {
	server := New()
	server.Use(logger(), errorHandle(), recovery())
	return server
}

// addRouter 添加路由
func (h *HttpServer) addRouter(method, path string, middlewares []Middleware, handler HandleFunc) {
	if !validHttpMethod(method) {
		panic("web：无效 HTTP Method")
	}
	h.trees.addRouter(method, path, middlewares, handler)
}

// ServeHTTP 处理 HTTP 请求，作为请求全局入口点
func (h *HttpServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	ctx := h.pool.Get().(*Context)
	ctx.reset()
	ctx.Request = request
	ctx.Writer = writer
	ctx.h = h

	// 回写响应 Middleware
	flush := func(next HandleFunc) HandleFunc {
		return func(ctx *Context) {
			defer func() {
				ctx.Writer.WriteHeader(ctx.RespStatus)
				if _, err := ctx.Writer.Write(ctx.RespData); err != nil {
					// TODO 将就用，打条日志出来就完事了
					log.Println("web：回写响应失败")
				}
			}()
			next(ctx)
		}
	}

	// 组合全局 middleware
	root := h.handleHttpRequest
	middlewares := append([]Middleware{flush}, h.middlewares...)
	for i := len(middlewares) - 1; i >= 0; i-- {
		root = middlewares[i](root)
	}
	root(ctx)

	// 归还 context
	h.pool.Put(ctx)
}

// handleHttpRequest 处理客户端请求
//  1. 从路由树中查找请求路由，未找到返回 404
//  2. 初始化 context
//  3. 处理器外部封装一层回写响应的 middleware
func (h *HttpServer) handleHttpRequest(c *Context) {
	// 校验请求方法
	if !validHttpMethod(c.Request.Method) {
		c.Status(http.StatusMethodNotAllowed)
		return
	}

	// 查找请求路由
	n, ok := h.trees.findRoute(c.Request)
	if !ok {
		c.Status(http.StatusNotFound)
		return
	}

	// 初始化 context
	c.Params = n.params
	c.handler = n.handler
	c.Route = n.fullPath

	// 组合局部 Middleware，并进行调用
	root := n.handler
	for i := len(n.middlewares) - 1; i >= 0; i-- {
		root = n.middlewares[i](root)
	}
	root(c)
}

// Use 注册 Middleware
func (h *HttpServer) Use(middlewares ...Middleware) {
	h.middlewares = append(h.middlewares, middlewares...)
}

// Run 根据传入的地址启动 HttpServer
func (h *HttpServer) Run(addr string) error {
	return http.ListenAndServe(addr, h)
}

// RunTLS 启动 HTTPS 服务
func (h *HttpServer) RunTLS(addr, certFile, keyFile string) error {
	return http.ListenAndServeTLS(addr, certFile, keyFile, h)
}

func (h *HttpServer) LoadHTMLFiles(files ...string) {
	h.templ = template.Must(template.New("").Delims("{{", "}}").ParseFiles(files...))
}
