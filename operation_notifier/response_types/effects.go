package response_types

import (
	"time"
)

type BaseEffects struct {
	*MatchEffects
	*DemurrageEffects
}

type DemurrageEffects struct {
	Asset      string    `json:"asset"`
	Amount     string    `json:"amount"`
	PeriodFrom time.Time `json:"period_from"`
	PeriodTo   time.Time `json:"period_to"`
}

type MatchEffects struct {
	BaseAsset  string         `json:"base_asset"`
	QuoteAsset string         `json:"quote_asset"`
	IsBuy      bool           `json:"is_buy"`
	Matches    []MatchDetails `json:"matches"`
}

type MatchDetails struct {
	BaseAmount  string `json:"base_amount"`
	QuoteAmount string `json:"quote_amount"`
	FeePaid     string `json:"fee_paid"`
	Price       string `json:"price"`
}
