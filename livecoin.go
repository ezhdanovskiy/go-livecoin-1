// Package Livecoin is an implementation of the Livecoin API in Golang.
package livecoin

import (
	"encoding/json"
	"errors"
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

// livecoin represent a livecoin client
type Livecoin struct {
	client *client
}

// set enable/disable http request/response dump
func (c *Livecoin) SetDebug(enable bool) {
	c.client.debug = enable
}

// Account

// GetBalances is used to retrieve all balances from your account
func (b *Livecoin) GetBalances() (balances []Balance, err error) {
	r, err := b.client.do("GET", "payment/balances", nil, true)
	if err != nil {
		return
	}
	//var response json.RawMessage
	//if err = json.Unmarshal(r, &response); err != nil {
	//	return
	//}
	//if err = handleErr(response); err != nil {
	//	return
	//}
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
	//var response jsonResponse
	//if err = json.Unmarshal(r, &response); err != nil {
	//	return
	//}
	//if err = handleErr(response); err != nil {
	//	return
	//}
	err = json.Unmarshal(r, &balance)
	return
}

// GetTransactions is used to retrieve your withdrawal and deposit history
// "Start" and "end" are given in UNIX timestamp format in miliseconds and used to specify the date range for the data returned.
func (b *Livecoin) GetTransactions(start uint64, end uint64) (transactions []Transaction, err error) {
	ressource := "payment/history/transactions"
	if end == 0 {
		end = uint64(time.Now().Unix()) * 1000
	}
	r, err := b.client.do("GET", ressource, map[string]string{"types": "DEPOSIT,WITHDRAWAL", "start": strconv.FormatUint(uint64(start), 10), "end": strconv.FormatUint(uint64(end), 10)}, true)
	if err != nil {
		return
	}
	//var response jsonResponse
	//if err = json.Unmarshal(r, &response); err != nil {
	//	return
	//}
	//if err = handleErr(response); err != nil {
	//	return
	//}
	err = json.Unmarshal(r, &transactions)
	return
}
