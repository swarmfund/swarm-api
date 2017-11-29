package config

import "html/template"

type EmailConfirmation struct {
	//*Base
	Template    *template.Template
	RedirectURL string
	ClientURL   string
}
