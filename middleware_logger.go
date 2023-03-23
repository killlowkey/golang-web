package web

import (
	"fmt"
	"log"
	"time"
)

type LogEntry struct {
	Method     string        // 请求方法
	Path       string        // 请求路径
	RespStatus int           // 响应状态码
	Latency    time.Duration // 处理耗时
}

type LogWriter interface {
	Write(entry *LogEntry) error
}

type LogMiddlewareBuild struct {
	writer LogWriter
}

func NewLogMiddlewareBuild(writer LogWriter) *LogMiddlewareBuild {
	if writer == nil {
		panic("web：AccessLog middleware 传入 writer 为 nil")
	}
	return &LogMiddlewareBuild{
		writer: writer,
	}
}

func (b *LogMiddlewareBuild) Build() Middleware {
	return func(next HandleFunc) HandleFunc {
		return func(ctx *Context) {
			start := time.Now()
			defer func() {
				e := &LogEntry{
					Method:     ctx.Request.Method,
					Path:       ctx.Request.URL.Path,
					Latency:    time.Now().Sub(start),
					RespStatus: ctx.RespStatus,
				}
				err := b.writer.Write(e)
				if err != nil {
					log.Println("web：[Log Middleware]写入 log 出错：" + err.Error())
				}
			}()
			next(ctx)
		}
	}
}

type ConsoleLogger struct {
}

func (c ConsoleLogger) Write(e *LogEntry) error {
	msg := fmt.Sprintf("[web]: %v | %#v %s %d %v",
		time.Now().Format("2006/01/02 - 15:04:05"),
		e.Path, e.Method, e.RespStatus, e.Latency)
	fmt.Println(msg)
	return nil
}

func logger() Middleware {
	build := NewLogMiddlewareBuild(new(ConsoleLogger))
	return build.Build()
}
