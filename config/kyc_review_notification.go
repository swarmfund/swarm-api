package config

import "html/template"

type KYCReviewNotification struct {
	//*Base
	Template *template.Template
}

//
//func (e *KYCReviewNotification) Init() error {
//	e.Template = e.getTemplate("kyc_review_notification")
//	return nil
//}
