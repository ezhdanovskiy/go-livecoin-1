package livecoin

import (
	"encoding/json"
	"time"
)

type Order struct {
	Id                   int64     `json:"id"`
	CurrencyPair         string    `json:"currencyPair"`
	GoodUntilTime        time.Time `json:"goodUntilTime"`
	Type                 string    `json:"type"`
	OrderStatus          string    `json:"orderStatus"`
	IssueTime            time.Time `json:"issueTime"`
	Price                float64   `json:"price"`
	Quantity             float64   `json:"quantity"`
	RemainingQuantity    float64   `json:"remainingQuantity"`
	Commission           float64   `json:"commission"`
	CommissionRate       float64   `json:"commissionRate"`
	LastModificationTime time.Time `json:"lastModificationTime"`
}

type ClientOrders struct {
	TotalRows uint64  `json:"totalRows"`
	StartRow  uint64  `json:"startRow"`
	EndRow    uint64  `json:"endRow"`
	Data      []Order `json:"data"`
}

func (o *Order) UnmarshalJSON(data []byte) error {
	var err error
	type Alias Order
	aux := &struct {
		GoodUntilTime        int64 `json:"goodUntilTime"`
		IssueTime            int64 `json:"issueTime"`
		LastModificationTime int64 `json:"lastModificationTime"`
		*Alias
	}{
		Alias: (*Alias)(o),
	}
	if err = json.Unmarshal(data, &aux); err != nil {
		return err
	}
	o.GoodUntilTime = time.Unix(0, aux.GoodUntilTime * int64(time.Millisecond))
	o.IssueTime = time.Unix(0, aux.IssueTime * int64(time.Millisecond))
	o.LastModificationTime = time.Unix(0, aux.LastModificationTime * int64(time.Millisecond))
	return nil
}
