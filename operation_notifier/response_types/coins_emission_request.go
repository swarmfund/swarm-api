package response_types

import (
	"encoding/json"
	"fmt"

	"gitlab.com/swarmfund/api/notificator"
)

type ManageCoinsEmissionRequest struct {
	OperationBase
	RequestID int64  `json:"request_id"`
	Amount    string `json:"amount"`
	Asset     string `json:"asset"`
}

type ReviewCoinsEmissionRequest struct {
	OperationBase
	RequestID  uint64  `json:"request_id"`
	Amount     string  `json:"amount"`
	Asset      string  `json:"asset"`
	IsApproved *bool   `json:"approved"`
	Issuer     string  `json:"issuer"`
	Reason     *string `json:"reason"`
}

func (p *ReviewCoinsEmissionRequest) FromBase(base *OperationBase, rawOperation []byte) error {
	err := json.Unmarshal(rawOperation, &p)
	if err != nil {
		return err
	}
	p.OperationBase = *base
	return nil
}

func (p *ReviewCoinsEmissionRequest) ToLetter(isManage bool) (letter *notificator.CoinsEmissionNoticeLetter) {
	letter = new(notificator.CoinsEmissionNoticeLetter)
	letter.Header = "{{ .Project }} | New Deposit"

	letter.Type = "Deposit"
	letter.Amount = fmt.Sprintf("%s %s", p.Amount, p.Asset)
	letter.Date = p.LedgerCloseTime.Format(notificator.TimeLayout)

	if p.Reason != nil {
		letter.RejectReason = *p.Reason
	} else {
		letter.RejectReason = "Reason is not specified"
	}

	for _, participant := range p.Participants {
		if participant.Email == "" {
			continue
		}

		letter.Id = fmt.Sprintf("%s;%s", p.ID, participant.BalanceID)
		letter.Email = participant.Email

		if participant.FullName != "" {
			letter.Addressee = participant.FullName
		} else {
			letter.Addressee = participant.Email
		}
		break
	}

	var m string

	if isManage {
		letter.Status = OpBaseStatuses[p.State]
		var b bool
		switch p.State {
		case OpBaseStatePending:
			p.IsApproved = nil
			break
		case OpBaseStateSuccess:
			b = true
			p.IsApproved = &b
			break
		case OpBaseStateFailed:
			m = "A deposit of %s failed."
			break
		case OpBaseStateRejected:
			b = false
			p.IsApproved = &b
			break
		}
	}

	if p.IsApproved == nil {
		letter.Status = "Pending"
		m = "A new deposit of %s has been requested to your account."
	} else if *p.IsApproved {
		letter.Status = "Success"
		m = "A new deposit of %s has been credited to your account."
	} else {
		letter.Status = "Rejected"
		m = "A deposit of %s has been rejected with reason: " + letter.RejectReason
	}

	letter.Message = fmt.Sprintf(m, letter.Amount)
	return letter

}
