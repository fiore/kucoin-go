package kucoin_test

import (
	"errors"
	"testing"

	kucoinGo "github.com/fiore/kucoin-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	// email: h2683915@nwytg.net
	// password: Test10000
	apiKey    = "5bd714589b5f4527c83152fb"
	apiSecret = "62d7852d-9515-49aa-905b-9e4d19220af9"
)

var (
	kucoin              *kucoinGo.Kucoin = kucoinGo.New(apiKey, apiSecret)
	defaultErrorMessage string           = "There should be no error"
)

func TestGetUserInfo(t *testing.T) {
	userInfo, err := kucoin.GetUserInfo()
	t.Logf("GetUserInfo : %#v\n", userInfo)
	require.NoError(t, err, defaultErrorMessage)
}

func TestGetSymbols(t *testing.T) {
	symbols, err := kucoin.GetSymbols()
	t.Logf("GetSymbols : %#v\n", symbols)
	require.NoError(t, err, defaultErrorMessage)
}

func TestGetCoinsPairs(t *testing.T) {
	coinPair, err := kucoin.GetCoinsPairs()
	t.Logf("GetCoinsPairs() : %#v\n", coinPair)
	require.NoError(t, err, defaultErrorMessage)
}

func TestGetUserSymbols(t *testing.T) {
	_, err := kucoin.GetUserSymbols("TEST", "KCS-BTC", "FAVOURITE")
	if assert.Error(t, err) {
		require.Equal(t, kucoinGo.ErrNonExistingMarket, err)
	}
	_, err = kucoin.GetUserSymbols("BTC", "TEST", "FAVOURITE")
	if assert.Error(t, err) {
		require.Equal(t, kucoinGo.ErrNonExistingSymbol, err)
	}
	_, err = kucoin.GetUserSymbols("BTC", "KCS-BTC", "TEST")
	if assert.Error(t, err) {
		require.Equal(t, errors.New("Entered invalid parameter. Accepted values: [FAVOURITE,STICK]"), err)
	}

	symbols, err := kucoin.GetUserSymbols("BTC", "KCS-BTC", "FAVOURITE")
	t.Logf("GetUserSymbols : %#v\n", symbols)
	require.NoError(t, err, defaultErrorMessage)
}

func TestGetSymbol(t *testing.T) {
	_, err := kucoin.GetSymbol("")
	if assert.Error(t, err) {
		require.Equal(t, kucoinGo.ErrSymbolRequired, err)
	}
	_, err = kucoin.GetSymbol("TEST")
	if assert.Error(t, err) {
		require.Equal(t, kucoinGo.ErrNonExistingSymbol, err)
	}

	symbol, err := kucoin.GetSymbol("KCS-BTC")
	t.Logf("GetSymbol : %#v\n", symbol)
	require.NoError(t, err, defaultErrorMessage)
}

func TestGetCoins(t *testing.T) {
	coins, err := kucoin.GetCoins()
	t.Logf("GetCoins : %#v\n", coins)
	require.NoError(t, err, defaultErrorMessage)
}

func TestGetCoin(t *testing.T) {
	_, err := kucoin.GetCoin("")
	if assert.Error(t, err) {
		require.Equal(t, kucoinGo.ErrSymbolRequired, err)
	}
	_, err = kucoin.GetCoin("TEST")
	if assert.Error(t, err) {
		require.Equal(t, kucoinGo.ErrNonExistingMarket, err)
	}

	coin, err := kucoin.GetCoin("BTC")
	t.Logf("GetCoin : %#v\n", coin)
	require.NoError(t, err, defaultErrorMessage)
}

func TestGetCoinBalance(t *testing.T) {
	_, err := kucoin.GetCoinBalance("")
	if assert.Error(t, err) {
		require.Equal(t, kucoinGo.ErrSymbolRequired, err)
	}
	_, err = kucoin.GetCoinBalance("TEST")
	if assert.Error(t, err) {
		require.Equal(t, kucoinGo.ErrNonExistingMarket, err)
	}

	coinBalance, err := kucoin.GetCoinBalance("BTC")
	t.Logf("GetCoinBalance : %#v\n", coinBalance)
	require.NoError(t, err, defaultErrorMessage)
}

func TestGetCoinDepositAddress(t *testing.T) {
	_, err := kucoin.GetCoinDepositAddress("")
	if assert.Error(t, err) {
		require.Equal(t, kucoinGo.ErrSymbolRequired, err)
	}
	_, err = kucoin.GetCoinDepositAddress("TEST")
	if assert.Error(t, err) {
		require.Equal(t, kucoinGo.ErrNonExistingMarket, err)
	}

	coinDepositAddress, err := kucoin.GetCoinDepositAddress("BTC")
	t.Logf("GetCoinDepositAddress : %#v\n", coinDepositAddress)
	require.NoError(t, err, defaultErrorMessage)
}

func TestListActiveMapOrders(t *testing.T) {
	_, err := kucoin.ListActiveMapOrders("", "")
	if assert.Error(t, err) {
		require.Equal(t, kucoinGo.ErrSymbolRequired, err)
	}
	_, err = kucoin.ListActiveMapOrders("TEST", "")
	if assert.Error(t, err) {
		require.Equal(t, kucoinGo.ErrNonExistingSymbol, err)
	}
	_, err = kucoin.ListActiveMapOrders("KCS-BTC", "TEST")
	if assert.Error(t, err) {
		require.Equal(t, errors.New("Entered invalid parameter. Accepted values: [BUY,SELL]"), err)
	}

	activeMapOrders, err := kucoin.ListActiveMapOrders("KCS-BTC", "BUY")
	t.Logf("ListActiveMapOrders : %#v\n", activeMapOrders)
	require.NoError(t, err, defaultErrorMessage)
}

func TestListActiveOrders(t *testing.T) {
	_, err := kucoin.ListActiveOrders("", "")
	if assert.Error(t, err) {
		require.Equal(t, kucoinGo.ErrSymbolRequired, err)
	}
	_, err = kucoin.ListActiveOrders("TEST", "")
	if assert.Error(t, err) {
		require.Equal(t, kucoinGo.ErrNonExistingSymbol, err)
	}
	_, err = kucoin.ListActiveOrders("KCS-BTC", "TEST")
	if assert.Error(t, err) {
		require.Equal(t, errors.New("Entered invalid parameter. Accepted values: [BUY,SELL]"), err)
	}

	activeOrders, err := kucoin.ListActiveOrders("KCS-BTC", "")
	t.Logf("ListActiveMapOrders : %#v\n", activeOrders)
	require.NoError(t, err, defaultErrorMessage)
}

func TestOrdersBook(t *testing.T) {
	_, err := kucoin.OrdersBook("", 0, 0, "")
	if assert.Error(t, err) {
		require.Equal(t, kucoinGo.ErrSymbolRequired, err)
	}
	_, err = kucoin.OrdersBook("TEST", 0, 0, "")
	if assert.Error(t, err) {
		require.Equal(t, kucoinGo.ErrNonExistingSymbol, err)
	}
	_, err = kucoin.OrdersBook("KCS-BTC", 0, 0, "TEST")
	if assert.Error(t, err) {
		require.Equal(t, errors.New("Entered invalid parameter. Accepted values: [BUY,SELL]"), err)
	}

	ordersBook, err := kucoin.OrdersBook("KCS-BTC", 0, 0, "BUY")
	t.Logf("OrdersBook : %#v\n", ordersBook)
	require.NoError(t, err, defaultErrorMessage)
}

func TestCreateOrder(t *testing.T) {
	_, err := kucoin.CreateOrder("", "", 0, 0)
	if assert.Error(t, err) {
		require.Equal(t, kucoinGo.ErrAllParamsRequired, err)
	}
	_, err = kucoin.CreateOrder("TEST", "BUY", 1, 1)
	if assert.Error(t, err) {
		require.Equal(t, kucoinGo.ErrNonExistingSymbol, err)
	}
	_, err = kucoin.CreateOrder("KCS-BTC", "TEST", 1, 1)
	if assert.Error(t, err) {
		require.Equal(t, errors.New("Entered invalid parameter. Accepted values: [BUY,SELL]"), err)
	}

	/*	Use the code below only with a real account	*/

	// 	orderOid, err = kucoin.CreateOrder("KCS-BTC", "BUY", 0.0001700, 1.5)
	// 	t.Logf("CreateOrder : %#v\n", orderOid)
	// 	require.NoError(t, err, defaultErrorMessage)
}

func TestCreateOrderByString(t *testing.T) {
	_, err := kucoin.CreateOrderByString("", "", "", "")
	if assert.Error(t, err) {
		require.Equal(t, kucoinGo.ErrAllParamsRequired, err)
	}
	_, err = kucoin.CreateOrderByString("TEST", "BUY", "1", "1")
	if assert.Error(t, err) {
		require.Equal(t, kucoinGo.ErrNonExistingSymbol, err)
	}
	_, err = kucoin.CreateOrderByString("KCS-BTC", "TEST", "1", "1")
	if assert.Error(t, err) {
		require.Equal(t, errors.New("Entered invalid parameter. Accepted values: [BUY,SELL]"), err)
	}
	/*	Use the code below only with a real account	*/

	// orderOid, err = kucoin.CreateOrderByString("KCS-BTC", "BUY", "0.0001700", "1.5")
	// t.Logf("CreateOrderByString : %#v\n", orderOid)
	// require.NoError(t, err, defaultErrorMessage)
}

func TestAccountHistory(t *testing.T) {
	_, err := kucoin.AccountHistory("", "", "", 0)
	if assert.Error(t, err) {
		require.Equal(t, kucoinGo.ErrAllParamsRequired, err)
	}
	_, err = kucoin.AccountHistory("TEST", "DEPOSIT", "FINISHED", 0)
	if assert.Error(t, err) {
		require.Equal(t, kucoinGo.ErrNonExistingMarket, err)
	}
	_, err = kucoin.AccountHistory("KCS", "TEST", "FINISHED", 0)
	if assert.Error(t, err) {
		require.Equal(t, errors.New("Entered invalid parameter. Accepted values: [DEPOSIT,WITHDRAW]"), err)
	}
	_, err = kucoin.AccountHistory("KCS", "DEPOSIT", "TEST", 0)
	if assert.Error(t, err) {
		require.Equal(t, errors.New("Entered invalid parameter. Accepted values: [FINISHED,CANCEL,PENDING]"), err)
	}

	/*	Use the code below only with a real account	*/

	// accountHistory, err = kucoin.AccountHistory("KCS-BTC", "DEPOSIT", "FINISHED", 0)
	// t.Logf("AccountHistory : %#v\n", accountHistory)
	// require.NoError(t, err, defaultErrorMessage)
}

func TestListSpecificDealtOrders(t *testing.T) {
	_, err := kucoin.ListSpecificDealtOrders("", "", 0, 0)
	if assert.Error(t, err) {
		require.Equal(t, kucoinGo.ErrSymbolRequired, err)
	}
	_, err = kucoin.ListSpecificDealtOrders("TEST", "", 0, 0)
	if assert.Error(t, err) {
		require.Equal(t, kucoinGo.ErrNonExistingSymbol, err)
	}
	_, err = kucoin.ListSpecificDealtOrders("KCS-BTC", "TEST", 0, 0)
	if assert.Error(t, err) {
		require.Equal(t, errors.New("Entered invalid parameter. Accepted values: [BUY,SELL]"), err)
	}

	specificDealtOrders, err := kucoin.ListSpecificDealtOrders("KCS-BTC", "", 0, 0)
	t.Logf("ListSpecificDealtOrders : %#v\n", specificDealtOrders)
	require.NoError(t, err, defaultErrorMessage)
}

func TestListMergedDealtOrders(t *testing.T) {
	_, err := kucoin.ListMergedDealtOrders("TEST", "", 0, 0, 0, 0)
	if assert.Error(t, err) {
		require.Equal(t, kucoinGo.ErrNonExistingSymbol, err)
	}
	_, err = kucoin.ListMergedDealtOrders("KCS-BTC", "TEST", 0, 0, 0, 0)
	if assert.Error(t, err) {
		require.Equal(t, errors.New("Entered invalid parameter. Accepted values: [BUY,SELL]"), err)
	}

	mergedDealtOrders, err := kucoin.ListMergedDealtOrders("", "", 0, 0, 0, 0)
	t.Logf("ListMergedDealtOrders : %#v\n", mergedDealtOrders)
	require.NoError(t, err, defaultErrorMessage)
}

func TestOrderDetails(t *testing.T) {
	_, err := kucoin.OrderDetails("", "", "", 0, 0)
	if assert.Error(t, err) {
		require.Equal(t, kucoinGo.ErrAllParamsRequired, err)
	}
	_, err = kucoin.OrderDetails("TEST", "BUY", "5969ddc96732d54312eb960e", 0, 0)
	if assert.Error(t, err) {
		require.Equal(t, kucoinGo.ErrNonExistingSymbol, err)
	}
	_, err = kucoin.OrderDetails("KCS-BTC", "TEST", "5969ddc96732d54312eb960e", 0, 0)
	if assert.Error(t, err) {
		require.Equal(t, errors.New("Entered invalid parameter. Accepted values: [BUY,SELL]"), err)
	}

	orderDetails, err := kucoin.OrderDetails("KCS-BTC", "BUY", "5969ddc96732d54312eb960e", 0, 0)
	t.Logf("OrderDetails : %#v\n", orderDetails)
	require.NoError(t, err, defaultErrorMessage)
}

func TestCreateWithdrawalApply(t *testing.T) {
	_, err := kucoin.CreateWithdrawalApply("", "", 0)
	if assert.Error(t, err) {
		require.Equal(t, kucoinGo.ErrAllParamsRequired, err)
	}
	_, err = kucoin.CreateWithdrawalApply("TEST", "5969ddc96732d54312eb960e", 1)
	if assert.Error(t, err) {
		require.Equal(t, kucoinGo.ErrNonExistingMarket, err)
	}

	withdrawalApply, err := kucoin.CreateWithdrawalApply("BTC", "5969ddc96732d54312eb960e", 1)
	t.Logf("CreateWithdrawalApply : %#v\n", withdrawalApply)
	require.NoError(t, err, defaultErrorMessage)
}

func TestCancelWithdrawal(t *testing.T) {
	_, err := kucoin.CancelWithdrawal("", "")
	if assert.Error(t, err) {
		require.Equal(t, kucoinGo.ErrAllParamsRequired, err)
	}
	_, err = kucoin.CancelWithdrawal("TEST", "5969ddc96732d54312eb960e")
	if assert.Error(t, err) {
		require.Equal(t, kucoinGo.ErrNonExistingMarket, err)
	}

	withdrawal, err := kucoin.CancelWithdrawal("BTC", "5969ddc96732d54312eb960e")
	t.Logf("CancelWithdrawal : %#v\n", withdrawal)
	require.NoError(t, err, defaultErrorMessage)
}

func TestCancelOrder(t *testing.T) {
	err := kucoin.CancelOrder("", "", "")
	if assert.Error(t, err) {
		require.Equal(t, kucoinGo.ErrAllParamsRequired, err)
	}
	err = kucoin.CancelOrder("TEST", "5969ddc96732d54312eb960e", "BUY")
	if assert.Error(t, err) {
		require.Equal(t, kucoinGo.ErrNonExistingSymbol, err)
	}
	err = kucoin.CancelOrder("KCS-BTC", "5969ddc96732d54312eb960e", "TEST")
	if assert.Error(t, err) {
		require.Equal(t, errors.New("Entered invalid parameter. Accepted values: [BUY,SELL]"), err)
	}

	err = kucoin.CancelOrder("KCS-BTC", "5969ddc96732d54312eb960e", "BUY")
	require.NoError(t, err, defaultErrorMessage)
}

func TestCancelAllOrders(t *testing.T) {
	err := kucoin.CancelAllOrders("", "")
	if assert.Error(t, err) {
		require.Equal(t, kucoinGo.ErrSymbolRequired, err)
	}
	err = kucoin.CancelAllOrders("TEST", "BUY")
	if assert.Error(t, err) {
		require.Equal(t, kucoinGo.ErrNonExistingSymbol, err)
	}
	err = kucoin.CancelAllOrders("KCS-BTC", "TEST")
	if assert.Error(t, err) {
		require.Equal(t, errors.New("Entered invalid parameter. Accepted values: [BUY,SELL]"), err)
	}

	err = kucoin.CancelAllOrders("KCS-BTC", "BUY")
	require.NoError(t, err, defaultErrorMessage)
}
