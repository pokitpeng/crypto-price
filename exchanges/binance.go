package exchanges

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

type Binance struct {
	client *http.Client
}

type BinanceResponse struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

func NewBinance(proxyURL string) *Binance {
	return &Binance{
		client: createHTTPClient(proxyURL),
	}
}

func (b *Binance) Name() string {
	return "Binance"
}

func (b *Binance) GetPrice(symbol string) (float64, error) {
	url := fmt.Sprintf("https://api.binance.com/api/v3/ticker/price?symbol=%sUSDT", symbol)
	resp, err := b.client.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var response BinanceResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return 0, err
	}

	return strconv.ParseFloat(response.Price, 64)
}
