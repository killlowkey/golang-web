package session

import (
	"context"
	"net/http"
)

// Session session 行为接口
type Session interface {
	// Get 获取 key 与之对应的 value
	Get(ctx context.Context, key string) (string, error)
	// Set 设置 key 响应的 value
	Set(ctx context.Context, key, value string) error
	// ID 获取 Session id 编号
	ID() string
}

// Store 用于存储 Session 会话
type Store interface {
	// Generate 生成 Session
	Generate(ctx context.Context, id string) (Session, error)
	// Refresh 刷新 Session
	Refresh(ctx context.Context, id string) error
	// Remove 移除 Session
	Remove(ctx context.Context, id string) error
	// Get 获取 Session
	Get(ctx context.Context, id string) (Session, error)
}

// Propagator 在请求和响应之间进行 Session 传播
type Propagator interface {
	// Inject 注入 Session id 到响应中
	Inject(id string, writer http.ResponseWriter) error
	// Extract 从请求提取出 Session id
	Extract(req *http.Request) (string, error)
	// Remove 移除 Session
	Remove(writer http.ResponseWriter) error
}
