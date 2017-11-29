package api

import (
	"errors"

	"gitlab.com/distributed_lab/skrill-go"
	"gitlab.com/swarmfund/api/render/hal"
	"gitlab.com/swarmfund/api/render/problem"
	"gitlab.com/swarmfund/horizon-connector"
)

type DepositMethod int32

const (
	_ DepositMethod = iota
	DepositMethodStripe
	DepositMethodSkrill
)

type DepositRequest struct {
	Receiver string               `json:"receiver" valid:"required"`
	Method   DepositMethod        `json:"method" valid:"required"`
	Amount   string               `json:"amount" valid:"required"`
	Stripe   StripeRequestDetails `json:"stripe" valid:"optional"`
}

type StripeRequestDetails struct {
	Token string `json:"token" valid:"required"`
}

type ChargeRequest struct {
	Token     string `json:"token"`
	Amount    string `json:"amount"`
	Reference string `json:"reference"`
	Receiver  string `json:"receiver"`
	Asset     string `json:"asset"`
}

type DepositAction struct {
	Action

	Request  DepositRequest
	Resource interface{}
}

func (action *DepositAction) JSON() {
	action.Do(
		action.ValidateBodyType,
		action.loadParams,
		action.validateReceiver,
		action.performRequest,
		func() {
			hal.Render(action.W, action.Resource)
		},
	)
}

func (action *DepositAction) loadParams() {
	action.UnmarshalBody(&action.Request)
	value, err := skrill.ParseAmount(action.Request.Amount, 4)
	if err != nil {
		action.SetInvalidField("amount", err)
		return
	}
	if value <= 5000 {
		action.SetInvalidField("amount", errors.New("too low"))
		return
	}
}

func (action *DepositAction) validateReceiver() {
	balance, err := horizon.ParseBalanceID(action.Request.Receiver)
	if err != nil {
		action.SetInvalidField("receiver", err)
		return
	}

	asset, err := action.App.horizon.BalanceAsset(balance.AsString())
	if err != nil {
		action.Log.WithError(err).Error("failed to get balance asset")
		action.Err = &problem.ServerError
		return
	}

	if asset == nil {
		action.SetInvalidField("receiver", errors.New("balance does not exists"))
		return
	}

	//if asset.Code != action.App.config.Deposit.Asset {
	//	action.SetInvalidField("receiver", fmt.Errorf("expected %s balance", action.App.config.Deposit.Asset))
	//	return
	//}
}

func (action *DepositAction) performRequest() {
	//reference, err := keypair.Random()
	//if err != nil {
	//	action.Log.WithError(err).Error("failed to generate reference")
	//	action.Err = &problem.ServerError
	//	return
	//}

	switch action.Request.Method {
	case DepositMethodSkrill:
		// craft return URL

		//returnURL := *action.App.config.ClientRouter
		//encodedPayload, err := redirect.SkrillDepositSuccess.Encode()
		//if err != nil {
		//	action.Log.WithError(err).Error("failed to encode return payload")
		//	action.Err = &problem.ServerError
		//	return
		//}
		//
		//query := returnURL.Query()
		//query.Set("action", encodedPayload)
		//returnURL.RawQuery = query.Encode()

		// craft cancel URL

		//cancelURL := *action.App.config.ClientRouter
		//encodedPayload, err = redirect.SkrillDepositSuccess.Encode()
		//if err != nil {
		//	action.Log.WithError(err).Error("failed to encode cancel payload")
		//	action.Err = &problem.ServerError
		//	return
		//}
		//
		//query = cancelURL.Query()
		//query.Set("action", encodedPayload)
		//cancelURL.RawQuery = query.Encode()

		// quick checkout session

		//sid, err := skrill.NewClient().QuickCheckoutSession(&skrill.QuickCheckoutParams{
		//	PayToEmail:    action.App.config.Deposit.Merchant,
		//	Amount:        action.Request.Amount,
		//	Currency:      action.App.config.Deposit.Currency,
		//	TransactionID: reference.Address(),
		//	ReturnURL:     returnURL.String(),
		//	CancelURL:     cancelURL.String(),
		//	CustomFields: map[string]string{
		//		"x-receiver": action.Request.Receiver,
		//		"x-asset":    action.App.config.Deposit.Asset,
		//	},
		//})

		//if err != nil {
		//	action.Log.WithError(err).Error("failed to get checkout session")
		//	action.Err = &problem.ServerError
		//	return
		//}

		//action.Resource = resource.Deposit{
		//	SID: sid,
		//}
	case DepositMethodStripe:
		//payload := ChargeRequest{
		//	Token:     action.Request.Stripe.Token,
		//	Amount:    action.Request.Amount,
		//	Reference: reference.Address(),
		//	Receiver:  action.Request.Receiver,
		//	Asset:     action.App.config.Deposit.Asset,
		//}
		//body, err := json.Marshal(&payload)
		//if err != nil {
		//	action.Log.WithError(err).Error("failed to marshal strip payload")
		//	action.Err = &problem.ServerError
		//	return
		//}
		//response, err := http.DefaultClient.Post(action.App.config.Deposit.StripeChargeEndpoint.String(), "application/json", bytes.NewReader(body))
		//if err != nil {
		//	action.Log.WithError(err).Error("failed to create request for stripe charger")
		//	action.Err = &problem.ServerError
		//	return
		//}
		//defer response.Body.Close()

		//switch response.StatusCode {
		//case 200:
		//	action.Err = &problem.Success
		//	return
		//case 400:
		//	action.SetInvalidField("token", errors.New("unable to complete payment"))
		//	return
		//default:
		//	action.Log.WithField("status", response.StatusCode).Error("failed to submit stripe charge")
		//	action.Err = &problem.ServerError
		//	return
		//}
	default:
		action.SetInvalidField("method", errors.New("unknown method"))
		return
	}
}
