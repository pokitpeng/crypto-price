package main

import (
	"crypto-price/config"
	"crypto-price/exchanges"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/getlantern/systray"
)

type App struct {
	config       *config.Config
	tokens       map[string]*systray.MenuItem
	exchanges    map[string]exchanges.Exchange
	refreshBtn   *systray.MenuItem
	refreshItems []*systray.MenuItem
}

func main() {
	systray.Run(onReady, onExit)
}

func onReady() {
	app := &App{
		tokens:    make(map[string]*systray.MenuItem),
		exchanges: make(map[string]exchanges.Exchange),
	}

	var err error
	app.config, err = config.LoadConfig()
	if err != nil {
		slog.Error("加载配置失败", "error", err)
		systray.Quit()
		return
	}

	if os.Getenv("PROXY_URL") != "" {
		app.config.ProxyURL = os.Getenv("PROXY_URL")
	}

	if app.config.ProxyURL != "" {
		slog.Info("使用代理", "proxy_url", app.config.ProxyURL)
	}

	exchanges.RegisterExchanges(app.config.ProxyURL)
	app.exchanges = exchanges.AvailableExchanges

	app.setupMenu()
	app.startPriceUpdate()
}

func (app *App) setupMenu() {
	systray.SetTitle("加载中...")

	// 添加代币选择菜单
	tokenMenu := systray.AddMenuItem("选择显示代币", "")
	for _, token := range app.config.Tokens {
		// 获取当前价格并格式化菜单项文本
		menuText := fmt.Sprintf("%s ($%.2f)", token.Name, token.Price)
		menuItem := tokenMenu.AddSubMenuItem(menuText, "")
		app.tokens[token.Symbol] = menuItem
		if token.Symbol == app.config.ActiveToken {
			menuItem.Check()
		}
	}

	// 添加新代币
	// _ = systray.AddMenuItem("添加新代币", "")

	// 设置刷新时间
	refreshTimeMenu := systray.AddMenuItem("设置刷新时间", "")
	refreshTimes := []int{15, 30, 60, 300}
	app.refreshItems = make([]*systray.MenuItem, len(refreshTimes))
	for i, t := range refreshTimes {
		app.refreshItems[i] = refreshTimeMenu.AddSubMenuItem(fmt.Sprintf("%d秒", t), "")
		if t == app.config.RefreshTime {
			app.refreshItems[i].Check()
		}
	}

	app.refreshBtn = systray.AddMenuItem("刷新", "立即刷新价格")

	systray.AddSeparator()
	quitBtn := systray.AddMenuItem("退出", "退出程序")

	// 处理刷新按钮
	go func() {
		for {
			<-app.refreshBtn.ClickedCh
			app.updatePrice()
		}
	}()

	// 处理退出按钮
	go func() {
		<-quitBtn.ClickedCh
		systray.Quit()
	}()

	// 处理代币选择
	for symbol, menuItem := range app.tokens {
		go func(s string, m *systray.MenuItem) {
			for {
				<-m.ClickedCh
				app.setActiveToken(s)
			}
		}(symbol, menuItem)
	}

	// 处理刷新时间设置
	for i, item := range app.refreshItems {
		go func(seconds int, m *systray.MenuItem) {
			for {
				<-m.ClickedCh
				app.setRefreshTime(seconds)
			}
		}(refreshTimes[i], item)
	}
}

func (app *App) startPriceUpdate() {
	go func() {
		for {
			app.updatePrice()
			time.Sleep(time.Duration(app.config.RefreshTime) * time.Second)
		}
	}()
}

func (app *App) updatePrice() {
	for i, token := range app.config.Tokens {
		if exchange, ok := app.exchanges[token.Exchange]; ok {
			price, err := exchange.GetPrice(token.Symbol)
			if err != nil {
				slog.Error("获取价格失败", "symbol", token.Symbol, "error", err)
				continue
			}
			app.config.Tokens[i].Price = price

			// 更新菜单项文本
			if menuItem, ok := app.tokens[token.Symbol]; ok {
				menuItem.SetTitle(fmt.Sprintf("%s ($%.3f)", token.Name, price))
			}

			if token.Symbol == app.config.ActiveToken {
				systray.SetTitle(fmt.Sprintf("%s $%.3f", token.Symbol, price))
			}
		}
	}
}

func (app *App) setActiveToken(symbol string) {
	for _, menuItem := range app.tokens {
		menuItem.Uncheck()
	}
	app.tokens[symbol].Check()
	app.config.ActiveToken = symbol
	config.SaveConfig(app.config)
	app.updatePrice()
}

func (app *App) setRefreshTime(seconds int) {
	// 取消所有刷新时间的对钩
	for _, item := range app.refreshItems {
		item.Uncheck()
	}

	// 为当前选中的刷新时间添加对钩
	refreshTimes := []int{15, 30, 60, 300, 600}
	for i, t := range refreshTimes {
		if t == seconds {
			app.refreshItems[i].Check()
			break
		}
	}

	app.config.RefreshTime = seconds
	config.SaveConfig(app.config)
}

func onExit() {
	slog.Info("正在退出程序...")
	os.Exit(0)
}
