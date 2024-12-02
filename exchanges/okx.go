package exchanges

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

type Okx struct {
	client *http.Client
}

type OkxResponse struct {
	Code string `json:"code"`
	Data []struct {
		Last string `json:"last"`
	} `json:"data"`
}

func NewOkx(proxyURL string) *Okx {
	return &Okx{
		client: createHTTPClient(proxyURL),
	}
}

func (o *Okx) Name() string {
	return "OKX"
}

func (o *Okx) GetPrice(symbol string) (float64, error) {
	url := fmt.Sprintf("https://www.okx.com/api/v5/market/ticker?instId=%s-USDT", symbol)
	resp, err := o.client.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var response OkxResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return 0, err
	}

	if len(response.Data) == 0 {
		return 0, fmt.Errorf("没有获取到okx %s-USDT价格数据", symbol)
	}

	return strconv.ParseFloat(response.Data[0].Last, 64)
}
