package config

import "html/template"

type RecoveryRequest struct {
	//*Base
	Template    *template.Template
	RedirectURL string
	ClientURL   string
}

//
//func NewRecoveryRequest(base *Base) RecoveryRequest {
//	rr := RecoveryRequest{
//		Base: base,
//	}
//	rr.DefineStructure()
//	return rr
//}
//
//func (c *RecoveryRequest) DefineStructure() {
//	c.bindEnv("redirect_url")
//	c.bindEnv("client_url")
//}
//
//func (c *RecoveryRequest) Init() error {
//	var err error
//
//	c.Template = c.getTemplate("recovery_request")
//
//	c.RedirectURL, err = c.getURLAsString("redirect_url")
//	if err != nil {
//		return err
//	}
//
//	c.ClientURL, err = c.getURLAsString("client_url")
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
