package websocket

import "encoding/json"

type wsReq struct {
	Id    uint64 `json:"id"`
	Type  string `json:"type"`
	Topic string `json:"topic"`
	Req   int    `json:"req"`
}

type wsResp struct {
	Success   bool            `json:"success"`
	Code      string          `json:"code"`
	Msg       string          `json:"msg"`
	Timestamp int64           `json:"timestamp"`
	Data      json.RawMessage `json:"data"`
}

type connData struct {
	BulletToken     string           `json:"bulletToken"`
	InstanceServers []instanceServer `json:"instanceServers"`
	HistoryServers  []historyServer  `json:"historyServers"`
}

type historyServer struct {
	Endpoint string `json:"endpoint"`
	Encrypt  bool   `json:"encrypt"`
	UserType string `json:"userType"`
}

type instanceServer struct {
	PingInterval int    `json:"pingInterval"`
	Endpoint     string `json:"endpoint"`
	Protocol     string `json:"protocol"`
	Encrypt      bool   `json:"encrypt"`
	PingTimeout  int    `json:"pingTimeout"`
	UserType     string `json:"userType"`
}

// Response represents server response after sending data.
type Response struct {
	Id   string `json:"id"`
	Type string `json:"type"`
}

// OrderBook is the type received from Orderbook subscription
type OrderBook struct {
	Volume float64 `json:"volume"`
	Price  float64 `json:"price"`
	Count  float64 `json:"count"`
	Action string  `json:"action"`
	Time   int64   `json:"time"`
	Type   string  `json:"type"`
}

// History is the type received from History subscription
type History struct {
	Id        string  `json:"oid"`
	Price     float64 `json:"price"`
	Count     float64 `json:"count"`
	Time      int64   `json:"time"`
	VolValue  float64 `json:"volValue"`
	Direction string  `json:"direction"`
}

// Market is the type received from Tick and Market subscription
type Market struct {
	CoinType      string  `json:"coinType"`
	Trading       bool    `json:"trading"`
	Symbol        string  `json:"symbol"`
	LastDealPrice float64 `json:"lastDealPrice"`
	Buy           float64 `json:"buy"`
	Sell          float64 `json:"sell"`
	Change        float64 `json:"change"`
	CoinTypePair  string  `json:"coinTypePair"`
	Sort          int     `json:"sort"`
	FeeRate       float64 `json:"feeRate"`
	VolValue      float64 `json:"volValue"`
	High          float64 `json:"high"`
	Datetime      int64   `json:"datetime"`
	Vol           float64 `json:"vol"`
	Low           float64 `json:"low"`
	ChangeRate    float64 `json:"changeRate"`
}
