package web

type PanicRecovery interface {
	Recovery(ctx *Context, err any)
}

type RecoveryMiddlewareBuild struct {
	recovery PanicRecovery
}

func NewRecoveryMiddlewareBuild(recovery PanicRecovery) *RecoveryMiddlewareBuild {
	if recovery == nil {
		panic("web：recovery middleware 传入 recovery 为 nil")
	}
	return &RecoveryMiddlewareBuild{
		recovery: recovery,
	}
}

func (m *RecoveryMiddlewareBuild) Build() Middleware {
	return func(next HandleFunc) HandleFunc {
		return func(ctx *Context) {
			defer func() {
				// 检测是否被 panic
				if err := recover(); err != nil {
					m.recovery.Recovery(ctx, err)
				}
			}()
			next(ctx)
		}
	}
}

type DefaultPanicRecovery int

func (d DefaultPanicRecovery) Recovery(ctx *Context, err any) {
	ctx.JSON(500, H{
		"code":  500,
		"path":  ctx.Request.URL.Path,
		"route": ctx.Route,
		"msg":   err,
	})
}

func recovery() Middleware {
	build := NewRecoveryMiddlewareBuild(new(DefaultPanicRecovery))
	return build.Build()
}
