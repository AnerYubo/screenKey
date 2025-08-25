package app

import (
	"context"

	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

var wailsContext *context.Context

type App struct {
	ctx          context.Context
	TOTPHandler  *TOTPHandler
	KnockHandler *KnockHandler
	// 未来还能加更多 Handler，比如 UserHandler
}

func NewApp() *App {
	a := &App{}
	ctx := context.Background()
	a.TOTPHandler = NewTOTPHandler(ctx)
	a.KnockHandler = NewKnockHandler(ctx)
	return a
}

func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
	a.TOTPHandler = NewTOTPHandler(ctx)
	a.KnockHandler = NewKnockHandler(ctx)
	wailsContext = &ctx

}

// OnSecondInstanceLaunch 阻止多开
func OnSecondInstanceLaunch(secondInstanceData options.SecondInstanceData) {
	runtime.WindowUnminimise(*wailsContext)
	runtime.Show(*wailsContext)
}

// ====== 窗口控制方法 ======
func (a *App) HideWindow() {
	runtime.Hide(a.ctx)
}

func (a *App) ShowWindow() {
	runtime.Show(a.ctx)
}
