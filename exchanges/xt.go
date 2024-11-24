package exchanges

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type XT struct {
	client *http.Client
}

type XTResponse struct {
	Rc     int    `json:"rc"`
	Mc     string `json:"mc"`
	Ma     []any  `json:"ma"`
	Result []struct {
		Symbol string `json:"s"`
		Time   int64  `json:"t"`
		Price  string `json:"p"`
	} `json:"result"`
}

func NewXT(proxyURL string) *XT {
	return &XT{
		client: createHTTPClient(proxyURL),
	}
}

func (b *XT) Name() string {
	return "XT"
}

func (b *XT) GetPrice(symbol string) (float64, error) {
	symbol_lower := strings.ToLower(symbol)
	url := fmt.Sprintf("https://www.xt.com/sapi/v4/market/public/ticker/price?symbol=%s_usdt", symbol_lower)
	resp, err := b.client.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var response XTResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return 0, err
	}

	return strconv.ParseFloat(response.Result[0].Price, 64)
}
