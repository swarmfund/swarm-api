package skrill

import (
	"strings"
)

func (c *Client) QuickCheckoutSession(params *QuickCheckoutParams) (string, error) {
	params.PrepareOnly = "1"
	body, err := c.post(skrillPayEndpoint, params)
	return string(body), err
}

type QuickCheckoutParams struct {
	PrepareOnly string `json:"prepare_only"`
	// required fields
	PayToEmail string `json:"pay_to_email"`
	Amount     string `json:"amount"`
	Currency   string `json:"currency"`
	// optional
	StatusURL     string `json:"status_url"`
	TransactionID string `json:"transaction_id"`
	ReturnURL     string `json:"return_url"`
	CancelURL     string `json:"cancel_url"`

	CustomFields map[string]string `json:"-"`
}

func (p *QuickCheckoutParams) MerchantFields() map[string]string {
	delete(p.CustomFields, "merchant_fields")
	keys := make([]string, len(p.CustomFields))
	for key, _ := range p.CustomFields {
		keys = append(keys, key)
	}
	p.CustomFields["merchant_fields"] = strings.Join(keys, ",")

	return p.CustomFields
}
