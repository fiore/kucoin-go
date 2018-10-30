package websocket

const (
	urlServers = "https://kitchen.kucoin.com/v1/bullet/usercenter/loginUser?protocol=websocket&encrypt=true"
	urlOrderBook = "/trade/%s_TRADE"   // Symbol
	urlHistory   = "/trade/%s_HISTORY" // Symbol
	urlTick      = "/market/%s_TICK"   // Symbol
	urlMarket    = "/market/%s"        // Coin
)

// Topic represents topic in which user can perform Type actions.
type Topic byte

const (
	TOrderBook Topic = iota
	THistory
	Tick
	TMarket
)

// Type represents actions that can be performed.
type Type string

const (
	Subscribe   Type = "subscribe"
	Unsubscribe Type = "unsubscribe"
	Ping        Type = "ping"
	Close       Type = "close"
)
