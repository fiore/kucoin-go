// Package kucoin is an implementation of the Kucoin API in Golang.
package kucoin

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	kucoinURL = "https://api.kucoin.com/v1/"
)

// Custom errors used when an input is required
var (
	ErrSymbolRequired    = errors.New("Symbol is required")
	ErrAllParamsRequired = errors.New("All parameters are required")
	ErrNonExistingSymbol = errors.New("Entered symbol doesn't exist in Kucoin")
	ErrNonExistingMarket = errors.New("Entered market doesn't exist in Kucoin")
)

var (
	coinsPairsList           = []CoinPair{}
	openMarketsList          = []string{}
	defaultMessageWrongInput = "Entered invalid parameter. Accepted values: [%s]"
)

func (k *Kucoin) getCoinsPairsList() []CoinPair {
	if len(coinsPairsList) == 0 {
		coinPair, err := k.GetCoinsPairs()
		if err == nil {
			coinsPairsList = coinPair
		}
	}
	return coinsPairsList
}

func (k *Kucoin) containsCoinsPairs(coinPair string) bool {
	coinsPairsList = k.getCoinsPairsList()
	for _, cp := range coinsPairsList {
		if cp.CoinPair == coinPair {
			return true
		}
	}
	return false
}

func (k *Kucoin) getOpenMarketsList() []string {
	if len(openMarketsList) == 0 {
		openMarket, err := k.GetOpenMarkets()
		if err == nil {
			openMarketsList = openMarket
		}
	}
	return openMarketsList
}

func (k *Kucoin) containsOpenMarkets(openMarket string) bool {
	openMarketsList = k.getOpenMarketsList()
	for _, om := range openMarketsList {
		if om == openMarket {
			return true
		}
	}
	return false
}

// New returns an instantiated Kucoin struct.
func New(apiKey, apiSecret string) *Kucoin {
	client := newClient(apiKey, apiSecret)
	return &Kucoin{client}
}

// NewCustomClient returns an instantiated Kucoin struct with custom http client.
func NewCustomClient(apiKey, apiSecret string, httpClient http.Client) *Kucoin {
	client := newClient(apiKey, apiSecret)
	client.httpClient = httpClient
	return &Kucoin{client}
}

// NewCustomTimeout returns an instantiated Kucoin struct with custom timeout.
func NewCustomTimeout(apiKey, apiSecret string, timeout time.Duration) *Kucoin {
	client := newClient(apiKey, apiSecret)
	client.httpClient.Timeout = timeout
	return &Kucoin{client}
}

func doArgs(args ...string) map[string]string {
	m := make(map[string]string)
	var lastK = ""
	for i, v := range args {
		if i&1 == 0 {
			lastK = v
		} else {
			m[lastK] = v
		}
	}
	return m
}

// handleErr gets JSON response from livecoin API en deal with error.
func handleErr(r interface{}) error {
	switch v := r.(type) {
	case map[string]interface{}:
		err := r.(map[string]interface{})["error"]
		if err != nil {
			switch v := err.(type) {
			case map[string]interface{}:
				errorMessage := err.(map[string]interface{})["message"]
				return errors.New(errorMessage.(string))
			default:
				return fmt.Errorf("don't recognized type %T", v)
			}
		}
	case []interface{}:
		return nil
	default:
		return fmt.Errorf("don't recognized type %T", v)
	}

	return nil
}

// Kucoin represent a Kucoin client.
type Kucoin struct {
	client *client
}

// SetDebug enables/disables http request/response dump.
func (k *Kucoin) SetDebug(enable bool) {
	k.client.debug = enable
}

// GetUserInfo is used to get the user information at Kucoin along with other meta data.
func (k *Kucoin) GetUserInfo() (userInfo UserInfo, err error) {
	r, err := k.client.do("GET", "user/info", nil, true)
	if err != nil {
		return
	}
	var response interface{}
	if err = json.Unmarshal(r, &response); err != nil {
		return
	}
	if err = handleErr(response); err != nil {
		return
	}
	var rawRes rawUserInfo
	err = json.Unmarshal(r, &rawRes)
	userInfo = rawRes.Data
	return
}

// GetSymbols is used to get the all open and available trading markets at Kucoin along with other meta data.
func (k *Kucoin) GetSymbols() (symbols []Symbol, err error) {
	r, err := k.client.do("GET", "market/open/symbols", nil, false)
	if err != nil {
		return
	}
	var response interface{}
	if err = json.Unmarshal(r, &response); err != nil {
		return
	}
	if err = handleErr(response); err != nil {
		return
	}
	var rawRes rawSymbols
	err = json.Unmarshal(r, &rawRes)
	symbols = rawRes.Data
	return

}

// GetCoinsPairs is used to get the all available trading markets at Kucoin.
func (k *Kucoin) GetCoinsPairs() (coinPair []CoinPair, err error) {
	r, err := k.client.do("GET", "market/open/coins-trending", nil, true)
	if err != nil {
		return
	}
	var response interface{}
	if err = json.Unmarshal(r, &response); err != nil {
		return
	}
	if err = handleErr(response); err != nil {
		return
	}
	var rawRes rawCoinPair
	err = json.Unmarshal(r, &rawRes)
	coinPair = rawRes.Data
	return
}

// GetUserSymbols is used to get the all open and available trading markets at Kucoin along with other meta data.
// The user should be logged to call this method.
// Example:
// - Market = BTC
// - Symbol = KCS-BTC
// - Filter = FAVOURITE | STICK
func (k *Kucoin) GetUserSymbols(market, symbol, filter string) (symbols []Symbol, err error) {
	if len(market) > 1 {
		if !k.containsOpenMarkets(strings.ToUpper(market)) {
			return symbols, ErrNonExistingMarket
		}
	}
	if len(symbol) > 1 {
		if !k.containsCoinsPairs(strings.ToUpper(symbol)) {
			return symbols, ErrNonExistingSymbol
		}
	}
	if len(filter) > 1 {
		if filter != "FAVOURITE" && filter != "STICK" {
			return symbols, fmt.Errorf(defaultMessageWrongInput, strings.Join([]string{"FAVOURITE", "STICK"}, ","))
		}
	}

	payload := map[string]string{
		"symbol": strings.ToUpper(symbol),
		"market": strings.ToUpper(market),
		"filter": strings.ToUpper(filter),
	}

	r, err := k.client.do("GET", "market/symbols", payload, true)
	if err != nil {
		return
	}

	var response interface{}
	if err = json.Unmarshal(r, &response); err != nil {
		return
	}
	if err = handleErr(response); err != nil {
		return
	}

	var rawRes rawSymbols
	err = json.Unmarshal(r, &rawRes)
	symbols = rawRes.Data
	return
}

// GetSymbol is used to get the open and available trading market at Kucoin along with other meta data.
// Trading symbol e.g. KCS-BTC. If not specified then you will get data of all symbols.
func (k *Kucoin) GetSymbol(s string) (symbol Symbol, err error) {
	if len(s) < 1 {
		return symbol, ErrSymbolRequired
	}
	if !k.containsCoinsPairs(strings.ToUpper(s)) {
		return symbol, ErrNonExistingSymbol
	}
	payload := map[string]string{
		"symbol": strings.ToUpper(s),
	}

	r, err := k.client.do("GET", "open/tick", payload, false)
	if err != nil {
		return
	}
	var response interface{}
	if err = json.Unmarshal(r, &response); err != nil {
		return
	}
	if err = handleErr(response); err != nil {
		return
	}
	var rawRes rawSymbol
	err = json.Unmarshal(r, &rawRes)
	symbol = rawRes.Data
	return
}

// GetOpenMarkets is used to get all open markets.
func (k *Kucoin) GetOpenMarkets() (markets []string, err error) {
	r, err := k.client.do("GET", "open/markets", nil, false)
	if err != nil {
		return
	}
	var response interface{}
	if err = json.Unmarshal(r, &response); err != nil {
		return
	}
	if err = handleErr(response); err != nil {
		return
	}
	var rawRes rawMarket
	err = json.Unmarshal(r, &rawRes)
	markets = rawRes.Data
	return
}

// GetCoins is used to get all open and available trading coins at Kucoin along with other meta data.
func (k *Kucoin) GetCoins() (coins []Coin, err error) {
	r, err := k.client.do("GET", "market/open/coins", nil, false)
	if err != nil {
		return
	}
	var response interface{}
	if err = json.Unmarshal(r, &response); err != nil {
		return
	}
	if err = handleErr(response); err != nil {
		return
	}
	var rawRes rawCoins
	err = json.Unmarshal(r, &rawRes)
	coins = rawRes.Data
	return
}

// GetCoin is used to get the open and available trading coin at Kucoin along with other meta data.
// Example:
// - Coin (required) = BTC
func (k *Kucoin) GetCoin(c string) (coin Coin, err error) {
	if len(c) < 1 {
		return coin, ErrSymbolRequired
	}
	if !k.containsOpenMarkets(strings.ToUpper(c)) {
		return coin, ErrNonExistingMarket
	}
	payload := map[string]string{
		"coin": strings.ToUpper(c),
	}

	r, err := k.client.do("GET", "market/open/coin-info", payload, false)
	if err != nil {
		return
	}
	var response interface{}
	if err = json.Unmarshal(r, &response); err != nil {
		return
	}
	if err = handleErr(response); err != nil {
		return
	}
	var rawRes rawCoin
	err = json.Unmarshal(r, &rawRes)
	coin = rawRes.Data
	return
}

// GetCoinBalance is used to get the balance at chosen coin at Kucoin along with other meta data.
// Example:
// - Coin (required) = BTC
func (k *Kucoin) GetCoinBalance(coin string) (coinBalance CoinBalance, err error) {
	if len(coin) < 1 {
		return coinBalance, ErrSymbolRequired
	}
	if !k.containsOpenMarkets(strings.ToUpper(coin)) {
		return coinBalance, ErrNonExistingMarket
	}

	r, err := k.client.do("GET", fmt.Sprintf("account/%s/balance", strings.ToUpper(coin)), nil, true)
	if err != nil {
		return
	}
	var response interface{}
	if err = json.Unmarshal(r, &response); err != nil {
		return
	}
	if err = handleErr(response); err != nil {
		return
	}
	var rawRes rawCoinBalance
	err = json.Unmarshal(r, &rawRes)
	coinBalance = rawRes.Data
	return
}

// GetCoinDepositAddress is used to get the address at chosen coin at Kucoin along with other meta data.
func (k *Kucoin) GetCoinDepositAddress(coin string) (coinDepositAddress CoinDepositAddress, err error) {
	if len(coin) < 1 {
		return coinDepositAddress, ErrSymbolRequired
	}
	if !k.containsOpenMarkets(strings.ToUpper(coin)) {
		return coinDepositAddress, ErrNonExistingMarket
	}

	r, err := k.client.do("GET", fmt.Sprintf("account/%s/wallet/address", strings.ToUpper(coin)), nil, true)
	if err != nil {
		return
	}
	var response interface{}
	if err = json.Unmarshal(r, &response); err != nil {
		return
	}
	if err = handleErr(response); err != nil {
		return
	}
	var rawRes rawCoinDepositAddress
	err = json.Unmarshal(r, &rawRes)
	coinDepositAddress = rawRes.Data
	return
}

// ListActiveMapOrders is used to get the information about active orders in user-friendly view
// at Kucoin along with other meta data.
// Example:
// - Symbol (required) = KCS-BTC
// - Type = BUY | SELL
func (k *Kucoin) ListActiveMapOrders(symbol, side string) (activeMapOrders ActiveMapOrder, err error) {
	if len(symbol) < 1 {
		return activeMapOrders, ErrSymbolRequired
	}
	if !k.containsCoinsPairs(strings.ToUpper(symbol)) {
		return activeMapOrders, ErrNonExistingSymbol
	}
	payload := make(map[string]string)
	payload["symbol"] = strings.ToUpper(symbol)
	if len(side) > 1 {
		if side != "BUY" && side != "SELL" {
			return activeMapOrders, fmt.Errorf(defaultMessageWrongInput, strings.Join([]string{"BUY", "SELL"}, ","))
		}
		payload["type"] = strings.ToUpper(side)
	}

	r, err := k.client.do("GET", "order/active-map", payload, true)
	if err != nil {
		return
	}
	var response interface{}
	if err = json.Unmarshal(r, &response); err != nil {
		return
	}
	if err = handleErr(response); err != nil {
		return
	}
	var rawRes rawActiveMapOrder
	err = json.Unmarshal(r, &rawRes)
	activeMapOrders = rawRes.Data
	return
}

// ListActiveOrders is used to get the information about active orders in array mode
// at Kucoin along with other meta data.
// Example:
// - Symbol (required) = KCS-BTC
// - Type = BUY | SELL
func (k *Kucoin) ListActiveOrders(symbol, side string) (activeOrders ActiveOrder, err error) {
	if len(symbol) < 1 {
		return activeOrders, ErrSymbolRequired
	}
	if !k.containsCoinsPairs(strings.ToUpper(symbol)) {
		return activeOrders, ErrNonExistingSymbol
	}
	payload := make(map[string]string)
	payload["symbol"] = strings.ToUpper(symbol)
	if len(side) > 1 {
		if side != "BUY" && side != "SELL" {
			return activeOrders, fmt.Errorf(defaultMessageWrongInput, strings.Join([]string{"BUY", "SELL"}, ","))
		}
		payload["type"] = strings.ToUpper(side)
	}

	r, err := k.client.do("GET", "order/active", payload, true)
	if err != nil {
		return
	}

	fmt.Println(string(r))

	var response interface{}
	if err = json.Unmarshal(r, &response); err != nil {
		return
	}
	if err = handleErr(response); err != nil {
		return
	}
	var rawRes rawActiveOrder
	err = json.Unmarshal(r, &rawRes)
	activeOrders = rawRes.Data
	return
}

// OrdersBook is used to get the information about active orders at Kucoin along with other meta data.
// Example:
// - Symbol (required) = KCS-BTC
// - Group
// - Limit
// - Direction = BUY | SELL
func (k *Kucoin) OrdersBook(symbol string, group, limit int, direction string) (ordersBook OrdersBook, err error) {
	if len(symbol) < 1 {
		return ordersBook, ErrSymbolRequired
	}
	if !k.containsCoinsPairs(strings.ToUpper(symbol)) {
		return ordersBook, ErrNonExistingSymbol
	}
	payload := make(map[string]string)
	payload["symbol"] = strings.ToUpper(symbol)
	if len(direction) > 1 {
		if direction != "BUY" && direction != "SELL" {
			return ordersBook, fmt.Errorf(defaultMessageWrongInput, strings.Join([]string{"BUY", "SELL"}, ","))
		}
		payload["direction"] = strings.ToUpper(direction)
	}
	if group > 0 {
		payload["group"] = fmt.Sprintf("%v", group)
	}
	if limit == 0 {
		payload["limit"] = fmt.Sprintf("%v", 1000)
	} else {
		payload["limit"] = fmt.Sprintf("%v", limit)
	}

	r, err := k.client.do("GET", "open/orders", payload, true)
	if err != nil {
		return
	}
	var response interface{}
	if err = json.Unmarshal(r, &response); err != nil {
		return
	}
	if err = handleErr(response); err != nil {
		return
	}
	var rawRes rawOrdersBook
	err = json.Unmarshal(r, &rawRes)
	ordersBook = rawRes.Data
	return
}

// CreateOrder is used to create order at Kucoin along with other meta data.
// Example:
// - Symbol (required) = KCS-BTC
// - Side (required) = BUY | SELL
// - Price (required) = 0.0001700
// - Amount (required) = 1.5
func (k *Kucoin) CreateOrder(symbol, side string, price, amount float64) (orderOid string, err error) {
	if len(symbol) < 1 || len(side) < 1 || price <= 0.0 || amount <= 0.0 {
		return orderOid, ErrAllParamsRequired
	}
	if !k.containsCoinsPairs(strings.ToUpper(symbol)) {
		return orderOid, ErrNonExistingSymbol
	}
	if side != "BUY" && side != "SELL" {
		return orderOid, fmt.Errorf(defaultMessageWrongInput, strings.Join([]string{"BUY", "SELL"}, ","))
	}
	payload := map[string]string{
		"symbol": strings.ToUpper(symbol),
		"amount": strconv.FormatFloat(amount, 'f', 8, 64),
		"price":  strconv.FormatFloat(price, 'f', 8, 64),
		"type":   strings.ToUpper(side),
	}

	r, err := k.client.do("POST", "order", payload, true)
	if err != nil {
		return
	}
	var response interface{}
	if err = json.Unmarshal(r, &response); err != nil {
		return
	}
	if err = handleErr(response); err != nil {
		return
	}
	var rawRes rawOrder
	err = json.Unmarshal(r, &rawRes)
	if err != nil {
		return
	}
	if !rawRes.Success {
		err = errors.New(string(r))
		return
	}
	orderOid = rawRes.Data.OrderOid
	return
}

// CreateOrderByString is used to create order at Kucoin along with other meta data.
// This ByString version is fix precise problem.
// Example:
// - Symbol (required) = KCS-BTC
// - Side (required) = BUY | SELL
// - Price (required) = 0.0001700
// - Amount (required) = 1.5
func (k *Kucoin) CreateOrderByString(symbol, side, price, amount string) (orderOid string, err error) {
	if len(symbol) < 1 || len(side) < 1 || len(price) < 1 || len(amount) < 1 {
		return orderOid, ErrAllParamsRequired
	}
	if !k.containsCoinsPairs(strings.ToUpper(symbol)) {
		return orderOid, ErrNonExistingSymbol
	}
	if side != "BUY" && side != "SELL" {
		return orderOid, fmt.Errorf(defaultMessageWrongInput, strings.Join([]string{"BUY", "SELL"}, ","))
	}
	payload := map[string]string{
		"symbol": strings.ToUpper(symbol),
		"amount": amount,
		"price":  price,
		"type":   strings.ToUpper(side),
	}

	r, err := k.client.do("POST", "order", payload, true)
	if err != nil {
		return
	}
	var response interface{}
	if err = json.Unmarshal(r, &response); err != nil {
		return
	}
	if err = handleErr(response); err != nil {
		return
	}
	var rawRes rawOrder
	err = json.Unmarshal(r, &rawRes)
	if err != nil {
		return
	}
	if !rawRes.Success {
		err = errors.New(string(r))
		return
	}
	orderOid = rawRes.Data.OrderOid
	return
}

// AccountHistory is used to get the information about list deposit & withdrawal
// at Kucoin along with other meta data. Coin, Side (type in Kucoin docs.)
// and Status are required parameters. Limit and page may be zeros.
// Example:
// - Coin (required) = KCS
// - Side (required) = DEPOSIT | WITHDRAW
// - Status (required) = FINISHED | CANCEL | PENDING
// - Page
func (k *Kucoin) AccountHistory(coin, side, status string, page int) (accountHistory AccountHistory, err error) {
	if len(coin) < 1 || len(side) < 1 || len(status) < 1 {
		return accountHistory, ErrAllParamsRequired
	}
	if !k.containsOpenMarkets(strings.ToUpper(coin)) {
		return accountHistory, ErrNonExistingMarket
	}
	if side != "DEPOSIT" && side != "WITHDRAW" {
		return accountHistory, fmt.Errorf(defaultMessageWrongInput, strings.Join([]string{"DEPOSIT", "WITHDRAW"}, ","))
	}
	if status != "FINISHED" && status != "CANCEL" && status != "PENDING" {
		return accountHistory, fmt.Errorf(defaultMessageWrongInput, strings.Join([]string{"FINISHED", "CANCEL", "PENDING"}, ","))
	}

	payload := map[string]string{
		"coin":   strings.ToUpper(coin),
		"type":   strings.ToUpper(side),
		"status": strings.ToUpper(status),
	}
	if page != 0 {
		payload["page"] = fmt.Sprintf("%v", page)
	}

	r, err := k.client.do("GET", fmt.Sprintf("account/%s/wallet/records", strings.ToUpper(coin)), payload, true)
	if err != nil {
		return
	}
	var response interface{}
	if err = json.Unmarshal(r, &response); err != nil {
		return
	}
	if err = handleErr(response); err != nil {
		return
	}
	var rawRes rawAccountHistory
	err = json.Unmarshal(r, &rawRes)
	accountHistory = rawRes.Data
	return
}

// ListSpecificDealtOrders is used to get the information about dealt orders for specific symbol at Kucoin along with other meta data.
// Symbol, Side (type in Kucoin docs.) are required parameters. Limit and page may be zeros.
// Example:
// - Symbol (required) = KCS-BTC
// - Side = BUY | SELL
// - Limit
// - Page
func (k *Kucoin) ListSpecificDealtOrders(symbol, side string, limit, page int) (specificDealtOrders SpecificDealtOrder, err error) {
	if len(symbol) < 1 {
		return specificDealtOrders, ErrSymbolRequired
	}
	if !k.containsCoinsPairs(strings.ToUpper(symbol)) {
		return specificDealtOrders, ErrNonExistingSymbol
	}
	payload := make(map[string]string)
	payload["symbol"] = strings.ToUpper(symbol)
	if len(side) > 1 {
		if side != "BUY" && side != "SELL" {
			return specificDealtOrders, fmt.Errorf(defaultMessageWrongInput, strings.Join([]string{"BUY", "SELL"}, ","))
		}
		payload["type"] = strings.ToUpper(side)
	}
	if limit == 0 {
		payload["limit"] = fmt.Sprintf("%v", 1000)
	} else {
		payload["limit"] = fmt.Sprintf("%v", limit)
	}
	if page != 0 {
		payload["page"] = fmt.Sprintf("%v", page)
	}

	r, err := k.client.do("GET", "deal-orders", payload, true)
	if err != nil {
		return
	}
	var response interface{}
	if err = json.Unmarshal(r, &response); err != nil {
		return
	}
	if err = handleErr(response); err != nil {
		return
	}
	var rawRes rawSpecificDealtOrder
	err = json.Unmarshal(r, &rawRes)
	specificDealtOrders = rawRes.Data
	return
}

// ListMergedDealtOrders is used to get the information about dealt orders for
// all symbols at Kucoin along with other meta data.
// All parameters are optional. Timestamp must be in milliseconds from Unix epoch.
// Example:
// - Symbol = KCS-BTC
// - Side = BUY | SELL
// - Limit
// - Page
// - Since
// - Before
func (k *Kucoin) ListMergedDealtOrders(symbol, side string, limit, page int, since, before int64) (mergedDealtOrders MergedDealtOrder, err error) {
	payload := make(map[string]string)
	if len(symbol) > 1 {
		if !k.containsCoinsPairs(strings.ToUpper(symbol)) {
			return mergedDealtOrders, ErrNonExistingSymbol
		}
		payload["symbol"] = strings.ToUpper(symbol)
	}
	if len(side) > 1 {
		if side != "BUY" && side != "SELL" {
			return mergedDealtOrders, fmt.Errorf(defaultMessageWrongInput, strings.Join([]string{"BUY", "SELL"}, ","))
		}
		payload["type"] = strings.ToUpper(side)
	}
	if (limit == 0 || limit > 100) && len(symbol) > 1 {
		payload["limit"] = fmt.Sprintf("%v", 100)
	} else if (limit == 0 || limit > 20) && len(symbol) < 1 {
		payload["limit"] = fmt.Sprintf("%v", 20)
	} else {
		payload["limit"] = fmt.Sprintf("%v", limit)
	}
	if page != 0 {
		payload["page"] = fmt.Sprintf("%v", page)
	}
	if since != 0 {
		payload["since"] = fmt.Sprintf("%v", since)
	}
	if before != 0 {
		payload["before"] = fmt.Sprintf("%v", before)
	}

	r, err := k.client.do("GET", "order/dealt", payload, true)
	if err != nil {
		return
	}
	var response interface{}
	if err = json.Unmarshal(r, &response); err != nil {
		return
	}
	if err = handleErr(response); err != nil {
		return
	}
	var rawRes rawMergedDealtOrder
	err = json.Unmarshal(r, &rawRes)
	mergedDealtOrders = rawRes.Data
	return
}

// OrderDetails is used to get the information about orders for specific symbol at Kucoin along with other meta data.
// Symbol, Side (type in Kucoin docs.) are required parameters.
// Limit may be zero, and not greater than 20. Page may be zero and by default is equal to 1.
// Example:
// - Symbol (required) = KCS-BTC
// - Side (required) = BUY | SELL
// - OrderOid (required)
// - Limit
// - Page
func (k *Kucoin) OrderDetails(symbol, side, orderOid string, limit, page int) (orderDetails OrderDetails, err error) {
	if len(symbol) < 1 || len(side) < 1 || len(orderOid) < 1 {
		return orderDetails, ErrAllParamsRequired
	}
	if !k.containsCoinsPairs(strings.ToUpper(symbol)) {
		return orderDetails, ErrNonExistingSymbol
	}
	if side != "BUY" && side != "SELL" {
		return orderDetails, fmt.Errorf(defaultMessageWrongInput, strings.Join([]string{"BUY", "SELL"}, ","))
	}
	payload := map[string]string{
		"symbol":   strings.ToUpper(symbol),
		"type":     strings.ToUpper(side),
		"orderOid": strings.ToUpper(orderOid),
	}
	if limit == 0 {
		payload["limit"] = fmt.Sprintf("%v", 20)
	} else {
		payload["limit"] = fmt.Sprintf("%v", limit)
	}
	if page == 0 {
		payload["page"] = fmt.Sprintf("%v", 1)
	} else {
		payload["page"] = fmt.Sprintf("%v", page)
	}

	r, err := k.client.do("GET", "order/detail", payload, true)
	if err != nil {
		return
	}
	var response interface{}
	if err = json.Unmarshal(r, &response); err != nil {
		return
	}
	if err = handleErr(response); err != nil {
		return
	}
	var rawRes rawOrderDetails
	err = json.Unmarshal(r, &rawRes)
	orderDetails = rawRes.Data
	return
}

// CreateWithdrawalApply is used to create withdrawal for specific coin
// at Kucoin along with other meta data.
// Example:
// - Coin (required) = KCS
// - Address (required) =
// - Amount (required) = 0.50
// Result:
// - Nothing.
func (k *Kucoin) CreateWithdrawalApply(coin, address string, amount float64) (withdrawalApply Withdrawal, err error) {
	if len(coin) < 1 || len(address) < 1 || amount <= 0.0 {
		return withdrawalApply, ErrAllParamsRequired
	}
	if !k.containsOpenMarkets(strings.ToUpper(coin)) {
		return withdrawalApply, ErrNonExistingMarket
	}
	payload := map[string]string{
		"address": address,
		"amount":  fmt.Sprintf("%v", amount),
	}

	r, err := k.client.do("POST", fmt.Sprintf("account/%s/withdraw/apply", strings.ToUpper(coin)), payload, true)
	if err != nil {
		return
	}
	var response interface{}
	if err = json.Unmarshal(r, &response); err != nil {
		return
	}
	if err = handleErr(response); err != nil {
		return
	}
	var rawRes rawWithdrawal
	err = json.Unmarshal(r, &rawRes)
	withdrawalApply = rawRes.Data
	return
}

// CancelWithdrawal used to cancel withdrawal for specific coin
// at Kucoin along with other meta data.
// Example:
// - Coin (required) = KCS
// - TxOid (required)
// Result:
// - Nothing.
func (k *Kucoin) CancelWithdrawal(coin, txOid string) (withdrawal Withdrawal, err error) {
	if len(coin) < 1 || len(txOid) < 1 {
		return withdrawal, ErrAllParamsRequired
	}
	if !k.containsOpenMarkets(strings.ToUpper(coin)) {
		return withdrawal, ErrNonExistingMarket
	}
	payload := map[string]string{
		"txOid": txOid,
	}

	r, err := k.client.do("POST", fmt.Sprintf("account/%s/withdraw/cancel", strings.ToUpper(coin)), payload, true)
	if err != nil {
		return
	}
	var response interface{}
	if err = json.Unmarshal(r, &response); err != nil {
		return
	}
	if err = handleErr(response); err != nil {
		return
	}
	var rawRes rawWithdrawal
	err = json.Unmarshal(r, &rawRes)
	withdrawal = rawRes.Data
	return
}

// CancelOrder is used to cancel execution of current order at Kucoin along with other meta data.
// Example:
// - Symbol (required) = KCS-BTC
// - OrderId (required)
// - Side (required) = BUY | SELL
func (k *Kucoin) CancelOrder(symbol, orderOid, side string) error {
	if len(symbol) < 1 || len(side) < 1 || len(orderOid) < 1 {
		return ErrAllParamsRequired
	}
	if !k.containsCoinsPairs(strings.ToUpper(symbol)) {
		return ErrNonExistingSymbol
	}
	if side != "BUY" && side != "SELL" {
		return fmt.Errorf(defaultMessageWrongInput, strings.Join([]string{"BUY", "SELL"}, ","))
	}
	payload := map[string]string{
		"symbol":   strings.ToUpper(symbol),
		"orderOid": orderOid,
		"type":     strings.ToUpper(side),
	}

	r, err := k.client.do("POST", "cancel-order", payload, true)
	if err != nil {
		return err
	}
	var response interface{}
	if err = json.Unmarshal(r, &response); err != nil {
		return err
	}
	return handleErr(response)
}

// CancelAllOrders is used to cancel execution of all orders at Kucoin along with other meta data.
// Example:
// - Symbol (required) = KCS-BTC
// - Side = BUY | SELL
func (k *Kucoin) CancelAllOrders(symbol, side string) error {
	if len(symbol) < 1 {
		return ErrSymbolRequired
	}
	if !k.containsCoinsPairs(strings.ToUpper(symbol)) {
		return ErrNonExistingSymbol
	}
	payload := make(map[string]string)
	payload["symbol"] = strings.ToUpper(symbol)
	if len(side) > 1 {
		if side != "BUY" && side != "SELL" {
			return fmt.Errorf(defaultMessageWrongInput, strings.Join([]string{"BUY", "SELL"}, ","))
		}
		payload["type"] = strings.ToUpper(side)
	}

	r, err := k.client.do("POST", "order/cancel-all", payload, true)
	if err != nil {
		return err
	}
	var response interface{}
	if err = json.Unmarshal(r, &response); err != nil {
		return err
	}
	return handleErr(response)
}
