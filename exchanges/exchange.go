package exchanges

import (
	"log/slog"
	"net/http"
	"net/url"
	"time"
)

type Exchange interface {
	GetPrice(symbol string) (float64, error)
	Name() string
}

var (
	AvailableExchanges = make(map[string]Exchange)
)

const (
	ExchangeBinance = "Binance"
	// ExchangeOKX     = "OKX"
	ExchangeXT     = "XT"
	ExchangeBitget = "Bitget"
)

func RegisterExchanges(proxyURL string) {
	AvailableExchanges[ExchangeBinance] = NewBinance(proxyURL)
	AvailableExchanges[ExchangeXT] = NewXT(proxyURL)
	AvailableExchanges[ExchangeBitget] = NewBitget(proxyURL)
}

// 创建带代理的 HTTP 客户端
func createHTTPClient(proxyURL string) *http.Client {
	if proxyURL != "" {
		proxy, err := url.Parse(proxyURL)
		if err != nil {
			slog.Error("代理地址解析错误", "error", err)
			return &http.Client{}
		}

		transport := &http.Transport{
			Proxy: http.ProxyURL(proxy),
		}

		return &http.Client{
			Transport: transport,
			Timeout:   10 * time.Second,
		}
	}

	return &http.Client{
		Timeout: 10 * time.Second,
	}
}
