package main

import (
	"github.com/doneland/yquotes"
	"testing"
)

func TestQuote(t *testing.T) {
	symbol := "AAPL"
	price, err := yquotes.GetPrice(symbol)
	if err == nil {
		t.Logf("%+v",price)
	} else {
		t.Error(err)
	}
}