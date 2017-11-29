package response_types

import (
	"encoding/json"
	"fmt"

	"gitlab.com/swarmfund/api/notificator"
)

type ManageInvoice struct {
	OperationBase
	Amount          string  `json:"amount"`
	Asset           string  `json:"asset"`
	ReceiverBalance string  `json:"receiver_balance"`
	Sender          string  `json:"sender"`
	InvoiceID       int64   `json:"invoice_id"`
	RejectReason    *string `json:"reject_reason"`
}

const (
	InvoiceStatePending   int32 = 1
	InvoiceStateSuccess         = 2
	InvoiceStateRejected        = 3
	InvoiceStateCancelled       = 4
	InvoiceStateFailed          = 5
)

var InvoiceStatuses = map[int32]string{
	InvoiceStatePending:   "Pending",
	InvoiceStateSuccess:   "Success",
	InvoiceStateRejected:  "Rejected",
	InvoiceStateCancelled: "Cancelled",
	InvoiceStateFailed:    "Failed",
}

func (p *ManageInvoice) FromBase(base *OperationBase, rawOperation []byte) error {
	err := json.Unmarshal(rawOperation, &p)
	if err != nil {
		return err
	}
	p.OperationBase = *base
	return nil
}

func (p *ManageInvoice) ToLetter(addressee string, isSender bool) (letter *notificator.InvoiceNoticeLetter) {
	letter = new(notificator.InvoiceNoticeLetter)
	letter.Header = "{{ .Project }} | New Request"
	letter.Type = "Payment Request"

	letter.Amount = fmt.Sprintf("%s %s", p.Amount, p.Asset)
	letter.Date = p.LedgerCloseTime.Format(notificator.TimeLayout)

	if p.RejectReason != nil {
		letter.RejectReason = *p.RejectReason
	} else {
		letter.RejectReason = "Reject Reason is not specified"
	}

	// For situation if len(p.Participants) == 0
	letter.Id = fmt.Sprintf("%s;%s;%s", p.ID, addressee, letter.Amount)

	for i, participant := range p.Participants {
		if participant.Email == "" {
			continue
		}

		if participant.AccountID == addressee {
			letter.Id = fmt.Sprintf("%s;%s;%d", p.ID, participant.BalanceID, i)
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

	letter.Status = InvoiceStatuses[p.State]

	var m string

	if isSender {
		letter.CounterpartyType = "Applicant"
		switch p.State {
		case InvoiceStatePending:
			m = "You have received a request for payment from %s."
			break
		case InvoiceStateSuccess:
			m = "A request for payment from %s fulfilled!"
			break
		case InvoiceStateRejected:
			m = "A request for payment from %s decline with reason: " + letter.RejectReason
			break
		case InvoiceStateCancelled:
			m = "%s cancel his request for payment."
			break
		case InvoiceStateFailed:
			m = "A request for payment from %s failed."
			break
		default:
			m = "Request for payment from %s"
		}
	} else {
		letter.CounterpartyType = "Sender"
		switch p.State {
		case InvoiceStatePending:
			m = "You sent a payment request to %s."
			break
		case InvoiceStateSuccess:
			m = "Your payment request to %s fulfilled!"
			break
		case InvoiceStateRejected:
			m = "%s decline your payment request with reason: " + letter.RejectReason
			break
		case InvoiceStateCancelled:
			m = "%s canceled."
			break
		case InvoiceStateFailed:
			m = "Your payment request to %s was failed"
			break
		default:
			m = "Request for payment to %s"
		}
	}

	letter.Message = fmt.Sprintf(m, letter.Counterparty)

	return letter
}
