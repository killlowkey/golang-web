package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/killlowkey/web/binding"
	"net/http"
	"net/url"
)

// Context 处理请求输入输出
type Context struct {
	Request *http.Request       // http 请求
	Writer  http.ResponseWriter // http 响应

	Params  Params     // http 路由参数
	handler HandleFunc // http 路由处理器

	Route      string // 路由信息
	RespStatus int    // 保存响应状态码
	RespData   []byte

	UserValues map[string]any

	h          *HttpServer
	queryCache url.Values
}

// reset 重置 Context，从 Context Pool 获取后需要进行重置
func (c *Context) reset() {
	c.Request = nil
	c.Writer = nil
	c.Params = c.Params[:0]
	c.handler = nil
	c.RespStatus = 200
	c.RespData = c.RespData[:0]
	c.Route = ""
	c.h = nil
	c.queryCache = nil
	c.UserValues = nil
}

// ===============================
// ====== Request register =======
// ===============================

// Param 获取请求路由对应的参数
// 绑定路由为 /user/:name，实际请求 /user/ray
// 方法参数传入 name，得到返回值为 ray
func (c *Context) Param(key string) string {
	res, ok := c.Params.Get(key)
	if !ok {
		return ""
	}
	return res
}

// Query 查询 url 中参数值
// GET /path?id=1234  => c.Query("id") == "1234"
func (c *Context) Query(key string) (value string) {
	value, _ = c.GetQuery(key)
	return
}
func (c *Context) GetQuery(key string) (string, bool) {
	if values, ok := c.GetQueryArray(key); ok {
		return values[0], ok
	}
	return "", false
}
func (c *Context) GetQueryArray(key string) (values []string, ok bool) {
	c.initQueryCache()
	values, ok = c.queryCache[key]
	return
}

func (c *Context) initQueryCache() {
	if c.queryCache == nil {
		if c.Request != nil {
			c.queryCache = c.Request.URL.Query()
		} else {
			c.queryCache = url.Values{}
		}
	}
}

// BindJSON 绑定请求 body
// TODO request 的 body 只能读取一次，此处仅为演示，后续需要修改
func (c *Context) BindJSON(val any) error {
	return binding.JSON.Bind(c.Request, val)
}

// BindXML 从请求 XML 进行绑定
func (c *Context) BindXML(obj any) error {
	return binding.XML.Bind(c.Request, obj)
}

// BindQuery  从请求的查询参数绑定
func (c *Context) BindQuery(obj any) error {
	return binding.QUERY.Bind(c.Request, obj)
}

// BindProtobuf 从请求的 protobuf 绑定
func (c *Context) BindProtobuf(obj any) error {
	return binding.PROTOBUF.Bind(c.Request, obj)
}

// =================================
// ======== Response register ======
// =================================

func (c *Context) Status(status int) {
	c.RespStatus = status
}

func (c *Context) Write(data []byte) {
	c.RespData = data
}

func (c *Context) WriteWithStatus(status int, data []byte) {
	c.RespStatus = status
	c.RespData = data
}

func (c *Context) Header(key, val string) {
	if val == "" {
		c.Writer.Header().Del(key)
		return
	}
	c.Writer.Header().Add(key, val)
}

func (c *Context) String(code int, format string, values ...any) {
	msg := format
	if len(values) != 0 {
		msg = fmt.Sprintf(format, values...)
	}
	c.WriteWithStatus(code, []byte(msg))
}

func (c *Context) JSON(status int, val any) {
	if val == nil {
		c.WriteWithStatus(http.StatusInternalServerError, []byte("val 为 nil"))
		return
	}

	if data, err := json.Marshal(val); err != nil {
		c.WriteWithStatus(http.StatusInternalServerError, []byte(err.Error()))
	} else {
		c.Header("Content-type", "application/json; charset=utf-8")
		c.WriteWithStatus(status, data)
	}
}

// HTML 渲染 HTML 模版
func (c *Context) HTML(code int, name string, obj any) {
	// TODO 错误示例，为了节省时间
	var data = &bytes.Buffer{}

	if err := c.h.templ.ExecuteTemplate(data, name, obj); err != nil {
		c.WriteWithStatus(http.StatusInternalServerError, []byte(err.Error()))
		return
	}

	c.Status(code)
	c.Header("Content-type", "text/html; charset=utf-8")
	c.Write(data.Bytes())
}

// File 文件服务器
func (c *Context) File(filepath string) {
	http.ServeFile(c.Writer, c.Request, filepath)
}
