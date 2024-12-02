package exchanges

import (
	"testing"
)

func TestBitget_GetPrice(t *testing.T) {
	bitget := NewBitget("http://127.0.0.1:7890")
	price, err := bitget.GetPrice("BTC")
	if err != nil {
		t.Errorf("Bitget.GetPrice() error = %v", err)
	}
	t.Logf("Bitget.GetPrice() = %v", price)
}
