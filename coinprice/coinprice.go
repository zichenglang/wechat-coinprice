package coinprice

import (
	"encoding/json"
	"fmt"
	"github.com/songtianyi/rrframework/logs"
	"github.com/songtianyi/wechat-go/wxweb"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type QuoteType struct {
	BestAsk        string    `json:"best_ask"`
	BestBid        string    `json:"best_bid"`
	InstrumentID   string    `json:"instrument_id"`
	ProductID      string    `json:"product_id"`
	Last           string    `json:"last"`
	Ask            string    `json:"ask"`
	Bid            string    `json:"bid"`
	Open24H        string    `json:"open_24h"`
	High24H        string    `json:"high_24h"`
	Low24H         string    `json:"low_24h"`
	BaseVolume24H  string    `json:"base_volume_24h"`
	Timestamp      time.Time `json:"timestamp"`
	QuoteVolume24H string    `json:"quote_volume_24h"`
}

var (
	quoteChan          = make(chan string, 100)
	getQuoteUrl string = "https://www.okex.com/api/spot/v3/instruments/%s_USDT/ticker"
	//okexClient =
)
// Register plugin
func Register(session *wxweb.Session) {
	session.HandlerRegister.Add(wxweb.MSG_TEXT, wxweb.Handler(listenCmd), "coinprice")
	if err := session.HandlerRegister.EnableByName("coinprice"); err != nil {
		logs.Error(err)
		return
	}
}

func getQuote(coin string) (string, error) {
	cli := http.Client{}
	cli.Timeout = time.Second * time.Duration(5)
	url := fmt.Sprintf(getQuoteUrl, coin)
	resp, err := cli.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	quote := QuoteType{}
	err = json.Unmarshal(result, &quote)
	if err != nil {
		return "", err
	}
	lastPrice := 0.0
	if lastPrice, err = strconv.ParseFloat(quote.Last, 32); err != nil {
		return "", err
	}
	openPrice := 0.0
	if openPrice, err = strconv.ParseFloat(quote.Open24H, 32); err != nil {
		return "", err
	}
	if openPrice < 0.000001 {
		openPrice = 0.000001
	}
	percent := (lastPrice - openPrice) / openPrice * 100
	r := fmt.Sprintf("币种 : %s \n"+
		"当前价格 : %s$ \n"+
		"24h最高价格 : %s$ \n"+
		"24h最低价格 : %s$ \n"+
		"涨幅 :  %.2f%%  ", quote.InstrumentID, quote.Last, quote.High24H, quote.Low24H, percent)

	return r, nil
}

func GetQuote(session *wxweb.Session, msg *wxweb.ReceivedMessage ) {
	res, err := getQuote(msg.Content)
	if err != nil {
		return
	}
	session.SendText(res, session.Bot.UserName, wxweb.RealTargetUserName(session, msg))
}
func listenCmd(session *wxweb.Session, msg *wxweb.ReceivedMessage) {

	if !msg.IsGroup || len(msg.Content) > 5 {
		return
	}

	logs.Info("coin type is : ", msg.Content)
	go GetQuote(session, msg)
//	logs.Info("coin type is : ", msg.Content)
//	select {
//	case <-time.After(time.Second * 5):
//		return
//	case quote := <-quoteChan:
//		session.SendText(quote, session.Bot.UserName, wxweb.RealTargetUserName(session, msg))
//	}
}
