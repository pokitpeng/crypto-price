package exchanges

import (
	"testing"
)

func TestXT_GetPrice(t *testing.T) {
	xt := NewXT("http://127.0.0.1:7890")
	price, err := xt.GetPrice("BTC")
	if err != nil {
		t.Errorf("XT.GetPrice() error = %v", err)
	}
	t.Logf("XT.GetPrice() = %v", price)
}
