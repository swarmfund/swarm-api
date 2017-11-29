package config

import "html/template"

type LoginNotification struct {
	//*Base
	Template *template.Template
}

//func (e *LoginNotification) Init() error {
//	e.Template = e.getTemplate("login_notification")
//	return nil
//}
