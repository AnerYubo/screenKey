package main

import (
	"embed"
	"fmt"
	"os"
	"screenKey/app"
	"screenKey/sqlite"

	"context"

	"github.com/getlantern/systray"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed assets/icon.ico
var iconFS embed.FS

// 全局保存 wails 上下文，用于托盘控制窗口
var mainCtx context.Context

func main() {
	// 初始化数据库
	if err := sqlite.InitDatabase(); err != nil {
		panic(err)
	}

	// 创建 app 实例
	appInstance := app.NewApp()
	TOTPHandler := appInstance.TOTPHandler
	KnockHandler := appInstance.KnockHandler

	// 启动系统托盘（独立 goroutine）
	go func() {
		systray.Run(onReady, onExit)
	}()

	// 启动 Wails 应用
	err := wails.Run(&options.App{
		Title:     "screenKey",
		Width:     1024,
		Height:    768,
		MinWidth:  500,
		MinHeight: 600,
		DragAndDrop: &options.DragAndDrop{
			EnableFileDrop: true,
		},
		Windows: &windows.Options{},
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		SingleInstanceLock: &options.SingleInstanceLock{
			UniqueId:               "c9c8fd93-6758-4144-87d1-34bdb0a8bd60",
			OnSecondInstanceLaunch: app.OnSecondInstanceLaunch,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},

		// 应用启动时保存上下文
		OnStartup: func(ctx context.Context) {
			mainCtx = ctx
			appInstance.Startup(ctx)
		},

		// 点击窗口 X 时只隐藏，不退出
		OnBeforeClose: func(ctx context.Context) (prevent bool) {
			runtime.WindowHide(ctx)
			return true
		},

		// 绑定 Go 方法给前端调用
		Bind: []interface{}{
			TOTPHandler, KnockHandler,
		},
	})

	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}
}

// 系统托盘初始化
func onReady() {
	// 设置托盘图标
	iconData, err := iconFS.ReadFile("assets/icon.ico")
	if err != nil {
		fmt.Println("加载图标失败:", err)
		return
	}
	systray.SetIcon(iconData)
	systray.SetTitle("ScreenKey")
	systray.SetTooltip("左键点击显示窗口，右键打开菜单")

	// ⚡ 默认菜单项（隐藏后作为左键触发器）
	mDefault := systray.AddMenuItem("打开", "显示主窗口")

	// 退出菜单
	mQuit := systray.AddMenuItem("退出", "退出应用程序")

	// 监听菜单点击
	go func() {
		for {
			select {
			case <-mDefault.ClickedCh:
				// 左键点击：显示主窗口
				if mainCtx != nil {
					runtime.WindowShow(mainCtx)
					runtime.WindowUnminimise(mainCtx)
				}
			case <-mQuit.ClickedCh:
				// 退出应用
				systray.Quit()
				os.Exit(0)
			}
		}
	}()
}

// 托盘退出回调
func onExit() {
	fmt.Println("系统托盘退出...")
}
