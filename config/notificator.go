package config

import (
	"html/template"

	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/swarmfund/api/assets"
)

const (
	notificatorConfigKey = "notificator"
)

var (
	notificatorConfig *Notificator
)

type Notificator struct {
	Endpoint     string
	Secret       string
	Public       string
	ClientRouter string

	EmailConfirmation *template.Template
	//KYCApproval            KYCApproval
	//RecoveryRequest        RecoveryRequest
	//LoginNotification      LoginNotification
	//OperationsNotification OperationsNotification
	//KYCReviewNotification  KYCReviewNotification
}

func (c *ViperConfig) Notificator() Notificator {
	if notificatorConfig == nil {
		notificatorConfig = &Notificator{}
		config := c.GetStringMap(notificatorConfigKey)
		if err := figure.Out(notificatorConfig).From(config).Please(); err != nil {
			panic(errors.Wrap(err, "failed to figure out notificator"))
		}
		notificatorConfig.EmailConfirmation = assets.Templates.Lookup("email_confirm")
	}
	return *notificatorConfig
}

//
//func (n *Notificator) DefineConfigStructure() {
//	n.bindEnv("endpoint")
//	n.bindEnv("secret")
//	n.bindEnv("public")
//	n.EmailConfirmation.Base = NewBase(n.Base, "email_confirm")
//	n.EmailConfirmation.DefineConfigStructure()
//
//	n.KYCApproval.Base = NewBase(n.Base, "kyc_approval")
//	n.RecoveryRequest = NewRecoveryRequest(NewBase(n.Base, "recovery_request"))
//	n.LoginNotification.Base = NewBase(n.Base, "login_notification")
//
//	n.OperationsNotification.Base = NewBase(n.Base, "op_notifications")
//	n.OperationsNotification.DefineConfigStructure()
//
//	n.KYCReviewNotification.Base = NewBase(n.Base, "kyc_review_notification")
//}
//
//func (n *Notificator) Init() error {
//	var err error
//	n.Endpoint, err = n.getURLAsString("endpoint")
//	if err != nil {
//		return err
//	}
//
//	n.Secret, err = n.getNonEmptyString("secret")
//	if err != nil {
//		return err
//	}
//
//	n.Public, err = n.getNonEmptyString("public")
//	if err != nil {
//		return err
//	}
//
//	err = n.EmailConfirmation.Init()
//	if err != nil {
//		return err
//	}
//
//	err = n.KYCApproval.Init()
//	if err != nil {
//		return err
//	}
//
//	err = n.RecoveryRequest.Init()
//	if err != nil {
//		return err
//	}
//
//	err = n.LoginNotification.Init()
//	if err != nil {
//		return err
//	}
//
//	err = n.OperationsNotification.Init()
//	if err != nil {
//		return err
//	}
//
//	if err = n.KYCReviewNotification.Init(); err != nil {
//		return err
//	}
//
//	return nil
//}
