package config

import (
	"html/template"
)

type OperationsNotification struct {
	//*Base
	Disable    bool
	LastCursor string

	CoinsEmission *template.Template
	Demurrage     *template.Template
	Forfeit       *template.Template
	Invoice       *template.Template
	Offer         *template.Template
	Payment       *template.Template
}

//func (e *OperationsNotification) DefineConfigStructure() {
//	e.setDefault("disable", false)
//
//	e.bindEnv("disable")
//	e.bindEnv("last_cursor")
//}
//
//func (e *OperationsNotification) Init() error {
//	e.Disable = e.getBool("disable")
//	e.LastCursor = e.getString("last_cursor")
//
//	e.CoinsEmission = e.getTemplate("operations_notifications/manage_coins_emission")
//	e.Demurrage = e.getTemplate("operations_notifications/demurrage")
//	e.Forfeit = e.getTemplate("operations_notifications/manage_forfeit")
//	e.Invoice = e.getTemplate("operations_notifications/manage_invoice")
//	e.Offer = e.getTemplate("operations_notifications/manage_offer")
//	e.Payment = e.getTemplate("operations_notifications/payment")
//	return nil
//}
