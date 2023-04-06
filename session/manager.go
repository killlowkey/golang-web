package session

import (
	"github.com/killlowkey/web"
)

// Manager 用于管理 Session
type Manager struct {
	Store
	Propagator
	SessCtxKey string
}

// GetSession 获取 Session
func (m *Manager) GetSession(ctx *web.Context) (Session, error) {
	if ctx.UserValues == nil {
		ctx.UserValues = make(map[string]any, 1)
	}

	// 先从缓存获取
	if v, ok := ctx.UserValues[m.SessCtxKey]; ok {
		return v.(Session), nil
	}

	// 从请求中提取到 session id
	id, err := m.Extract(ctx.Request)
	if err != nil {
		return nil, err
	}
	// 从底层存储中（数据库、Redis、内存）获取 Session
	if session, err := m.Store.Get(ctx.Request.Context(), id); err != nil {
		return nil, err
	} else {
		// 进行缓存
		ctx.UserValues[m.SessCtxKey] = session
		return session, nil
	}
}

// InitSession 初始化 session
func (m *Manager) InitSession(ctx web.Context, id string) (Session, error) {
	session, err := m.Generate(ctx.Request.Context(), id)
	if err != nil {
		return nil, err
	}

	// 注入 Session 到响应中
	err = m.Inject(id, ctx.Writer)
	if err != nil {
		return nil, err
	}

	return session, nil
}

// RefreshSession 刷新 Session
func (m *Manager) RefreshSession(ctx *web.Context) (Session, error) {
	sess, err := m.GetSession(ctx)
	if err != nil {
		return nil, err
	}
	// 刷新存储的过期时间
	err = m.Refresh(ctx.Request.Context(), sess.ID())
	if err != nil {
		return nil, err
	}
	// 重新注入 HTTP 里面
	if err = m.Inject(sess.ID(), ctx.Writer); err != nil {
		return nil, err
	}
	return sess, nil
}

// RemoveSession 删除 Session
func (m *Manager) RemoveSession(ctx *web.Context) error {
	sess, err := m.GetSession(ctx)
	if err != nil {
		return err
	}
	err = m.Store.Remove(ctx.Request.Context(), sess.ID())
	if err != nil {
		return err
	}
	return m.Propagator.Remove(ctx.Writer)
}
