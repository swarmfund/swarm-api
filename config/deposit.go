package config

//type Deposit struct {
//	*Base
//
//	Merchant string
//	Currency string
//	Asset    string
//
//	StripeChargeEndpoint *url.URL
//	StripePK             string
//}
//
//func (n *Deposit) DefineConfigStructure() {
//	n.bindEnv("merchant")
//	n.bindEnv("currency")
//	n.bindEnv("asset")
//
//	n.bindEnv("stripe_charge_endpoint")
//	n.bindEnv("stripe_pk")
//}
//
//func (n *Deposit) Init() error {
//	var err error
//
//	n.Merchant, err = n.getNonEmptyString("merchant")
//	if err != nil {
//		return err
//	}
//
//	n.Currency, err = n.getNonEmptyString("currency")
//	if err != nil {
//		return err
//	}
//
//	n.Asset, err = n.getNonEmptyString("asset")
//	if err != nil {
//		return err
//	}
//
//	n.StripeChargeEndpoint, err = n.getParsedURL("stripe_charge_endpoint")
//	if err != nil {
//		return err
//	}
//
//	n.StripePK, err = n.getNonEmptyString("stripe_pk")
//	if err != nil {
//		return err
//	}
//
//	return err
//}
