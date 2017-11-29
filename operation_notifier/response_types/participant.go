package response_types

import (
	"errors"
	"fmt"
	"time"

	"gitlab.com/swarmfund/api/db2/api"
	"gitlab.com/swarmfund/api/notificator"
	"gitlab.com/swarmfund/go/amount"
)

type Participant struct {
	AccountID string      `json:"account_id,omitempty"`
	BalanceID string      `json:"balance_id,omitempty"`
	Email     string      `json:"email,omitempty"`
	FullName  string      `json:"full_name,omitempty"`
	Nickname  string      `json:"nickname,omitempty"`
	Effects   BaseEffects `json:"effects,omitempty"`
}

func (p *Participant) ToApiParticipant() (ap api.Participant) {
	ap.AccountID = p.AccountID
	ap.BalanceID = p.BalanceID
	ap.Nickname = p.Nickname
	ap.Email = &p.Email
	return
}

func (p *Participant) FromApiParticipant(ap *api.Participant) {
	p.AccountID = ap.AccountID
	p.BalanceID = ap.BalanceID
	p.Nickname = ap.Nickname

	if ap.Email != nil {
		p.Email = *ap.Email
	}

	if ap.Name != nil {
		p.FullName = *ap.Name
		return
	}

	name := ap.Details.DisplayName()
	if name != nil {
		p.FullName = *name
	}
}

func (p *Participant) ToDemurrageLetter(opId string, date time.Time) (letter *notificator.DemurrageNoticeLetter, err error) {
	if p.Effects.DemurrageEffects == nil {
		return nil, errors.New("Demurrage effect empty")
	}
	letter = new(notificator.DemurrageNoticeLetter)

	effects := p.Effects.DemurrageEffects
	demurrageAmount, err := amount.Parse(effects.Amount)
	if err != nil {
		return nil, err
	}

	if demurrageAmount == 0 {
		// Doesn't need to send letter of zero demurrage
		// not initialized letter will be ignored by notificator
		return letter, nil
	}

	letter.Id = fmt.Sprintf("%s;%s", opId, p.BalanceID)
	letter.Header = "{{ .Project }} | Demurrage"
	letter.Type = "Demurrage"
	letter.Email = p.Email
	letter.Asset = effects.Asset
	letter.Amount = fmt.Sprintf("%s %s", effects.Amount, effects.Asset)
	letter.Message = fmt.Sprintf("The commission for the storage of the %s was automatically withdrawn from Your account.", effects.Asset)

	if p.FullName != "" {
		letter.Addressee = p.FullName
	} else {
		letter.Addressee = p.Email
	}

	letter.PeriodTo = effects.PeriodTo.Format(notificator.TimeLayout)
	letter.PeriodFrom = effects.PeriodFrom.Format(notificator.TimeLayout)
	letter.Date = date.Format(notificator.TimeLayout)

	return letter, nil
}

func (p *Participant) ToOfferLetters(opId string, date time.Time) (lts []notificator.OfferNoticeLetter, err error) {
	if p.Effects.MatchEffects == nil {
		return nil, errors.New("Match effects empty")
	}
	effects := p.Effects.MatchEffects

	var letter = notificator.OfferNoticeLetter{}
	letter.Type = "Trade"
	letter.Email = p.Email
	letter.Header = "{{ .Project }} | New Match"
	letter.Date = date.Format(notificator.TimeLayout)

	if p.FullName != "" {
		letter.Addressee = p.FullName
	} else {
		letter.Addressee = p.Email
	}

	letters := make([]notificator.OfferNoticeLetter, len(effects.Matches))
	for i, m := range effects.Matches {
		letters[i].Id = fmt.Sprintf("%s;%s;%d", opId, p.BalanceID, i)
		letters[i].Type = "Trade"
		letters[i].Addressee = letter.Addressee
		letters[i].Date = letter.Date
		letters[i].Email = letter.Email
		letters[i].Header = letter.Header

		letters[i].Amount = fmt.Sprintf("%s %s", m.BaseAmount, effects.BaseAsset)
		letters[i].Fee = fmt.Sprintf("%s %s", m.FeePaid, effects.QuoteAsset)
		letters[i].Price = fmt.Sprintf("%s %s", m.Price, effects.QuoteAsset)
		letters[i].QuoteAmount = fmt.Sprintf("%s %s", m.QuoteAmount, effects.QuoteAsset)

		qa, err := amount.Parse(m.QuoteAmount)
		if err != nil {
			return lts, err
		}
		fee, err := amount.Parse(m.FeePaid)
		if err != nil {
			return lts, err
		}

		var action string
		if effects.IsBuy {
			action = "buy"
			letters[i].OrderPrice = fmt.Sprintf("%s %s", amount.String(qa+fee), effects.QuoteAsset)
		} else {
			action = "sell"
			letters[i].OrderPrice = fmt.Sprintf("%s %s", amount.String(qa-fee), effects.QuoteAsset)
		}

		letters[i].Message = fmt.Sprintf("Your order for %s %s at a price of %s is fulfilled.",
			action, letters[i].Amount, letters[i].Price)
	}

	return letters, nil
}
