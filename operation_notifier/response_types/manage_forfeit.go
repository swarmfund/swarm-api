package response_types

import (
	"encoding/json"
	"fmt"

	"gitlab.com/swarmfund/api/notificator"
)

type ManageForfeitRequest struct {
	OperationBase
	Action      int32  `json:"action"`
	RequestID   uint64 `json:"request_id"`
	Amount      string `json:"amount"`
	Asset       string `json:"asset"`
	UserDetails string `json:"user_details"`
}

type ReviewForfeitRequest struct {
	OperationBase
	RequestID uint64 `json:"request_id"`
	Balance   string `json:"balance"`
	Accept    bool   `json:"accept"`
	Amount    string `json:"amount"`
	Asset     string `json:"asset"`
}

func (p *ManageForfeitRequest) FromBase(base *OperationBase, rawOperation []byte) error {
	err := json.Unmarshal(rawOperation, &p)
	if err != nil {
		return err
	}
	p.OperationBase = *base
	return nil
}

func (p *ManageForfeitRequest) ToLetter() (letter *notificator.ForfeitNoticeLetter) {
	letter = new(notificator.ForfeitNoticeLetter)

	if p.Asset == "USD" {
		letter.Type = "Withdraw"

	} else {
		letter.Type = "Redemption"
	}

	letter.Header = "{{ .Project }} | New " + letter.Type
	letter.Amount = fmt.Sprintf("%s %s", p.Amount, p.Asset)
	letter.Date = p.LedgerCloseTime.Format(notificator.TimeLayout)

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
	letter.Status = OpBaseStatuses[p.State]

	switch p.State {
	case OpBaseStatePending:
		m = "%s request of %s created."
		break
	case OpBaseStateSuccess:
		m = "%s request of %s successfully done."
		break
	case OpBaseStateFailed:
		m = "%s request of %s failed."
		break
	case OpBaseStateRejected:
		m = "%s request of %s was declined."
		break
	default:
		m = "%s request of %s."
	}

	letter.Message = fmt.Sprintf(m, letter.Type, letter.Amount)
	return letter

}
