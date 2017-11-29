package response_types

import (
	"encoding/json"
	"fmt"

	"gitlab.com/swarmfund/api/notificator"
	"gitlab.com/swarmfund/go/amount"
)

type BasePayment struct {
	From                  string `json:"from"`
	To                    string `json:"to"`
	FromBalance           string `json:"from_balance"`
	ToBalance             string `json:"to_balance"`
	Amount                string `json:"amount"`
	UserDetails           string `json:"user_details"`
	Asset                 string `json:"asset"`
	SourcePaymentFee      string `json:"source_payment_fee"`
	DestinationPaymentFee string `json:"destination_payment_fee"`
	SourceFixedFee        string `json:"source_fixed_fee"`
	DestinationFixedFee   string `json:"destination_fixed_fee"`
	SourcePaysForDest     bool   `json:"source_pays_for_dest"`

	AmountI                int64
	SourcePaymentFeeI      int64
	DestinationPaymentFeeI int64
	SourceFixedFeeI        int64
	DestinationFixedFeeI   int64
}

type Payment struct {
	OperationBase
	BasePayment
	Subject   string `json:"subject"`
	Reference string `json:"reference"`
	Asset     string `json:"asset"`
}

func (p *Payment) FromBase(base *OperationBase, rawOperation []byte) error {
	err := json.Unmarshal(rawOperation, &p)
	if err != nil {
		return err
	}
	p.OperationBase = *base

	p.AmountI, err = amount.Parse(p.Amount)
	if err != nil {
		return err
	}

	p.SourcePaymentFeeI, err = amount.Parse(p.SourcePaymentFee)
	if err != nil {
		return err
	}

	p.DestinationPaymentFeeI, err = amount.Parse(p.DestinationPaymentFee)
	if err != nil {
		return err
	}

	p.SourceFixedFeeI, err = amount.Parse(p.SourceFixedFee)
	if err != nil {
		return err
	}

	p.DestinationFixedFeeI, err = amount.Parse(p.DestinationFixedFee)
	if err != nil {
		return err
	}
	return nil
}

func (p *Payment) ToLetter(addressee string, isSender bool) (letter *notificator.PaymentNoticeLetter) {
	letter = new(notificator.PaymentNoticeLetter)
	letter.Header = "{{ .Project }} | New Transaction"

	letter.Amount = fmt.Sprintf("%s %s", p.Amount, p.Asset)
	letter.Date = p.LedgerCloseTime.Format(notificator.TimeLayout)

	// For situation if len(p.Participants) == 0
	letter.Id = fmt.Sprintf("%s;%s;%s", p.ID, addressee, letter.Amount)
	for _, participant := range p.Participants {
		if participant.Email == "" {
			continue
		}

		if participant.AccountID == addressee {
			letter.Id = fmt.Sprintf("%s;%s", p.ID, participant.BalanceID)
			letter.Email = participant.Email

			if participant.FullName != "" {
				letter.Addressee = participant.FullName
			} else {
				letter.Addressee = participant.Email
			}
		} else {
			if participant.Nickname != "" {
				letter.Counterparty = participant.Nickname
			} else if participant.FullName != "" {
				letter.Counterparty = participant.FullName
			} else {
				letter.Counterparty = participant.AccountID
			}
		}
	}

	var fee int64
	if isSender {
		if p.SourcePaysForDest {
			fee = p.SourceFixedFeeI + p.SourcePaymentFeeI + p.DestinationFixedFeeI + p.DestinationPaymentFeeI
		} else {
			fee = p.SourceFixedFeeI + p.SourcePaymentFeeI
		}
		letter.Fee = fmt.Sprintf("%s %s", amount.String(fee), p.Asset)
		letter.FullAmount = fmt.Sprintf("%s %s", amount.String(p.AmountI+fee), p.Asset)
	} else {
		if p.SourcePaysForDest {
			letter.Fee = "Sender paid"
			letter.FullAmount = letter.Amount
		} else {
			fee = p.DestinationFixedFeeI + p.DestinationPaymentFeeI
			letter.Fee = fmt.Sprintf("%s %s", amount.String(fee), p.Asset)
			letter.FullAmount = fmt.Sprintf("%s %s", amount.String(p.AmountI-fee), p.Asset)
		}
	}

	sjPrefix := "No message"

	if len(p.Subject) > 4 {
		sjPrefix = p.Subject[0:4]
		letter.Reference = p.Subject[4:]
	} else {
		letter.Reference = p.Subject
	}

	switch sjPrefix {
	case "gf: ":
		letter.Type = "Gift"
	case "in: ":
		letter.Type = "Invoice"
	case "tf: ":
		letter.Type = "Transfer"
	default:
		letter.Type = "Transfer"
	}

	if isSender {
		letter.Action = "Paid"
		letter.CounterpartyType = "Receiver"
		letter.Message = fmt.Sprintf("You sent a %s to %s", letter.Amount, letter.Counterparty)
	} else {
		letter.Action = "Received"
		letter.CounterpartyType = "Sender"
		letter.Message = fmt.Sprintf("You received a %s from %s", letter.Amount, letter.Counterparty)
	}

	return letter
}
