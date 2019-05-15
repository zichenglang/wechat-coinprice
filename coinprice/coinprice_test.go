package coinprice

import "testing"

func TestGetQuote(t *testing.T) {
	c,er := getQuote("btc")
	t.Log(c,er)
}
