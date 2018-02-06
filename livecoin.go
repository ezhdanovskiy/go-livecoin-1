// Package Livecoin is an implementation of the Livecoin API in Golang.
package livecoin

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
	"strconv"
	"strings"
)

const (
	API_BASE    = "https://api.livecoin.net" // Livecoin API endpoint
)

// New returns an instantiated livecoin struct
func New(apiKey, apiSecret string) *Livecoin {
	client := NewClient(apiKey, apiSecret)
	return &Livecoin{client}
}

// NewWithCustomHttpClient returns an instantiated livecoin struct with custom http client
func NewWithCustomHttpClient(apiKey, apiSecret string, httpClient *http.Client) *Livecoin {
	client := NewClientWithCustomHttpConfig(apiKey, apiSecret, httpClient)
	return &Livecoin{client}
}

// NewWithCustomTimeout returns an instantiated livecoin struct with custom timeout
func NewWithCustomTimeout(apiKey, apiSecret string, timeout time.Duration) *Livecoin {
	client := NewClientWithCustomTimeout(apiKey, apiSecret, timeout)
	return &Livecoin{client}
}

// handleErr gets JSON response from livecoin API en deal with error
func handleErr(r interface{}) error {
	switch v := r.(type) {
	case map[string]interface{}:
		errorMessage := r.(map[string]interface{})["errorMessage"]
		if errorMessage != nil && errorMessage.(string) != "" {
			return errors.New(errorMessage.(string))
		}
	case []interface{}:
		return nil
	default:
		return fmt.Errorf("I don't know about type %T!\n", v)
	}

	return nil
}

// livecoin represent a livecoin client
type Livecoin struct {
	client *client
}

// set enable/disable http request/response dump
func (b *Livecoin) SetDebug(enable bool) {
	b.client.debug = enable
}

// Account

// GetBalances is used to retrieve all balances from your account
func (b *Livecoin) GetBalances() (balances []Balance, err error) {
	r, err := b.client.do("GET", "payment/balances", nil, true)
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
	err = json.Unmarshal(r, &balances)
	return
}

// Getbalance is used to retrieve the balance from your account for a specific currency.
// currency: a string literal for the currency (ex: LTC)
func (b *Livecoin) GetBalance(currency string) (balance Balance, err error) {
	r, err := b.client.do("GET", "payment/balance", map[string]string{"currency": strings.ToUpper(currency)}, true)
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
	err = json.Unmarshal(r, &balance)
	return
}

// GetTrades used to retrieve your trade history.
// market string literal for the market (ie. BTC/LTC). If set to "all", will return for all market
func (b *Livecoin) GetTrades(currencyPair string) (trades []Trade, err error) {
	payload := make(map[string]string)
	if currencyPair != "all" {
		payload["currencyPair"] = currencyPair
	}
	r, err := b.client.do("GET", "exchange/trades", payload, true)
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
	if res, ok := response.(map[string]interface{}); ok {
		if exception, ok := res["exception"]; ok && exception == "Data not found" {
			return
		}
	}
	err = json.Unmarshal(r, &trades)
	return
}

// GetTransactions is used to retrieve your withdrawal and deposit history
// "Start" and "end" are given in UNIX timestamp format in miliseconds and used to specify the date range for the data returned.
func (b *Livecoin) GetTransactions(start uint64, end uint64) (transactions []Transaction, err error) {
	if end == 0 {
		end = uint64(time.Now().Unix()) * 1000
	}
	r, err := b.client.do("GET", "payment/history/transactions", map[string]string{"types": "DEPOSIT,WITHDRAWAL", "start": strconv.FormatUint(uint64(start), 10), "end": strconv.FormatUint(uint64(end), 10)}, true)
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
	err = json.Unmarshal(r, &transactions)
	return
}

func (b *Livecoin) GetClientOrders(openClosed string) (clientOrders ClientOrders, err error) {
	r, err := b.client.do("GET", "exchange/client_orders", map[string]string{"openClosed": openClosed}, true)
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
	err = json.Unmarshal(r, &clientOrders)
	return
}
