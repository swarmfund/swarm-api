package resource

import (
	horizon "gitlab.com/swarmfund/horizon-connector"
)

type Balance struct {
	BalanceID  string `json:"balance_id"`
	ExchangeID string `json:"exchange_id"`
	Asset      string `json:"asset"`
}

func (b *Balance) Populate(h *horizon.Balance) {
	b.BalanceID = h.BalanceID
	b.ExchangeID = h.ExchangeID
	b.Asset = h.Asset
}

type Balances struct {
	Records []Balance `json:"records"`
}

func (b *Balances) Populate(records []horizon.Balance) {
	b.Records = []Balance{}
	for _, record := range records {
		var r Balance
		r.Populate(&record)
		b.Records = append(b.Records, r)
	}
}
