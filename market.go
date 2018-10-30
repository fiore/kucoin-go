package kucoin

type rawMarket struct {
	Success   bool     `json:"success,omitempty"`
	Code      string   `json:"code,omitempty"`
	Msg       string   `json:"msg,omitempty"`
	Timestamp int64    `json:"timestamp,omitempty"`
	Data      []string `json:"data,omitempty"`
}
