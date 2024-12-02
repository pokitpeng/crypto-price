package exchanges

import (
	"testing"
)

func TestOkx_GetPrice(t *testing.T) {
	okx := NewOkx("http://127.0.0.1:7890")
	price, err := okx.GetPrice("ACT")
	if err != nil {
		t.Errorf("Okx.GetPrice() error = %v", err)
	}
	t.Logf("Okx.GetPrice() = %v", price)
}
