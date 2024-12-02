package exchanges

import (
	"testing"
)

func TestBinance_GetPrice(t *testing.T) {
	binance := NewBinance("http://127.0.0.1:7890")
	price, err := binance.GetPrice("BTC")
	if err != nil {
		t.Errorf("Binance.GetPrice() error = %v", err)
	}
	t.Logf("Binance.GetPrice() = %v", price)
}
