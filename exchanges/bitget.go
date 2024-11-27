package exchanges

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

type Bitget struct {
	client *http.Client
}

type BitgetResponse struct {
	Code        string `json:"code"`
	Msg         string `json:"msg"`
	RequestTime int64  `json:"requestTime"`
	Data        []struct {
		Open         string `json:"open"`
		Symbol       string `json:"symbol"`
		High24H      string `json:"high24h"`
		Low24H       string `json:"low24h"`
		LastPr       string `json:"lastPr"`
		QuoteVolume  string `json:"quoteVolume"`
		BaseVolume   string `json:"baseVolume"`
		UsdtVolume   string `json:"usdtVolume"`
		Ts           string `json:"ts"`
		BidPr        string `json:"bidPr"`
		AskPr        string `json:"askPr"`
		BidSz        string `json:"bidSz"`
		AskSz        string `json:"askSz"`
		OpenUtc      string `json:"openUtc"`
		ChangeUtc24H string `json:"changeUtc24h"`
		Change24H    string `json:"change24h"`
	} `json:"data"`
}

func NewBitget(proxyURL string) *Bitget {
	return &Bitget{
		client: &http.Client{},
	}
}

func (b *Bitget) Name() string {
	return "Bitget"
}

func (b *Bitget) GetPrice(symbol string) (float64, error) {
	url := fmt.Sprintf("https://api.bitget.com/api/v2/spot/market/tickers?symbol=%sUSDT", symbol)
	resp, err := b.client.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var response BitgetResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return 0, err
	}

	return strconv.ParseFloat(response.Data[0].LastPr, 64)
}
